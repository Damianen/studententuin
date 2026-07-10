package docker

import (
	"bytes"
	"encoding/binary"
	"strings"
	"testing"
	"time"

	"github.com/moby/moby/api/pkg/stdcopy"
)

// muxFrame writes one stdcopy frame: 8-byte header (stream type, three zero
// bytes, big-endian payload length) followed by the payload — the daemon's
// wire format for non-TTY container logs.
func muxFrame(buf *bytes.Buffer, stream stdcopy.StdType, payload string) {
	header := [8]byte{byte(stream)}
	binary.BigEndian.PutUint32(header[4:], uint32(len(payload)))
	buf.Write(header[:])
	buf.WriteString(payload)
}

// muxStream builds a multiplexed log stream of timestamp-prefixed lines.
func muxStream(t *testing.T, stdout, stderr []string) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	for _, line := range stdout {
		muxFrame(&buf, stdcopy.Stdout, line+"\n")
	}
	for _, line := range stderr {
		muxFrame(&buf, stdcopy.Stderr, line+"\n")
	}
	return &buf
}

func TestDemuxLogs(t *testing.T) {
	t.Run("merges streams sorted by timestamp", func(t *testing.T) {
		src := muxStream(t,
			[]string{
				"2026-06-13T10:00:01.000000001Z first out",
				"2026-06-13T10:00:03.000000001Z second out",
			},
			[]string{"2026-06-13T10:00:02.000000001Z first err"},
		)

		lines, err := demuxLogs(src)
		if err != nil {
			t.Fatalf("demuxLogs: %v", err)
		}
		if len(lines) != 3 {
			t.Fatalf("len = %d, want 3; lines %+v", len(lines), lines)
		}

		wantMessages := []string{"first out", "first err", "second out"}
		wantStreams := []string{"stdout", "stderr", "stdout"}
		for i := range lines {
			if lines[i].Message != wantMessages[i] || lines[i].Stream != wantStreams[i] {
				t.Errorf("lines[%d] = %+v, want message %q stream %q", i, lines[i], wantMessages[i], wantStreams[i])
			}
			if lines[i].Timestamp.IsZero() {
				t.Errorf("lines[%d] has zero timestamp", i)
			}
		}
		if !lines[0].Timestamp.Equal(time.Date(2026, 6, 13, 10, 0, 1, 1, time.UTC)) {
			t.Errorf("timestamp = %v, want parsed RFC3339Nano", lines[0].Timestamp)
		}
	})

	t.Run("keeps malformed lines with zero timestamp", func(t *testing.T) {
		src := muxStream(t, []string{"no timestamp here"}, nil)

		lines, err := demuxLogs(src)
		if err != nil {
			t.Fatalf("demuxLogs: %v", err)
		}
		if len(lines) != 1 {
			t.Fatalf("len = %d, want 1", len(lines))
		}
		if lines[0].Message != "no timestamp here" || !lines[0].Timestamp.IsZero() {
			t.Errorf("line = %+v, want full message and zero timestamp", lines[0])
		}
	})

	t.Run("empty stream", func(t *testing.T) {
		lines, err := demuxLogs(&bytes.Buffer{})
		if err != nil {
			t.Fatalf("demuxLogs: %v", err)
		}
		if len(lines) != 0 {
			t.Errorf("len = %d, want 0", len(lines))
		}
	})
}

func TestParseLogLines(t *testing.T) {
	t.Run("trims carriage returns and skips blank lines", func(t *testing.T) {
		buf := bytes.NewBufferString("2026-06-13T10:00:01Z windows line\r\n\n2026-06-13T10:00:02Z plain\n")
		lines := parseLogLines(buf, "stdout")
		if len(lines) != 2 {
			t.Fatalf("len = %d, want 2; lines %+v", len(lines), lines)
		}
		if lines[0].Message != "windows line" {
			t.Errorf("message = %q, want trailing \\r trimmed", lines[0].Message)
		}
	})

	t.Run("long lines survive", func(t *testing.T) {
		long := strings.Repeat("x", 200*1024)
		buf := bytes.NewBufferString("2026-06-13T10:00:01Z " + long + "\n")
		lines := parseLogLines(buf, "stdout")
		if len(lines) != 1 || lines[0].Message != long {
			t.Fatalf("long line not preserved (len %d)", len(lines))
		}
	})

	t.Run("timestamp-only line keeps stream tag", func(t *testing.T) {
		buf := bytes.NewBufferString("2026-06-13T10:00:01Z \n")
		lines := parseLogLines(buf, "stderr")
		if len(lines) != 1 || lines[0].Stream != "stderr" || lines[0].Message != "" {
			t.Fatalf("lines = %+v", lines)
		}
	})
}
