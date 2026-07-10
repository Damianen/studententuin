package deployments

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"servermanager/internal/app/apps"
	deploysvc "servermanager/internal/app/deployments"
	"servermanager/internal/domain"
	"servermanager/internal/mocks"

	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

const testAppID = "8b9f5e9e-9a3a-4b7e-9a59-1d2f3a4b5c6d"

type okRunner struct{}

func (okRunner) Validate(apps.RunInput) error { return nil }
func (okRunner) Execute(context.Context, apps.RunInput) (string, error) {
	return "cid-new", nil
}

type harness struct {
	router  *gin.Engine
	fetcher *mocks.MockSourceFetcher
	runtime *mocks.MockContainerRuntime
}

// setupRouter wires the real service behind the controller with mocked infra
// — pipelines started by POSTs run for real, so tests either block them or
// wait for the terminal state.
func setupRouter(t *testing.T) harness {
	t.Helper()
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)

	h := harness{
		fetcher: mocks.NewMockSourceFetcher(ctrl),
		runtime: mocks.NewMockContainerRuntime(ctrl),
	}
	router := gin.New()
	v1 := router.Group("/v1")
	SetupRouter(deploysvc.Dependencies{
		Fetcher:  h.fetcher,
		Builder:  mocks.NewMockImageBuilder(ctrl),
		Runtime:  h.runtime,
		Images:   mocks.NewMockImageStore(ctrl),
		Runner:   okRunner{},
		Clock:    fixedClock{},
		GitHosts: []string{"github.com"},
		Budgets: deploysvc.Budgets{
			CloneTimeout: time.Second,
			BuildTimeout: time.Second,
			HealthGrace:  time.Millisecond,
		},
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

func deployBody() string {
	return `{"repository_url":"https://github.com/user/repo","branch":"main","port":3000}`
}

// waitJobTerminal polls the GET endpoint until the async pipeline settles,
// so gomock's Finish never races the detached goroutine.
func waitJobTerminal(t *testing.T, router *gin.Engine, id string) DeploymentResponse {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for {
		w := request(router, http.MethodGet, "/v1/deployments/"+id, "")
		if w.Code != http.StatusOK {
			t.Fatalf("GET deployment = %d: %s", w.Code, w.Body.String())
		}
		var resp DeploymentResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("decoding deployment: %v", err)
		}
		if resp.Status == string(domain.DeploymentStatusRunning) || resp.Status == string(domain.DeploymentStatusFailed) {
			return resp
		}
		if time.Now().After(deadline) {
			t.Fatalf("deployment stuck in %s", resp.Status)
		}
		time.Sleep(2 * time.Millisecond)
	}
}

func TestDeployAcceptedAndPollable(t *testing.T) {
	h := setupRouter(t)
	// The pipeline runs for real; fail it at the clone so only fetcher is hit.
	h.fetcher.EXPECT().Fetch(gomock.Any(), "https://github.com/user/repo", "main").
		Return(nil, errors.New("repository not found")).AnyTimes()

	w := request(h.router, http.MethodPost, "/v1/apps/"+testAppID+"/deploy", deployBody())
	if w.Code != http.StatusAccepted {
		t.Fatalf("POST deploy = %d: %s", w.Code, w.Body.String())
	}
	var resp DeployResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil || resp.DeploymentID == "" {
		t.Fatalf("deploy response %q: %v", w.Body.String(), err)
	}

	job := waitJobTerminal(t, h.router, resp.DeploymentID)
	if job.Status != string(domain.DeploymentStatusFailed) || job.Error != "repository not found" {
		t.Errorf("job = %+v", job)
	}
	if job.AppID != testAppID || job.ID != resp.DeploymentID {
		t.Errorf("job identity = %+v", job)
	}
}

func TestDeployValidation(t *testing.T) {
	h := setupRouter(t)

	cases := []struct {
		name, path, body string
		wantCode         int
		wantErr          string
	}{
		{"bad app uuid", "/v1/apps/not-a-uuid/deploy", deployBody(), http.StatusBadRequest, "id must be a UUID"},
		{"bad body", "/v1/apps/" + testAppID + "/deploy", "{not json", http.StatusBadRequest, "invalid request body"},
		{"missing repo url", "/v1/apps/" + testAppID + "/deploy", `{"port":3000}`, http.StatusBadRequest, "repository url is required"},
		{"http url", "/v1/apps/" + testAppID + "/deploy", `{"repository_url":"http://github.com/u/r"}`, http.StatusBadRequest, "must use https"},
		{"bad branch", "/v1/apps/" + testAppID + "/deploy", `{"repository_url":"https://github.com/u/r","branch":"a b"}`, http.StatusBadRequest, "branch"},
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
}

func TestDeployConflict(t *testing.T) {
	h := setupRouter(t)
	release := make(chan struct{})
	h.fetcher.EXPECT().Fetch(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(context.Context, string, string) (*domain.SourceCheckout, error) {
			<-release
			return nil, errors.New("released")
		}).AnyTimes()

	first := request(h.router, http.MethodPost, "/v1/apps/"+testAppID+"/deploy", deployBody())
	if first.Code != http.StatusAccepted {
		t.Fatalf("first deploy = %d", first.Code)
	}
	var accepted DeployResponse
	if err := json.Unmarshal(first.Body.Bytes(), &accepted); err != nil {
		t.Fatal(err)
	}

	second := request(h.router, http.MethodPost, "/v1/apps/"+testAppID+"/deploy", deployBody())
	if second.Code != http.StatusConflict {
		t.Fatalf("second deploy = %d, want 409: %s", second.Code, second.Body.String())
	}
	if !strings.Contains(second.Body.String(), "already in flight") {
		t.Errorf("conflict body = %s", second.Body.String())
	}

	close(release)
	waitJobTerminal(t, h.router, accepted.DeploymentID)
}

func TestGetDeploymentErrors(t *testing.T) {
	h := setupRouter(t)

	if w := request(h.router, http.MethodGet, "/v1/deployments/not-a-uuid", ""); w.Code != http.StatusBadRequest {
		t.Errorf("bad uuid = %d, want 400", w.Code)
	}
	w := request(h.router, http.MethodGet, "/v1/deployments/"+testAppID, "")
	if w.Code != http.StatusNotFound || !strings.Contains(w.Body.String(), "no such deployment") {
		t.Errorf("unknown deployment = %d: %s", w.Code, w.Body.String())
	}
}
