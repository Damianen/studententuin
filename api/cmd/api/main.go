package main

import (
	"api/internal/api/http/user"
	appUser "api/internal/app/user"
	"api/internal/infra/auth"
	"api/internal/infra/postgres"
	"api/internal/infra/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/heath", func (c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"server": "running",
		})
	})

	db, err := postgres.ConnectDB()
	if err != nil {
	 	panic(err)
	}
	userRepo := postgres.GormUserRepo{ DB: db }
	passwordRepo := postgres.GormPasswordRepo{ DB: db }
	clock := utils.SystemClock{}
	hasher := auth.NewBcryptHasher(10)
	userDepdencies := appUser.Dependencies{
		UserRepo: &userRepo,
		PasswordRepo: &passwordRepo,
		Clock: clock,
		Hasher: hasher,
	}

	user.SetupRouter(userDepdencies, router)

	router.Run(":8080")
}
