package servermanager

import (
	"context"
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
