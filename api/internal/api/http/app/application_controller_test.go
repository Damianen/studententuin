package app

import (
	"api/internal/api/middlewares"
	"api/internal/app/ports"
	"api/internal/domain"
	"api/internal/infra/utils"
	"api/internal/mocks"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	subdomainhttp "api/internal/api/http/subdomain"
	appapp "api/internal/app/app"
	appsubdomain "api/internal/app/subdomain"
	infraauth "api/internal/infra/auth"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func init() {
	gin.SetMode(gin.TestMode)
}

const testSecret = "test-secret-key"

type appMocks struct {
	applicationRepo *mocks.MockApplicationRepo
	subdomainRepo   *mocks.MockSubdomainRepo
	userRepo        *mocks.MockUserRepo
	serverManager   *mocks.MockServerManagerClient
	clock           *mocks.MockClock
}

type envelope struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func newAppRouter(t *testing.T) (*gin.Engine, appMocks) {
	ctrl := gomock.NewController(t)
	m := appMocks{
		applicationRepo: mocks.NewMockApplicationRepo(ctrl),
		subdomainRepo:   mocks.NewMockSubdomainRepo(ctrl),
		userRepo:        mocks.NewMockUserRepo(ctrl),
		serverManager:   mocks.NewMockServerManagerClient(ctrl),
		clock:           mocks.NewMockClock(ctrl),
	}
	sdDeps := appsubdomain.Dependencies{
		SubdomainRepo: m.subdomainRepo,
		UserRepo:      m.userRepo,
		Clock:         m.clock,
	}
	appDeps := appapp.Dependencies{
		ApplicationRepo: m.applicationRepo,
		ServerManager:   m.serverManager,
		Clock:           m.clock,
	}
	mw := middlewares.AuthMiddleware{
		JwtTokenizer: infraauth.JwtTokenizer{Clock: utils.SystemClock{}, SecretKey: testSecret},
	}
	r := gin.New()
	group := subdomainhttp.SetupRouter(sdDeps, mw, r)
	SetupRouter(appDeps, sdDeps, mw, group)
	return r, m
}

func authCookie(t *testing.T, userID string) *http.Cookie {
	t.Helper()
	tokenizer := infraauth.JwtTokenizer{Clock: utils.SystemClock{}, SecretKey: testSecret}
	token, err := tokenizer.CreateToken(userID)
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}
	return &http.Cookie{Name: "AuthToken", Value: token}
}

func doJSON(r *gin.Engine, method, path, body string, cookie *http.Cookie) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if cookie != nil {
		req.AddCookie(cookie)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestCreateApplication(t *testing.T) {
	ownerID := uuid.New()
	subID := uuid.New()
	owned := &domain.Subdomain{ID: subID, UserID: ownerID}
	path := "/subdomain/" + subID.String() + "/application"

	t.Run("success normalizes type and stores pending application", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.applicationRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
			func(app *domain.Application, ctx any) error {
				if app.SubdomainID != subID {
					t.Errorf("expected subdomain %s, got %s", subID, app.SubdomainID)
				}
				if app.Type != domain.ApplicationTypeNodejs {
					t.Errorf("expected normalized type Nodejs, got %s", app.Type)
				}
				if app.Status != domain.ApplicationStatusPending {
					t.Errorf("expected status pending, got %s", app.Status)
				}
				if app.Name == nil || *app.Name != "my-app" {
					t.Errorf("expected name my-app, got %v", app.Name)
				}
				return nil
			})

		// "node" must be accepted and normalized to the canonical "Nodejs".
		body := `{"name":"my-app","type":"node","repo_url":"https://github.com/me/my-app","branch":"main","build_command":"npm run build","start_command":"npm start"}`
		w := doJSON(r, http.MethodPost, path, body, authCookie(t, ownerID.String()))
		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d (body %s)", w.Code, w.Body.String())
		}
	})

	t.Run("unsupported type returns 400", func(t *testing.T) {
		r, _ := newAppRouter(t)
		body := `{"name":"my-app","type":"python","repo_url":"https://github.com/me/my-app","branch":"main","build_command":"make","start_command":"./run"}`
		w := doJSON(r, http.MethodPost, path, body, authCookie(t, ownerID.String()))
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("missing required field returns 400", func(t *testing.T) {
		r, _ := newAppRouter(t)
		body := `{"name":"my-app","type":"Nodejs","branch":"main","build_command":"npm run build","start_command":"npm start"}`
		w := doJSON(r, http.MethodPost, path, body, authCookie(t, ownerID.String()))
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("non-owner gets 403", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)

		intruder := uuid.New()
		body := `{"name":"my-app","type":"Nodejs","repo_url":"https://github.com/me/my-app","branch":"main","build_command":"npm run build","start_command":"npm start"}`
		w := doJSON(r, http.MethodPost, path, body, authCookie(t, intruder.String()))
		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", w.Code)
		}
	})

	t.Run("no cookie returns 401", func(t *testing.T) {
		r, _ := newAppRouter(t)
		w := doJSON(r, http.MethodPost, path, `{}`, nil)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})
}

