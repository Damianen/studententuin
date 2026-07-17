package databases

import (
	"servermanager/internal/app/databases"

	"github.com/gin-gonic/gin"
)

// SetupRouter registers the database lifecycle routes on the (already
// authenticated) /v1 group. Package databases, routes /dbs (§3.2).
func SetupRouter(d databases.Dependencies, group *gin.RouterGroup) {
	c := NewController(d)

	dbGroup := group.Group("/dbs/:id")
	{
		dbGroup.POST("/provision", c.Provision)
		dbGroup.GET("", c.Status)
		dbGroup.DELETE("", c.Delete)
		dbGroup.GET("/logs", c.Logs)
		dbGroup.GET("/logs/stream", c.StreamLogs)
		dbGroup.GET("/metrics", c.Metrics)
	}
}
