package db

import (
	"api/internal/api/middlewares"
	"api/internal/app/db"
	"api/internal/app/subdomain"

	"github.com/gin-gonic/gin"
)

func SetupRouter(d db.Dependencies, sd subdomain.Dependencies,  middleware middlewares.AuthMiddleware, group *gin.RouterGroup) {
	c := NewController(d, sd)

	group.POST("/:subId/database", middleware.Auth, c.Create)
	group.DELETE("/:subId/database/:dbId", middleware.Auth, c.Delete)
	group.PATCH("/:subId/database/:dbId", middleware.Auth, c.Update)
	group.GET("/:subId/database/:dbId", middleware.Auth, c.Get)
}
