package app

import (
	"api/internal/api/middlewares"
	"api/internal/app/app"
	"api/internal/app/subdomain"

	"github.com/gin-gonic/gin"
)

func SetupRouter(d app.Dependencies, sd subdomain.Dependencies, middleware middlewares.AuthMiddleware, group *gin.RouterGroup) {
	c := NewController(d, sd)

	appGroup := group.Group(":id/application")
	{
		appGroup.POST("", middleware.Auth, c.Create)
		appGroup.DELETE("/:appId", middleware.Auth, c.Delete)
		appGroup.PATCH("/:appId", middleware.Auth, c.Updates)
		appGroup.GET("/:appId", middleware.Auth, c.Get)
		appGroup.GET("/:appId/logs", middleware.Auth, c.GetLogs)
		appGroup.GET("/:appId/logs/stream", middleware.Auth, c.StreamLogs)
		appGroup.GET("/:appId/metrics", middleware.Auth, c.GetMetrics)
		appGroup.POST("/:appId/deploy", middleware.Auth, c.Deploy)
		appGroup.GET("/:appId/deployment/:deploymentId", middleware.Auth, c.GetDeployment)
		appGroup.GET("/:appId/deployments", middleware.Auth, c.ListDeployments)
		appGroup.POST("/:appId/start", middleware.Auth, c.Start)
		appGroup.POST("/:appId/stop", middleware.Auth, c.Stop)
	}
}
