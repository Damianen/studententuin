package apps

import "time"

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
