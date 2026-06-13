package domain

import "time"

// ContainerSpec is everything the runtime adapter needs to create a hardened
// user container (SERVERMANAGEMENT.md §3.3).
type ContainerSpec struct {
	AppID          string // application UUID from the api; names are derived from it
	Image          string
	Env            map[string]string
	Port           int      // container port the app listens on
	Cmd            []string // optional start command override
	MemoryBytes    int64
	NanoCPUs       int64
	PidsLimit      int64
	Runtime        string   // must pass ValidRuntime
	Volumes        []string // named volumes only — never host bind mounts
	ReadonlyRootfs bool
}

// ContainerState is the actual Docker-side state of a container.
type ContainerState struct {
	Exists       bool
	Running      bool
	StartedAt    time.Time
	RestartCount int
	Status       string
}

type LogOptions struct {
	Tail  int
	Since time.Time
}

type LogLine struct {
	Timestamp time.Time
	Stream    string // "stdout" or "stderr"
	Message   string
}
