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
