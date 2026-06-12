import type { DatabaseDto } from '@/dtos/database_dtos';

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
}

/*
 * The API has no telemetry endpoints yet, so the detail pages grow on
 * sample data. Everything is seeded by resource id so values stay stable
 * across navigation and reloads instead of reshuffling on every render.
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

export function seededInt(seed: string, min: number, max: number): number {
	const rng = mulberry32(hashSeed(seed));
	return Math.floor(min + rng() * (max - min + 1));
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

function int(rng: () => number, min: number, max: number): number {
	return Math.floor(min + rng() * (max - min + 1));
}

function hex(rng: () => number, length: number): string {
	return Array.from({ length }, () => Math.floor(rng() * 16).toString(16)).join(
		'',
	);
}

type LogFactory = (rng: () => number) => { level: LogLevel; message: string };

const APP_LOG_LINES: LogFactory[] = [
	(rng) => ({ level: 'info', message: `GET / 200 ${int(rng, 3, 40)}ms` }),
	(rng) => ({
		level: 'info',
		message: `GET /api/items 200 ${int(rng, 12, 140)}ms`,
	}),
	(rng) => ({
		level: 'info',
		message: `POST /api/auth/login 200 ${int(rng, 80, 240)}ms`,
	}),
	(rng) => ({
		level: 'info',
		message: `GET /assets/index-${hex(rng, 8)}.js 200 ${int(rng, 2, 18)}ms`,
	}),
	() => ({ level: 'info', message: 'health check ok' }),
	(rng) => ({
		level: 'debug',
		message: `cache hit for key session:${hex(rng, 6)}`,
	}),
	(rng) => ({
		level: 'debug',
		message: `websocket ping rtt=${int(rng, 8, 48)}ms`,
	}),
	(rng) => ({
		level: 'warn',
		message: `slow request: GET /api/search took ${int(rng, 900, 2100)}ms`,
	}),
	(rng) => ({
		level: 'warn',
		message: `rate limit hit for 145.97.${int(rng, 0, 255)}.${int(rng, 0, 255)}`,
	}),
	() => ({
		level: 'error',
		message: 'upstream timeout after 5000ms, retrying (1/3)',
	}),
];

const DB_LOG_LINES: LogFactory[] = [
	(rng) => ({
		level: 'info',
		message: `checkpoint complete: wrote ${int(rng, 40, 420)} buffers (${int(rng, 1, 9)}.${int(rng, 0, 9)}%)`,
	}),
	(rng) => ({
		level: 'info',
		message: `connection received: host=10.0.${int(rng, 0, 9)}.${int(rng, 2, 250)} port=${int(rng, 40000, 60000)}`,
	}),
	() => ({
		level: 'info',
		message: 'connection authorized: user=app database=main',
	}),
	(rng) => ({
		level: 'info',
		message: `automatic analyze of table "main.public.sessions" took ${int(rng, 30, 400)}ms`,
	}),
	(rng) => ({
		level: 'debug',
		message: `autovacuum: processed "public.events" in ${int(rng, 60, 900)}ms`,
	}),
	(rng) => ({
		level: 'warn',
		message: `slow query (${int(rng, 1100, 3400)}ms): SELECT * FROM orders WHERE created_at > $1`,
	}),
	(rng) => ({
		level: 'warn',
		message: `temporary file: size ${int(rng, 8, 96)} MB for sort operation`,
	}),
	() => ({
		level: 'error',
		message: 'could not receive data from client: connection reset by peer',
	}),
];

function formatTimestamp(date: Date): string {
	const pad = (value: number) => String(value).padStart(2, '0');
	return (
		`${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ` +
		`${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
	);
}

export function makeLogs(seed: string, kind: ResourceKind): LogEntry[] {
	const rng = mulberry32(hashSeed(`${seed}:logs`));
	const pool = kind === 'application' ? APP_LOG_LINES : DB_LOG_LINES;
	const count = 28;
	const gaps = Array.from({ length: count }, () => 20 + rng() * 160);
	let elapsed = 0;
	const offsets = gaps.map((gap) => (elapsed += gap));
	const now = Date.now();
	return offsets.map((offset, index) => {
		const line = pool[Math.floor(rng() * pool.length)](rng);
		return {
			id: String(index),
			timestamp: formatTimestamp(new Date(now - (elapsed - offset) * 1000)),
			level: line.level,
			message: line.message,
		};
	});
}

const DEPLOY_MESSAGES = [
	'Fix login redirect after session expiry',
	'Add caching headers to static assets',
	'Bump dependencies and patch advisories',
	'Tidy the navbar on small screens',
	'Add health check endpoint',
	'Refactor request logging middleware',
	'Speed up image builds with layer caching',
	'Handle empty states on the projects page',
];

const DEPLOY_AGES = [
	'2 hours ago',
	'yesterday',
	'2 days ago',
	'4 days ago',
	'5 days ago',
	'last week',
];

export function makeDeployments(
	seed: string,
	branch: string,
): DeploymentRecord[] {
	const rng = mulberry32(hashSeed(`${seed}:deploys`));
	const messages = [...DEPLOY_MESSAGES].sort(() => rng() - 0.5);
	return DEPLOY_AGES.map((timeAgo, index) => {
		const failed = index === 3;
		const minutes = int(rng, 1, 4);
		const seconds = int(rng, 0, 59);
		return {
			id: String(index),
			status: failed ? 'failed' : index === 0 && rng() > 0.6 ? 'in_progress' : 'success',
			commit: hex(rng, 7),
			branch,
			message: messages[index % messages.length],
			author: rng() > 0.8 ? 'dependabot[bot]' : 'damian',
			timeAgo,
			duration: failed
				? `${seconds}s`
				: `${minutes}m ${String(seconds).padStart(2, '0')}s`,
		};
	});
}

export function defaultEnvVars(
	database?: DatabaseDto,
): Array<{ key: string; value: string }> {
	const vars = [
		{ key: 'NODE_ENV', value: 'production' },
		{ key: 'PORT', value: '3000' },
	];
	if (database?.connection_string) {
		vars.push({ key: 'DATABASE_URL', value: database.connection_string });
	}
	return vars;
}
