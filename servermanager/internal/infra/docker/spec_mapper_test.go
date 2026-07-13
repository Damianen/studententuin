package docker

import (
	"errors"
	"reflect"
	"testing"

	"servermanager/internal/domain"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/api/types/network"
)

const testAppID = "8b9f5e9e-9a3a-4b7e-9a59-1d2f3a4b5c6d"

func fullSpec() domain.ContainerSpec {
	return domain.ContainerSpec{
		AppID:          testAppID,
		Image:          "busybox:1.37",
		Env:            map[string]string{"B": "2", "A": "1"},
		Port:           3000,
		Cmd:            []string{"sleep", "3600"},
		MemoryBytes:    256 * 1024 * 1024,
		NanoCPUs:       500_000_000,
		PidsLimit:      256,
		Runtime:        domain.RuntimeRunc,
		Volumes:        []string{"data:/var/lib/data"},
		ReadonlyRootfs: true,
	}
}

// TestBuildCreateOptionsHardening is the §7 security-checklist test: every
// container gets limits + pids cap + no-new-privileges + cap-drop ALL +
// isolated network + no binds + log rotation, asserted on the generated
// HostConfig without a daemon.
func TestBuildCreateOptionsHardening(t *testing.T) {
	spec := fullSpec()
	opts, err := buildCreateOptions(spec)
	if err != nil {
		t.Fatalf("buildCreateOptions: %v", err)
	}

	if opts.Name != "stt-app-"+testAppID {
		t.Errorf("Name = %q, want stt-app-%s", opts.Name, testAppID)
	}

	hc := opts.HostConfig
	if hc.Memory != spec.MemoryBytes {
		t.Errorf("Memory = %d, want %d", hc.Memory, spec.MemoryBytes)
	}
	if hc.MemorySwap != spec.MemoryBytes {
		t.Errorf("MemorySwap = %d, want %d (no swap)", hc.MemorySwap, spec.MemoryBytes)
	}
	if hc.NanoCPUs != spec.NanoCPUs {
		t.Errorf("NanoCPUs = %d, want %d", hc.NanoCPUs, spec.NanoCPUs)
	}
	if hc.PidsLimit == nil || *hc.PidsLimit != spec.PidsLimit {
		t.Errorf("PidsLimit = %v, want %d", hc.PidsLimit, spec.PidsLimit)
	}
	if want := []string{"no-new-privileges:true"}; !reflect.DeepEqual(hc.SecurityOpt, want) {
		t.Errorf("SecurityOpt = %v, want %v", hc.SecurityOpt, want)
	}
	if want := []string{"ALL"}; !reflect.DeepEqual(hc.CapDrop, want) {
		t.Errorf("CapDrop = %v, want %v", hc.CapDrop, want)
	}
	if len(hc.CapAdd) != 0 {
		t.Errorf("CapAdd = %v, want empty", hc.CapAdd)
	}
	if hc.Privileged {
		t.Error("Privileged = true, want false")
	}
	if want := container.NetworkMode("stt-net-" + testAppID); hc.NetworkMode != want {
		t.Errorf("NetworkMode = %q, want %q", hc.NetworkMode, want)
	}
	if len(hc.Binds) != 0 {
		t.Errorf("Binds = %v, want empty — bind mounts are forbidden", hc.Binds)
	}
	if len(hc.PortBindings) != 0 {
		t.Errorf("PortBindings = %v, want empty — no published ports", hc.PortBindings)
	}
	if hc.PidMode != "" || hc.IpcMode != "" {
		t.Errorf("PidMode/IpcMode = %q/%q, want defaults", hc.PidMode, hc.IpcMode)
	}
	wantMounts := []mount.Mount{{Type: mount.TypeVolume, Source: "data", Target: "/var/lib/data"}}
	if !reflect.DeepEqual(hc.Mounts, wantMounts) {
		t.Errorf("Mounts = %v, want %v", hc.Mounts, wantMounts)
	}
	if !hc.ReadonlyRootfs {
		t.Error("ReadonlyRootfs = false, want true")
	}
	if hc.Tmpfs["/tmp"] == "" {
		t.Errorf("Tmpfs = %v, want /tmp entry", hc.Tmpfs)
	}
	if hc.RestartPolicy.Name != container.RestartPolicyOnFailure || hc.RestartPolicy.MaximumRetryCount != maxRestartRetries {
		t.Errorf("RestartPolicy = %+v, want on-failure/%d", hc.RestartPolicy, maxRestartRetries)
	}
	if hc.Runtime != domain.RuntimeRunc {
		t.Errorf("Runtime = %q, want runc", hc.Runtime)
	}
	wantLog := container.LogConfig{Type: "json-file", Config: map[string]string{"max-size": "10m", "max-file": "3"}}
	if !reflect.DeepEqual(hc.LogConfig, wantLog) {
		t.Errorf("LogConfig = %+v, want %+v", hc.LogConfig, wantLog)
	}

	cfg := opts.Config
	if cfg.Image != spec.Image {
		t.Errorf("Image = %q, want %q", cfg.Image, spec.Image)
	}
	if want := []string{"A=1", "B=2"}; !reflect.DeepEqual(cfg.Env, want) {
		t.Errorf("Env = %v, want sorted %v", cfg.Env, want)
	}
	if !reflect.DeepEqual(cfg.Cmd, spec.Cmd) {
		t.Errorf("Cmd = %v, want %v", cfg.Cmd, spec.Cmd)
	}
	if len(cfg.ExposedPorts) != 1 {
		t.Errorf("ExposedPorts = %v, want one 3000/tcp entry", cfg.ExposedPorts)
	} else if p := network.MustParsePort("3000/tcp"); cfg.ExposedPorts[p] != struct{}{} {
		t.Errorf("ExposedPorts = %v, want 3000/tcp", cfg.ExposedPorts)
	}

	netName := "stt-net-" + testAppID
	if opts.NetworkingConfig == nil || opts.NetworkingConfig.EndpointsConfig[netName] == nil {
		t.Errorf("NetworkingConfig missing endpoint for %s", netName)
	}

	// App specs must not accidentally grow database-only settings.
	if cfg.User != "" {
		t.Errorf("User = %q, want empty for app specs", cfg.User)
	}
	if cfg.Healthcheck != nil {
		t.Errorf("Healthcheck = %+v, want nil for app specs", cfg.Healthcheck)
	}
	if hc.ShmSize != 0 {
		t.Errorf("ShmSize = %d, want 0 (daemon default) for app specs", hc.ShmSize)
	}
}

