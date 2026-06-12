import { makeSeries } from '@/lib/mock_telemetry';
import type { ResourceKind } from '@/lib/mock_telemetry';

export const CHART_COLORS = {
	chartreuse: 'var(--primary)',
	terracotta: 'var(--terracotta)',
	amber: 'var(--status-pending)',
	moss: 'var(--status-running)',
} as const;

export interface ChartSpec {
	key: string;
	title: string;
	unit: string;
	min: number;
	max: number;
	color: string;
}

export const CHART_SPECS: Record<ResourceKind, ChartSpec[]> = {
	application: [
		{ key: 'cpu', title: 'CPU usage', unit: '%', min: 8, max: 82, color: CHART_COLORS.chartreuse },
		{ key: 'mem', title: 'Memory', unit: ' MB', min: 310, max: 960, color: CHART_COLORS.terracotta },
		{ key: 'resp', title: 'Response time', unit: ' ms', min: 28, max: 210, color: CHART_COLORS.amber },
		{ key: 'req', title: 'Requests', unit: '/min', min: 60, max: 880, color: CHART_COLORS.moss },
	],
	database: [
		{ key: 'conn', title: 'Active connections', unit: '', min: 2, max: 46, color: CHART_COLORS.chartreuse },
		{ key: 'qps', title: 'Queries', unit: '/s', min: 8, max: 230, color: CHART_COLORS.moss },
		{ key: 'cpu', title: 'CPU usage', unit: '%', min: 3, max: 58, color: CHART_COLORS.amber },
		{ key: 'disk', title: 'Disk usage', unit: ' MB', min: 180, max: 540, color: CHART_COLORS.terracotta },
	],
};

export function formatMetricValue(value: number): string {
	return value >= 100
		? String(Math.round(value))
		: String(Math.round(value * 10) / 10);
}

/** The single chart the overview tab previews, per resource kind. */
export function overviewSpec(kind: ResourceKind): ChartSpec {
	return kind === 'application'
		? CHART_SPECS.application[3]
		: CHART_SPECS.database[0];
}

export function headlineStats(
	kind: ResourceKind,
	seedId: string,
): Array<{ label: string; value: string }> {
	return CHART_SPECS[kind].map((spec) => {
		const data = makeSeries(`${seedId}:${spec.key}`, spec);
		return {
			label: spec.title,
			value: `${formatMetricValue(data[data.length - 1].value)}${spec.unit}`,
		};
	});
}
