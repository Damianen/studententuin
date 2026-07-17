
export type ResourceKind = 'application' | 'database';
export type LogLevel = 'info' | 'warn' | 'error' | 'debug';
export type DeploymentStatus = 'success' | 'failed' | 'in_progress';

export interface MetricPoint {
	time: string;
	value: number;
}

export interface LogEntry {
	id: string;
	timestamp: string;
	level: LogLevel;
	message: string;
}

export interface DeploymentRecord {
	id: string;
	status: DeploymentStatus;
	commit: string;
	branch: string;
	message: string;
	author: string;
	timeAgo: string;
	duration: string;
	durationSeconds?: number;
}

/*
 * Only the resp/req application charts still grow on sample data — they can
 * only come from the edge proxy, which lands with Traefik in phase 7. The
 * series are seeded by resource id so values stay stable across navigation
 * and reloads instead of reshuffling on every render.
 */

function hashSeed(input: string): number {
	let hash = 2166136261;
	for (let index = 0; index < input.length; index++) {
		hash ^= input.charCodeAt(index);
		hash = Math.imul(hash, 16777619);
	}
	return hash >>> 0;
}

function mulberry32(seed: number): () => number {
	let state = seed;
	return () => {
		state = (state + 0x6d2b79f5) | 0;
		let t = Math.imul(state ^ (state >>> 15), 1 | state);
		t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
		return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
	};
}

function hourLabel(hoursAgo: number): string {
	const date = new Date(Date.now() - hoursAgo * 3_600_000);
	return `${String(date.getHours()).padStart(2, '0')}:00`;
}

export function makeSeries(
	seed: string,
	{ min, max, points = 25 }: { min: number; max: number; points?: number },
): MetricPoint[] {
	const rng = mulberry32(hashSeed(seed));
	const range = max - min;
	let value = min + range * (0.3 + rng() * 0.4);
	const series: MetricPoint[] = [];
	for (let i = points - 1; i >= 0; i--) {
		// Random walk with a soft pull back toward the middle of the range.
		value += (min + range / 2 - value) * 0.15 + (rng() - 0.5) * range * 0.3;
		value = Math.min(max, Math.max(min, value));
		series.push({ time: hourLabel(i), value: Math.round(value * 10) / 10 });
	}
	return series;
}
