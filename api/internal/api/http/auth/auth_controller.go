package auth

import (
	"api/internal/api/dtos"
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
		ginc.JSON(400, gin.H{"error": "invalid JSON or missing values"})
	}

	loginInput := auth.LoginInput{
		Email: req.Email,
		Password: req.Password,
	}

	token, err := c.service.Login.Execute(context, loginInput)
	if err != nil {
		fmt.Println(err.Error())
		ginc.JSON(401, gin.H{"error": "email or password not correct!"})
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
	ginc.JSON(200, gin.H{"success": "login was successful"})
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
	ginc.JSON(200, gin.H{"success": "logout was successful"})
}
