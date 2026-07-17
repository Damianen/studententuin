package apps

import (
	"servermanager/internal/app/metrics"
	"servermanager/internal/ports"
)

// Limits carries the server-side resource caps and defaults from config, so
// the app layer doesn't depend on the config package.
type Limits struct {
	DefaultMemoryBytes int64
	MaxMemoryBytes     int64
	DefaultNanoCPUs    int64
	MaxNanoCPUs        int64
	DefaultPidsLimit   int64
	DefaultRuntime     string
}

type Dependencies struct {
	Runtime ports.ContainerRuntime
	// Metrics is the shared collector-fed store; nil only in tests that never
	// hit the metrics endpoint.
	Metrics *metrics.Store
	Limits  Limits
}

type Service struct {
	Run     *RunApp
	Start   *StartApp
	Stop    *StopApp
	Remove  *RemoveApp
	Status  *GetAppStatus
	Logs    *GetAppLogs
	Follow  *FollowAppLogs
	Metrics *GetAppMetrics
}

func NewService(d Dependencies) *Service {
	return &Service{
		Run: &RunApp{
			runtime: d.Runtime,
			limits:  d.Limits,
		},
		Start: &StartApp{
			runtime: d.Runtime,
		},
		Stop: &StopApp{
			runtime: d.Runtime,
		},
		Remove: &RemoveApp{
			runtime: d.Runtime,
		},
		Status: &GetAppStatus{
			runtime: d.Runtime,
		},
		Logs: &GetAppLogs{
			runtime: d.Runtime,
		},
		Follow: &FollowAppLogs{
			runtime: d.Runtime,
		},
		Metrics: &GetAppMetrics{
			store: d.Metrics,
		},
	}
}