func TestGetApplication(t *testing.T) {
	ownerID := uuid.New()
	subID := uuid.New()
	appID := uuid.New()
	owned := &domain.Subdomain{ID: subID, UserID: ownerID}
	path := "/subdomain/" + subID.String() + "/application/" + appID.String()

	t.Run("success maps application to DTO", func(t *testing.T) {
		r, m := newAppRouter(t)
		name := "my-app"
		repoURL := "https://github.com/me/my-app"
		branch := "main"
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.applicationRepo.EXPECT().FindByID(appID.String(), gomock.Any()).Return(&domain.Application{
			ID:                   appID,
			Name:                 &name,
			Type:                 domain.ApplicationTypeNodejs,
			Status:               domain.ApplicationStatusRunning,
			RepositoryURL:        &repoURL,
			Branch:               &branch,
			EnvironmentVariables: map[string]string{"NODE_ENV": "production"},
		}, nil)

		w := doJSON(r, http.MethodGet, path, "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d (body %s)", w.Code, w.Body.String())
		}

		var resp envelope
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		var data map[string]any
		if err := json.Unmarshal(resp.Data, &data); err != nil {
			t.Fatalf("failed to decode data: %v", err)
		}
		if data["id"] != appID.String() || data["name"] != name || data["type"] != "Nodejs" ||
			data["status"] != "running" || data["repo_url"] != repoURL || data["branch"] != branch {
			t.Errorf("unexpected application mapping: %v", data)
		}
		env, ok := data["environment_variables"].(map[string]any)
		if !ok || env["NODE_ENV"] != "production" {
			t.Errorf("expected environment_variables in response, got %v", data["environment_variables"])
		}
	})

	t.Run("missing application returns 404", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.applicationRepo.EXPECT().FindByID(appID.String(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

		w := doJSON(r, http.MethodGet, path, "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("non-owner gets 403", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)

		intruder := uuid.New()
		w := doJSON(r, http.MethodGet, path, "", authCookie(t, intruder.String()))
		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", w.Code)
		}
	})
}

