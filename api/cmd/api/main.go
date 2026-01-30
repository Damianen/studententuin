package main

import (
	"api/internal/api/http/user"
	"api/internal/api/http/auth"
	appUser "api/internal/app/user"
	appAuth "api/internal/app/auth"
	authUtils "api/internal/infra/auth"
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
	hasher := authUtils.NewBcryptHasher(10)
	jwtTokenizer := authUtils.JwtTokenizer{
		
	}
	userDeps := appUser.Dependencies{
		UserRepo: &userRepo,
		PasswordRepo: &passwordRepo,
		Clock: clock,
		Hasher: hasher,
	}
	authDeps := appAuth.Dependencies{
		UserRepo: &userRepo,
		PasswordRepo: &passwordRepo,
		Clock: clock,
		Hasher: hasher,
		JwtTokenizer: &jwtTokenizer,
	}

	user.SetupRouter(userDeps, router)
	auth.SetupRouter(authDeps, router)


	router.Run(":8080")
}
