package config

import (
	"fmt"
	"os"
	"strconv"

	"servermanager/internal/domain"
)

type Config struct {
	Port           string
	BindAddr       string
	Token          string
	DefaultRuntime string

	DefaultMemoryBytes int64
	MaxMemoryBytes     int64
	DefaultNanoCPUs    int64
	MaxNanoCPUs        int64
	DefaultPidsLimit   int64
}

// Load reads the SM_* environment (SERVERMANAGEMENT.md §9) and fails on a
// missing token or any invalid value so the manager never starts half-configured.
func Load() (*Config, error) {
	cfg := &Config{
		Port:           getenv("SM_PORT", "8080"),
		BindAddr:       getenv("SM_BIND_ADDR", "127.0.0.1"),
		Token:          os.Getenv("SM_TOKEN"),
		DefaultRuntime: getenv("SM_DEFAULT_RUNTIME", domain.RuntimeRunc),
	}

	if cfg.Token == "" {
		return nil, fmt.Errorf("SM_TOKEN is required")
	}
	if !domain.ValidRuntime(cfg.DefaultRuntime) {
		return nil, fmt.Errorf("SM_DEFAULT_RUNTIME %q is not on the runtime allowlist", cfg.DefaultRuntime)
	}

	var err error
	if cfg.DefaultMemoryBytes, err = domain.ParseMemoryLimit(getenv("SM_DEFAULT_MEMORY", "256m")); err != nil {
		return nil, fmt.Errorf("SM_DEFAULT_MEMORY: %w", err)
	}
	if cfg.MaxMemoryBytes, err = domain.ParseMemoryLimit(getenv("SM_MAX_MEMORY", "1g")); err != nil {
		return nil, fmt.Errorf("SM_MAX_MEMORY: %w", err)
	}
	if cfg.DefaultNanoCPUs, err = domain.ParseCPULimit(getenv("SM_DEFAULT_CPU", "0.5")); err != nil {
		return nil, fmt.Errorf("SM_DEFAULT_CPU: %w", err)
	}
	if cfg.MaxNanoCPUs, err = domain.ParseCPULimit(getenv("SM_MAX_CPU", "2")); err != nil {
		return nil, fmt.Errorf("SM_MAX_CPU: %w", err)
	}
	if cfg.DefaultPidsLimit, err = strconv.ParseInt(getenv("SM_DEFAULT_PIDS", "256"), 10, 64); err != nil || cfg.DefaultPidsLimit <= 0 {
		return nil, fmt.Errorf("SM_DEFAULT_PIDS must be a positive integer")
	}

	if cfg.DefaultMemoryBytes > cfg.MaxMemoryBytes {
		return nil, fmt.Errorf("SM_DEFAULT_MEMORY exceeds SM_MAX_MEMORY")
	}
	if cfg.DefaultNanoCPUs > cfg.MaxNanoCPUs {
		return nil, fmt.Errorf("SM_DEFAULT_CPU exceeds SM_MAX_CPU")
	}

	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
