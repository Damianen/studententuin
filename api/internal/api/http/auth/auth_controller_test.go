package auth

import (
	"api/internal/domain"
	"api/internal/mocks"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	appauth "api/internal/app/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type authMocks struct {
	userRepo     *mocks.MockUserRepo
	passwordRepo *mocks.MockPasswordRepo
	clock        *mocks.MockClock
	hasher       *mocks.MockPasswordHasher
	jwt          *mocks.MockJwtTokenizer
}

func newAuthRouter(t *testing.T) (*gin.Engine, authMocks) {
	ctrl := gomock.NewController(t)
	m := authMocks{
		userRepo:     mocks.NewMockUserRepo(ctrl),
		passwordRepo: mocks.NewMockPasswordRepo(ctrl),
		clock:        mocks.NewMockClock(ctrl),
		hasher:       mocks.NewMockPasswordHasher(ctrl),
		jwt:          mocks.NewMockJwtTokenizer(ctrl),
	}
	deps := appauth.Dependencies{
		UserRepo:     m.userRepo,
		PasswordRepo: m.passwordRepo,
		Clock:        m.clock,
		Hasher:       m.hasher,
		JwtTokenizer: m.jwt,
	}
	r := gin.New()
	SetupRouter(deps, r)
	return r, m
}

func findCookie(t *testing.T, w *httptest.ResponseRecorder, name string) *http.Cookie {
	t.Helper()
	for _, c := range w.Result().Cookies() {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func TestLogin(t *testing.T) {
	userID := uuid.New()
	user := &domain.User{ID: userID, Email: "test@example.com"}
	cred := &domain.PasswordCredential{UserId: userID, PasswordHash: "hashed"}

	tests := []struct {
		name       string
		body       string
		setup      func(m authMocks)
		wantStatus int
		wantCookie string
	}{
		{
			name: "success sets auth cookie",
			body: `{"email":"test@example.com","password":"secret123"}`,
			setup: func(m authMocks) {
				m.userRepo.EXPECT().FindByEmail("test@example.com", gomock.Any()).Return(user, nil)
				m.passwordRepo.EXPECT().FindById(userID.String(), gomock.Any()).Return(cred, nil)
				m.hasher.EXPECT().Compare("secret123", "hashed").Return(true, nil)
				m.jwt.EXPECT().CreateToken(userID.String()).Return("tok123", nil)
				m.clock.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
			},
			wantStatus: http.StatusOK,
			wantCookie: "tok123",
		},
		{
			name: "wrong password returns 401 without cookie",
			body: `{"email":"test@example.com","password":"wrong"}`,
			setup: func(m authMocks) {
				m.userRepo.EXPECT().FindByEmail("test@example.com", gomock.Any()).Return(user, nil)
				m.passwordRepo.EXPECT().FindById(userID.String(), gomock.Any()).Return(cred, nil)
				m.hasher.EXPECT().Compare("wrong", "hashed").Return(false, nil)
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "unknown user returns 401",
			body: `{"email":"nope@example.com","password":"secret123"}`,
			setup: func(m authMocks) {
				m.userRepo.EXPECT().FindByEmail("nope@example.com", gomock.Any()).Return(nil, errors.New("user not found"))
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "malformed JSON returns 400",
			body:       `{"email":`,
			setup:      func(m authMocks) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing password returns 400",
			body:       `{"email":"test@example.com"}`,
			setup:      func(m authMocks) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid email format returns 400",
			body:       `{"email":"not-an-email","password":"secret123"}`,
			setup:      func(m authMocks) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, m := newAuthRouter(t)
			tt.setup(m)

			req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Fatalf("expected %d, got %d (body %s)", tt.wantStatus, w.Code, w.Body.String())
			}

			cookie := findCookie(t, w, "AuthToken")
			if tt.wantCookie == "" {
				if cookie != nil {
					t.Fatalf("expected no AuthToken cookie, got %q", cookie.Value)
				}
				return
			}
			if cookie == nil {
				t.Fatal("expected AuthToken cookie, got none")
			}
			if cookie.Value != tt.wantCookie {
				t.Errorf("expected cookie value %q, got %q", tt.wantCookie, cookie.Value)
			}
			if !cookie.HttpOnly {
				t.Error("expected cookie to be HttpOnly")
			}
			if cookie.MaxAge != 86400 {
				t.Errorf("expected Max-Age 86400, got %d", cookie.MaxAge)
			}
			if cookie.Path != "/" {
				t.Errorf("expected Path /, got %q", cookie.Path)
			}
		})
	}
}

func TestLogout_ClearsCookie(t *testing.T) {
	r, _ := newAuthRouter(t)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	cookie := findCookie(t, w, "AuthToken")
	if cookie == nil {
		t.Fatal("expected clearing AuthToken cookie, got none")
	}
	if cookie.Value != "" {
		t.Errorf("expected empty cookie value, got %q", cookie.Value)
	}
	// Go serializes MaxAge -1 as Max-Age=0, which parses back as MaxAge < 0.
	if cookie.MaxAge >= 0 {
		t.Errorf("expected expired cookie (MaxAge < 0), got %d", cookie.MaxAge)
	}
}
