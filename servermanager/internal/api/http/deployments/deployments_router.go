package deployments

import (
	"servermanager/internal/app/deployments"

	"github.com/gin-gonic/gin"
)

// SetupRouter mounts the deploy trigger next to the other app routes and the
// job resource under its own prefix. Both live behind the group's bearer
// auth, like everything under /v1.
func SetupRouter(d deployments.Dependencies, group *gin.RouterGroup) {
	c := NewController(d)
	group.POST("/apps/:id/deploy", c.Deploy)
	group.GET("/deployments/:id", c.Get)
}
