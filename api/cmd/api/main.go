package main

import (
	"api/internal/api/http/auth"
	"api/internal/api/http/subdomain"
	"api/internal/api/http/user"
	database "api/internal/api/http/db"
	"api/internal/api/http/app"

	"api/internal/api/middlewares"

	appAuth "api/internal/app/auth"
	appUser "api/internal/app/user"
	authUtils "api/internal/infra/auth"
	appSubdomain "api/internal/app/subdomain"
	appDb "api/internal/app/db"
	appApp "api/internal/app/app"

	"api/internal/infra/postgres"
	"api/internal/infra/utils"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	router.Use(gin.Logger(), gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

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
	subdomainRepo := postgres.GormSubdomainRepo{ DB: db }
	dbRepo := postgres.GormDatabaseRepo{ DB: db }
	appRepo := postgres.GormApplicationRepo{ DB: db }
	clock := utils.SystemClock{}
	hasher := authUtils.NewBcryptHasher(10)
	jwtToken := os.Getenv("JWT_TOKEN")
	jwtTokenizer := authUtils.JwtTokenizer{
		SecretKey: jwtToken,
		Clock: clock,
	}
	middleware := middlewares.AuthMiddleware{
		JwtTokenizer: jwtTokenizer,
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
	subdomainDeps := appSubdomain.Dependencies{
		UserRepo: &userRepo,
		SubdomainRepo: &subdomainRepo,
		Clock: &clock,
	}
	dbDeps := appDb.Dependencies{
		DatabaseRepo: &dbRepo,
		Clock: &clock,
	}
	appDeps := appApp.Dependencies{
		ApplicationRepo: &appRepo,
		Clock: &clock,
	}

	user.SetupRouter(userDeps, middleware, router)
	auth.SetupRouter(authDeps, router)
	subGroup := subdomain.SetupRouter(subdomainDeps, middleware, router)
	database.SetupRouter(dbDeps, subdomainDeps, middleware, subGroup)
	app.SetupRouter(appDeps, subdomainDeps, middleware, subGroup)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