// TestBuildCreateOptionsDBSpec is the database twin of the hardening test: a
// KindDB spec derives db names, carries the postgres-specific settings, and
// keeps every §3.3 default that matters.
func TestBuildCreateOptionsDBSpec(t *testing.T) {
	spec := domain.ContainerSpec{
		AppID: testAppID,
		Kind:  domain.KindDB,
		Image: "postgres:16",
		Env:   map[string]string{"POSTGRES_USER": "app"},
		Port:  5432,
		User:  "postgres",
		Healthcheck: &domain.Healthcheck{
			Test:     []string{"CMD", "pg_isready", "-h", "127.0.0.1", "-U", "app"},
			Interval: 2_000_000_000,
			Timeout:  3_000_000_000,
			Retries:  30,
		},
		MemoryBytes:  512 * 1024 * 1024,
		NanoCPUs:     500_000_000,
		PidsLimit:    256,
		ShmSizeBytes: 256 * 1024 * 1024,
		Runtime:      domain.RuntimeRunc,
		Volumes:      []string{"stt-db-data-" + testAppID + ":/var/lib/postgresql/data"},
	}
	opts, err := buildCreateOptions(spec)
	if err != nil {
		t.Fatalf("buildCreateOptions: %v", err)
	}

	if opts.Name != "stt-db-"+testAppID {
		t.Errorf("Name = %q, want stt-db-%s", opts.Name, testAppID)
	}
	hc := opts.HostConfig
	if want := container.NetworkMode("stt-dbnet-" + testAppID); hc.NetworkMode != want {
		t.Errorf("NetworkMode = %q, want %q", hc.NetworkMode, want)
	}
	if opts.NetworkingConfig == nil || opts.NetworkingConfig.EndpointsConfig["stt-dbnet-"+testAppID] == nil {
		t.Errorf("NetworkingConfig missing endpoint for stt-dbnet-%s", testAppID)
	}
	if opts.Config.User != "postgres" {
		t.Errorf("User = %q, want postgres", opts.Config.User)
	}
	gotHC := opts.Config.Healthcheck
	if gotHC == nil || !reflect.DeepEqual(gotHC.Test, spec.Healthcheck.Test) ||
		gotHC.Interval != spec.Healthcheck.Interval || gotHC.Timeout != spec.Healthcheck.Timeout ||
		gotHC.Retries != spec.Healthcheck.Retries {
		t.Errorf("Healthcheck = %+v, want %+v", gotHC, spec.Healthcheck)
	}
	if hc.ShmSize != spec.ShmSizeBytes {
		t.Errorf("ShmSize = %d, want %d", hc.ShmSize, spec.ShmSizeBytes)
	}
	if hc.ReadonlyRootfs {
		t.Error("ReadonlyRootfs = true, want false for postgres")
	}

	// The §3.3 hardening must hold for databases exactly as for apps.
	if want := []string{"ALL"}; !reflect.DeepEqual(hc.CapDrop, want) {
		t.Errorf("CapDrop = %v, want %v", hc.CapDrop, want)
	}
	if want := []string{"no-new-privileges:true"}; !reflect.DeepEqual(hc.SecurityOpt, want) {
		t.Errorf("SecurityOpt = %v, want %v", hc.SecurityOpt, want)
	}
	if len(hc.Binds) != 0 || len(hc.PortBindings) != 0 {
		t.Errorf("Binds/PortBindings = %v/%v, want empty", hc.Binds, hc.PortBindings)
	}
	if hc.PidsLimit == nil || *hc.PidsLimit != spec.PidsLimit {
		t.Errorf("PidsLimit = %v, want %d", hc.PidsLimit, spec.PidsLimit)
	}
	if hc.Memory != spec.MemoryBytes || hc.MemorySwap != spec.MemoryBytes {
		t.Errorf("Memory/MemorySwap = %d/%d, want %d", hc.Memory, hc.MemorySwap, spec.MemoryBytes)
	}
	wantLog := container.LogConfig{Type: "json-file", Config: map[string]string{"max-size": "10m", "max-file": "3"}}
	if !reflect.DeepEqual(hc.LogConfig, wantLog) {
		t.Errorf("LogConfig = %+v, want %+v", hc.LogConfig, wantLog)
	}
}

