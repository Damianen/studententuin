package apps

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	appsvc "servermanager/internal/app/apps"
	"servermanager/internal/domain"
	"servermanager/internal/mocks"

	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

const testAppID = "8b9f5e9e-9a3a-4b7e-9a59-1d2f3a4b5c6d"

func testLimits() appsvc.Limits {
	return appsvc.Limits{
		DefaultMemoryBytes: 256 * 1024 * 1024,
		MaxMemoryBytes:     1024 * 1024 * 1024,
		DefaultNanoCPUs:    500_000_000,
		MaxNanoCPUs:        2_000_000_000,
		DefaultPidsLimit:   256,
		DefaultRuntime:     domain.RuntimeRunc,
	}
}

func setupRouter(t *testing.T) (*gin.Engine, *mocks.MockContainerRuntime) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	rt := mocks.NewMockContainerRuntime(ctrl)

	router := gin.New()
	v1 := router.Group("/v1")
	SetupRouter(appsvc.Dependencies{Runtime: rt, Limits: testLimits()}, v1)
	return router, rt
}

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

func TestRunEndpoint(t *testing.T) {
	runBody := `{"image":"busybox:1.37","start_command":["sleep","3600"]}`

	t.Run("created", func(t *testing.T) {
		router, rt := setupRouter(t)
		rt.EXPECT().Create(gomock.Any(), gomock.Any()).Return("cid-1", nil)
		rt.EXPECT().Start(gomock.Any(), "cid-1").Return(nil)

		w := request(router, http.MethodPost, "/v1/apps/"+testAppID+"/run", runBody)
		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, want 201; body %s", w.Code, w.Body)
		}
		var resp RunResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("invalid JSON response: %v", err)
		}
		if resp.ContainerID != "cid-1" || resp.Name != "stt-app-"+testAppID {
			t.Errorf("response = %+v", resp)
		}
	})

	t.Run("non-uuid id", func(t *testing.T) {
		router, _ := setupRouter(t)
		w := request(router, http.MethodPost, "/v1/apps/not-a-uuid/run", runBody)
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("malformed json", func(t *testing.T) {
		router, _ := setupRouter(t)
		w := request(router, http.MethodPost, "/v1/apps/"+testAppID+"/run", `{"image":`)
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("missing image", func(t *testing.T) {
		router, _ := setupRouter(t)
		w := request(router, http.MethodPost, "/v1/apps/"+testAppID+"/run", `{}`)
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("validation error from service", func(t *testing.T) {
		router, _ := setupRouter(t)
		w := request(router, http.MethodPost, "/v1/apps/"+testAppID+"/run",
			`{"image":"busybox:1.37","memory_limit":"99g"}`)
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400; body %s", w.Code, w.Body)
		}
	})

	t.Run("image missing on host", func(t *testing.T) {
		router, rt := setupRouter(t)
		rt.EXPECT().Create(gomock.Any(), gomock.Any()).Return("", domain.ErrImageMissing)
		w := request(router, http.MethodPost, "/v1/apps/"+testAppID+"/run", runBody)
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("conflict", func(t *testing.T) {
		router, rt := setupRouter(t)
		rt.EXPECT().Create(gomock.Any(), gomock.Any()).Return("", domain.ErrConflict)
		w := request(router, http.MethodPost, "/v1/apps/"+testAppID+"/run", runBody)
		if w.Code != http.StatusConflict {
			t.Errorf("status = %d, want 409", w.Code)
		}
	})

	t.Run("docker failure is a generic 500", func(t *testing.T) {
		router, rt := setupRouter(t)
		rt.EXPECT().Create(gomock.Any(), gomock.Any()).Return("", errors.New("socket exploded"))
		w := request(router, http.MethodPost, "/v1/apps/"+testAppID+"/run", runBody)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500", w.Code)
		}
		if strings.Contains(w.Body.String(), "socket exploded") {
			t.Errorf("response leaks internal error detail: %s", w.Body)
		}
	})
}

func TestLifecycleEndpoints(t *testing.T) {
	name := domain.AppContainerName(testAppID)

	t.Run("start ok", func(t *testing.T) {
		router, rt := setupRouter(t)
		rt.EXPECT().Start(gomock.Any(), name).Return(nil)
		if w := request(router, http.MethodPost, "/v1/apps/"+testAppID+"/start", ""); w.Code != http.StatusOK {
			t.Errorf("status = %d, want 200", w.Code)
		}
	})

	t.Run("start missing container", func(t *testing.T) {
		router, rt := setupRouter(t)
		rt.EXPECT().Start(gomock.Any(), name).Return(domain.ErrNotFound)
		if w := request(router, http.MethodPost, "/v1/apps/"+testAppID+"/start", ""); w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want 404", w.Code)
		}
	})

	t.Run("stop ok", func(t *testing.T) {
		router, rt := setupRouter(t)
		rt.EXPECT().Stop(gomock.Any(), name).Return(nil)
		if w := request(router, http.MethodPost, "/v1/apps/"+testAppID+"/stop", ""); w.Code != http.StatusOK {
			t.Errorf("status = %d, want 200", w.Code)
		}
	})

	t.Run("remove is idempotent", func(t *testing.T) {
		router, rt := setupRouter(t)
		rt.EXPECT().Remove(gomock.Any(), name).Return(domain.ErrNotFound)
		if w := request(router, http.MethodDelete, "/v1/apps/"+testAppID, ""); w.Code != http.StatusOK {
			t.Errorf("status = %d, want 200", w.Code)
		}
	})
}

func TestStatusEndpoint(t *testing.T) {
	name := domain.AppContainerName(testAppID)

	t.Run("running container", func(t *testing.T) {
		router, rt := setupRouter(t)
		rt.EXPECT().Inspect(gomock.Any(), name).Return(&domain.ContainerState{
			Exists: true, Running: true, Status: "running", RestartCount: 1,
		}, nil)

		w := request(router, http.MethodGet, "/v1/apps/"+testAppID, "")
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
		var resp ContainerStateResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("invalid JSON response: %v", err)
		}
		if !resp.Exists || !resp.Running || resp.Status != "running" || resp.RestartCount != 1 {
			t.Errorf("response = %+v", resp)
		}
		if resp.StartedAt != nil {
			t.Errorf("StartedAt = %v, want omitted for zero time", resp.StartedAt)
		}
	})

	t.Run("missing container is exists false, not 404", func(t *testing.T) {
		router, rt := setupRouter(t)
		rt.EXPECT().Inspect(gomock.Any(), name).Return(nil, domain.ErrNotFound)

		w := request(router, http.MethodGet, "/v1/apps/"+testAppID, "")
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
		var resp ContainerStateResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("invalid JSON response: %v", err)
		}
		if resp.Exists || resp.Running {
			t.Errorf("response = %+v, want exists false", resp)
		}
	})
}
