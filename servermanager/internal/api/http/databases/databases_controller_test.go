package databases

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"servermanager/internal/app/apps"
	dbsvc "servermanager/internal/app/databases"
	"servermanager/internal/app/metrics"
	"servermanager/internal/domain"
	"servermanager/internal/mocks"

	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

const testDBID = "6f1d2c3b-4a5e-4f60-8b7a-9c0d1e2f3a4b"

type harness struct {
	router  *gin.Engine
	runtime *mocks.MockContainerRuntime
	metrics *metrics.Store
}

func setupRouter(t *testing.T) harness {
	t.Helper()
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)

	h := harness{
		runtime: mocks.NewMockContainerRuntime(ctrl),
		metrics: metrics.NewStore(30*time.Second, 24*time.Hour, fixedClock{}),
	}
	router := gin.New()
	v1 := router.Group("/v1")
	SetupRouter(dbsvc.Dependencies{
		Runtime: h.runtime,
		Metrics: h.metrics,
		Limits: apps.Limits{
			DefaultMemoryBytes: 256 * 1024 * 1024,
			MaxMemoryBytes:     1024 * 1024 * 1024,
			DefaultNanoCPUs:    500_000_000,
			MaxNanoCPUs:        2_000_000_000,
			DefaultPidsLimit:   256,
			DefaultRuntime:     domain.RuntimeRunc,
		},
		Clock:   fixedClock{},
		Budgets: dbsvc.Budgets{PullTimeout: time.Second, HealthBudget: time.Second},
	}, v1)
	h.router = router
	return h
}

type fixedClock struct{}

func (fixedClock) Now() time.Time { return time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC) }

func request(router *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func provisionBody() string {
	return `{"type":"postgres","version":"16","db_name":"appdb","db_user":"app","db_password":"0123456789abcdef0123456789abcdef"}` // gitleaks:allow — fixture, not a secret
}

// waitProvisionTerminal polls the status endpoint until the async pipeline
// settles, so gomock's Finish never races the detached goroutine.
func waitProvisionTerminal(t *testing.T, router *gin.Engine) DatabaseStatusResponse {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for {
		w := request(router, http.MethodGet, "/v1/dbs/"+testDBID, "")
		if w.Code != http.StatusOK {
			t.Fatalf("GET db = %d: %s", w.Code, w.Body.String())
		}
		var resp DatabaseStatusResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("decoding status: %v", err)
		}
		if p := resp.Provision; p != nil &&
			(p.Status == string(domain.ProvisionStatusRunning) || p.Status == string(domain.ProvisionStatusFailed)) {
			return resp
		}
		if time.Now().After(deadline) {
			t.Fatalf("provision never terminal: %s", w.Body.String())
		}
		time.Sleep(2 * time.Millisecond)
	}
}

func TestProvisionAcceptedAndPollable(t *testing.T) {
	h := setupRouter(t)
	name := domain.DBContainerName(testDBID)
	// Pre-check plus every status poll; the container never comes to exist
	// because the pull fails.
	h.runtime.EXPECT().Inspect(gomock.Any(), name).
		Return(nil, fmt.Errorf("inspect: %w", domain.ErrNotFound)).AnyTimes()
	// Fail at the pull so the pipeline stays one mock deep.
	h.runtime.EXPECT().PullImage(gomock.Any(), "postgres:16").
		Return(errors.New("registry unreachable")).AnyTimes()

	w := request(h.router, http.MethodPost, "/v1/dbs/"+testDBID+"/provision", provisionBody())
	if w.Code != http.StatusAccepted {
		t.Fatalf("POST provision = %d: %s", w.Code, w.Body.String())
	}
	var resp ProvisionResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil || resp.ProvisionID == "" {
		t.Fatalf("provision response %q: %v", w.Body.String(), err)
	}

	status := waitProvisionTerminal(t, h.router)
	if status.Provision.Status != string(domain.ProvisionStatusFailed) ||
		!strings.Contains(status.Provision.Error, "could not pull") {
		t.Errorf("provision = %+v", status.Provision)
	}
	// The password from the request body must never surface in the status.
	if strings.Contains(strings.ToLower(mustJSON(t, status)), "0123456789abcdef") {
		t.Error("status response leaks the password")
	}
}

func TestProvisionValidationAndConflict(t *testing.T) {
	h := setupRouter(t)

	cases := []struct {
		name, path, body string
		wantCode         int
		wantErr          string
	}{
		{"bad uuid", "/v1/dbs/not-a-uuid/provision", provisionBody(), http.StatusBadRequest, "id must be a UUID"},
		{"bad body", "/v1/dbs/" + testDBID + "/provision", "{not json", http.StatusBadRequest, "invalid request body"},
		{"missing fields", "/v1/dbs/" + testDBID + "/provision", `{"type":"postgres"}`, http.StatusBadRequest, "invalid request body"},
		{"type not allowlisted", "/v1/dbs/" + testDBID + "/provision",
			`{"type":"mysql","version":"8.0","db_name":"a","db_user":"a","db_password":"0123456789abcdef0123456789abcdef"}`, // gitleaks:allow — fixture, not a secret
			http.StatusBadRequest, "not supported"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := request(h.router, http.MethodPost, tc.path, tc.body)
			if w.Code != tc.wantCode {
				t.Fatalf("code = %d, want %d: %s", w.Code, tc.wantCode, w.Body.String())
			}
			if !strings.Contains(w.Body.String(), tc.wantErr) {
				t.Errorf("body %q missing %q", w.Body.String(), tc.wantErr)
			}
		})
	}

	t.Run("existing container conflicts", func(t *testing.T) {
		h.runtime.EXPECT().Inspect(gomock.Any(), domain.DBContainerName(testDBID)).
			Return(&domain.ContainerState{Exists: true}, nil)
		w := request(h.router, http.MethodPost, "/v1/dbs/"+testDBID+"/provision", provisionBody())
		if w.Code != http.StatusConflict {
			t.Fatalf("code = %d, want 409: %s", w.Code, w.Body.String())
		}
	})
}

