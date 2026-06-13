package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	middleware := AuthMiddleware{Token: "secret-token"}
	router.GET("/protected", middleware.Auth, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	return router
}

func request(router *gin.Engine, authHeader string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestAuthValidToken(t *testing.T) {
	w := request(setupRouter(), "Bearer secret-token")
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestAuthRejects(t *testing.T) {
	cases := map[string]string{
		"missing header": "",
		"wrong token":    "Bearer wrong-token",
		"wrong scheme":   "Basic secret-token",
		"empty bearer":   "Bearer ",
		"token only":     "secret-token",
	}
	router := setupRouter()
	for name, header := range cases {
		t.Run(name, func(t *testing.T) {
			if w := request(router, header); w.Code != http.StatusUnauthorized {
				t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
			}
		})
	}
}
