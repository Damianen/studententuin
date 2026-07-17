//go:build integration

package docker

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"servermanager/internal/domain"

	"github.com/google/uuid"
)

func TestIntegrationExecCapture(t *testing.T) {
	r := integrationRuntime(t)
	ctx := context.Background()
	appID := uuid.NewString()
	spec := integrationSpec(appID)
	cleanupApp(t, r, appID)

	cid, err := r.Create(ctx, spec)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if err := r.Start(ctx, cid); err != nil {
		t.Fatalf("Start: %v", err)
	}

	t.Run("captures stdout", func(t *testing.T) {
		out, err := r.ExecCapture(ctx, cid, []string{"echo", "hello"})
		if err != nil {
			t.Fatalf("ExecCapture: %v", err)
		}
		if out != "hello\n" {
			t.Errorf("stdout = %q, want %q", out, "hello\n")
		}
	})

	t.Run("non-zero exit carries the stderr tail", func(t *testing.T) {
		_, err := r.ExecCapture(ctx, cid, []string{"sh", "-c", "echo oops >&2; exit 3"})
		if err == nil {
			t.Fatal("ExecCapture = nil error, want exit failure")
		}
		if !strings.Contains(err.Error(), "exit 3") || !strings.Contains(err.Error(), "oops") {
			t.Errorf("err = %v, want exit code and stderr tail", err)
		}
	})

	t.Run("ctx deadline unblocks a non-exiting exec", func(t *testing.T) {
		// The hijacked connection ignores ctx after dialing, so this pins the
		// watchdog: without it a wedged psql would stall the collector forever.
		shortCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()

		start := time.Now()
		_, err := r.ExecCapture(shortCtx, cid, []string{"sleep", "60"})
		if err == nil {
			t.Fatal("ExecCapture = nil error, want deadline failure")
		}
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("err = %v, want context.DeadlineExceeded", err)
		}
		if elapsed := time.Since(start); elapsed > 5*time.Second {
			t.Errorf("ExecCapture blocked %v past a 1s deadline", elapsed)
		}
	})

	t.Run("container env returns only requested keys", func(t *testing.T) {
		env, err := r.ContainerEnv(ctx, cid, "FIXTURE", "NO_SUCH_KEY")
		if err != nil {
			t.Fatalf("ContainerEnv: %v", err)
		}
		if env["FIXTURE"] != "1" {
			t.Errorf("FIXTURE = %q, want 1", env["FIXTURE"])
		}
		if _, ok := env["NO_SUCH_KEY"]; ok {
			t.Error("NO_SUCH_KEY present, want absent")
		}
		if _, ok := env["PATH"]; ok {
			t.Error("PATH present, want only requested keys")
		}
	})
}

func TestIntegrationExecCaptureStoppedContainer(t *testing.T) {
	r := integrationRuntime(t)
	ctx := context.Background()
	appID := uuid.NewString()
	cleanupApp(t, r, appID)

	cid, err := r.Create(ctx, integrationSpec(appID))
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	// Created but never started: exec must fail cleanly, not hang.
	if _, err := r.ExecCapture(ctx, cid, []string{"echo", "hi"}); err == nil {
		t.Error("ExecCapture on a stopped container = nil, want error")
	}
}

func TestIntegrationExecCaptureUnknownContainer(t *testing.T) {
	r := integrationRuntime(t)
	_, err := r.ExecCapture(context.Background(), domain.AppContainerName(uuid.NewString()), []string{"true"})
	if err == nil {
		t.Error("ExecCapture on unknown container = nil, want error")
	}
}
