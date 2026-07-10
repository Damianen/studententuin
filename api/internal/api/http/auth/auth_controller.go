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
	clock   ports.Clock
}

func NewController(deps auth.Dependencies) *Controller {
	return &Controller{service: auth.NewService(deps), clock: deps.Clock}
}

func (c *Controller) Login(ginc *gin.Context) {
	var req dtos.LoginUserRequest
	if !middlewares.BindJSON(ginc, &req) {
		return
	}

	loginInput := auth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	token, err := c.service.Login.Execute(ginc.Request.Context(), loginInput)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusUnauthorized, "email or password not correct!", nil)
		return
	}

	ginc.SetCookieData(&http.Cookie{
		Name:     "AuthToken",
		Value:    token,
		Path:     "/",
		Expires:  c.clock.Now().Add(24 * time.Hour),
		MaxAge:   86400,
		// Browsers exempt localhost from the Secure-over-HTTPS rule, so this
		// stays compatible with local dev and the compose-based e2e stack.
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	middlewares.Respond(ginc, http.StatusOK, "login was successful", nil)
}

func (c *Controller) Logout(ginc *gin.Context) {
	ginc.SetCookieData(&http.Cookie{
		Name:     "AuthToken",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	middlewares.Respond(ginc, http.StatusOK, "logout was successful", nil)
}
