package docker

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"servermanager/internal/domain"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"
)

// maxLogLine bounds a single scanned log line. The json-file driver splits
// container lines above 16KB, but stay defensive against other drivers.
const maxLogLine = 1024 * 1024

// Logs reads the container's log tail and returns the stdout/stderr lines
// merged in timestamp order. The manager never creates TTY containers
// (see spec_mapper), so the stream is always multiplexed and stdcopy can
// split it — which is what makes the stdout→info / stderr→error mapping
// reliable without parsing line content.
func (r *Runtime) Logs(ctx context.Context, nameOrID string, opts domain.LogOptions) ([]domain.LogLine, error) {
	clOpts := client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
	}
	if opts.Tail > 0 {
		clOpts.Tail = strconv.Itoa(opts.Tail)
	}
	if !opts.Since.IsZero() {
		clOpts.Since = opts.Since.Format(time.RFC3339Nano)
	}

	rc, err := r.cli.ContainerLogs(ctx, nameOrID, clOpts)
	if err != nil {
		return nil, mapErr("logs "+nameOrID, err)
	}
	defer rc.Close()

	lines, err := demuxLogs(rc)
	if err != nil {
		return nil, mapErr("logs demux "+nameOrID, err)
	}
	return lines, nil
}

// FollowLogs streams log lines as the container produces them. The returned
// channel closes when the container's log stream ends or ctx is canceled
// (the SDK closes the underlying reader on cancel, which unwinds the
// demuxing goroutines).
func (r *Runtime) FollowLogs(ctx context.Context, nameOrID string, opts domain.LogOptions) (<-chan domain.LogLine, error) {
	clOpts := client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Follow:     true,
		Tail:       "0", // history is served by Logs; follow only pushes new lines
	}
	if opts.Tail > 0 {
		clOpts.Tail = strconv.Itoa(opts.Tail)
	}
	if !opts.Since.IsZero() {
		clOpts.Since = opts.Since.Format(time.RFC3339Nano)
	}

	rc, err := r.cli.ContainerLogs(ctx, nameOrID, clOpts)
	if err != nil {
		return nil, mapErr("follow logs "+nameOrID, err)
	}

	lines := make(chan domain.LogLine)
	outR, outW := io.Pipe()
	errR, errW := io.Pipe()

	go func() {
		// StdCopy returns when rc ends (container gone, daemon closed it, or
		// ctx canceled). Propagating the close unwinds both scanners.
		_, err := stdcopy.StdCopy(outW, errW, rc)
		_ = outW.CloseWithError(err)
		_ = errW.CloseWithError(err)
		_ = rc.Close()
	}()

	var wg sync.WaitGroup
	for _, s := range []struct {
		r      io.Reader
		stream string
	}{{outR, "stdout"}, {errR, "stderr"}} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			scanLogLines(ctx, s.r, s.stream, lines)
		}()
	}
	go func() {
		wg.Wait()
		close(lines)
	}()

	return lines, nil
}

// scanLogLines feeds parsed lines from one demuxed stream into out until the
// reader ends or ctx is canceled.
func scanLogLines(ctx context.Context, r io.Reader, stream string, out chan<- domain.LogLine) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), maxLogLine)
	for scanner.Scan() {
		raw := strings.TrimSuffix(scanner.Text(), "\r")
		if raw == "" {
			continue
		}
		select {
		case out <- parseLogLine(raw, stream):
		case <-ctx.Done():
			return
		}
	}
}

// demuxLogs splits a multiplexed log stream into per-stream lines and merges
// them in timestamp order.
func demuxLogs(src io.Reader) ([]domain.LogLine, error) {
	var stdout, stderr bytes.Buffer
	if _, err := stdcopy.StdCopy(&stdout, &stderr, src); err != nil {
		return nil, err
	}

	lines := parseLogLines(&stdout, "stdout")
	lines = append(lines, parseLogLines(&stderr, "stderr")...)

	// Docker's tail counts lines of the combined log, so tail=N holds for the
	// merged view. Stable sort keeps within-stream order; cross-stream lines
	// with identical timestamps can swap — unrecoverable after demux, cosmetic.
	sort.SliceStable(lines, func(i, j int) bool {
		return lines[i].Timestamp.Before(lines[j].Timestamp)
	})
	return lines, nil
}

// parseLogLines splits a demuxed stream into LogLines. With Timestamps on,
// the daemon prefixes every line with RFC3339Nano + space. A line that
// doesn't parse is kept whole with a zero timestamp rather than dropped —
// missing log lines are worse than missing timestamps.
func parseLogLines(buf *bytes.Buffer, stream string) []domain.LogLine {
	var lines []domain.LogLine
	scanner := bufio.NewScanner(buf)
	scanner.Buffer(make([]byte, 0, 64*1024), maxLogLine)
	for scanner.Scan() {
		raw := strings.TrimSuffix(scanner.Text(), "\r")
		if raw == "" {
			continue
		}
		lines = append(lines, parseLogLine(raw, stream))
	}
	return lines
}

func parseLogLine(raw, stream string) domain.LogLine {
	line := domain.LogLine{Stream: stream, Message: raw}
	if prefix, rest, ok := strings.Cut(raw, " "); ok {
		if ts, err := time.Parse(time.RFC3339Nano, prefix); err == nil {
			line.Timestamp = ts
			line.Message = rest
		}
	}
	return line
}
