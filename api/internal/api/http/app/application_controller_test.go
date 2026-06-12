package app

import (
	"api/internal/api/middlewares"
	"api/internal/domain"
	"api/internal/infra/utils"
	"api/internal/mocks"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	subdomainhttp "api/internal/api/http/subdomain"
	appapp "api/internal/app/app"
	appsubdomain "api/internal/app/subdomain"
	infraauth "api/internal/infra/auth"

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
		clock:           mocks.NewMockClock(ctrl),
	}
	sdDeps := appsubdomain.Dependencies{
		SubdomainRepo: m.subdomainRepo,
		UserRepo:      m.userRepo,
		Clock:         m.clock,
	}
	appDeps := appapp.Dependencies{
		ApplicationRepo: m.applicationRepo,
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
			ID:            appID,
			Name:          &name,
			Type:          domain.ApplicationTypeNodejs,
			Status:        domain.ApplicationStatusRunning,
			RepositoryURL: &repoURL,
			Branch:        &branch,
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
