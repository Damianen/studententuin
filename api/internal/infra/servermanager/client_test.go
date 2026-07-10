package servermanager

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"api/internal/app/ports"
)

const testToken = "test-sm-token"

func TestClientLogs(t *testing.T) {
	appID := "8b9f5e9e-9a3a-4b7e-9a59-1d2f3a4b5c6d"

	t.Run("requests the manager and decodes entries", func(t *testing.T) {
		var gotPath, gotAuth, gotTail, gotSince string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			gotAuth = r.Header.Get("Authorization")
			gotTail = r.URL.Query().Get("tail")
			gotSince = r.URL.Query().Get("since")
			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write([]byte(`{"logs":[{"id":"1-0","timestamp":"2026-06-13T10:00:01Z","level":"info","message":"tick"}]}`)); err != nil {
				t.Errorf("writing response: %v", err)
			}
		}))
		defer srv.Close()

		client := NewClient(srv.URL+"/", testToken) // trailing slash must not double up
		since := time.Date(2026, 6, 13, 9, 0, 0, 0, time.UTC)
		entries, err := client.Logs(context.Background(), appID, ports.LogOptions{Tail: 50, Since: since})
		if err != nil {
			t.Fatalf("Logs: %v", err)
		}

		if gotPath != "/v1/apps/"+appID+"/logs" {
			t.Errorf("path = %q", gotPath)
		}
		if gotAuth != "Bearer "+testToken {
			t.Errorf("auth header = %q", gotAuth)
		}
		if gotTail != "50" || gotSince != "2026-06-13T09:00:00Z" {
			t.Errorf("query tail=%q since=%q", gotTail, gotSince)
		}
		if len(entries) != 1 || entries[0].Level != "info" || entries[0].Message != "tick" {
			t.Errorf("entries = %+v", entries)
		}
		if !entries[0].Timestamp.Equal(time.Date(2026, 6, 13, 10, 0, 1, 0, time.UTC)) {
			t.Errorf("timestamp = %v", entries[0].Timestamp)
		}
	})

	t.Run("404 means not deployed", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, `{"error":"no container for this application"}`, http.StatusNotFound)
		}))
		defer srv.Close()

		_, err := NewClient(srv.URL, testToken).Logs(context.Background(), appID, ports.LogOptions{})
		if !errors.Is(err, ports.ErrAppNotDeployed) {
			t.Errorf("err = %v, want ErrAppNotDeployed", err)
		}
	})

	t.Run("401 surfaces a token hint without the token", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer srv.Close()

		_, err := NewClient(srv.URL, testToken).Logs(context.Background(), appID, ports.LogOptions{})
		if err == nil || !strings.Contains(err.Error(), "token") {
			t.Errorf("err = %v, want token hint", err)
		}
		if err != nil && strings.Contains(err.Error(), testToken) {
			t.Errorf("error leaks the token: %v", err)
		}
	})

	t.Run("other statuses surface the manager error body", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		}))
		defer srv.Close()

		_, err := NewClient(srv.URL, testToken).Logs(context.Background(), appID, ports.LogOptions{})
		if err == nil || !strings.Contains(err.Error(), "internal error") || !strings.Contains(err.Error(), "500") {
			t.Errorf("err = %v, want status and detail", err)
		}
	})

	t.Run("unreachable manager is an error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		srv.Close() // closed listener: connection refused

		if _, err := NewClient(srv.URL, testToken).Logs(context.Background(), appID, ports.LogOptions{}); err == nil {
			t.Error("err = nil, want transport error")
		}
	})
}

