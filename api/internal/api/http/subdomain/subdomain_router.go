package subdomain

import (
	"api/internal/api/middlewares"
	"api/internal/app/subdomain"

	"github.com/gin-gonic/gin"
)

func SetupRouter(d subdomain.Dependencies, middleware middlewares.AuthMiddleware, r *gin.Engine) {
	c := NewController(d)

	subdomain := r.Group("/subdomain")
	{
		subdomain.GET("/:id", middleware.Auth, c.Get)
		subdomain.DELETE("/:id", middleware.Auth, c.Delete)
		subdomain.PATCH("/:id", middleware.Auth, c.Update)
		subdomain.POST("", middleware.Auth, c.Create)
	}
}
