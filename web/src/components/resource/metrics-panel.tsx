import { useMemo } from 'react';
import { makeSeries } from '@/lib/mock_telemetry';
import type { MetricPoint, ResourceKind } from '@/lib/mock_telemetry';
import { CHART_SPECS, formatMetricValue } from '@/lib/metric_specs';
import type { ChartSpec } from '@/lib/metric_specs';
import { useMetrics } from '@/hooks/use_metrics';
import {
	Card,
	CardAction,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { AreaChart } from '@/components/resource/area-chart';
import { DemoBadge } from '@/components/resource/demo-badge';

// resp/req can only come from the edge proxy, so they stay seeded (and
// badged) until Traefik lands in phase 7.
const SEEDED_KEYS = new Set(['resp', 'req']);

interface PanelChart {
	spec: ChartSpec;
	data: MetricPoint[];
	seeded: boolean;
}

export function MetricsPanel({
	kind,
	subdomainId,
	resourceId,
}: {
	kind: ResourceKind;
	subdomainId: string;
	resourceId: string;
}) {
	const { series, loading, error } = useMetrics(kind, subdomainId, resourceId);
	const noun = kind === 'application' ? 'application' : 'database';

	const charts: PanelChart[] = useMemo(
		() =>
			CHART_SPECS[kind].map((spec) => {
				const seeded = kind === 'application' && SEEDED_KEYS.has(spec.key);
				return {
					spec,
					seeded,
					data: seeded
						? makeSeries(`${resourceId}:${spec.key}`, spec)
						: (series[spec.key] ?? []),
				};
			}),
		[kind, resourceId, series],
	);

	return (
		<div className="space-y-5">
			<p className="font-mono text-xs uppercase tracking-[0.22em] text-primary">
				Vital signs — last 24 hours
			</p>

			{error && (
				<p className="text-sm text-destructive">
					Could not fetch metrics: {error}
				</p>
			)}

			<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
				{charts.map(({ spec, data }) => {
					const current =
						data.length > 0 ? data[data.length - 1].value : null;
					const average =
						data.length > 0
							? data.reduce((sum, point) => sum + point.value, 0) /
								data.length
							: null;
					return (
						<Card key={spec.key} className="gap-2 py-5">
							<CardHeader>
								<CardDescription className="font-mono text-[10px] uppercase tracking-[0.25em]">
									{spec.title}
								</CardDescription>
								<CardTitle className="font-display text-3xl font-semibold tracking-tight">
									{current === null ? (
										'—'
									) : (
										<>
											{formatMetricValue(current)}
											<span className="ml-1 text-base font-normal text-muted-foreground">
												{spec.unit.trim()}
											</span>
										</>
									)}
								</CardTitle>
							</CardHeader>
							<CardContent className="text-xs text-muted-foreground">
								{average === null
									? 'no samples yet'
									: `avg ${formatMetricValue(average)}${spec.unit} over 24h`}
							</CardContent>
						</Card>
					);
				})}
			</div>

			<div className="grid gap-5 lg:grid-cols-2">
				{charts.map(({ spec, data, seeded }) => (
					<Card key={spec.key} className="gap-4">
						<CardHeader>
							<CardTitle className="font-mono text-xs uppercase tracking-[0.22em] text-primary">
								{spec.title}
							</CardTitle>
							{seeded && (
								<CardAction>
									<DemoBadge />
								</CardAction>
							)}
						</CardHeader>
						<CardContent>
							{!seeded && loading ? (
								<Skeleton className="h-44 w-full" />
							) : data.length >= 2 ? (
								// AreaChart cannot render fewer than two points.
								<AreaChart data={data} color={spec.color} unit={spec.unit} />
							) : (
								<div className="flex h-44 items-center justify-center rounded-lg border border-dashed px-6 text-center text-sm text-muted-foreground">
									{error
										? 'Metrics are unavailable right now.'
										: `No metrics yet — they appear once the ${noun} is running.`}
								</div>
							)}
						</CardContent>
					</Card>
				))}
			</div>
		</div>
	);
}
