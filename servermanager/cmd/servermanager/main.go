package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	appsHTTP "servermanager/internal/api/http/apps"
	"servermanager/internal/api/http/middlewares"
	"servermanager/internal/app/apps"
	"servermanager/internal/config"
	"servermanager/internal/infra/docker"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	// The manager is useless without Docker — fail fast like a config error
	// instead of surfacing per-request 500s.
	runtime, err := docker.New(context.Background())
	if err != nil {
		slog.Error("docker unreachable", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			slog.Warn("closing docker client", "error", err)
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(middlewares.RequestLogger(), gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"server": "running"})
	})

	auth := middlewares.AuthMiddleware{Token: cfg.Token}
	v1 := router.Group("/v1", auth.Auth)
	appsHTTP.SetupRouter(apps.Dependencies{
		Runtime: runtime,
		Limits: apps.Limits{
			DefaultMemoryBytes: cfg.DefaultMemoryBytes,
			MaxMemoryBytes:     cfg.MaxMemoryBytes,
			DefaultNanoCPUs:    cfg.DefaultNanoCPUs,
			MaxNanoCPUs:        cfg.MaxNanoCPUs,
			DefaultPidsLimit:   cfg.DefaultPidsLimit,
			DefaultRuntime:     cfg.DefaultRuntime,
		},
	}, v1)

	addr := net.JoinHostPort(cfg.BindAddr, cfg.Port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}
	slog.Info("servermanager listening", "addr", addr, "default_runtime", cfg.DefaultRuntime)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