func TestGetApplicationLogs(t *testing.T) {
	ownerID := uuid.New()
	subID := uuid.New()
	appID := uuid.New()
	owned := &domain.Subdomain{ID: subID, UserID: ownerID}
	ownedApp := &domain.Application{ID: appID, SubdomainID: subID}
	path := "/subdomain/" + subID.String() + "/application/" + appID.String() + "/logs"

	t.Run("success returns log entries", func(t *testing.T) {
		r, m := newAppRouter(t)
		ts := time.Date(2026, 6, 13, 10, 0, 1, 0, time.UTC)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.applicationRepo.EXPECT().FindByID(appID.String(), gomock.Any()).Return(ownedApp, nil)
		m.serverManager.EXPECT().Logs(gomock.Any(), appID.String(), ports.LogOptions{Tail: 200}).Return([]ports.LogEntry{
			{ID: "1-0", Timestamp: ts, Level: "info", Message: "tick"},
			{ID: "2-1", Timestamp: ts.Add(time.Second), Level: "error", Message: "tock"},
		}, nil)

		w := doJSON(r, http.MethodGet, path, "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d (body %s)", w.Code, w.Body.String())
		}

		var resp envelope
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		var data []map[string]any
		if err := json.Unmarshal(resp.Data, &data); err != nil {
			t.Fatalf("failed to decode data: %v", err)
		}
		if len(data) != 2 {
			t.Fatalf("expected 2 entries, got %d", len(data))
		}
		if data[0]["level"] != "info" || data[0]["message"] != "tick" || data[0]["timestamp"] != "2026-06-13T10:00:01Z" {
			t.Errorf("unexpected entry mapping: %v", data[0])
		}
		if data[1]["level"] != "error" {
			t.Errorf("expected stderr entry to stay error, got %v", data[1])
		}
	})

	t.Run("tail and since are forwarded", func(t *testing.T) {
		r, m := newAppRouter(t)
		since := time.Date(2026, 6, 13, 9, 0, 0, 0, time.UTC)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.applicationRepo.EXPECT().FindByID(appID.String(), gomock.Any()).Return(ownedApp, nil)
		m.serverManager.EXPECT().Logs(gomock.Any(), appID.String(), ports.LogOptions{Tail: 50, Since: since}).Return([]ports.LogEntry{}, nil)

		w := doJSON(r, http.MethodGet, path+"?tail=50&since=2026-06-13T09:00:00Z", "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d (body %s)", w.Code, w.Body.String())
		}
	})

	t.Run("oversized tail clamps to the cap", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.applicationRepo.EXPECT().FindByID(appID.String(), gomock.Any()).Return(ownedApp, nil)
		m.serverManager.EXPECT().Logs(gomock.Any(), appID.String(), ports.LogOptions{Tail: 1000}).Return([]ports.LogEntry{}, nil)

		w := doJSON(r, http.MethodGet, path+"?tail=5000", "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("invalid tail returns 400", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil).Times(2)

		for _, tail := range []string{"abc", "0"} {
			w := doJSON(r, http.MethodGet, path+"?tail="+tail, "", authCookie(t, ownerID.String()))
			if w.Code != http.StatusBadRequest {
				t.Errorf("tail=%s: expected 400, got %d", tail, w.Code)
			}
		}
	})

	t.Run("invalid since returns 400", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)

		w := doJSON(r, http.MethodGet, path+"?since=yesterday", "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("app in another subdomain returns 404", func(t *testing.T) {
		r, m := newAppRouter(t)
		foreign := &domain.Application{ID: appID, SubdomainID: uuid.New()}
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.applicationRepo.EXPECT().FindByID(appID.String(), gomock.Any()).Return(foreign, nil)

		w := doJSON(r, http.MethodGet, path, "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("missing application returns 404", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.applicationRepo.EXPECT().FindByID(appID.String(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

		w := doJSON(r, http.MethodGet, path, "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("non-owner gets 403 and the manager is never called", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)

		intruder := uuid.New()
		w := doJSON(r, http.MethodGet, path, "", authCookie(t, intruder.String()))
		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", w.Code)
		}
	})

	t.Run("manager failure returns 502", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.applicationRepo.EXPECT().FindByID(appID.String(), gomock.Any()).Return(ownedApp, nil)
		m.serverManager.EXPECT().Logs(gomock.Any(), appID.String(), gomock.Any()).Return(nil, errors.New("connection refused"))

		w := doJSON(r, http.MethodGet, path, "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusBadGateway {
			t.Fatalf("expected 502, got %d", w.Code)
		}
		if strings.Contains(w.Body.String(), "connection refused") {
			t.Errorf("response leaks upstream error detail: %s", w.Body)
		}
	})
}