func TestStatusShape(t *testing.T) {
	h := setupRouter(t)
	started := time.Date(2026, 7, 10, 11, 0, 0, 0, time.UTC)
	h.runtime.EXPECT().Inspect(gomock.Any(), domain.DBContainerName(testDBID)).
		Return(&domain.ContainerState{
			Exists: true, Running: true, StartedAt: started, Status: "running", Health: "healthy",
		}, nil)

	w := request(h.router, http.MethodGet, "/v1/dbs/"+testDBID, "")
	if w.Code != http.StatusOK {
		t.Fatalf("GET = %d: %s", w.Code, w.Body.String())
	}
	var resp DatabaseStatusResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if !resp.Exists || !resp.Running || resp.Health != "healthy" || resp.StartedAt == nil {
		t.Errorf("status = %+v", resp)
	}
	if resp.Provision != nil {
		t.Errorf("provision = %+v, want absent without a job", resp.Provision)
	}
}

func TestDeleteMapsConflict(t *testing.T) {
	h := setupRouter(t)
	h.runtime.EXPECT().Stop(gomock.Any(), gomock.Any()).Return(nil)
	h.runtime.EXPECT().RemoveContainer(gomock.Any(), gomock.Any()).Return(nil)
	h.runtime.EXPECT().RemoveNetwork(gomock.Any(), gomock.Any()).Return(nil)
	h.runtime.EXPECT().RemoveVolume(gomock.Any(), gomock.Any()).
		Return(fmt.Errorf("volume in use: %w", domain.ErrConflict))

	w := request(h.router, http.MethodDelete, "/v1/dbs/"+testDBID, "")
	if w.Code != http.StatusConflict {
		t.Fatalf("DELETE = %d, want 409: %s", w.Code, w.Body.String())
	}
}

func TestLogsEndpoint(t *testing.T) {
	h := setupRouter(t)
	h.runtime.EXPECT().Logs(gomock.Any(), domain.DBContainerName(testDBID), gomock.Any()).
		Return([]domain.LogLine{
			{Timestamp: time.Date(2026, 7, 10, 11, 0, 0, 0, time.UTC), Stream: "stderr", Message: "ready to accept connections"},
		}, nil)

	w := request(h.router, http.MethodGet, "/v1/dbs/"+testDBID+"/logs?tail=50", "")
	if w.Code != http.StatusOK {
		t.Fatalf("GET logs = %d: %s", w.Code, w.Body.String())
	}
	var resp LogsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Logs) != 1 || resp.Logs[0].Level != "error" || resp.Logs[0].Message != "ready to accept connections" {
		t.Errorf("logs = %+v", resp.Logs)
	}

	if w := request(h.router, http.MethodGet, "/v1/dbs/"+testDBID+"/logs?tail=zero", ""); w.Code != http.StatusBadRequest {
		t.Errorf("bad tail = %d, want 400", w.Code)
	}
}

func TestMetricsEndpoint(t *testing.T) {
	base := "/v1/dbs/" + testDBID + "/metrics"

	t.Run("bad range", func(t *testing.T) {
		h := setupRouter(t)
		if w := request(h.router, http.MethodGet, base+"?range=7d", ""); w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("non-uuid id", func(t *testing.T) {
		h := setupRouter(t)
		if w := request(h.router, http.MethodGet, "/v1/dbs/not-a-uuid/metrics", ""); w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("no samples is 200 with all four series empty", func(t *testing.T) {
		h := setupRouter(t)
		w := request(h.router, http.MethodGet, base, "")
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body)
		}
		var resp MetricsResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("invalid JSON response: %v", err)
		}
		if resp.Range != "24h" {
			t.Errorf("range = %q, want default 24h", resp.Range)
		}
		for _, key := range []string{"conn", "qps", "cpu", "disk"} {
			if points, ok := resp.Series[key]; !ok || len(points) != 0 {
				t.Errorf("series[%q] = %v, want present and empty", key, points)
			}
		}
		if strings.Contains(w.Body.String(), "null") {
			t.Errorf("body contains null: %s", w.Body)
		}
	})

	t.Run("seeded series round-trips", func(t *testing.T) {
		h := setupRouter(t)
		h.metrics.Record(domain.DBContainerName(testDBID), fixedClock{}.Now().Add(-5*time.Minute),
			map[string]float64{"conn": 3, "disk": 42})

		w := request(h.router, http.MethodGet, base+"?range=1h", "")
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body)
		}
		var resp MetricsResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("invalid JSON response: %v", err)
		}
		if len(resp.Series["conn"]) != 1 || resp.Series["conn"][0].Value != 3 {
			t.Errorf("conn = %+v, want one 3 point", resp.Series["conn"])
		}
		if len(resp.Series["disk"]) != 1 || resp.Series["disk"][0].Value != 42 {
			t.Errorf("disk = %+v, want one 42 point", resp.Series["disk"])
		}
		if len(resp.Series["qps"]) != 0 || len(resp.Series["cpu"]) != 0 {
			t.Errorf("qps/cpu = %+v/%+v, want empty", resp.Series["qps"], resp.Series["cpu"])
		}
	})
}

func mustJSON(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
