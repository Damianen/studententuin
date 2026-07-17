package apps

import (
	"servermanager/internal/app/apps"

	"github.com/gin-gonic/gin"
)

// SetupRouter registers the app lifecycle routes on the (already
// authenticated) /v1 group.
func SetupRouter(d apps.Dependencies, group *gin.RouterGroup) {
	c := NewController(d)

	appGroup := group.Group("/apps/:id")
	{
		appGroup.POST("/run", c.Run)
		appGroup.POST("/start", c.Start)
		appGroup.POST("/stop", c.Stop)
		appGroup.DELETE("", c.Remove)
		appGroup.GET("", c.Status)
		appGroup.GET("/logs", c.Logs)
		appGroup.GET("/logs/stream", c.StreamLogs)
		appGroup.GET("/metrics", c.Metrics)
	}
}
