package subdomain

import (
	"api/internal/api/middlewares"
	"api/internal/domain"
	"api/internal/infra/postgres"
	"api/internal/infra/utils"
	"api/internal/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

type subdomainMocks struct {
	subdomainRepo   *mocks.MockSubdomainRepo
	userRepo        *mocks.MockUserRepo
	applicationRepo *mocks.MockApplicationRepo
	databaseRepo    *mocks.MockDatabaseRepo
	serverManager   *mocks.MockServerManagerClient
	clock           *mocks.MockClock
}

type envelope struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func newSubdomainRouter(t *testing.T) (*gin.Engine, subdomainMocks) {
	ctrl := gomock.NewController(t)
	m := subdomainMocks{
		subdomainRepo:   mocks.NewMockSubdomainRepo(ctrl),
		userRepo:        mocks.NewMockUserRepo(ctrl),
		applicationRepo: mocks.NewMockApplicationRepo(ctrl),
		databaseRepo:    mocks.NewMockDatabaseRepo(ctrl),
		serverManager:   mocks.NewMockServerManagerClient(ctrl),
		clock:           mocks.NewMockClock(ctrl),
	}
	deps := appsubdomain.Dependencies{
		SubdomainRepo:   m.subdomainRepo,
		UserRepo:        m.userRepo,
		ApplicationRepo: m.applicationRepo,
		DatabaseRepo:    m.databaseRepo,
		ServerManager:   m.serverManager,
		Clock:           m.clock,
	}
	mw := middlewares.AuthMiddleware{
		JwtTokenizer: infraauth.JwtTokenizer{Clock: utils.SystemClock{}, SecretKey: testSecret},
	}
	r := gin.New()
	SetupRouter(deps, mw, r)
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

func TestGetAllSubdomains(t *testing.T) {
	ownerID := uuid.New()

	t.Run("no cookie returns 401", func(t *testing.T) {
		r, _ := newSubdomainRouter(t)
		w := doJSON(r, http.MethodGet, "/subdomain", "", nil)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("maps nested database and application", func(t *testing.T) {
		r, m := newSubdomainRouter(t)

		appName := "my-app"
		repoURL := "https://github.com/me/my-app"
		branch := "main"
		connStr := "postgres://user:pass@host/db"
		bare := domain.Subdomain{
			ID:         uuid.New(),
			Name:       "bare",
			FullDomain: "bare.studententuin.com",
			UserID:     ownerID,
			IsActive:   true,
		}
		full := domain.Subdomain{
			ID:         uuid.New(),
			Name:       "full",
			FullDomain: "full.studententuin.com",
			UserID:     ownerID,
			IsActive:   false,
			Database: &domain.Database{
				ID:               uuid.New(),
				Name:             "db",
				Type:             domain.DatabaseTypePostgres,
				Version:          "17",
				Status:           domain.DatabaseStatusRunning,
				ConnectionString: &connStr,
			},
			Application: &domain.Application{
				ID:            uuid.New(),
				Name:          &appName,
				Type:          domain.ApplicationTypeNodejs,
				Status:        domain.ApplicationStatusRunning,
				RepositoryURL: &repoURL,
				Branch:        &branch,
			},
		}
		m.subdomainRepo.EXPECT().FindAllByUserID(ownerID.String(), gomock.Any()).Return([]domain.Subdomain{bare, full}, nil)

		w := doJSON(r, http.MethodGet, "/subdomain", "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d (body %s)", w.Code, w.Body.String())
		}

		var resp envelope
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		var items []map[string]any
		if err := json.Unmarshal(resp.Data, &items); err != nil {
			t.Fatalf("failed to decode data: %v", err)
		}
		if len(items) != 2 {
			t.Fatalf("expected 2 items, got %d", len(items))
		}
		if items[0]["name"] != "bare" || items[0]["database"] != nil || items[0]["application"] != nil {
			t.Errorf("expected bare subdomain without resources, got %v", items[0])
		}
		db, _ := items[1]["database"].(map[string]any)
		if db == nil || db["type"] != "postgres" || db["connection_string"] != connStr {
			t.Errorf("expected nested database, got %v", items[1]["database"])
		}
		app, _ := items[1]["application"].(map[string]any)
		if app == nil || app["name"] != appName || app["repo_url"] != repoURL {
			t.Errorf("expected nested application, got %v", items[1]["application"])
		}
	})

	t.Run("repo error returns 500", func(t *testing.T) {
		r, m := newSubdomainRouter(t)
		m.subdomainRepo.EXPECT().FindAllByUserID(ownerID.String(), gomock.Any()).Return(nil, errors.New("db down"))

		w := doJSON(r, http.MethodGet, "/subdomain", "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})
}

func TestCreateSubdomain(t *testing.T) {
	ownerID := uuid.New()

	tests := []struct {
		name       string
		body       string
		setup      func(m subdomainMocks)
		wantStatus int
	}{
		{
			name: "success",
			body: `{"name":"myapp","fullDomain":"myapp.studententuin.com"}`,
			setup: func(m subdomainMocks) {
				m.subdomainRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
					func(s *domain.Subdomain, ctx any) error {
						if s.UserID != ownerID {
							t.Errorf("expected owner %s, got %s", ownerID, s.UserID)
						}
						if !s.IsActive {
							t.Error("expected new subdomain to be active")
						}
						return nil
					})
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "duplicate domain returns 409",
			body: `{"name":"myapp","fullDomain":"taken.studententuin.com"}`,
			setup: func(m subdomainMocks) {
				m.subdomainRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(postgres.ErrFullDomainAlreadyInUse)
			},
			wantStatus: http.StatusConflict,
		},
		{
			name:       "missing fullDomain returns 400",
			body:       `{"name":"myapp"}`,
			setup:      func(m subdomainMocks) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "repo error returns 500",
			body: `{"name":"myapp","fullDomain":"myapp.studententuin.com"}`,
			setup: func(m subdomainMocks) {
				m.subdomainRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("db down"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, m := newSubdomainRouter(t)
			tt.setup(m)

			w := doJSON(r, http.MethodPost, "/subdomain", tt.body, authCookie(t, ownerID.String()))
			if w.Code != tt.wantStatus {
				t.Fatalf("expected %d, got %d (body %s)", tt.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestGetSubdomain_Ownership(t *testing.T) {
	ownerID := uuid.New()
	subID := uuid.New()
	owned := &domain.Subdomain{
		ID:         subID,
		Name:       "myapp",
		FullDomain: "myapp.studententuin.com",
		UserID:     ownerID,
		IsActive:   true,
	}

	t.Run("owner gets mapped subdomain", func(t *testing.T) {
		r, m := newSubdomainRouter(t)
		// Once for CheckOwnership, once for the Get usecase.
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil).Times(2)

		w := doJSON(r, http.MethodGet, "/subdomain/"+subID.String(), "", authCookie(t, ownerID.String()))
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
		if data["id"] != subID.String() || data["name"] != "myapp" || data["fullDomain"] != "myapp.studententuin.com" || data["isActive"] != true {
			t.Errorf("unexpected subdomain mapping: %v", data)
		}
	})

	t.Run("non-owner gets 403", func(t *testing.T) {
		r, m := newSubdomainRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)

		intruder := uuid.New()
		w := doJSON(r, http.MethodGet, "/subdomain/"+subID.String(), "", authCookie(t, intruder.String()))
		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d (body %s)", w.Code, w.Body.String())
		}
	})

	t.Run("missing subdomain returns 404", func(t *testing.T) {
		r, m := newSubdomainRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

		w := doJSON(r, http.MethodGet, "/subdomain/"+subID.String(), "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("repo error returns 500", func(t *testing.T) {
		r, m := newSubdomainRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(nil, errors.New("db down"))

		w := doJSON(r, http.MethodGet, "/subdomain/"+subID.String(), "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("no cookie returns 401", func(t *testing.T) {
		r, _ := newSubdomainRouter(t)
		w := doJSON(r, http.MethodGet, "/subdomain/"+subID.String(), "", nil)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})
}

func TestUpdateSubdomain(t *testing.T) {
	ownerID := uuid.New()
	subID := uuid.New()
	owned := &domain.Subdomain{ID: subID, UserID: ownerID}
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Run("owner updates provided fields only", func(t *testing.T) {
		r, m := newSubdomainRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.clock.EXPECT().Now().Return(now)
		m.subdomainRepo.EXPECT().Update(subID.String(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(id string, updates map[string]any, ctx any) error {
				if updates["name"] != "renamed" {
					t.Errorf("expected name renamed, got %v", updates["name"])
				}
				if len(updates) != 2 { // name, updated_at
					t.Errorf("expected 2 update keys, got %v", updates)
				}
				return nil
			})

		w := doJSON(r, http.MethodPatch, "/subdomain/"+subID.String(), `{"name":"renamed"}`, authCookie(t, ownerID.String()))
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d (body %s)", w.Code, w.Body.String())
		}
	})

	t.Run("non-owner gets 403 without update", func(t *testing.T) {
		r, m := newSubdomainRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)

		intruder := uuid.New()
		w := doJSON(r, http.MethodPatch, "/subdomain/"+subID.String(), `{"name":"hacked"}`, authCookie(t, intruder.String()))
		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", w.Code)
		}
	})

	t.Run("update on missing record returns 404", func(t *testing.T) {
		r, m := newSubdomainRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.clock.EXPECT().Now().Return(now)
		m.subdomainRepo.EXPECT().Update(subID.String(), gomock.Any(), gomock.Any()).Return(gorm.ErrRecordNotFound)

		w := doJSON(r, http.MethodPatch, "/subdomain/"+subID.String(), `{"name":"renamed"}`, authCookie(t, ownerID.String()))
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})
}

func TestDeleteSubdomain(t *testing.T) {
	ownerID := uuid.New()
	subID := uuid.New()
	owned := &domain.Subdomain{ID: subID, UserID: ownerID}

	t.Run("owner deletes subdomain", func(t *testing.T) {
		r, m := newSubdomainRouter(t)
		// Once for CheckOwnership, once inside the delete usecase.
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil).Times(2)
		m.applicationRepo.EXPECT().FindBySubdomainID(subID.String(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
		m.databaseRepo.EXPECT().FindBySubdomainID(subID.String(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
		m.subdomainRepo.EXPECT().Delete(owned, gomock.Any()).Return(nil)

		w := doJSON(r, http.MethodDelete, "/subdomain/"+subID.String(), "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d (body %s)", w.Code, w.Body.String())
		}
	})

	t.Run("non-owner gets 403 without delete", func(t *testing.T) {
		r, m := newSubdomainRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)

		intruder := uuid.New()
		w := doJSON(r, http.MethodDelete, "/subdomain/"+subID.String(), "", authCookie(t, intruder.String()))
		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", w.Code)
		}
	})

	t.Run("missing subdomain returns 404", func(t *testing.T) {
		r, m := newSubdomainRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

		w := doJSON(r, http.MethodDelete, "/subdomain/"+subID.String(), "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})
}
