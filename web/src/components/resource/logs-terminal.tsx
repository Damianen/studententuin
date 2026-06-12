import { useEffect, useMemo, useRef, useState } from 'react';
import { Download, Search, X } from 'lucide-react';
import { makeLogs } from '@/lib/mock_telemetry';
import type { LogLevel, ResourceKind } from '@/lib/mock_telemetry';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
	Card,
	CardAction,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';
import { DemoBadge } from '@/components/resource/demo-badge';

const LEVELS: LogLevel[] = ['info', 'warn', 'error', 'debug'];

const LEVEL_TEXT: Record<LogLevel, string> = {
	info: 'text-primary',
	warn: 'text-status-pending',
	error: 'text-status-failed',
	debug: 'text-muted-foreground',
};

export function LogsTerminal({
	kind,
	seedId,
	resourceName,
}: {
	kind: ResourceKind;
	seedId: string;
	resourceName: string;
}) {
	const logs = useMemo(() => makeLogs(seedId, kind), [seedId, kind]);
	const [search, setSearch] = useState('');
	const [levels, setLevels] = useState<Set<LogLevel>>(new Set(LEVELS));
	const terminalRef = useRef<HTMLDivElement>(null);

	const filtered = logs.filter(
		(log) =>
			levels.has(log.level) &&
			log.message.toLowerCase().includes(search.toLowerCase()),
	);

	useEffect(() => {
		const terminal = terminalRef.current;
		if (terminal) terminal.scrollTop = terminal.scrollHeight;
	}, [filtered.length]);

	const toggleLevel = (level: LogLevel) => {
		setLevels((current) => {
			const next = new Set(current);
			if (next.has(level)) {
				next.delete(level);
			} else {
				next.add(level);
			}
			return next;
		});
	};

	const downloadLogs = () => {
		const text = logs
			.map((log) => `[${log.timestamp}] ${log.level.toUpperCase()}: ${log.message}`)
			.join('\n');
		const url = URL.createObjectURL(new Blob([text], { type: 'text/plain' }));
		const anchor = document.createElement('a');
		anchor.href = url;
		anchor.download = `${resourceName}-logs.txt`;
		anchor.click();
		URL.revokeObjectURL(url);
	};

	return (
		<Card className="gap-0 overflow-hidden pb-0">
			<CardHeader className="pb-5">
				<CardTitle className="font-mono text-xs uppercase tracking-[0.22em] text-primary">
					Runtime logs
				</CardTitle>
				<CardDescription>
					What {resourceName} has been murmuring lately.
				</CardDescription>
				<CardAction className="flex items-center gap-2">
					<DemoBadge />
					<Button variant="outline" size="sm" onClick={downloadLogs}>
						<Download className="size-3.5" />
						Download
					</Button>
				</CardAction>
			</CardHeader>
			<div className="flex flex-wrap items-center gap-2 px-6 pb-5">
				<div className="relative min-w-48 flex-1">
					<Search className="absolute left-3 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
					<Input
						value={search}
						onChange={(event) => setSearch(event.target.value)}
						placeholder="Search logs…"
						className="pl-9 pr-9"
					/>
					{search && (
						<Button
							variant="ghost"
							size="icon"
							className="absolute right-0 top-0 h-full text-muted-foreground"
							onClick={() => setSearch('')}
							aria-label="Clear search"
						>
							<X className="size-3.5" />
						</Button>
					)}
				</div>
				<div className="flex gap-1.5">
					{LEVELS.map((level) => (
						<button
							key={level}
							type="button"
							onClick={() => toggleLevel(level)}
							className={cn(
								'rounded-full border px-3 py-1 font-mono text-[10px] uppercase tracking-widest transition-colors',
								levels.has(level)
									? 'border-primary/50 bg-primary/10 text-primary'
									: 'border-border text-muted-foreground hover:text-foreground',
							)}
						>
							{level}
						</button>
					))}
				</div>
			</div>
			<div
				ref={terminalRef}
				className="max-h-[26rem] overflow-y-auto border-t bg-[oklch(0.125_0.02_155)] px-5 py-4 font-mono text-xs leading-6"
			>
				{filtered.map((log) => (
					<div
						key={log.id}
						className="flex gap-3 rounded px-1 hover:bg-foreground/5"
					>
						<span className="shrink-0 text-muted-foreground/60">
							{log.timestamp}
						</span>
						<span
							className={cn(
								'w-12 shrink-0 font-semibold uppercase',
								LEVEL_TEXT[log.level],
							)}
						>
							{log.level}
						</span>
						<span className="whitespace-pre-wrap break-all text-foreground/85">
							{log.message}
						</span>
					</div>
				))}
				{filtered.length === 0 && (
					<p className="py-8 text-center text-muted-foreground">
						No log lines match — the garden is quiet.
					</p>
				)}
				<span className="terminal-cursor mt-1 inline-block px-1 text-primary">
					▌
				</span>
			</div>
		</Card>
	);
}
