package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"servermanager/internal/domain"
)

// maxHealthGrace keeps a typo'd SM_HEALTH_GRACE from stalling every deploy.
const maxHealthGrace = 60 * time.Second

// maxDBHealthBudget bounds SM_DB_HEALTH_BUDGET the same way.
const maxDBHealthBudget = 5 * time.Minute

// Metrics bounds (§6 phase 6): the interval keeps the sampler cheap but the
// deltas meaningful; retention bounds the in-memory store.
const (
	minMetricsInterval  = 5 * time.Second
	maxMetricsInterval  = 5 * time.Minute
	minMetricsRetention = time.Hour
	maxMetricsRetention = 168 * time.Hour
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

	GitHosts      []string
	CloneTimeout  time.Duration
	CloneMaxBytes int64
	BuildTimeout  time.Duration
	HealthGrace   time.Duration
	NixpacksBin   string

	DBPullTimeout  time.Duration
	DBHealthBudget time.Duration

	MetricsInterval  time.Duration
	MetricsRetention time.Duration
	CgroupRoot       string
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

	for _, host := range strings.Split(getenv("SM_GIT_HOSTS", "github.com,gitlab.com"), ",") {
		if host = strings.ToLower(strings.TrimSpace(host)); host != "" {
			cfg.GitHosts = append(cfg.GitHosts, host)
		}
	}
	if len(cfg.GitHosts) == 0 {
		return nil, fmt.Errorf("SM_GIT_HOSTS must list at least one host")
	}

	if cfg.CloneTimeout, err = positiveDuration("SM_CLONE_TIMEOUT", "120s"); err != nil {
		return nil, err
	}
	// A byte size, but ParseMemoryLimit is exactly the parser it needs.
	if cfg.CloneMaxBytes, err = domain.ParseMemoryLimit(getenv("SM_CLONE_MAX_SIZE", "200m")); err != nil {
		return nil, fmt.Errorf("SM_CLONE_MAX_SIZE: %w", err)
	}
	if cfg.BuildTimeout, err = positiveDuration("SM_BUILD_TIMEOUT", "10m"); err != nil {
		return nil, err
	}
	if cfg.HealthGrace, err = positiveDuration("SM_HEALTH_GRACE", "3s"); err != nil {
		return nil, err
	}
	if cfg.HealthGrace > maxHealthGrace {
		return nil, fmt.Errorf("SM_HEALTH_GRACE exceeds the %s maximum", maxHealthGrace)
	}
	if cfg.NixpacksBin = getenv("SM_NIXPACKS_BIN", "nixpacks"); cfg.NixpacksBin == "" {
		return nil, fmt.Errorf("SM_NIXPACKS_BIN must not be empty")
	}

	// Database provisioning (§6 phase 5): the first pull of an official image
	// dominates; the health budget covers initdb on slow disks.
	if cfg.DBPullTimeout, err = positiveDuration("SM_DB_PULL_TIMEOUT", "5m"); err != nil {
		return nil, err
	}
	if cfg.DBHealthBudget, err = positiveDuration("SM_DB_HEALTH_BUDGET", "60s"); err != nil {
		return nil, err
	}
	if cfg.DBHealthBudget > maxDBHealthBudget {
		return nil, fmt.Errorf("SM_DB_HEALTH_BUDGET exceeds the %s maximum", maxDBHealthBudget)
	}

	// Metrics (§6 phase 6). The default cgroup root works uncontainerized in
	// local dev; compose mounts the host tree at /host/cgroup.
	if cfg.MetricsInterval, err = positiveDuration("SM_METRICS_INTERVAL", "30s"); err != nil {
		return nil, err
	}
	if cfg.MetricsInterval < minMetricsInterval || cfg.MetricsInterval > maxMetricsInterval {
		return nil, fmt.Errorf("SM_METRICS_INTERVAL must be between %s and %s", minMetricsInterval, maxMetricsInterval)
	}
	if cfg.MetricsRetention, err = positiveDuration("SM_METRICS_RETENTION", "24h"); err != nil {
		return nil, err
	}
	if cfg.MetricsRetention < minMetricsRetention || cfg.MetricsRetention > maxMetricsRetention {
		return nil, fmt.Errorf("SM_METRICS_RETENTION must be between %s and %s", minMetricsRetention, maxMetricsRetention)
	}
	if cfg.MetricsRetention <= cfg.MetricsInterval {
		return nil, fmt.Errorf("SM_METRICS_RETENTION must exceed SM_METRICS_INTERVAL")
	}
	if cfg.CgroupRoot = getenv("SM_CGROUP_ROOT", "/sys/fs/cgroup"); cfg.CgroupRoot == "" {
		return nil, fmt.Errorf("SM_CGROUP_ROOT must not be empty")
	}

	return cfg, nil
}

func positiveDuration(key, fallback string) (time.Duration, error) {
	d, err := time.ParseDuration(getenv(key, fallback))
	if err != nil || d <= 0 {
		return 0, fmt.Errorf("%s must be a positive duration like %q", key, fallback)
	}
	return d, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
