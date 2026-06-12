package user

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

	appuser "api/internal/app/user"
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

type userMocks struct {
	userRepo     *mocks.MockUserRepo
	passwordRepo *mocks.MockPasswordRepo
	clock        *mocks.MockClock
	hasher       *mocks.MockPasswordHasher
}

type envelope struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func newUserRouter(t *testing.T) (*gin.Engine, userMocks) {
	ctrl := gomock.NewController(t)
	m := userMocks{
		userRepo:     mocks.NewMockUserRepo(ctrl),
		passwordRepo: mocks.NewMockPasswordRepo(ctrl),
		clock:        mocks.NewMockClock(ctrl),
		hasher:       mocks.NewMockPasswordHasher(ctrl),
	}
	deps := appuser.Dependencies{
		UserRepo:     m.userRepo,
		PasswordRepo: m.passwordRepo,
		Clock:        m.clock,
		Hasher:       m.hasher,
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
	var reader *strings.Reader
	if body == "" {
		reader = strings.NewReader("")
	} else {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	if cookie != nil {
		req.AddCookie(cookie)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestRegister(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		body       string
		setup      func(m userMocks)
		wantStatus int
	}{
		{
			name: "success",
			body: `{"email":"new@example.com","password":"secret123","name":"Alice"}`,
			setup: func(m userMocks) {
				m.clock.EXPECT().Now().Return(now)
				m.hasher.EXPECT().Hash("secret123").Return("hashed", nil)
				m.userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				m.passwordRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "duplicate email returns 409",
			body: `{"email":"taken@example.com","password":"secret123","name":"Bob"}`,
			setup: func(m userMocks) {
				m.clock.EXPECT().Now().Return(now)
				m.hasher.EXPECT().Hash("secret123").Return("hashed", nil)
				m.userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(postgres.ErrEmailAlreadyInUse)
			},
			wantStatus: http.StatusConflict,
		},
		{
			name: "repo error returns 500",
			body: `{"email":"new@example.com","password":"secret123","name":"Alice"}`,
			setup: func(m userMocks) {
				m.clock.EXPECT().Now().Return(now)
				m.hasher.EXPECT().Hash("secret123").Return("hashed", nil)
				m.userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("db down"))
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "invalid email returns 400",
			body:       `{"email":"not-an-email","password":"secret123","name":"Alice"}`,
			setup:      func(m userMocks) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "short password returns 400",
			body:       `{"email":"new@example.com","password":"short","name":"Alice"}`,
			setup:      func(m userMocks) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing name returns 400",
			body:       `{"email":"new@example.com","password":"secret123"}`,
			setup:      func(m userMocks) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "malformed JSON returns 400",
			body:       `{"email":`,
			setup:      func(m userMocks) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, m := newUserRouter(t)
			tt.setup(m)

			w := doJSON(r, http.MethodPost, "/user/register", tt.body, nil)
			if w.Code != tt.wantStatus {
				t.Fatalf("expected %d, got %d (body %s)", tt.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	userID := uuid.New()

	t.Run("no cookie returns 401", func(t *testing.T) {
		r, _ := newUserRouter(t)
		w := doJSON(r, http.MethodGet, "/user", "", nil)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("invalid token returns 401", func(t *testing.T) {
		r, _ := newUserRouter(t)
		w := doJSON(r, http.MethodGet, "/user", "", &http.Cookie{Name: "AuthToken", Value: "garbage.token.value"})
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("success maps user to DTO", func(t *testing.T) {
		r, m := newUserRouter(t)
		m.userRepo.EXPECT().FindByID(userID.String(), gomock.Any()).Return(&domain.User{
			ID:          userID,
			Email:       "test@example.com",
			DisplayName: "Alice",
			Status:      "active",
		}, nil)

		w := doJSON(r, http.MethodGet, "/user", "", authCookie(t, userID.String()))
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
		if data["email"] != "test@example.com" {
			t.Errorf("expected email test@example.com, got %v", data["email"])
		}
		if data["name"] != "Alice" {
			t.Errorf("expected name Alice, got %v", data["name"])
		}
		if data["status"] != "active" {
			t.Errorf("expected status active, got %v", data["status"])
		}
	})

	t.Run("not found returns 404", func(t *testing.T) {
		r, m := newUserRouter(t)
		m.userRepo.EXPECT().FindByID(userID.String(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

		w := doJSON(r, http.MethodGet, "/user", "", authCookie(t, userID.String()))
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("repo error returns 500", func(t *testing.T) {
		r, m := newUserRouter(t)
		m.userRepo.EXPECT().FindByID(userID.String(), gomock.Any()).Return(nil, errors.New("db down"))

		w := doJSON(r, http.MethodGet, "/user", "", authCookie(t, userID.String()))
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})
}

func TestUpdateUser(t *testing.T) {
	userID := uuid.New()
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Run("success updates provided fields only", func(t *testing.T) {
		r, m := newUserRouter(t)
		m.clock.EXPECT().Now().Return(now)
		m.userRepo.EXPECT().Update(userID.String(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(id string, updates map[string]any, ctx any) error {
				if updates["email"] != "new@example.com" {
					t.Errorf("expected lowercased email, got %v", updates["email"])
				}
				if updates["display_name"] != "New Name" {
					t.Errorf("expected display_name New Name, got %v", updates["display_name"])
				}
				if len(updates) != 3 { // email, display_name, updated_at
					t.Errorf("expected 3 update keys, got %v", updates)
				}
				return nil
			})

		w := doJSON(r, http.MethodPatch, "/user", `{"email":"New@Example.com","name":"New Name"}`, authCookie(t, userID.String()))
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d (body %s)", w.Code, w.Body.String())
		}
	})

	t.Run("empty body returns 500 (no fields to update)", func(t *testing.T) {
		r, m := newUserRouter(t)
		m.clock.EXPECT().Now().Return(now)

		w := doJSON(r, http.MethodPatch, "/user", `{}`, authCookie(t, userID.String()))
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("not found returns 404", func(t *testing.T) {
		r, m := newUserRouter(t)
		m.clock.EXPECT().Now().Return(now)
		m.userRepo.EXPECT().Update(userID.String(), gomock.Any(), gomock.Any()).Return(gorm.ErrRecordNotFound)

		w := doJSON(r, http.MethodPatch, "/user", `{"name":"New Name"}`, authCookie(t, userID.String()))
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("no cookie returns 401", func(t *testing.T) {
		r, _ := newUserRouter(t)
		w := doJSON(r, http.MethodPatch, "/user", `{"name":"New Name"}`, nil)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})
}

func TestDeleteUser(t *testing.T) {
	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		r, m := newUserRouter(t)
		user := &domain.User{ID: userID}
		m.userRepo.EXPECT().FindByID(userID.String(), gomock.Any()).Return(user, nil)
		m.userRepo.EXPECT().Delete(user, gomock.Any()).Return(nil)

		w := doJSON(r, http.MethodDelete, "/user", "", authCookie(t, userID.String()))
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d (body %s)", w.Code, w.Body.String())
		}
	})

	t.Run("repo error returns 500", func(t *testing.T) {
		r, m := newUserRouter(t)
		m.userRepo.EXPECT().FindByID(userID.String(), gomock.Any()).Return(nil, errors.New("db down"))

		w := doJSON(r, http.MethodDelete, "/user", "", authCookie(t, userID.String()))
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("no cookie returns 401", func(t *testing.T) {
		r, _ := newUserRouter(t)
		w := doJSON(r, http.MethodDelete, "/user", "", nil)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})
}
