package servermanager

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"api/internal/app/ports"
)

// Client implements ports.ServerManagerClient over the manager's internal
// HTTP API (bearer token, private network).
type Client struct {
	baseURL string
	token   string
	http    *http.Client
	// streaming requests live as long as the log follow does — a global
	// Timeout would cut the tail off, so they get their own client.
	streamHTTP *http.Client
}

var _ ports.ServerManagerClient = (*Client)(nil)

func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		token:      token,
		http:       &http.Client{Timeout: 10 * time.Second},
		streamHTTP: &http.Client{},
	}
}

func (c *Client) Logs(ctx context.Context, appID string, opts ports.LogOptions) ([]ports.LogEntry, error) {
	return c.logsAt(ctx, "/v1/apps/"+url.PathEscape(appID)+"/logs", opts, ports.ErrAppNotDeployed)
}

// DatabaseLogs mirrors Logs for the database container.
func (c *Client) DatabaseLogs(ctx context.Context, dbID string, opts ports.LogOptions) ([]ports.LogEntry, error) {
	return c.logsAt(ctx, "/v1/dbs/"+url.PathEscape(dbID)+"/logs", opts, ports.ErrDatabaseNotProvisioned)
}

// logsAt reads one log endpoint; apps and databases share the response
// contract and differ only in path and 404 meaning.
func (c *Client) logsAt(ctx context.Context, path string, opts ports.LogOptions, notFound error) ([]ports.LogEntry, error) {
	endpoint := c.baseURL + path
	if query := logQuery(opts); len(query) > 0 {
		endpoint += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("servermanager logs request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	res, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("servermanager logs: %w", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		var body struct {
			Logs []ports.LogEntry `json:"logs"`
		}
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return nil, fmt.Errorf("servermanager logs response: %w", err)
		}
		return body.Logs, nil
	case http.StatusNotFound:
		return nil, fmt.Errorf("servermanager logs: %w", notFound)
	case http.StatusUnauthorized:
		// Operator-facing: the shared token doesn't match. Never echo it.
		return nil, fmt.Errorf("servermanager logs: bearer token rejected (check SERVERMANAGER_TOKEN/SM_TOKEN)")
	default:
		return nil, fmt.Errorf("servermanager logs: status %d: %s", res.StatusCode, managerError(res))
	}
}

func logQuery(opts ports.LogOptions) url.Values {
	query := url.Values{}
	if opts.Tail > 0 {
		query.Set("tail", strconv.Itoa(opts.Tail))
	}
	if !opts.Since.IsZero() {
		query.Set("since", opts.Since.Format(time.RFC3339))
	}
	return query
}

// StreamLogs follows the manager's NDJSON log stream and re-emits each line
// as a LogEntry. The channel closes when the manager ends the stream or ctx
// is canceled (which also closes the underlying response body).
func (c *Client) StreamLogs(ctx context.Context, appID string, opts ports.LogOptions) (<-chan ports.LogEntry, error) {
	return c.streamAt(ctx, "/v1/apps/"+url.PathEscape(appID)+"/logs/stream", opts, ports.ErrAppNotDeployed)
}

// StreamDatabaseLogs mirrors StreamLogs for the database container.
func (c *Client) StreamDatabaseLogs(ctx context.Context, dbID string, opts ports.LogOptions) (<-chan ports.LogEntry, error) {
	return c.streamAt(ctx, "/v1/dbs/"+url.PathEscape(dbID)+"/logs/stream", opts, ports.ErrDatabaseNotProvisioned)
}

