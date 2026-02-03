package middlewares

import (
	"api/internal/infra/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(secret []byte) gin.HandlerFunc {
	return func(c * gin.Context) {
		tokenString, err := c.Cookie("AuthToken")
		if err != nil || tokenString == "" {
			Respond(c, http.StatusUnauthorized, "authentication required", nil)
			c.Abort()
			return
		}

		claims := &auth.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secret, nil
		})

		if err != nil || !token.Valid || claims.UserID == "" {
			Respond(c, http.StatusUnauthorized, "ivalid token", nil)
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)

		c.Next()
	}
}
