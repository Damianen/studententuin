package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// ManagedContainer is one running container the manager owns, discovered by
// its derived name. ID is the Docker container ID (changes at cutover), Name
// the stable derived name, OwnerID the app/database UUID the api knows.
type ManagedContainer struct {
	ID      string
	Name    string
	Kind    ContainerKind
	OwnerID string
}

// StatsSample is one host-level resource reading for a container (§6 phase 6:
// cgroup reads, not `docker stats`).
type StatsSample struct {
	CPUUsageUsec       uint64  // cumulative cpu time from cpu.stat
	CPULimitCores      float64 // from cpu.max; 0 = unlimited
	MemWorkingSetBytes int64   // memory.current - inactive_file, clamped at 0
}

// MetricPoint is one value on a metrics series.
type MetricPoint struct {
	Time  time.Time
	Value float64
}

// ParseManagedName maps a Docker container name onto the manager's naming
// scheme. Strict on purpose: Docker's name filter matches substrings, so the
// discovery path must re-validate that the remainder after the stt-app-/
// stt-db- prefix is exactly a canonical UUID (volumes like stt-db-data-<id>
// and free-named lookalikes fall out here). A leading '/' (Docker's list API
// form) is stripped.
func ParseManagedName(name string) (ManagedContainer, bool) {
	name = strings.TrimPrefix(name, "/")

	var kind ContainerKind
	var rest string
	switch {
	case strings.HasPrefix(name, "stt-app-"):
		kind, rest = KindApp, strings.TrimPrefix(name, "stt-app-")
	case strings.HasPrefix(name, "stt-db-"):
		kind, rest = KindDB, strings.TrimPrefix(name, "stt-db-")
	default:
		return ManagedContainer{}, false
	}

	// uuid.Parse also accepts urn/braced/32-hex forms; require the canonical
	// 36-character form the manager generates names from.
	if len(rest) != 36 {
		return ManagedContainer{}, false
	}
	id, err := uuid.Parse(rest)
	if err != nil {
		return ManagedContainer{}, false
	}
	return ManagedContainer{Name: name, Kind: kind, OwnerID: id.String()}, true
}