func TestBuildCreateOptionsNoPort(t *testing.T) {
	spec := fullSpec()
	spec.Port = 0
	opts, err := buildCreateOptions(spec)
	if err != nil {
		t.Fatalf("buildCreateOptions: %v", err)
	}
	if len(opts.Config.ExposedPorts) != 0 {
		t.Errorf("ExposedPorts = %v, want empty", opts.Config.ExposedPorts)
	}
}

func TestBuildCreateOptionsReadonlyOff(t *testing.T) {
	spec := fullSpec()
	spec.ReadonlyRootfs = false
	opts, err := buildCreateOptions(spec)
	if err != nil {
		t.Fatalf("buildCreateOptions: %v", err)
	}
	if opts.HostConfig.ReadonlyRootfs {
		t.Error("ReadonlyRootfs = true, want false")
	}
	if opts.HostConfig.Tmpfs["/tmp"] == "" {
		t.Error("Tmpfs /tmp should be set regardless of readonly rootfs")
	}
}

func TestBuildCreateOptionsRejects(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*domain.ContainerSpec)
	}{
		{"runtime not on allowlist", func(s *domain.ContainerSpec) { s.Runtime = "gvisor" }},
		{"empty runtime", func(s *domain.ContainerSpec) { s.Runtime = "" }},
		{"bind mount volume", func(s *domain.ContainerSpec) { s.Volumes = []string{"/host:/data"} }},
		{"malformed volume", func(s *domain.ContainerSpec) { s.Volumes = []string{"data"} }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			spec := fullSpec()
			tc.mutate(&spec)
			if _, err := buildCreateOptions(spec); !errors.Is(err, domain.ErrInvalid) {
				t.Errorf("buildCreateOptions error = %v, want ErrInvalid", err)
			}
		})
	}
}