func TestStreamApplicationLogs(t *testing.T) {
	ownerID := uuid.New()
	subID := uuid.New()
	appID := uuid.New()
	owned := &domain.Subdomain{ID: subID, UserID: ownerID}
	ownedApp := &domain.Application{ID: appID, SubdomainID: subID}
	streamPath := "/subdomain/" + subID.String() + "/application/" + appID.String() + "/logs/stream"

	dial := func(t *testing.T, r *gin.Engine, cookie *http.Cookie) (*websocket.Conn, *http.Response, error) {
		t.Helper()
		srv := httptest.NewServer(r)
		t.Cleanup(srv.Close)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		t.Cleanup(cancel)
		headers := http.Header{}
		if cookie != nil {
			headers.Set("Cookie", cookie.String())
		}
		wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + streamPath
		return websocket.Dial(ctx, wsURL, &websocket.DialOptions{HTTPHeader: headers}) //nolint:bodyclose
	}

	t.Run("streams entries as websocket messages", func(t *testing.T) {
		r, m := newAppRouter(t)
		ts := time.Date(2026, 6, 13, 10, 0, 1, 0, time.UTC)
		ch := make(chan ports.LogEntry, 2)
		ch <- ports.LogEntry{ID: "1-0", Timestamp: ts, Level: "info", Message: "tick"}
		ch <- ports.LogEntry{ID: "2-1", Timestamp: ts.Add(time.Second), Level: "error", Message: "tock"}
		close(ch)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.applicationRepo.EXPECT().FindByID(appID.String(), gomock.Any()).Return(ownedApp, nil)
		m.serverManager.EXPECT().StreamLogs(gomock.Any(), appID.String(), ports.LogOptions{Tail: 200}).Return((<-chan ports.LogEntry)(ch), nil)

		conn, _, err := dial(t, r, authCookie(t, ownerID.String()))
		if err != nil {
			t.Fatalf("websocket dial: %v", err)
		}
		defer conn.Close(websocket.StatusNormalClosure, "done")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var got []map[string]any
		for range 2 {
			_, payload, err := conn.Read(ctx)
			if err != nil {
				t.Fatalf("read after %d messages: %v", len(got), err)
			}
			var entry map[string]any
			if err := json.Unmarshal(payload, &entry); err != nil {
				t.Fatalf("message not JSON: %v", err)
			}
			got = append(got, entry)
		}
		if got[0]["message"] != "tick" || got[0]["level"] != "info" {
			t.Errorf("first = %v", got[0])
		}
		if got[1]["message"] != "tock" || got[1]["level"] != "error" {
			t.Errorf("second = %v", got[1])
		}

		// Channel closed → server closes the socket normally.
		if _, _, err := conn.Read(ctx); websocket.CloseStatus(err) != websocket.StatusNormalClosure {
			t.Errorf("close status = %v (err %v), want normal closure", websocket.CloseStatus(err), err)
		}
	})

	t.Run("unauthenticated handshake is rejected", func(t *testing.T) {
		r, _ := newAppRouter(t)
		_, res, err := dial(t, r, nil)
		if err == nil {
			t.Fatal("dial succeeded without auth cookie")
		}
		if res != nil && res.StatusCode != http.StatusUnauthorized {
			t.Errorf("handshake status = %d, want 401", res.StatusCode)
		}
	})

	t.Run("not deployed rejects the handshake with 404", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.applicationRepo.EXPECT().FindByID(appID.String(), gomock.Any()).Return(ownedApp, nil)
		m.serverManager.EXPECT().StreamLogs(gomock.Any(), appID.String(), gomock.Any()).Return(nil, ports.ErrAppNotDeployed)

		_, res, err := dial(t, r, authCookie(t, ownerID.String()))
		if err == nil {
			t.Fatal("dial succeeded for undeployed app")
		}
		if res != nil && res.StatusCode != http.StatusNotFound {
			t.Errorf("handshake status = %d, want 404", res.StatusCode)
		}
	})

	t.Run("foreign app rejects the handshake with 404", func(t *testing.T) {
		r, m := newAppRouter(t)
		foreign := &domain.Application{ID: appID, SubdomainID: uuid.New()}
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.applicationRepo.EXPECT().FindByID(appID.String(), gomock.Any()).Return(foreign, nil)

		_, res, err := dial(t, r, authCookie(t, ownerID.String()))
		if err == nil {
			t.Fatal("dial succeeded for foreign app")
		}
		if res != nil && res.StatusCode != http.StatusNotFound {
			t.Errorf("handshake status = %d, want 404", res.StatusCode)
		}
	})
}

