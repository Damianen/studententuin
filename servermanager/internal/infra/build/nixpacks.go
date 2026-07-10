// Package build implements ports.ImageBuilder with nixpacks: the binary only
// *generates* .nixpacks/Dockerfile from the checkout (§3.4 auto-detection),
// and the actual build runs through the Docker daemon, so user build steps
// execute in daemon-managed containers — never on the manager host.
package build

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"servermanager/internal/domain"
	"servermanager/internal/ports"

	"github.com/moby/moby/api/types/build"
	"github.com/moby/moby/client"
)

// dockerBuildAPI is the one client method the builder needs, kept narrow so
// unit tests can fake it without a daemon.
type dockerBuildAPI interface {
	ImageBuild(ctx context.Context, buildContext io.Reader, options client.ImageBuildOptions) (client.ImageBuildResult, error)
}

type Nixpacks struct {
	bin string
	api dockerBuildAPI
}

var _ ports.ImageBuilder = (*Nixpacks)(nil)

// New resolves and probes the nixpacks binary up front, so a missing builder
// fails at startup like any other bad configuration.
func New(ctx context.Context, bin string, api dockerBuildAPI) (*Nixpacks, error) {
	path, err := exec.LookPath(bin)
	if err != nil {
		return nil, fmt.Errorf("nixpacks binary %q not found: %w", bin, err)
	}

	probeCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	// #nosec G204 -- fixed argv: the binary path comes from server config
	// (SM_NIXPACKS_BIN), never from a request.
	out, err := exec.CommandContext(probeCtx, path, "--version").Output()
	if err != nil {
		return nil, fmt.Errorf("probing %s --version: %w", path, err)
	}
	slog.Info("nixpacks builder ready", "binary", path, "version", strings.TrimSpace(string(out)))

	return &Nixpacks{bin: path, api: api}, nil
}

// Build generates the build plan with nixpacks, streams the checkout as a
// tar context to the daemon, and builds in.ImageRef with the legacy builder.
// All generator and daemon output is written to logSink as it arrives.
func (n *Nixpacks) Build(ctx context.Context, in domain.BuildInput, logSink io.Writer) error {
	configName, buildArgs, err := writeBuildConfig(in)
	if err != nil {
		return err
	}

	if err := n.generate(ctx, in.Dir, configName, logSink); err != nil {
		return err
	}

	dockerfile := filepath.Join(".nixpacks", "Dockerfile")
	if _, err := os.Stat(filepath.Join(in.Dir, dockerfile)); err != nil {
		// The canary for nixpacks flag/behavior drift: --out must leave the
		// Dockerfile in place for the daemon build below.
		return fmt.Errorf("nixpacks did not produce %s (version drift?): %w", dockerfile, err)
	}

	pr, pw := io.Pipe()
	go func() { pw.CloseWithError(writeTar(pw, in.Dir)) }()

	res, err := n.api.ImageBuild(ctx, pr, client.ImageBuildOptions{
		Tags:        []string{in.ImageRef},
		Dockerfile:  dockerfile,
		BuildArgs:   buildArgs,
		Remove:      true,
		ForceRemove: true,
		// The generated Dockerfile is plain (no BuildKit-only directives with
		// --no-cache), and BuilderV1 needs no session plumbing. Pin it so a
		// buildkit-default daemon can't reroute the request.
		Version: build.BuilderV1,
		Labels:  map[string]string{"studententuin.managed": "true"},
	})
	if err != nil {
		// Unblock the tar goroutine if the daemon rejected the request.
		_ = pr.CloseWithError(err)
		return fmt.Errorf("starting image build: %w", err)
	}

	return decodeBuildStream(res.Body, logSink)
}

// generate runs the nixpacks binary over the checkout. User-influenced values
// never enter argv — they travel via the config file writeBuildConfig wrote.
func (n *Nixpacks) generate(ctx context.Context, dir, configName string, logSink io.Writer) error {
	args := []string{"build", dir, "--out", dir, "--no-cache"}
	if configName != "" {
		args = append(args, "--config", configName)
	}

	// #nosec G204 -- fixed argv: binary from server config, fixed flags, and
	// manager-generated paths; user data travels via the config file.
	cmd := exec.CommandContext(ctx, n.bin, args...)
	cmd.Dir = dir
	// A minimal, fixed environment: never os.Environ(), so manager secrets
	// (SM_TOKEN) cannot leak into the generator or its output.
	cmd.Env = []string{"PATH=" + os.Getenv("PATH"), "HOME=" + dir}
	cmd.Stdout = logSink
	cmd.Stderr = logSink

	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("nixpacks plan generation timed out")
		}
		return fmt.Errorf("nixpacks plan generation failed: %w", err)
	}
	return nil
}
