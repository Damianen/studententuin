package user

import (
	"api/internal/api/middlewares"
	"api/internal/app/user"

	"github.com/gin-gonic/gin"
)

func SetupRouter(d user.Dependencies, middleware middlewares.AuthMiddleware, r *gin.Engine) {
	c := NewController(d)

	user := r.Group("/user")
	{
		user.POST("/register", c.Create)
		user.GET("", middleware.Auth, c.Get)
		user.DELETE("", middleware.Auth, c.Delete)
		user.PATCH("", middleware.Auth, c.Update)
	}
}
