package docker

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"servermanager/internal/ports"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"
)

var _ ports.ContainerExec = (*Runtime)(nil)

// maxExecOutputBytes bounds each captured exec stream — metrics commands
// return one short line, anything bigger is a bug, not data to keep.
const maxExecOutputBytes = 64 * 1024

// maxExecStderrTail bounds the stderr excerpt carried on a non-zero-exit
// error, so a chatty failure can't bloat logs or job records.
const maxExecStderrTail = 512

// ExecCapture runs argv in the container (no TTY, so the output stream stays
// multiplexed and stdcopy can split it, like logs.go) and returns stdout.
// A non-zero exit becomes an error with a short stderr tail — callers must
// never feed error text into a metrics series.
func (r *Runtime) ExecCapture(ctx context.Context, nameOrID string, argv []string) (string, error) {
	created, err := r.cli.ExecCreate(ctx, nameOrID, client.ExecCreateOptions{
		Cmd:          argv,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return "", mapErr("exec create "+nameOrID, err)
	}

	attach, err := r.cli.ExecAttach(ctx, created.ID, client.ExecAttachOptions{})
	if err != nil {
		return "", mapErr("exec attach "+nameOrID, err)
	}
	defer attach.Close()

	// The hijacked connection honors ctx only while dialing (moby client
	// v0.4.1), so without a watchdog StdCopy would block until the exec'd
	// process exits — one wedged psql would stall the collector forever.
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-ctx.Done():
			attach.Close()
		case <-done:
		}
	}()

	stdout := &cappedBuffer{max: maxExecOutputBytes}
	stderr := &cappedBuffer{max: maxExecOutputBytes}
	if _, err := stdcopy.StdCopy(stdout, stderr, attach.Reader); err != nil {
		if ctx.Err() != nil {
			return "", fmt.Errorf("exec read %s: %w", nameOrID, ctx.Err())
		}
		return "", mapErr("exec read "+nameOrID, err)
	}

	info, err := r.cli.ExecInspect(ctx, created.ID, client.ExecInspectOptions{})
	if err != nil {
		return "", mapErr("exec inspect "+nameOrID, err)
	}
	if info.ExitCode != 0 {
		return "", fmt.Errorf("exec %s: exit %d: %s", nameOrID, info.ExitCode, stderrTail(stderr.String()))
	}
	return stdout.String(), nil
}

// ContainerEnv returns only the requested keys from the container's
// configured environment. The values are potential secrets — they go to the
// caller and nowhere else, never into logs or errors.
func (r *Runtime) ContainerEnv(ctx context.Context, nameOrID string, keys ...string) (map[string]string, error) {
	res, err := r.cli.ContainerInspect(ctx, nameOrID, client.ContainerInspectOptions{})
	if err != nil {
		return nil, mapErr("inspect env "+nameOrID, err)
	}

	want := make(map[string]bool, len(keys))
	for _, key := range keys {
		want[key] = true
	}
	out := make(map[string]string, len(keys))
	if cfg := res.Container.Config; cfg != nil {
		for _, kv := range cfg.Env {
			if k, v, ok := strings.Cut(kv, "="); ok && want[k] {
				out[k] = v
			}
		}
	}
	return out, nil
}

// cappedBuffer keeps the first max bytes and silently drops the rest, so a
// runaway command can't grow memory while StdCopy drains the stream.
type cappedBuffer struct {
	buf bytes.Buffer
	max int
}

func (b *cappedBuffer) Write(p []byte) (int, error) {
	if room := b.max - b.buf.Len(); room > 0 {
		if len(p) > room {
			p = p[:room]
		}
		b.buf.Write(p)
	}
	return len(p), nil
}

func (b *cappedBuffer) String() string { return b.buf.String() }

// stderrTail returns the last maxExecStderrTail bytes, single-line.
func stderrTail(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > maxExecStderrTail {
		s = s[len(s)-maxExecStderrTail:]
	}
	return strings.ReplaceAll(s, "\n", " ")
}
