package auth

import (
	"api/internal/app/auth"

	"github.com/gin-gonic/gin"
)

func SetupRouter(d auth.Dependencies, r *gin.Engine) {
	c := NewController(d)

	auth := r.Group("/auth")
	{
		auth.POST("login", c.Login)
		auth.POST("logout", c.Logout)
	}
}