func (c *Client) streamAt(ctx context.Context, path string, opts ports.LogOptions, notFound error) (<-chan ports.LogEntry, error) {
	endpoint := c.baseURL + path
	if query := logQuery(opts); len(query) > 0 {
		endpoint += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("servermanager stream request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	res, err := c.streamHTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("servermanager stream: %w", err)
	}

	switch res.StatusCode {
	case http.StatusOK:
		// fall through to streaming below
	case http.StatusNotFound:
		closeBody(res)
		return nil, fmt.Errorf("servermanager stream: %w", notFound)
	case http.StatusUnauthorized:
		closeBody(res)
		return nil, fmt.Errorf("servermanager stream: bearer token rejected (check SERVERMANAGER_TOKEN/SM_TOKEN)")
	default:
		defer closeBody(res)
		return nil, fmt.Errorf("servermanager stream: status %d: %s", res.StatusCode, managerError(res))
	}

	entries := make(chan ports.LogEntry)
	go func() {
		defer close(entries)
		defer closeBody(res)

		scanner := bufio.NewScanner(res.Body)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for scanner.Scan() {
			line := scanner.Bytes()
			if len(line) == 0 {
				continue
			}
			var entry ports.LogEntry
			if err := json.Unmarshal(line, &entry); err != nil {
				continue // skip torn/malformed lines rather than killing the tail
			}
			select {
			case entries <- entry:
			case <-ctx.Done():
				return
			}
		}
	}()
	return entries, nil
}

// Metrics reads the app's sampled series from the manager.
func (c *Client) Metrics(ctx context.Context, appID string, rng string) (*ports.MetricsResponse, error) {
	return c.metricsAt(ctx, "/v1/apps/"+url.PathEscape(appID)+"/metrics", rng, ports.ErrAppNotDeployed)
}

// DatabaseMetrics mirrors Metrics for the database container.
func (c *Client) DatabaseMetrics(ctx context.Context, dbID string, rng string) (*ports.MetricsResponse, error) {
	return c.metricsAt(ctx, "/v1/dbs/"+url.PathEscape(dbID)+"/metrics", rng, ports.ErrDatabaseNotProvisioned)
}

// metricsAt reads one metrics endpoint; apps and databases share the response
// envelope and differ only in path and 404 meaning.
func (c *Client) metricsAt(ctx context.Context, path, rng string, notFound error) (*ports.MetricsResponse, error) {
	if rng != "" {
		query := url.Values{}
		query.Set("range", rng)
		path += "?" + query.Encode()
	}

	res, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("servermanager metrics: %w", err)
	}
	defer closeBody(res)

	switch res.StatusCode {
	case http.StatusOK:
		var body ports.MetricsResponse
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return nil, fmt.Errorf("servermanager metrics response: %w", err)
		}
		return &body, nil
	case http.StatusNotFound:
		return nil, fmt.Errorf("servermanager metrics: %w", notFound)
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("servermanager metrics: bearer token rejected (check SERVERMANAGER_TOKEN/SM_TOKEN)")
	default:
		return nil, fmt.Errorf("servermanager metrics: status %d: %s", res.StatusCode, managerError(res))
	}
}

// Deploy queues the async job on the manager and returns its deployment ID.
func (c *Client) Deploy(ctx context.Context, appID string, spec ports.DeploySpec) (string, error) {
	payload, err := json.Marshal(spec)
	if err != nil {
		return "", fmt.Errorf("servermanager deploy payload: %w", err)
	}

	res, err := c.do(ctx, http.MethodPost, "/v1/apps/"+url.PathEscape(appID)+"/deploy", payload)
	if err != nil {
		return "", fmt.Errorf("servermanager deploy: %w", err)
	}
	defer closeBody(res)

	switch res.StatusCode {
	case http.StatusAccepted:
		var body struct {
			DeploymentID string `json:"deployment_id"`
		}
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil || body.DeploymentID == "" {
			return "", fmt.Errorf("servermanager deploy response: %w", err)
		}
		return body.DeploymentID, nil
	case http.StatusConflict:
		return "", fmt.Errorf("servermanager deploy: %w", ports.ErrDeployInFlight)
	case http.StatusBadRequest:
		return "", fmt.Errorf("%w: %s", ports.ErrSpecRejected, managerError(res))
	case http.StatusUnauthorized:
		return "", fmt.Errorf("servermanager deploy: bearer token rejected (check SERVERMANAGER_TOKEN/SM_TOKEN)")
	default:
		return "", fmt.Errorf("servermanager deploy: status %d: %s", res.StatusCode, managerError(res))
	}
}

// DeploymentStatus polls the manager's job resource.
func (c *Client) DeploymentStatus(ctx context.Context, deploymentID string) (*ports.DeploymentStatus, error) {
	res, err := c.do(ctx, http.MethodGet, "/v1/deployments/"+url.PathEscape(deploymentID), nil)
	if err != nil {
		return nil, fmt.Errorf("servermanager deployment status: %w", err)
	}
	defer closeBody(res)

	switch res.StatusCode {
	case http.StatusOK:
		var status ports.DeploymentStatus
		if err := json.NewDecoder(res.Body).Decode(&status); err != nil {
			return nil, fmt.Errorf("servermanager deployment status response: %w", err)
		}
		return &status, nil
	case http.StatusNotFound:
		return nil, fmt.Errorf("servermanager deployment status: %w", ports.ErrDeploymentNotFound)
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("servermanager deployment status: bearer token rejected (check SERVERMANAGER_TOKEN/SM_TOKEN)")
	default:
		return nil, fmt.Errorf("servermanager deployment status: status %d: %s", res.StatusCode, managerError(res))
	}
}

func (c *Client) Start(ctx context.Context, appID string) error {
	return c.lifecycle(ctx, http.MethodPost, appID, "/start", "start")
}

func (c *Client) Stop(ctx context.Context, appID string) error {
	return c.lifecycle(ctx, http.MethodPost, appID, "/stop", "stop")
}

