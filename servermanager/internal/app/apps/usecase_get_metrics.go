package apps

import (
	"servermanager/internal/app/metrics"
	"servermanager/internal/domain"
)

type GetAppMetrics struct {
	store *metrics.Store
}

// Execute returns the app's bucketed series over the window. Store-local
// truth: a valid app the collector never sampled (new, stopped, or the
// manager just restarted) yields empty series, never an error.
func (u *GetAppMetrics) Execute(appID string, rng metrics.Range) map[string][]domain.MetricPoint {
	return u.store.Query(domain.AppContainerName(appID), metrics.AppSeriesKeys, rng.Duration)
}
