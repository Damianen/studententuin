package deployments

import "servermanager/internal/domain"

type GetDeployment struct {
	jobs *JobStore
}

// Execute returns a snapshot copy of the job; unknown IDs (including jobs
// lost to a manager restart) map to ErrNotFound.
func (u *GetDeployment) Execute(id string) (domain.DeploymentJob, error) {
	return u.jobs.Get(id)
}
