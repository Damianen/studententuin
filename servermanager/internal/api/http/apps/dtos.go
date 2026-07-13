package apps

import (
	"strconv"
	"time"

	"servermanager/internal/domain"
)

// RunRequest mirrors the §3.2 deploy body minus the git/build fields.
// start_command is an argv array, not a string — the manager never
// shell-splits user input.
type RunRequest struct {
	Image          string            `json:"image" binding:"required"`
	Env            map[string]string `json:"env"`
	Port           int               `json:"port"`
	StartCommand   []string          `json:"start_command"`
	MemoryLimit    string            `json:"memory_limit"`
	CpuLimit       string            `json:"cpu_limit"`
	Runtime        string            `json:"runtime"`
	Volumes        []string          `json:"volumes"`
	ReadonlyRootfs bool              `json:"readonly_rootfs"`
	DatabaseID     string            `json:"database_id"`
}

type RunResponse struct {
	ContainerID string `json:"container_id"`
	Name        string `json:"name"`
}

type ContainerStateResponse struct {
	Exists       bool       `json:"exists"`
	Running      bool       `json:"running"`
	StartedAt    *time.Time `json:"started_at,omitempty"`
	RestartCount int        `json:"restart_count"`
	Status       string     `json:"status"`
}

// LogEntryResponse is the frontend's LogEntry shape (§3.2): the MVP level
// mapping is stdout→info, stderr→error — reliable without parsing because
// Docker multiplexes the two streams.
type LogEntryResponse struct {
	ID        string `json:"id"`        // "<unixnano>-<index>", unique within a response
	Timestamp string `json:"timestamp"` // RFC3339Nano, UTC; empty when the line had no parsable prefix
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// LogsResponse wraps the entries in an object so fields (e.g. truncated)
// can be added later without breaking the api's client.
type LogsResponse struct {
	Logs []LogEntryResponse `json:"logs"`
}

func toLogsResponse(lines []domain.LogLine) LogsResponse {
	resp := LogsResponse{Logs: make([]LogEntryResponse, 0, len(lines))}
	for i, line := range lines {
		resp.Logs = append(resp.Logs, toLogEntryResponse(line, i))
	}
	return resp
}

func toLogEntryResponse(line domain.LogLine, index int) LogEntryResponse {
	entry := LogEntryResponse{
		ID:      strconv.FormatInt(line.Timestamp.UnixNano(), 10) + "-" + strconv.Itoa(index),
		Level:   "info",
		Message: line.Message,
	}
	if line.Stream == "stderr" {
		entry.Level = "error"
	}
	if !line.Timestamp.IsZero() {
		entry.Timestamp = line.Timestamp.UTC().Format(time.RFC3339Nano)
	}
	return entry
}