func TestClientStreamLogs(t *testing.T) {
	appID := "8b9f5e9e-9a3a-4b7e-9a59-1d2f3a4b5c6d"

	t.Run("re-emits ndjson lines until the stream ends", func(t *testing.T) {
		var gotPath, gotAuth string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			gotAuth = r.Header.Get("Authorization")
			w.Header().Set("Content-Type", "application/x-ndjson")
			_, _ = w.Write([]byte(`{"id":"1-0","timestamp":"2026-06-13T10:00:01Z","level":"info","message":"tick"}` + "\n"))
			_, _ = w.Write([]byte("not json, skipped\n"))
			_, _ = w.Write([]byte(`{"id":"2-1","timestamp":"2026-06-13T10:00:02Z","level":"error","message":"tock"}` + "\n"))
		}))
		defer srv.Close()

		entries, err := NewClient(srv.URL, testToken).StreamLogs(context.Background(), appID, ports.LogOptions{Tail: 100})
		if err != nil {
			t.Fatalf("StreamLogs: %v", err)
		}

		var got []ports.LogEntry
		for entry := range entries {
			got = append(got, entry)
		}
		if gotPath != "/v1/apps/"+appID+"/logs/stream" {
			t.Errorf("path = %q", gotPath)
		}
		if gotAuth != "Bearer "+testToken {
			t.Errorf("auth header = %q", gotAuth)
		}
		if len(got) != 2 || got[0].Message != "tick" || got[1].Level != "error" {
			t.Errorf("entries = %+v", got)
		}
	})

	t.Run("404 means not deployed", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, `{"error":"no container for this application"}`, http.StatusNotFound)
		}))
		defer srv.Close()

		_, err := NewClient(srv.URL, testToken).StreamLogs(context.Background(), appID, ports.LogOptions{})
		if !errors.Is(err, ports.ErrAppNotDeployed) {
			t.Errorf("err = %v, want ErrAppNotDeployed", err)
		}
	})

	t.Run("context cancel ends the stream", func(t *testing.T) {
		release := make(chan struct{})
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/x-ndjson")
			_, _ = w.Write([]byte(`{"id":"1-0","timestamp":"2026-06-13T10:00:01Z","level":"info","message":"tick"}` + "\n"))
			w.(http.Flusher).Flush()
			<-release // hold the stream open like a real follow
		}))
		defer srv.Close()
		defer close(release)

		ctx, cancel := context.WithCancel(context.Background())
		entries, err := NewClient(srv.URL, testToken).StreamLogs(ctx, appID, ports.LogOptions{})
		if err != nil {
			t.Fatalf("StreamLogs: %v", err)
		}

		if entry := <-entries; entry.Message != "tick" {
			t.Fatalf("first entry = %+v", entry)
		}
		cancel()

		select {
		case _, open := <-entries:
			if open {
				t.Error("expected closed channel after cancel")
			}
		case <-time.After(5 * time.Second):
			t.Error("stream did not end after context cancel")
		}
	})
}

func TestClientDeploy(t *testing.T) {
	appID := "8b9f5e9e-9a3a-4b7e-9a59-1d2f3a4b5c6d"
	spec := ports.DeploySpec{
		RepositoryURL: "https://github.com/user/repo",
		Branch:        "main",
		Type:          "Nodejs",
		StartCommand:  "node index.js",
		Env:           map[string]string{"FOO": "bar"},
		Port:          3000,
		MemoryLimit:   "256m",
	}

	t.Run("posts the spec and returns the deployment id", func(t *testing.T) {
		var gotPath, gotAuth, gotContentType string
		var gotBody map[string]any
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			gotAuth = r.Header.Get("Authorization")
			gotContentType = r.Header.Get("Content-Type")
			if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
				t.Errorf("decoding deploy body: %v", err)
			}
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte(`{"deployment_id":"dep-123"}`)) //nolint:errcheck
		}))
		defer srv.Close()

		id, err := NewClient(srv.URL, testToken).Deploy(context.Background(), appID, spec)
		if err != nil {
			t.Fatalf("Deploy: %v", err)
		}
		if id != "dep-123" {
			t.Errorf("deployment id = %q", id)
		}
		if gotPath != "/v1/apps/"+appID+"/deploy" || gotAuth != "Bearer "+testToken || gotContentType != "application/json" {
			t.Errorf("request path=%q auth=%q content-type=%q", gotPath, gotAuth, gotContentType)
		}
		if gotBody["repository_url"] != spec.RepositoryURL || gotBody["start_command"] != spec.StartCommand ||
			gotBody["port"] != float64(3000) || gotBody["memory_limit"] != "256m" {
			t.Errorf("deploy body = %v", gotBody)
		}
	})

	t.Run("409 means a deploy is in flight", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, `{"error":"a deployment is already in flight for this application"}`, http.StatusConflict)
		}))
		defer srv.Close()

		if _, err := NewClient(srv.URL, testToken).Deploy(context.Background(), appID, spec); !errors.Is(err, ports.ErrDeployInFlight) {
			t.Errorf("err = %v, want ErrDeployInFlight", err)
		}
	})

	t.Run("400 surfaces the manager's reason", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, `{"error":"repository host \"evil.io\" is not allowed"}`, http.StatusBadRequest)
		}))
		defer srv.Close()

		_, err := NewClient(srv.URL, testToken).Deploy(context.Background(), appID, spec)
		if !errors.Is(err, ports.ErrSpecRejected) || !strings.Contains(err.Error(), "evil.io") {
			t.Errorf("err = %v, want ErrSpecRejected with the manager's detail", err)
		}
	})
}

