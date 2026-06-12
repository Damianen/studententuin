package db

import (
	"api/internal/api/middlewares"
	"api/internal/domain"
	"api/internal/infra/utils"
	"api/internal/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	subdomainhttp "api/internal/api/http/subdomain"
	appdb "api/internal/app/db"
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

type dbMocks struct {
	databaseRepo  *mocks.MockDatabaseRepo
	subdomainRepo *mocks.MockSubdomainRepo
	userRepo      *mocks.MockUserRepo
	clock         *mocks.MockClock
}

type envelope struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func newDbRouter(t *testing.T) (*gin.Engine, dbMocks) {
	ctrl := gomock.NewController(t)
	m := dbMocks{
		databaseRepo:  mocks.NewMockDatabaseRepo(ctrl),
		subdomainRepo: mocks.NewMockSubdomainRepo(ctrl),
		userRepo:      mocks.NewMockUserRepo(ctrl),
		clock:         mocks.NewMockClock(ctrl),
	}
	sdDeps := appsubdomain.Dependencies{
		SubdomainRepo: m.subdomainRepo,
		UserRepo:      m.userRepo,
		Clock:         m.clock,
	}
	dbDeps := appdb.Dependencies{
		DatabaseRepo: m.databaseRepo,
		Clock:        m.clock,
	}
	mw := middlewares.AuthMiddleware{
		JwtTokenizer: infraauth.JwtTokenizer{Clock: utils.SystemClock{}, SecretKey: testSecret},
	}
	r := gin.New()
	group := subdomainhttp.SetupRouter(sdDeps, mw, r)
	SetupRouter(dbDeps, sdDeps, mw, group)
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

func TestCreateDatabase(t *testing.T) {
	ownerID := uuid.New()
	subID := uuid.New()
	owned := &domain.Subdomain{ID: subID, UserID: ownerID}
	validBody := `{"name":"mydb","type":"postgres","version":"17","db_name":"app","db_password":"secret123"}`

	t.Run("success returns 201 and stores provisioning database", func(t *testing.T) {
		r, m := newDbRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.databaseRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *domain.Database, ctx any) error {
				if db.SubdomainID != subID {
					t.Errorf("expected subdomain %s, got %s", subID, db.SubdomainID)
				}
				if db.Type != domain.DatabaseTypePostgres {
					t.Errorf("expected type postgres, got %s", db.Type)
				}
				if db.Status != domain.DatabaseStatusProvisioning {
					t.Errorf("expected status provisioning, got %s", db.Status)
				}
				return nil
			})

		w := doJSON(r, http.MethodPost, "/subdomain/"+subID.String()+"/database", validBody, authCookie(t, ownerID.String()))
		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d (body %s)", w.Code, w.Body.String())
		}
	})

	t.Run("invalid engine type returns 400", func(t *testing.T) {
		r, _ := newDbRouter(t)
		body := `{"name":"mydb","type":"redis","version":"7","db_name":"app","db_password":"secret123"}`
		w := doJSON(r, http.MethodPost, "/subdomain/"+subID.String()+"/database", body, authCookie(t, ownerID.String()))
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("missing field returns 400", func(t *testing.T) {
		r, _ := newDbRouter(t)
		body := `{"name":"mydb","type":"postgres","version":"17"}`
		w := doJSON(r, http.MethodPost, "/subdomain/"+subID.String()+"/database", body, authCookie(t, ownerID.String()))
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("non-owner gets 403", func(t *testing.T) {
		r, m := newDbRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)

		intruder := uuid.New()
		w := doJSON(r, http.MethodPost, "/subdomain/"+subID.String()+"/database", validBody, authCookie(t, intruder.String()))
		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", w.Code)
		}
	})

	t.Run("no cookie returns 401", func(t *testing.T) {
		r, _ := newDbRouter(t)
		w := doJSON(r, http.MethodPost, "/subdomain/"+subID.String()+"/database", validBody, nil)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})
}

func TestGetDatabase(t *testing.T) {
	ownerID := uuid.New()
	subID := uuid.New()
	dbID := uuid.New()
	owned := &domain.Subdomain{ID: subID, UserID: ownerID}
	path := "/subdomain/" + subID.String() + "/database/" + dbID.String()

	t.Run("success maps database to DTO", func(t *testing.T) {
		r, m := newDbRouter(t)
		connStr := "postgres://user:pass@host/db"
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.databaseRepo.EXPECT().FindByID(dbID.String(), gomock.Any()).Return(&domain.Database{
			ID:               dbID,
			Name:             "mydb",
			Type:             domain.DatabaseTypePostgres,
			Version:          "17",
			Status:           domain.DatabaseStatusRunning,
			ConnectionString: &connStr,
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
		if data["id"] != dbID.String() || data["name"] != "mydb" || data["type"] != "postgres" ||
			data["version"] != "17" || data["status"] != "running" || data["connection_string"] != connStr {
			t.Errorf("unexpected database mapping: %v", data)
		}
	})

	t.Run("missing database returns 404", func(t *testing.T) {
		r, m := newDbRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.databaseRepo.EXPECT().FindByID(dbID.String(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

		w := doJSON(r, http.MethodGet, path, "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("non-owner gets 403", func(t *testing.T) {
		r, m := newDbRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)

		intruder := uuid.New()
		w := doJSON(r, http.MethodGet, path, "", authCookie(t, intruder.String()))
		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", w.Code)
		}
	})
}

func TestUpdateDatabase(t *testing.T) {
	ownerID := uuid.New()
	subID := uuid.New()
	dbID := uuid.New()
	owned := &domain.Subdomain{ID: subID, UserID: ownerID}
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	path := "/subdomain/" + subID.String() + "/database/" + dbID.String()

	t.Run("updates provided fields only", func(t *testing.T) {
		r, m := newDbRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.clock.EXPECT().Now().Return(now)
		m.databaseRepo.EXPECT().Update(dbID.String(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(id string, updates map[string]any, ctx any) error {
				if updates["name"] != "renamed" {
					t.Errorf("expected name renamed, got %v", updates["name"])
				}
				if updates["version"] != "18" {
					t.Errorf("expected version 18, got %v", updates["version"])
				}
				if len(updates) != 3 { // name, version, updated_at
					t.Errorf("expected 3 update keys, got %v", updates)
				}
				return nil
			})

		w := doJSON(r, http.MethodPatch, path, `{"name":"renamed","version":"18"}`, authCookie(t, ownerID.String()))
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d (body %s)", w.Code, w.Body.String())
		}
	})

	t.Run("invalid engine type returns 400", func(t *testing.T) {
		r, _ := newDbRouter(t)
		w := doJSON(r, http.MethodPatch, path, `{"type":"redis"}`, authCookie(t, ownerID.String()))
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("missing database returns 404", func(t *testing.T) {
		r, m := newDbRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.clock.EXPECT().Now().Return(now)
		m.databaseRepo.EXPECT().Update(dbID.String(), gomock.Any(), gomock.Any()).Return(gorm.ErrRecordNotFound)

		w := doJSON(r, http.MethodPatch, path, `{"name":"renamed"}`, authCookie(t, ownerID.String()))
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})
}

func TestDeleteDatabase(t *testing.T) {
	ownerID := uuid.New()
	subID := uuid.New()
	dbID := uuid.New()
	owned := &domain.Subdomain{ID: subID, UserID: ownerID}
	path := "/subdomain/" + subID.String() + "/database/" + dbID.String()

	t.Run("success", func(t *testing.T) {
		r, m := newDbRouter(t)
		database := &domain.Database{ID: dbID}
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.databaseRepo.EXPECT().FindByID(dbID.String(), gomock.Any()).Return(database, nil)
		m.databaseRepo.EXPECT().Delete(database, gomock.Any()).Return(nil)

		w := doJSON(r, http.MethodDelete, path, "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d (body %s)", w.Code, w.Body.String())
		}
	})

	t.Run("missing database returns 404", func(t *testing.T) {
		r, m := newDbRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.databaseRepo.EXPECT().FindByID(dbID.String(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

		w := doJSON(r, http.MethodDelete, path, "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("repo error returns 500", func(t *testing.T) {
		r, m := newDbRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.databaseRepo.EXPECT().FindByID(dbID.String(), gomock.Any()).Return(nil, errors.New("db down"))

		w := doJSON(r, http.MethodDelete, path, "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})
}
