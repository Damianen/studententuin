import { useId, useMemo, useState } from 'react';
import type { MetricPoint } from '@/lib/mock_telemetry';
import { CHART_COLORS, formatMetricValue } from '@/lib/metric_specs';
import { cn } from '@/lib/utils';

const WIDTH = 600;
const HEIGHT = 210;
const PAD_X = 6;
const PAD_TOP = 18;
const PAD_BOTTOM = 26;
const INNER_WIDTH = WIDTH - PAD_X * 2;
const INNER_HEIGHT = HEIGHT - PAD_TOP - PAD_BOTTOM;

function niceCeil(value: number): number {
	if (value <= 0) return 1;
	const power = 10 ** Math.floor(Math.log10(value));
	for (const step of [1, 2, 2.5, 5, 10]) {
		if (value <= step * power) return step * power;
	}
	return 10 * power;
}

export function AreaChart({
	data,
	color = CHART_COLORS.chartreuse,
	unit = '',
	className,
}: {
	data: MetricPoint[];
	color?: string;
	unit?: string;
	className?: string;
}) {
	const gradientId = useId();
	const [hovered, setHovered] = useState<number | null>(null);

	const { points, linePath, areaPath, ceiling } = useMemo(() => {
		const ceiling = niceCeil(Math.max(...data.map((point) => point.value)));
		const points = data.map((point, index) => ({
			...point,
			x: PAD_X + (index / Math.max(1, data.length - 1)) * INNER_WIDTH,
			y: PAD_TOP + INNER_HEIGHT - (point.value / ceiling) * INNER_HEIGHT,
		}));
		const linePath = points
			.map(
				(point, index) =>
					`${index === 0 ? 'M' : 'L'}${point.x.toFixed(1)},${point.y.toFixed(1)}`,
			)
			.join(' ');
		const baseline = HEIGHT - PAD_BOTTOM;
		const areaPath =
			`${linePath} L${points[points.length - 1].x.toFixed(1)},${baseline}` +
			` L${points[0].x.toFixed(1)},${baseline} Z`;
		return { points, linePath, areaPath, ceiling };
	}, [data]);

	const handleMove = (event: React.MouseEvent<SVGSVGElement>) => {
		const rect = event.currentTarget.getBoundingClientRect();
		const x = ((event.clientX - rect.left) / rect.width) * WIDTH;
		const index = Math.round(((x - PAD_X) / INNER_WIDTH) * (data.length - 1));
		setHovered(Math.min(data.length - 1, Math.max(0, index)));
	};

	const active = hovered === null ? null : points[hovered];
	const tooltipLeft = active !== null && active.x > WIDTH * 0.72;

	return (
		<svg
			viewBox={`0 0 ${WIDTH} ${HEIGHT}`}
			className={cn('w-full', className)}
			role="img"
			onMouseMove={handleMove}
			onMouseLeave={() => setHovered(null)}
		>
			<defs>
				<linearGradient id={gradientId} x1="0" y1="0" x2="0" y2="1">
					<stop offset="0%" stopColor={color} stopOpacity="0.32" />
					<stop offset="100%" stopColor={color} stopOpacity="0.02" />
				</linearGradient>
			</defs>
			{[0.5, 1].map((fraction) => {
				const y = PAD_TOP + INNER_HEIGHT - fraction * INNER_HEIGHT;
				return (
					<g key={fraction}>
						<line
							x1={PAD_X}
							x2={WIDTH - PAD_X}
							y1={y}
							y2={y}
							stroke="var(--border)"
							strokeDasharray="3 5"
						/>
						<text
							x={PAD_X + 2}
							y={y - 5}
							fontSize="9"
							fontFamily="var(--font-mono)"
							fill="var(--muted-foreground)"
							opacity="0.8"
						>
							{formatMetricValue(fraction * ceiling)}
							{unit}
						</text>
					</g>
				);
			})}
			<line
				x1={PAD_X}
				x2={WIDTH - PAD_X}
				y1={HEIGHT - PAD_BOTTOM}
				y2={HEIGHT - PAD_BOTTOM}
				stroke="var(--border)"
			/>
			<path d={areaPath} fill={`url(#${gradientId})`} />
			<path
				d={linePath}
				fill="none"
				stroke={color}
				strokeWidth="2"
				strokeLinejoin="round"
			/>
			{points.map((point, index) =>
				index % 6 === 0 ? (
					<text
						key={`${point.time}-${index}`}
						x={point.x}
						y={HEIGHT - 8}
						fontSize="9"
						textAnchor={
							index === 0 ? 'start' : index === points.length - 1 ? 'end' : 'middle'
						}
						fontFamily="var(--font-mono)"
						fill="var(--muted-foreground)"
					>
						{point.time}
					</text>
				) : null,
			)}
			{active && (
				<g pointerEvents="none">
					<line
						x1={active.x}
						x2={active.x}
						y1={PAD_TOP}
						y2={HEIGHT - PAD_BOTTOM}
						stroke={color}
						strokeOpacity="0.35"
					/>
					<circle
						cx={active.x}
						cy={active.y}
						r="3.5"
						fill={color}
						stroke="var(--background)"
						strokeWidth="1.5"
					/>
					<g
						transform={`translate(${tooltipLeft ? active.x - 122 : active.x + 12}, ${Math.max(
							PAD_TOP,
							active.y - 38,
						)})`}
					>
						<rect
							width="110"
							height="36"
							rx="9"
							fill="var(--popover)"
							stroke="var(--border)"
						/>
						<text
							x="10"
							y="15"
							fontSize="9"
							fontFamily="var(--font-mono)"
							fill="var(--muted-foreground)"
						>
							{active.time}
						</text>
						<text
							x="10"
							y="28"
							fontSize="11"
							fontFamily="var(--font-mono)"
							fill="var(--foreground)"
						>
							{formatMetricValue(active.value)}
							{unit}
						</text>
					</g>
				</g>
			)}
		</svg>
	);
}
