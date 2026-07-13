package databases

import (
	"context"
	"errors"

	"servermanager/internal/domain"
	"servermanager/internal/ports"
)

type GetDatabaseStatus struct {
	runtime ports.ContainerRuntime
	jobs    *JobStore
}

// DatabaseStatus combines the actual Docker-side state with the latest
// provision job (nil after a manager restart — the api's poller treats
// "no container, no job" as a lost provision).
type DatabaseStatus struct {
	State     domain.ContainerState
	Provision *domain.ProvisionJob
}

// Execute reports the database's Docker state. A missing container is not an
// error — that's what ContainerState.Exists is for.
func (u *GetDatabaseStatus) Execute(ctx context.Context, dbID string) (*DatabaseStatus, error) {
	state, err := u.runtime.Inspect(ctx, domain.DBContainerName(dbID))
	if errors.Is(err, domain.ErrNotFound) {
		state = &domain.ContainerState{Exists: false}
	} else if err != nil {
		return nil, err
	}

	status := &DatabaseStatus{State: *state}
	if job, ok := u.jobs.Latest(dbID); ok {
		status.Provision = &job
	}
	return status, nil
}
