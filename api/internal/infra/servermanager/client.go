package servermanager

import (
	"bufio"
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
	query := url.Values{}
	if opts.Tail > 0 {
		query.Set("tail", strconv.Itoa(opts.Tail))
	}
	if !opts.Since.IsZero() {
		query.Set("since", opts.Since.Format(time.RFC3339))
	}

	endpoint := c.baseURL + "/v1/apps/" + url.PathEscape(appID) + "/logs"
	if len(query) > 0 {
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
		return nil, fmt.Errorf("servermanager logs: %w", ports.ErrAppNotDeployed)
	case http.StatusUnauthorized:
		// Operator-facing: the shared token doesn't match. Never echo it.
		return nil, fmt.Errorf("servermanager logs: bearer token rejected (check SERVERMANAGER_TOKEN/SM_TOKEN)")
	default:
		return nil, fmt.Errorf("servermanager logs: status %d: %s", res.StatusCode, managerError(res))
	}
}

// StreamLogs follows the manager's NDJSON log stream and re-emits each line
// as a LogEntry. The channel closes when the manager ends the stream or ctx
// is canceled (which also closes the underlying response body).
func (c *Client) StreamLogs(ctx context.Context, appID string, opts ports.LogOptions) (<-chan ports.LogEntry, error) {
	query := url.Values{}
	if opts.Tail > 0 {
		query.Set("tail", strconv.Itoa(opts.Tail))
	}
	if !opts.Since.IsZero() {
		query.Set("since", opts.Since.Format(time.RFC3339))
	}

	endpoint := c.baseURL + "/v1/apps/" + url.PathEscape(appID) + "/logs/stream"
	if len(query) > 0 {
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
		return nil, fmt.Errorf("servermanager stream: %w", ports.ErrAppNotDeployed)
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
