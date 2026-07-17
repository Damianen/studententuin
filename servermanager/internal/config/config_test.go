package config

import (
	"testing"
	"time"
)

// clearEnv blanks every SM_* var so tests don't inherit the developer's shell.
func clearEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"SM_PORT", "SM_BIND_ADDR", "SM_TOKEN", "SM_DEFAULT_RUNTIME",
		"SM_DEFAULT_MEMORY", "SM_MAX_MEMORY", "SM_DEFAULT_CPU", "SM_MAX_CPU",
		"SM_DEFAULT_PIDS",
		"SM_GIT_HOSTS", "SM_CLONE_TIMEOUT", "SM_CLONE_MAX_SIZE",
		"SM_BUILD_TIMEOUT", "SM_HEALTH_GRACE", "SM_NIXPACKS_BIN",
		"SM_METRICS_INTERVAL", "SM_METRICS_RETENTION", "SM_CGROUP_ROOT",
	} {
		t.Setenv(key, "")
	}
}

func TestLoadMissingToken(t *testing.T) {
	clearEnv(t)
	if _, err := Load(); err == nil {
		t.Fatal("expected error when SM_TOKEN is missing")
	}
}

func TestLoadDefaults(t *testing.T) {
	clearEnv(t)
	t.Setenv("SM_TOKEN", "test-token")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}
	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want 8080", cfg.Port)
	}
	if cfg.BindAddr != "127.0.0.1" {
		t.Errorf("BindAddr = %q, want 127.0.0.1", cfg.BindAddr)
	}
	if cfg.DefaultRuntime != "runc" {
		t.Errorf("DefaultRuntime = %q, want runc", cfg.DefaultRuntime)
	}
	if cfg.DefaultMemoryBytes != 256*1024*1024 {
		t.Errorf("DefaultMemoryBytes = %d, want 256MiB", cfg.DefaultMemoryBytes)
	}
	if cfg.MaxMemoryBytes != 1024*1024*1024 {
		t.Errorf("MaxMemoryBytes = %d, want 1GiB", cfg.MaxMemoryBytes)
	}
	if cfg.DefaultNanoCPUs != 500_000_000 {
		t.Errorf("DefaultNanoCPUs = %d, want 500000000", cfg.DefaultNanoCPUs)
	}
	if cfg.MaxNanoCPUs != 2_000_000_000 {
		t.Errorf("MaxNanoCPUs = %d, want 2000000000", cfg.MaxNanoCPUs)
	}
	if cfg.DefaultPidsLimit != 256 {
		t.Errorf("DefaultPidsLimit = %d, want 256", cfg.DefaultPidsLimit)
	}
	if len(cfg.GitHosts) != 2 || cfg.GitHosts[0] != "github.com" || cfg.GitHosts[1] != "gitlab.com" {
		t.Errorf("GitHosts = %v, want [github.com gitlab.com]", cfg.GitHosts)
	}
	if cfg.CloneTimeout != 120*time.Second {
		t.Errorf("CloneTimeout = %v, want 120s", cfg.CloneTimeout)
	}
	if cfg.CloneMaxBytes != 200*1024*1024 {
		t.Errorf("CloneMaxBytes = %d, want 200MiB", cfg.CloneMaxBytes)
	}
	if cfg.BuildTimeout != 10*time.Minute {
		t.Errorf("BuildTimeout = %v, want 10m", cfg.BuildTimeout)
	}
	if cfg.HealthGrace != 3*time.Second {
		t.Errorf("HealthGrace = %v, want 3s", cfg.HealthGrace)
	}
	if cfg.NixpacksBin != "nixpacks" {
		t.Errorf("NixpacksBin = %q, want nixpacks", cfg.NixpacksBin)
	}
	if cfg.MetricsInterval != 30*time.Second {
		t.Errorf("MetricsInterval = %v, want 30s", cfg.MetricsInterval)
	}
	if cfg.MetricsRetention != 24*time.Hour {
		t.Errorf("MetricsRetention = %v, want 24h", cfg.MetricsRetention)
	}
	if cfg.CgroupRoot != "/sys/fs/cgroup" {
		t.Errorf("CgroupRoot = %q, want /sys/fs/cgroup", cfg.CgroupRoot)
	}
}

