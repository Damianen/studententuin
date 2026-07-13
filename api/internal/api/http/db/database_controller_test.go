package db

import (
	"api/internal/api/middlewares"
	"api/internal/app/ports"
	"api/internal/domain"
	"api/internal/infra/utils"
	"api/internal/mocks"
	"encoding/json"
	"errors"
	"fmt"
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
	serverManager *mocks.MockServerManagerClient
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
		serverManager: mocks.NewMockServerManagerClient(ctrl),
		clock:         mocks.NewMockClock(ctrl),
	}
	m.clock.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)).AnyTimes()
	sdDeps := appsubdomain.Dependencies{
		SubdomainRepo: m.subdomainRepo,
		UserRepo:      m.userRepo,
		Clock:         m.clock,
	}
	dbDeps := appdb.Dependencies{
		DatabaseRepo:  m.databaseRepo,
		Clock:         m.clock,
		ServerManager: m.serverManager,
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
	validBody := `{"name":"mydb","type":"postgres","version":"17","db_name":"app"}`

	t.Run("success returns 201, stores provisioning database, fires provision", func(t *testing.T) {
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
		m.serverManager.EXPECT().ProvisionDatabase(gomock.Any(), gomock.Any(), gomock.Any()).Return("prov-1", nil)
		// The poller starts polling; park it on a non-terminal status.
		m.serverManager.EXPECT().DatabaseStatus(gomock.Any(), gomock.Any()).
			Return(&ports.DBStatus{Provision: &ports.DBProvisionState{Status: "pulling"}}, nil).AnyTimes()

		w := doJSON(r, http.MethodPost, "/subdomain/"+subID.String()+"/database", validBody, authCookie(t, ownerID.String()))
		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d (body %s)", w.Code, w.Body.String())
		}
	})

	t.Run("engine outside the allowlist returns 400", func(t *testing.T) {
		r, _ := newDbRouter(t)
		for _, body := range []string{
			`{"name":"mydb","type":"redis","version":"7","db_name":"app"}`,
			`{"name":"mydb","type":"mysql","version":"8.0","db_name":"app"}`,
			`{"name":"mydb","type":"postgres","version":"15","db_name":"app"}`,
			`{"name":"mydb","type":"postgres","version":"17","db_name":"Bad Name"}`,
		} {
			w := doJSON(r, http.MethodPost, "/subdomain/"+subID.String()+"/database", body, authCookie(t, ownerID.String()))
			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected 400 for %s, got %d", body, w.Code)
			}
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
			SubdomainID:      subID,
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

	t.Run("a foreign database behind your own subdomain is 404, not a leak", func(t *testing.T) {
		r, m := newDbRouter(t)
		connStr := "postgres://user:secret@foreign/db"
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.databaseRepo.EXPECT().FindByID(dbID.String(), gomock.Any()).Return(&domain.Database{
			ID:               dbID,
			SubdomainID:      uuid.New(), // someone else's subdomain
			ConnectionString: &connStr,
		}, nil)

		w := doJSON(r, http.MethodGet, path, "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d (body %s)", w.Code, w.Body.String())
		}
		if strings.Contains(w.Body.String(), "secret") {
			t.Error("foreign connection string leaked")
		}
	})
}

func TestUpdateDatabase(t *testing.T) {
	ownerID := uuid.New()
	subID := uuid.New()
	dbID := uuid.New()
	owned := &domain.Subdomain{ID: subID, UserID: ownerID}
	path := "/subdomain/" + subID.String() + "/database/" + dbID.String()

	t.Run("renames the database", func(t *testing.T) {
		r, m := newDbRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.databaseRepo.EXPECT().FindByID(dbID.String(), gomock.Any()).
			Return(&domain.Database{ID: dbID, SubdomainID: subID}, nil)
		m.databaseRepo.EXPECT().Update(dbID.String(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(id string, updates map[string]any, ctx any) error {
				if updates["name"] != "renamed" {
					t.Errorf("expected name renamed, got %v", updates["name"])
				}
				if len(updates) != 2 { // name, updated_at
					t.Errorf("expected 2 update keys, got %v", updates)
				}
				return nil
			})

		w := doJSON(r, http.MethodPatch, path, `{"name":"renamed"}`, authCookie(t, ownerID.String()))
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d (body %s)", w.Code, w.Body.String())
		}
	})

	t.Run("engine changes are rejected", func(t *testing.T) {
		r, _ := newDbRouter(t)
		for _, body := range []string{`{"type":"redis"}`, `{"version":"18"}`} {
			w := doJSON(r, http.MethodPatch, path, body, authCookie(t, ownerID.String()))
			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected 400 for %s, got %d", body, w.Code)
			}
		}
	})

	t.Run("missing database returns 404", func(t *testing.T) {
		r, m := newDbRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.databaseRepo.EXPECT().FindByID(dbID.String(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

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

	t.Run("success removes the container before the row", func(t *testing.T) {
		r, m := newDbRouter(t)
		database := &domain.Database{ID: dbID, SubdomainID: subID}
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		gomock.InOrder(
			m.databaseRepo.EXPECT().FindByID(dbID.String(), gomock.Any()).Return(database, nil),
			m.serverManager.EXPECT().RemoveDatabase(gomock.Any(), dbID.String()).Return(nil),
			m.databaseRepo.EXPECT().Delete(database, gomock.Any()).Return(nil),
		)

		w := doJSON(r, http.MethodDelete, path, "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d (body %s)", w.Code, w.Body.String())
		}
	})

	t.Run("provision in flight returns 409", func(t *testing.T) {
		r, m := newDbRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.databaseRepo.EXPECT().FindByID(dbID.String(), gomock.Any()).
			Return(&domain.Database{ID: dbID, SubdomainID: subID}, nil)
		m.serverManager.EXPECT().RemoveDatabase(gomock.Any(), dbID.String()).
			Return(fmt.Errorf("servermanager remove database: %w", ports.ErrProvisionInFlight))

		w := doJSON(r, http.MethodDelete, path, "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d (body %s)", w.Code, w.Body.String())
		}
	})

	t.Run("a foreign database is 404 and never reaches the manager", func(t *testing.T) {
		r, m := newDbRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.databaseRepo.EXPECT().FindByID(dbID.String(), gomock.Any()).
			Return(&domain.Database{ID: dbID, SubdomainID: uuid.New()}, nil)

		w := doJSON(r, http.MethodDelete, path, "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
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
