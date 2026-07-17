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
	"api/internal/infra/servermanager"
	"api/internal/infra/utils"
	"context"
	"fmt"
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
	deploymentRepo := postgres.GormDeploymentRepo{ DB: db }
	// Polling state is in-memory and single-instance: at boot every in_flight
	// history row is orphaned by definition. Best-effort — history never
	// blocks startup.
	if err := deploymentRepo.FailInFlight("interrupted by an api restart", context.Background()); err != nil {
		fmt.Println("sweeping in-flight deployments:", err.Error())
	}
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
	// The servermanager is root-equivalent on the hosting server; refuse to
	// start half-configured rather than failing per request.
	smURL := os.Getenv("SERVERMANAGER_URL")
	smToken := os.Getenv("SERVERMANAGER_TOKEN")
	if smURL == "" || smToken == "" {
		panic("SERVERMANAGER_URL and SERVERMANAGER_TOKEN are required")
	}
	smClient := servermanager.NewClient(smURL, smToken)

	subdomainDeps := appSubdomain.Dependencies{
		UserRepo: &userRepo,
		SubdomainRepo: &subdomainRepo,
		ApplicationRepo: &appRepo,
		DatabaseRepo: &dbRepo,
		ServerManager: smClient,
		Clock: &clock,
	}
	dbDeps := appDb.Dependencies{
		DatabaseRepo: &dbRepo,
		Clock: &clock,
		ServerManager: smClient,
	}
	appDeps := appApp.Dependencies{
		ApplicationRepo: &appRepo,
		DatabaseRepo: &dbRepo,
		DeploymentRepo: &deploymentRepo,
		ServerManager: smClient,
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
	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}
