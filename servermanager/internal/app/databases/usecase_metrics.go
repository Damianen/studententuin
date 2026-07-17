package databases

import (
	"servermanager/internal/app/metrics"
	"servermanager/internal/domain"
)

type GetDatabaseMetrics struct {
	store *metrics.Store
}

// Execute returns the database's bucketed series over the window. Store-local
// truth: a valid database the collector never sampled yields empty series,
// never an error — mirrors the apps metrics usecase.
func (u *GetDatabaseMetrics) Execute(dbID string, rng metrics.Range) map[string][]domain.MetricPoint {
	return u.store.Query(domain.DBContainerName(dbID), metrics.DBSeriesKeys, rng.Duration)
}
