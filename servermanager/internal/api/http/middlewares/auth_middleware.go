package middlewares

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	Token string
}

// Auth requires "Authorization: Bearer <token>". Both sides are hashed before
// the constant-time compare so the token length doesn't leak either.
func (a *AuthMiddleware) Auth(c *gin.Context) {
	const prefix = "Bearer "
	header := c.GetHeader("Authorization")
	if len(header) <= len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	got := sha256.Sum256([]byte(header[len(prefix):]))
	want := sha256.Sum256([]byte(a.Token))
	if subtle.ConstantTimeCompare(got[:], want[:]) != 1 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.Next()
}
