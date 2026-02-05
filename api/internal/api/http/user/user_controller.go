package user

import (
	"api/internal/api/dtos"
	"api/internal/api/middlewares"
	"api/internal/app/user"
	"api/internal/infra/postgres"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	context := ginc.Request.Context()

	var req dtos.CreateUserRequest
	err := ginc.ShouldBindBodyWithJSON(&req)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 400, "invalid JSON or missing values", nil)
		return
	}

	userInput := user.UserInput{
		Email: req.Email,
		Name: req.Name,
		Password: req.Password,
	}

	err = c.service.Create.Execute(context, userInput)
	if errors.Is(err, postgres.ErrEmailAlreadyInUse) {
    	middlewares.Respond(ginc, 409, "email already in use", nil)
    	return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "failed to create user", nil)
		return
	}

	middlewares.Respond(ginc, 201, "success", nil)
}

func (c *Controller) Update(ginc *gin.Context) {
	context := ginc.Request.Context()

	var req dtos.UpdateUserRequest
	err := ginc.ShouldBindBodyWithJSON(&req)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 400, "invalid JSON or missing values", nil)
		return
	}

	userID := ginc.GetString("userID")
	userInput := user.UserUpdateInput{
		ID: userID,
		Email: req.Email,
		Name: req.Name,
	}

	err = c.service.Update.Execute(context, userInput)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		middlewares.Respond(ginc, 404, "user not found", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "failed to update user", nil)
		return
	}

	middlewares.Respond(ginc, 204, "success", nil)
}


func (c *Controller) Delete(ginc *gin.Context) {
	context := ginc.Request.Context()

	userID := ginc.GetString("userID")

	err := c.service.Delete.Execute(context, userID)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "failed to delete user", nil)
		return
	}

	middlewares.Respond(ginc, 204, "success", nil)
}

func (c *Controller) Get(ginc *gin.Context) {
	context := ginc.Request.Context()

	userID := ginc.GetString("userID")

	user, err := c.service.Get.Execute(context, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		middlewares.Respond(ginc, 404, "user not found", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "failed to get user", nil)
		return
	}

	userResponse := dtos.UserResponse{
		Email:       user.Email,
    	DisplayName: user.DisplayName,
    	Status:      user.Status,
	}

	middlewares.Respond(ginc, 200, "success", userResponse)
}
