package user

import (
	"api/internal/api/dtos"
	"api/internal/api/middlewares"
	"api/internal/app/user"
	"api/internal/infra/postgres"
	"errors"
	"fmt"
	"net/http"

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
	var req dtos.CreateUserRequest
	if !middlewares.BindJSON(ginc, &req) {
		return
	}

	userInput := user.UserInput{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	}

	err := c.service.Create.Execute(ginc.Request.Context(), userInput)
	if errors.Is(err, postgres.ErrEmailAlreadyInUse) {
		middlewares.Respond(ginc, http.StatusConflict, "email already in use", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusInternalServerError, "failed to create user", nil)
		return
	}

	middlewares.Respond(ginc, http.StatusCreated, "success", nil)
}

func (c *Controller) Update(ginc *gin.Context) {
	var req dtos.UpdateUserRequest
	if !middlewares.BindJSON(ginc, &req) {
		return
	}

	userID := ginc.GetString("userID")
	userInput := user.UserUpdateInput{
		ID:    userID,
		Email: req.Email,
		Name:  req.Name,
	}

	err := c.service.Update.Execute(ginc.Request.Context(), userInput)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		middlewares.Respond(ginc, http.StatusNotFound, "user not found", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusInternalServerError, "failed to update user", nil)
		return
	}

	middlewares.Respond(ginc, http.StatusOK, "success", nil)
}

func (c *Controller) Delete(ginc *gin.Context) {
	userID := ginc.GetString("userID")

	err := c.service.Delete.Execute(ginc.Request.Context(), userID)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusInternalServerError, "failed to delete user", nil)
		return
	}

	middlewares.Respond(ginc, http.StatusOK, "success", nil)
}

func (c *Controller) Get(ginc *gin.Context) {
	userID := ginc.GetString("userID")

	user, err := c.service.Get.Execute(ginc.Request.Context(), userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		middlewares.Respond(ginc, http.StatusNotFound, "user not found", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusInternalServerError, "failed to get user", nil)
		return
	}

	userResponse := dtos.UserResponse{
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Status:      user.Status,
	}

	middlewares.Respond(ginc, http.StatusOK, "success", userResponse)
}
