package user

import (
	"api/internal/app/user"

	"github.com/gin-gonic/gin"
)

func SetupRouter(d user.Dependencies, r *gin.Engine) {
	c := NewController(d)

	user := r.Group("/user")
	{
		user.POST("/register", c.Create)
	}
}