func TestClientDeploymentStatus(t *testing.T) {
	t.Run("decodes the job resource", func(t *testing.T) {
		var gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			w.Write([]byte(`{"id":"dep-123","app_id":"app-1","status":"running","commit_sha":"abc","image":"stt-app-app-1:dep","container_id":"cid","container_name":"stt-app-app-1","created_at":"2026-07-10T12:00:00Z","updated_at":"2026-07-10T12:01:00Z"}`)) //nolint:errcheck
		}))
		defer srv.Close()

		status, err := NewClient(srv.URL, testToken).DeploymentStatus(context.Background(), "dep-123")
		if err != nil {
			t.Fatalf("DeploymentStatus: %v", err)
		}
		if gotPath != "/v1/deployments/dep-123" {
			t.Errorf("path = %q", gotPath)
		}
		if status.Status != "running" || status.ContainerID != "cid" || status.Image != "stt-app-app-1:dep" || !status.Terminal() {
			t.Errorf("status = %+v", status)
		}
	})

	t.Run("404 means unknown deployment", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, `{"error":"no such deployment"}`, http.StatusNotFound)
		}))
		defer srv.Close()

		if _, err := NewClient(srv.URL, testToken).DeploymentStatus(context.Background(), "dep-x"); !errors.Is(err, ports.ErrDeploymentNotFound) {
			t.Errorf("err = %v, want ErrDeploymentNotFound", err)
		}
	})
}

func TestClientLifecycle(t *testing.T) {
	appID := "8b9f5e9e-9a3a-4b7e-9a59-1d2f3a4b5c6d"

	cases := []struct {
		name, wantMethod, wantPath string
		call                       func(*Client) error
	}{
		{"start", http.MethodPost, "/v1/apps/" + appID + "/start", func(c *Client) error { return c.Start(context.Background(), appID) }},
		{"stop", http.MethodPost, "/v1/apps/" + appID + "/stop", func(c *Client) error { return c.Stop(context.Background(), appID) }},
		{"remove", http.MethodDelete, "/v1/apps/" + appID, func(c *Client) error { return c.Remove(context.Background(), appID) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var gotMethod, gotPath string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotMethod, gotPath = r.Method, r.URL.Path
				w.WriteHeader(http.StatusOK)
			}))
			defer srv.Close()

			if err := tc.call(NewClient(srv.URL, testToken)); err != nil {
				t.Fatalf("%s: %v", tc.name, err)
			}
			if gotMethod != tc.wantMethod || gotPath != tc.wantPath {
				t.Errorf("request = %s %s, want %s %s", gotMethod, gotPath, tc.wantMethod, tc.wantPath)
			}
		})

		t.Run(tc.name+" 404 means not deployed", func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				http.Error(w, `{"error":"no container for this application"}`, http.StatusNotFound)
			}))
			defer srv.Close()

			if err := tc.call(NewClient(srv.URL, testToken)); !errors.Is(err, ports.ErrAppNotDeployed) {
				t.Errorf("err = %v, want ErrAppNotDeployed", err)
			}
		})
	}
}