func (c *Client) Remove(ctx context.Context, appID string) error {
	return c.lifecycle(ctx, http.MethodDelete, appID, "", "remove")
}

// lifecycle covers the body-less container operations that share one
// response contract: 200 ok, 404 not deployed.
func (c *Client) lifecycle(ctx context.Context, method, appID, suffix, op string) error {
	res, err := c.do(ctx, method, "/v1/apps/"+url.PathEscape(appID)+suffix, nil)
	if err != nil {
		return fmt.Errorf("servermanager %s: %w", op, err)
	}
	defer closeBody(res)

	switch res.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("servermanager %s: %w", op, ports.ErrAppNotDeployed)
	case http.StatusUnauthorized:
		return fmt.Errorf("servermanager %s: bearer token rejected (check SERVERMANAGER_TOKEN/SM_TOKEN)", op)
	default:
		return fmt.Errorf("servermanager %s: status %d: %s", op, res.StatusCode, managerError(res))
	}
}

// ProvisionDatabase queues the async provision job on the manager.
func (c *Client) ProvisionDatabase(ctx context.Context, dbID string, spec ports.DBProvisionSpec) (string, error) {
	payload, err := json.Marshal(spec)
	if err != nil {
		return "", fmt.Errorf("servermanager provision payload: %w", err)
	}

	res, err := c.do(ctx, http.MethodPost, "/v1/dbs/"+url.PathEscape(dbID)+"/provision", payload)
	if err != nil {
		return "", fmt.Errorf("servermanager provision: %w", err)
	}
	defer closeBody(res)

	switch res.StatusCode {
	case http.StatusAccepted:
		var body struct {
			ProvisionID string `json:"provision_id"`
		}
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil || body.ProvisionID == "" {
			return "", fmt.Errorf("servermanager provision response: %w", err)
		}
		return body.ProvisionID, nil
	case http.StatusConflict:
		return "", fmt.Errorf("servermanager provision: %w", ports.ErrProvisionInFlight)
	case http.StatusBadRequest:
		return "", fmt.Errorf("%w: %s", ports.ErrSpecRejected, managerError(res))
	case http.StatusUnauthorized:
		return "", fmt.Errorf("servermanager provision: bearer token rejected (check SERVERMANAGER_TOKEN/SM_TOKEN)")
	default:
		return "", fmt.Errorf("servermanager provision: status %d: %s", res.StatusCode, managerError(res))
	}
}

// DatabaseStatus polls the manager's database resource.
func (c *Client) DatabaseStatus(ctx context.Context, dbID string) (*ports.DBStatus, error) {
	res, err := c.do(ctx, http.MethodGet, "/v1/dbs/"+url.PathEscape(dbID), nil)
	if err != nil {
		return nil, fmt.Errorf("servermanager database status: %w", err)
	}
	defer closeBody(res)

	switch res.StatusCode {
	case http.StatusOK:
		var status ports.DBStatus
		if err := json.NewDecoder(res.Body).Decode(&status); err != nil {
			return nil, fmt.Errorf("servermanager database status response: %w", err)
		}
		return &status, nil
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("servermanager database status: bearer token rejected (check SERVERMANAGER_TOKEN/SM_TOKEN)")
	default:
		return nil, fmt.Errorf("servermanager database status: status %d: %s", res.StatusCode, managerError(res))
	}
}

// RemoveDatabase sweeps the database's container, network, and volume. The
// manager treats an already-gone database as success, so there is no 404 case.
func (c *Client) RemoveDatabase(ctx context.Context, dbID string) error {
	res, err := c.do(ctx, http.MethodDelete, "/v1/dbs/"+url.PathEscape(dbID), nil)
	if err != nil {
		return fmt.Errorf("servermanager remove database: %w", err)
	}
	defer closeBody(res)

	switch res.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusConflict:
		return fmt.Errorf("servermanager remove database: %w", ports.ErrProvisionInFlight)
	case http.StatusUnauthorized:
		return fmt.Errorf("servermanager remove database: bearer token rejected (check SERVERMANAGER_TOKEN/SM_TOKEN)")
	default:
		return fmt.Errorf("servermanager remove database: status %d: %s", res.StatusCode, managerError(res))
	}
}

// do issues an authenticated JSON request on the short-timeout client.
func (c *Client) do(ctx context.Context, method, path string, payload []byte) (*http.Response, error) {
	var body *bytes.Reader
	if payload != nil {
		body = bytes.NewReader(payload)
	} else {
		body = bytes.NewReader(nil)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.http.Do(req)
}

func closeBody(res *http.Response) {
	_ = res.Body.Close()
}

// managerError extracts the manager's {"error": "..."} body, if any.
func managerError(res *http.Response) string {
	var body struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil || body.Error == "" {
		return "no detail"
	}
	return body.Error
}