func TestUpdateApplication(t *testing.T) {
	ownerID := uuid.New()
	subID := uuid.New()
	appID := uuid.New()
	owned := &domain.Subdomain{ID: subID, UserID: ownerID}
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	path := "/subdomain/" + subID.String() + "/application/" + appID.String()

	t.Run("updates provided fields and normalizes type", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.clock.EXPECT().Now().Return(now)
		m.applicationRepo.EXPECT().Update(appID.String(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(id string, updates map[string]any, ctx any) error {
				if updates["name"] != "renamed" {
					t.Errorf("expected name renamed, got %v", updates["name"])
				}
				if updates["type"] != domain.ApplicationTypeNodejs {
					t.Errorf("expected normalized type Nodejs, got %v", updates["type"])
				}
				if len(updates) != 3 { // name, type, updated_at
					t.Errorf("expected 3 update keys, got %v", updates)
				}
				return nil
			})

		w := doJSON(r, http.MethodPatch, path, `{"name":"renamed","type":"node"}`, authCookie(t, ownerID.String()))
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d (body %s)", w.Code, w.Body.String())
		}
	})

	t.Run("unsupported type returns 400", func(t *testing.T) {
		r, _ := newAppRouter(t)
		w := doJSON(r, http.MethodPatch, path, `{"type":"python"}`, authCookie(t, ownerID.String()))
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("persists environment variables", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.clock.EXPECT().Now().Return(now)
		m.applicationRepo.EXPECT().Update(appID.String(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(id string, updates map[string]any, ctx any) error {
				env, ok := updates["environment_variables"].(map[string]string)
				if !ok {
					t.Errorf("expected environment_variables in updates, got %v", updates)
				} else if env["NODE_ENV"] != "production" || len(env) != 1 {
					t.Errorf("unexpected env vars: %v", env)
				}
				return nil
			})

		body := `{"environment_variables":{"NODE_ENV":"production"}}`
		w := doJSON(r, http.MethodPatch, path, body, authCookie(t, ownerID.String()))
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d (body %s)", w.Code, w.Body.String())
		}
	})

	t.Run("invalid env var key returns 400 naming the key only", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.clock.EXPECT().Now().Return(now)

		body := `{"environment_variables":{"1BAD-KEY":"topsecret"}}`
		w := doJSON(r, http.MethodPatch, path, body, authCookie(t, ownerID.String()))
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d (body %s)", w.Code, w.Body.String())
		}
		if !strings.Contains(w.Body.String(), "1BAD-KEY") {
			t.Errorf("expected the offending key in the message, got %s", w.Body.String())
		}
		if strings.Contains(w.Body.String(), "topsecret") {
			t.Errorf("env value leaked into the error response: %s", w.Body.String())
		}
	})

	t.Run("missing application returns 404", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.clock.EXPECT().Now().Return(now)
		m.applicationRepo.EXPECT().Update(appID.String(), gomock.Any(), gomock.Any()).Return(gorm.ErrRecordNotFound)

		w := doJSON(r, http.MethodPatch, path, `{"name":"renamed"}`, authCookie(t, ownerID.String()))
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})
}

func TestDeleteApplication(t *testing.T) {
	ownerID := uuid.New()
	subID := uuid.New()
	appID := uuid.New()
	owned := &domain.Subdomain{ID: subID, UserID: ownerID}
	path := "/subdomain/" + subID.String() + "/application/" + appID.String()

	t.Run("success", func(t *testing.T) {
		r, m := newAppRouter(t)
		application := &domain.Application{ID: appID}
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.applicationRepo.EXPECT().FindByID(appID.String(), gomock.Any()).Return(application, nil)
		m.serverManager.EXPECT().Remove(gomock.Any(), appID.String()).Return(nil)
		m.applicationRepo.EXPECT().Delete(application, gomock.Any()).Return(nil)

		w := doJSON(r, http.MethodDelete, path, "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d (body %s)", w.Code, w.Body.String())
		}
	})

	t.Run("missing application returns 404", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.applicationRepo.EXPECT().FindByID(appID.String(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

		w := doJSON(r, http.MethodDelete, path, "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("non-owner gets 403 without delete", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)

		intruder := uuid.New()
		w := doJSON(r, http.MethodDelete, path, "", authCookie(t, intruder.String()))
		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", w.Code)
		}
	})
}
