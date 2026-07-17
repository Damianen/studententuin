export interface MetricPointDto {
	time: string; // RFC3339
	value: number;
}

// The manager's series envelope, proxied verbatim by the api (phase 6).
// Applications carry cpu/mem, databases conn/qps/cpu/disk; keys may be
// missing or empty while nothing has sampled yet.
export interface MetricsDto {
	range: string;
	series?: Record<string, MetricPointDto[]>;
}
