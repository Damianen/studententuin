package apps

import "servermanager/internal/ports"

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
	Limits  Limits
}

type Service struct {
	Run    *RunApp
	Start  *StartApp
	Stop   *StopApp
	Remove *RemoveApp
	Status *GetAppStatus
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
	}
}
