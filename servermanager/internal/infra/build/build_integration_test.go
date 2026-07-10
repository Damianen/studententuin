//go:build integration

package build

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"servermanager/internal/domain"

	"github.com/google/uuid"
	"github.com/moby/moby/client"
)

// TestIntegrationNixpacksBuild is the canary for nixpacks flag drift and
// legacy-builder compatibility: real binary, real daemon, real image.
func TestIntegrationNixpacksBuild(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	cli, err := client.New(client.FromEnv, client.WithUserAgent("stt-build-test"))
	if err != nil {
		t.Fatalf("docker client: %v", err)
	}
	defer cli.Close()
	if _, err := cli.Ping(ctx, client.PingOptions{}); err != nil {
		t.Fatalf("docker daemon unreachable (integration tests need one): %v", err)
	}

	builder, err := New(ctx, "nixpacks", cli)
	if err != nil {
		t.Fatalf("New: %v (integration tests need the nixpacks binary)", err)
	}

	// Copy the fixture: the generator writes .nixpacks/ and the config file
	// into the build dir, and testdata must stay pristine.
	dir := t.TempDir()
	for _, name := range []string{"package.json", "index.js"} {
		data, err := os.ReadFile(filepath.Join("testdata", "nodeapp", name))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, name), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	appID := uuid.NewString()
	imageRef := domain.AppImageRef(appID, uuid.NewString())
	t.Cleanup(func() {
		_, _ = cli.ImageRemove(context.Background(), imageRef, client.ImageRemoveOptions{PruneChildren: true})
	})

	var log strings.Builder
	err = builder.Build(ctx, domain.BuildInput{
		Dir:          dir,
		ImageRef:     imageRef,
		Env:          map[string]string{"FIXTURE_MARKER": "stt"},
		StartCommand: "node index.js",
		Port:         3000,
	}, &log)
	if err != nil {
		t.Fatalf("Build: %v\n--- build log ---\n%s", err, log.String())
	}

	if log.Len() == 0 {
		t.Error("build log sink is empty")
	}
	if _, err := os.Stat(filepath.Join(dir, ".nixpacks", "Dockerfile")); err != nil {
		t.Errorf(".nixpacks/Dockerfile missing after generate: %v", err)
	}

	res, err := cli.ImageList(ctx, client.ImageListOptions{
		Filters: client.Filters{}.Add("reference", imageRef),
	})
	if err != nil {
		t.Fatalf("ImageList: %v", err)
	}
	if len(res.Items) != 1 {
		t.Fatalf("image %s not found after build (%d matches)", imageRef, len(res.Items))
	}

	// The request env must be baked into the image (nixpacks ARG+ENV plus
	// our BuildArgs); a lost build arg yields an empty value here.
	inspect, err := cli.ImageInspect(ctx, imageRef)
	if err != nil {
		t.Fatalf("ImageInspect: %v", err)
	}
	var marker, port bool
	for _, env := range inspect.Config.Env {
		if env == "FIXTURE_MARKER=stt" {
			marker = true
		}
		if env == "PORT=3000" {
			port = true
		}
	}
	if !marker || !port {
		t.Errorf("image env missing FIXTURE_MARKER/PORT: %v", inspect.Config.Env)
	}
}

func TestIntegrationNixpacksBuildFailure(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	cli, err := client.New(client.FromEnv, client.WithUserAgent("stt-build-test"))
	if err != nil {
		t.Fatalf("docker client: %v", err)
	}
	defer cli.Close()
	if _, err := cli.Ping(ctx, client.PingOptions{}); err != nil {
		t.Fatalf("docker daemon unreachable: %v", err)
	}

	builder, err := New(ctx, "nixpacks", cli)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	dir := t.TempDir()
	data, err := os.ReadFile(filepath.Join("testdata", "nodeapp", "package.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "package.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "index.js"), []byte("ok\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var log strings.Builder
	err = builder.Build(ctx, domain.BuildInput{
		Dir:          dir,
		ImageRef:     domain.AppImageRef(uuid.NewString(), uuid.NewString()),
		BuildCommand: "exit 7", // the build step itself fails inside the container
	}, &log)
	if err == nil {
		t.Fatalf("Build succeeded, want failure\n--- log ---\n%s", log.String())
	}
	if !strings.Contains(err.Error(), "build failed") {
		t.Errorf("err = %v, want a daemon build failure", err)
	}
}
