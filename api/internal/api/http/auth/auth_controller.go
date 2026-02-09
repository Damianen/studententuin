package auth

import (
	"api/internal/api/dtos"
	"api/internal/api/middlewares"
	"api/internal/app/auth"
	"api/internal/app/ports"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *auth.Service
	clock ports.Clock
}

func NewController(deps auth.Dependencies) (*Controller) {
	return &Controller{ service: auth.NewService(deps), clock: deps.Clock}
}

func (c *Controller) Login(ginc *gin.Context) {
	context := ginc.Request.Context()

	var req dtos.LoginUserRequest
	err := ginc.ShouldBindBodyWithJSON(&req)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusBadRequest, "invalid JSON or missing values", nil)
		return
	}

	loginInput := auth.LoginInput{
		Email: req.Email,
		Password: req.Password,
	}

	token, err := c.service.Login.Execute(context, loginInput)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusUnauthorized, "email or password not correct!", nil)
		return
	}

	ginc.SetCookieData(&http.Cookie{
		Name: "AuthToken",
		Value: token,
		Path: "/",
		Expires: c.clock.Now().Add(24 * time.Hour),
		MaxAge: 86400,
		Secure: false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	middlewares.Respond(ginc, http.StatusOK, "login was successful", nil)
}

func (c *Controller) Logout(ginc *gin.Context) {
	ginc.SetCookieData(&http.Cookie{
		Name: "AuthToken",
		Value: "",
		Path: "/",
		MaxAge: -1,
		Secure: false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	middlewares.Respond(ginc, http.StatusOK, "logout was successful", nil)
}
