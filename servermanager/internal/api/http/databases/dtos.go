package databases

import (
	"strconv"
	"time"

	"servermanager/internal/domain"
)

// ProvisionRequest is the api → manager provision body. The image is chosen
// server-side from type+version (allowlist §7); credentials are api-generated
// and pass through to the container env only.
type ProvisionRequest struct {
	Type        string `json:"type" binding:"required"`
	Version     string `json:"version" binding:"required"`
	DBName      string `json:"db_name" binding:"required"`
	DBUser      string `json:"db_user" binding:"required"`
	DBPassword  string `json:"db_password" binding:"required"`
	MemoryLimit string `json:"memory_limit"`
	CpuLimit    string `json:"cpu_limit"`
	Runtime     string `json:"runtime"`
}

type ProvisionResponse struct {
	ProvisionID string `json:"provision_id"`
}

// ProvisionStateResponse mirrors the deployments poll shape for provisions.
type ProvisionStateResponse struct {
	ID            string    `json:"id"`
	Status        string    `json:"status"`
	Error         string    `json:"error,omitempty"`
	Image         string    `json:"image,omitempty"`
	ContainerID   string    `json:"container_id,omitempty"`
	ContainerName string    `json:"container_name,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// DatabaseStatusResponse is Docker state plus the latest provision job. The
// provision block is absent when the manager lost the job (restart) — the
// api's poller treats "no container, no provision" as a lost provision.
type DatabaseStatusResponse struct {
	Exists       bool                    `json:"exists"`
	Running      bool                    `json:"running"`
	StartedAt    *time.Time              `json:"started_at,omitempty"`
	RestartCount int                     `json:"restart_count"`
	Status       string                  `json:"status"`
	Health       string                  `json:"health,omitempty"`
	Provision    *ProvisionStateResponse `json:"provision,omitempty"`
}

// LogEntryResponse / LogsResponse duplicate the apps package's shapes on
// purpose: the two endpoints must be able to evolve independently, and the
// mapping is trivial.
type LogEntryResponse struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

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
