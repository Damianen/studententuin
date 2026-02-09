package db

import (
	"api/internal/api/middlewares"
	"api/internal/app/db"
	"api/internal/app/subdomain"

	"github.com/gin-gonic/gin"
)

func SetupRouter(d db.Dependencies, sd subdomain.Dependencies,  middleware middlewares.AuthMiddleware, group *gin.RouterGroup) {
	c := NewController(d, sd)

	dbGroup := group.Group(":id/database")
	{
		dbGroup.POST("", middleware.Auth, c.Create)
		dbGroup.DELETE("/:dbId", middleware.Auth, c.Delete)
		dbGroup.PATCH("/:dbId", middleware.Auth, c.Update)
		dbGroup.GET("/:dbId", middleware.Auth, c.Get)
	}
}