func TestLoadOverrides(t *testing.T) {
	clearEnv(t)
	t.Setenv("SM_TOKEN", "test-token")
	t.Setenv("SM_PORT", "8095")
	t.Setenv("SM_BIND_ADDR", "0.0.0.0")
	t.Setenv("SM_DEFAULT_RUNTIME", "kata")
	t.Setenv("SM_DEFAULT_MEMORY", "512m")
	t.Setenv("SM_MAX_MEMORY", "2g")
	t.Setenv("SM_DEFAULT_CPU", "1")
	t.Setenv("SM_MAX_CPU", "4")
	t.Setenv("SM_DEFAULT_PIDS", "128")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}
	if cfg.Port != "8095" || cfg.BindAddr != "0.0.0.0" || cfg.DefaultRuntime != "kata" {
		t.Errorf("unexpected config: %+v", cfg)
	}
	if cfg.DefaultMemoryBytes != 512*1024*1024 || cfg.MaxMemoryBytes != 2*1024*1024*1024 {
		t.Errorf("unexpected memory limits: %+v", cfg)
	}
	if cfg.DefaultNanoCPUs != 1_000_000_000 || cfg.MaxNanoCPUs != 4_000_000_000 || cfg.DefaultPidsLimit != 128 {
		t.Errorf("unexpected cpu/pids limits: %+v", cfg)
	}
}

func TestLoadInvalidValues(t *testing.T) {
	cases := map[string]string{
		"SM_DEFAULT_RUNTIME": "gvisor",
		"SM_DEFAULT_MEMORY":  "abc",
		"SM_MAX_MEMORY":      "-1g",
		"SM_DEFAULT_CPU":     "0",
		"SM_MAX_CPU":         "NaN",
		"SM_DEFAULT_PIDS":    "-5",
		"SM_GIT_HOSTS":       " , ,",
		"SM_CLONE_TIMEOUT":   "soon",
		"SM_CLONE_MAX_SIZE":  "big",
		"SM_BUILD_TIMEOUT":   "-5m",
		"SM_HEALTH_GRACE":    "5m", // above the 60s sanity cap
	}
	for key, value := range cases {
		t.Run(key, func(t *testing.T) {
			clearEnv(t)
			t.Setenv("SM_TOKEN", "test-token")
			t.Setenv(key, value)
			if _, err := Load(); err == nil {
				t.Errorf("expected error for %s=%q", key, value)
			}
		})
	}
}

func TestLoadMetricsBounds(t *testing.T) {
	cases := []struct {
		name     string
		env      map[string]string
		wantFail bool
	}{
		{"interval below floor", map[string]string{"SM_METRICS_INTERVAL": "1s"}, true},
		{"interval above ceiling", map[string]string{"SM_METRICS_INTERVAL": "10m"}, true},
		{"interval garbage", map[string]string{"SM_METRICS_INTERVAL": "often"}, true},
		{"retention below floor", map[string]string{"SM_METRICS_RETENTION": "30m"}, true},
		{"retention above ceiling", map[string]string{"SM_METRICS_RETENTION": "200h"}, true},
		{"valid overrides", map[string]string{"SM_METRICS_INTERVAL": "10s", "SM_METRICS_RETENTION": "48h", "SM_CGROUP_ROOT": "/host/cgroup"}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clearEnv(t)
			t.Setenv("SM_TOKEN", "test-token")
			for key, value := range tc.env {
				t.Setenv(key, value)
			}
			cfg, err := Load()
			if tc.wantFail {
				if err == nil {
					t.Errorf("expected error for %v", tc.env)
				}
				return
			}
			if err != nil {
				t.Fatalf("Load() returned error: %v", err)
			}
			if cfg.MetricsInterval != 10*time.Second || cfg.MetricsRetention != 48*time.Hour || cfg.CgroupRoot != "/host/cgroup" {
				t.Errorf("metrics config = %v %v %q", cfg.MetricsInterval, cfg.MetricsRetention, cfg.CgroupRoot)
			}
		})
	}
}

func TestLoadDefaultAboveMax(t *testing.T) {
	clearEnv(t)
	t.Setenv("SM_TOKEN", "test-token")
	t.Setenv("SM_DEFAULT_MEMORY", "2g")

	if _, err := Load(); err == nil {
		t.Error("expected error when SM_DEFAULT_MEMORY exceeds SM_MAX_MEMORY")
	}
}
