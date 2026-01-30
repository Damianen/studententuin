package user

import (
	"api/internal/api/dtos"
	"api/internal/app/user"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *user.Service
}

func NewController(d user.Dependencies) *Controller {
	return &Controller{
		service: user.NewService(d),
	}
}

func (c *Controller) Create(ginc *gin.Context) {
	context := ginc.Request.Context();

	var req dtos.CreateUserRequest
	err := ginc.ShouldBindBodyWithJSON(&req)
	if err != nil {
		fmt.Println(err.Error())
		ginc.JSON(400, gin.H{"error": "invalid JSON or missing values"})
	}

	userInput := user.UserInput{
		Email: req.Email,
		Name: req.Name,
		Password: req.Password,
	}

	err = c.service.Create.Execute(context, userInput)
	if err != nil {
		fmt.Println(err.Error())
		ginc.JSON(500, gin.H{"error": "failed to create user"})
		return
	}

	ginc.JSON(200, gin.H{"success": "user was created!"})
}
