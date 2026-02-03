package main

import (
	"api/internal/api/http/auth"
	"api/internal/api/http/user"
	appAuth "api/internal/app/auth"
	appUser "api/internal/app/user"
	authUtils "api/internal/infra/auth"
	"api/internal/infra/postgres"
	"api/internal/infra/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health", func (c *gin.Context) {
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
	jwtToken := os.Getenv("JWT_TOKEN")
	jwtTokenizer := authUtils.JwtTokenizer{
		SecretKey: jwtToken,
		Clock: clock,
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
