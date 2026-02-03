package middlewares

import (
	"api/internal/infra/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	JwtTokenizer auth.JwtTokenizer
}

func (a *AuthMiddleware) Auth(c *gin.Context) {
	tokenString, err := c.Cookie("AuthToken")
	if err != nil || tokenString == "" {
		Respond(c, http.StatusUnauthorized, "authentication required", nil)
		c.Abort()
		return
	}

	userID, err := a.JwtTokenizer.VerifyToken(tokenString)
	if err != nil {
		Respond(c, http.StatusUnauthorized, "authentication required", nil)
		c.Abort()
		return
	}

	c.Set("userID", userID)

	c.Next()
}
