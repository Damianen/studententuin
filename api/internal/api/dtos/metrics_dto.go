package dtos

// MetricPointResponse is one sampled point. Time is RFC3339 UTC; the web
// formats display labels in the browser's timezone.
type MetricPointResponse struct {
	Time  string  `json:"time"`
	Value float64 `json:"value"`
}

// MetricsResponse passes the manager's envelope through verbatim: every
// requested series key is present with a (possibly empty) array, never null.
type MetricsResponse struct {
	Range  string                           `json:"range"`
	Series map[string][]MetricPointResponse `json:"series"`
}
