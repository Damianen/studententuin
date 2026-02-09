package middlewares

import (
	"api/internal/infra/auth"
	"api/internal/infra/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestMiddleware() (*AuthMiddleware, *auth.JwtTokenizer) {
	tokenizer := &auth.JwtTokenizer{
		Clock:     utils.SystemClock{},
		SecretKey: "test-secret-key",
	}
	return &AuthMiddleware{JwtTokenizer: *tokenizer}, tokenizer
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	mw, tokenizer := newTestMiddleware()

	token, err := tokenizer.CreateToken("user-123")
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	var capturedUserID any
	r := gin.New()
	r.Use(mw.Auth)
	r.GET("/test", func(c *gin.Context) {
		capturedUserID, _ = c.Get("userID")
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{Name: "AuthToken", Value: token})
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	userIDStr, ok := capturedUserID.(string)
	if !ok {
		t.Fatal("expected userID to be a string")
	}
	if userIDStr != "user-123" {
		t.Fatalf("expected user-123, got %s", userIDStr)
	}
}

func TestAuthMiddleware_MissingCookie(t *testing.T) {
	mw, _ := newTestMiddleware()

	r := gin.New()
	r.Use(mw.Auth)
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	mw, _ := newTestMiddleware()

	r := gin.New()
	r.Use(mw.Auth)
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{Name: "AuthToken", Value: "invalid.token.here"})
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
