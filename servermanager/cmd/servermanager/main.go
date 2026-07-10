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
	deploymentsHTTP "servermanager/internal/api/http/deployments"
	"servermanager/internal/api/http/middlewares"
	"servermanager/internal/app/apps"
	"servermanager/internal/app/deployments"
	"servermanager/internal/config"
	"servermanager/internal/infra/build"
	"servermanager/internal/infra/clock"
	"servermanager/internal/infra/docker"
	gitinfra "servermanager/internal/infra/git"

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

	// Like Docker, a missing nixpacks binary is a startup failure, not a
	// per-deploy 500.
	builder, err := build.New(context.Background(), cfg.NixpacksBin, runtime)
	if err != nil {
		slog.Error("nixpacks unavailable", "error", err)
		os.Exit(1)
	}

	appsDeps := apps.Dependencies{
		Runtime: runtime,
		Limits: apps.Limits{
			DefaultMemoryBytes: cfg.DefaultMemoryBytes,
			MaxMemoryBytes:     cfg.MaxMemoryBytes,
			DefaultNanoCPUs:    cfg.DefaultNanoCPUs,
			MaxNanoCPUs:        cfg.MaxNanoCPUs,
			DefaultPidsLimit:   cfg.DefaultPidsLimit,
			DefaultRuntime:     cfg.DefaultRuntime,
		},
	}

	auth := middlewares.AuthMiddleware{Token: cfg.Token}
	v1 := router.Group("/v1", auth.Auth)
	appsHTTP.SetupRouter(appsDeps, v1)
	deploymentsHTTP.SetupRouter(deployments.Dependencies{
		Fetcher:  gitinfra.NewFetcher(cfg.GitHosts, cfg.CloneMaxBytes),
		Builder:  builder,
		Runtime:  runtime,
		Images:   runtime,
		Runner:   apps.NewService(appsDeps).Run,
		Clock:    clock.System{},
		GitHosts: cfg.GitHosts,
		Budgets: deployments.Budgets{
			CloneTimeout: cfg.CloneTimeout,
			BuildTimeout: cfg.BuildTimeout,
			HealthGrace:  cfg.HealthGrace,
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
