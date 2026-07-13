package domain

import "time"

// ContainerKind selects the container/network naming scheme. The zero value
// is an application container, so existing specs keep their behaviour.
type ContainerKind string

const (
	KindApp ContainerKind = ""   // stt-app-<id> on stt-net-<id>
	KindDB  ContainerKind = "db" // stt-db-<id> on stt-dbnet-<id>
)

// Healthcheck is a Docker healthcheck in exec (argv) form — never a shell
// string, for the same injection reasons as start commands.
type Healthcheck struct {
	Test        []string // e.g. ["CMD", "pg_isready", ...]
	Interval    time.Duration
	Timeout     time.Duration
	StartPeriod time.Duration
	Retries     int
}

// ContainerSpec is everything the runtime adapter needs to create a hardened
// user container (SERVERMANAGEMENT.md §3.3).
type ContainerSpec struct {
	AppID          string        // owner UUID from the api (app or database); names are derived from it
	Kind           ContainerKind // app (default) or db
	Image          string
	Env            map[string]string
	Port           int      // container port the app listens on
	Cmd            []string // optional start command override
	User           string   // optional non-root user, e.g. "postgres"
	MemoryBytes    int64
	NanoCPUs       int64
	PidsLimit      int64
	ShmSizeBytes   int64    // 0 keeps the daemon default (64m)
	Runtime        string   // must pass ValidRuntime
	Volumes        []string // named volumes only — never host bind mounts
	ReadonlyRootfs bool
	Healthcheck    *Healthcheck
}

// ContainerName derives the Docker container name for this spec's kind.
func (s ContainerSpec) ContainerName() string {
	if s.Kind == KindDB {
		return DBContainerName(s.AppID)
	}
	return AppContainerName(s.AppID)
}

// NetworkName derives the private network this spec's container lives on.
func (s ContainerSpec) NetworkName() string {
	if s.Kind == KindDB {
		return DBNetworkName(s.AppID)
	}
	return AppNetworkName(s.AppID)
}

// ContainerState is the actual Docker-side state of a container.
type ContainerState struct {
	Exists       bool
	Running      bool
	StartedAt    time.Time
	RestartCount int
	Status       string
	Health       string // "starting" | "healthy" | "unhealthy", "" without a healthcheck
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
