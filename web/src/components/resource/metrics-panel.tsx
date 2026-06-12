import { useMemo } from 'react';
import { makeSeries, seededInt } from '@/lib/mock_telemetry';
import type { ResourceKind } from '@/lib/mock_telemetry';
import { CHART_SPECS, formatMetricValue } from '@/lib/metric_specs';
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';
import { AreaChart } from '@/components/resource/area-chart';
import { DemoBadge } from '@/components/resource/demo-badge';

interface Fact {
	label: string;
	value: string;
	note: string;
}

function appFacts(seedId: string): Fact[] {
	const errors = seededInt(`${seedId}:errors`, 0, 8);
	const requests = seededInt(`${seedId}:requests`, 6000, 30000);
	return [
		{
			label: 'Error rate',
			value: `${((errors / requests) * 100).toFixed(2)}%`,
			note: `${errors} of ${requests.toLocaleString()} requests this week`,
		},
		{
			label: 'Network I/O',
			value: `${seededInt(`${seedId}:net-in`, 2, 38)} MB/s in`,
			note: `${seededInt(`${seedId}:net-out`, 2, 30)} MB/s out`,
		},
		{
			label: 'Open connections',
			value: String(seededInt(`${seedId}:open-conn`, 12, 240)),
			note: 'Currently held open by clients',
		},
	];
}

function dbFacts(seedId: string): Fact[] {
	return [
		{
			label: 'Cache hit ratio',
			value: `${(95 + seededInt(`${seedId}:cache`, 0, 49) / 10).toFixed(1)}%`,
			note: 'Shared buffers, last 24 hours',
		},
		{
			label: 'Disk used',
			value: `${seededInt(`${seedId}:disk-used`, 120, 720)} MB`,
			note: 'of 1 GB volume',
		},
		{
			label: 'Slow queries',
			value: String(seededInt(`${seedId}:slow`, 0, 7)),
			note: 'Slower than 1s in the last 24 hours',
		},
	];
}

export function MetricsPanel({
	kind,
	seedId,
}: {
	kind: ResourceKind;
	seedId: string;
}) {
	const charts = useMemo(
		() =>
			CHART_SPECS[kind].map((spec) => ({
				spec,
				data: makeSeries(`${seedId}:${spec.key}`, spec),
			})),
		[kind, seedId],
	);
	const facts = useMemo(
		() => (kind === 'application' ? appFacts(seedId) : dbFacts(seedId)),
		[kind, seedId],
	);

	return (
		<div className="space-y-5">
			<div className="flex flex-wrap items-center justify-between gap-3">
				<p className="font-mono text-xs uppercase tracking-[0.22em] text-primary">
					Vital signs — last 24 hours
				</p>
				<DemoBadge />
			</div>

			<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
				{charts.map(({ spec, data }) => {
					const current = data[data.length - 1].value;
					const average =
						data.reduce((sum, point) => sum + point.value, 0) / data.length;
					return (
						<Card key={spec.key} className="gap-2 py-5">
							<CardHeader>
								<CardDescription className="font-mono text-[10px] uppercase tracking-[0.25em]">
									{spec.title}
								</CardDescription>
								<CardTitle className="font-display text-3xl font-semibold tracking-tight">
									{formatMetricValue(current)}
									<span className="ml-1 text-base font-normal text-muted-foreground">
										{spec.unit.trim()}
									</span>
								</CardTitle>
							</CardHeader>
							<CardContent className="text-xs text-muted-foreground">
								avg {formatMetricValue(average)}
								{spec.unit} over 24h
							</CardContent>
						</Card>
					);
				})}
			</div>

			<div className="grid gap-5 lg:grid-cols-2">
				{charts.map(({ spec, data }) => (
					<Card key={spec.key} className="gap-4">
						<CardHeader>
							<CardTitle className="font-mono text-xs uppercase tracking-[0.22em] text-primary">
								{spec.title}
							</CardTitle>
						</CardHeader>
						<CardContent>
							<AreaChart data={data} color={spec.color} unit={spec.unit} />
						</CardContent>
					</Card>
				))}
			</div>

			<div className="grid gap-4 md:grid-cols-3">
				{facts.map((fact) => (
					<Card key={fact.label} className="gap-2 py-5">
						<CardHeader>
							<CardDescription className="font-mono text-[10px] uppercase tracking-[0.25em]">
								{fact.label}
							</CardDescription>
							<CardTitle className="font-display text-2xl font-semibold tracking-tight">
								{fact.value}
							</CardTitle>
						</CardHeader>
						<CardContent className="text-xs text-muted-foreground">
							{fact.note}
						</CardContent>
					</Card>
				))}
			</div>
		</div>
	);
}
