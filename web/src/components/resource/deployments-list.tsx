import { useCallback, useMemo } from 'react';
import type { LucideIcon } from 'lucide-react';
import { Check, GitBranch, RefreshCw, Rocket, User, X } from 'lucide-react';
import { toast } from 'sonner';
import type { DeploymentStatus } from '@/lib/mock_telemetry';
import { cn } from '@/lib/utils';
import { formatDuration } from '@/lib/format_time';
import ApplicationController from '@/controllers/application_controller';
import { useAsync } from '@/hooks/use_async';
import { STAGE_LABELS, useDeployment } from '@/hooks/use_deployment';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import {
	Card,
	CardAction,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';

const STATUS_META: Record<
	DeploymentStatus,
	{ label: string; icon: LucideIcon; circle: string; text: string }
> = {
	success: {
		label: 'Success',
		icon: Check,
		circle: 'bg-status-running/10 text-status-running',
		text: 'text-status-running',
	},
	failed: {
		label: 'Failed',
		icon: X,
		circle: 'bg-status-failed/10 text-status-failed',
		text: 'text-status-failed',
	},
	in_progress: {
		label: 'Deploying',
		icon: RefreshCw,
		circle: 'bg-status-pending/10 text-status-pending',
		text: 'text-status-pending',
	},
};

export function DeploymentsList({
	branch,
	resourceName,
	subdomainId,
	appId,
	repoUrl,
	onChanged,
}: {
	branch: string;
	resourceName: string;
	subdomainId: string;
	appId: string;
	repoUrl: string;
	onChanged?: () => void;
}) {
	const {
		data,
		loading,
		error: listError,
		reload,
	} = useAsync(
		useCallback(
			() => ApplicationController.getDeployments(subdomainId, appId),
			[subdomainId, appId],
		),
	);
	const deployments = useMemo(() => data ?? [], [data]);

	const { deploying, stage, error, buildLog, deploy } = useDeployment(
		subdomainId,
		appId,
		(ok) => {
			if (ok) {
				toast.success(`${resourceName} is live`);
			} else {
				toast.error('Deployment failed');
			}
			void reload();
			onChanged?.();
		},
	);

	const stats = useMemo(() => {
		if (deployments.length === 0) return [];
		const terminal = deployments.filter(
			(deployment) => deployment.status !== 'in_progress',
		);
		const succeeded = terminal.filter(
			(deployment) => deployment.status === 'success',
		);
		const durations = succeeded
			.map((deployment) => deployment.durationSeconds)
			.filter((seconds): seconds is number => seconds != null);
		return [
			{
				label: 'Success rate',
				value:
					terminal.length > 0
						? `${Math.round((succeeded.length / terminal.length) * 100)}%`
						: '—',
				note: `${succeeded.length} of ${terminal.length} finished deploys`,
			},
			{
				label: 'Avg deploy time',
				value:
					durations.length > 0
						? formatDuration(
								durations.reduce((sum, seconds) => sum + seconds, 0) /
									durations.length,
							)
						: '—',
				note: 'Successful deploys, clone to live',
			},
			{
				label: 'Last deployment',
				value: deployments[0].timeAgo,
				note: STATUS_META[deployments[0].status].label,
			},
		];
	}, [deployments]);

	return (
		<div className="space-y-5">
			<Card className="gap-0 overflow-hidden pb-0">
				<CardHeader className="pb-6">
					<CardTitle className="font-mono text-xs uppercase tracking-[0.22em] text-primary">
						Deployments
					</CardTitle>
					<CardDescription>
						Every push to {branch} replants {resourceName} fresh.
					</CardDescription>
					<CardAction>
						<Button
							size="sm"
							disabled={deploying || !repoUrl}
							title={!repoUrl ? 'Set a repository URL first' : undefined}
							onClick={deploy}
						>
							{deploying ? (
								<RefreshCw className="size-3.5 animate-spin" />
							) : (
								<Rocket className="size-3.5" />
							)}
							{deploying ? 'Deploying…' : 'Deploy now'}
						</Button>
					</CardAction>
				</CardHeader>
				<div className="divide-y border-t">
					{deploying && (
						<div className="flex items-start gap-4 bg-status-pending/5 px-6 py-5">
							<span className="mt-0.5 flex size-9 shrink-0 items-center justify-center rounded-full bg-status-pending/10 text-status-pending">
								<RefreshCw className="size-4 animate-spin" />
							</span>
							<div className="min-w-0 flex-1 space-y-1.5">
								<div className="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
									<span className="font-medium text-status-pending">
										Deploying
									</span>
									<span>·</span>
									<span>{STAGE_LABELS[stage ?? ''] ?? stage}</span>
								</div>
								<p className="truncate text-sm font-medium">
									Building {resourceName} from {branch}
								</p>
							</div>
						</div>
					)}
					{!deploying && error && (
						<div className="space-y-3 bg-status-failed/5 px-6 py-5">
							<div className="flex items-start gap-4">
								<span className="mt-0.5 flex size-9 shrink-0 items-center justify-center rounded-full bg-status-failed/10 text-status-failed">
									<X className="size-4" />
								</span>
								<div className="min-w-0 flex-1 space-y-1.5">
									<p className="text-xs font-medium text-status-failed">
										Last deployment failed
									</p>
									<p className="text-sm">{error}</p>
								</div>
							</div>
							{buildLog && (
								<details className="text-xs">
									<summary className="cursor-pointer text-muted-foreground hover:text-foreground">
										Build log
									</summary>
									<pre className="mt-2 max-h-64 overflow-auto whitespace-pre-wrap rounded bg-muted/60 p-3 font-mono text-[11px] leading-relaxed">
										{buildLog}
									</pre>
								</details>
							)}
						</div>
					)}
					{loading &&
						[0, 1, 2].map((row) => (
							<div key={row} className="flex items-start gap-4 px-6 py-5">
								<Skeleton className="size-9 rounded-full" />
								<div className="flex-1 space-y-2 py-1">
									<Skeleton className="h-3 w-44" />
									<Skeleton className="h-4 w-72" />
								</div>
							</div>
						))}
					{!loading && listError && (
						<p className="px-6 py-10 text-center text-sm text-status-failed">
							Could not fetch deployments: {listError}
						</p>
					)}
					{!loading && !listError && deployments.length === 0 && (
						<p className="px-6 py-10 text-center text-sm text-muted-foreground">
							No deployments yet — plant the first one with Deploy now.
						</p>
					)}
					{deployments.map((deployment) => {
						const meta = STATUS_META[deployment.status];
						const StatusIcon = meta.icon;
						return (
							<div
								key={deployment.id}
								className="flex items-start gap-4 px-6 py-5 transition-colors hover:bg-muted/40"
							>
								<span
									className={cn(
										'mt-0.5 flex size-9 shrink-0 items-center justify-center rounded-full',
										meta.circle,
									)}
								>
									<StatusIcon
										className={cn(
											'size-4',
											deployment.status === 'in_progress' && 'animate-spin',
										)}
									/>
								</span>
								<div className="min-w-0 flex-1 space-y-1.5">
									<div className="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
										<span className={cn('font-medium', meta.text)}>
											{meta.label}
										</span>
										<span>·</span>
										<span>{deployment.timeAgo}</span>
										<span>·</span>
										<span>{deployment.duration}</span>
									</div>
									<p className="truncate text-sm font-medium">
										{deployment.message}
									</p>
									<div className="flex flex-wrap items-center gap-4 text-xs text-muted-foreground">
										<span className="flex items-center gap-1.5">
											<GitBranch className="size-3.5" />
											{deployment.branch}
										</span>
										<code className="rounded bg-muted/60 px-1.5 py-0.5 font-mono text-[11px]">
											{deployment.commit}
										</code>
										<span className="flex items-center gap-1.5">
											<User className="size-3.5" />
											{deployment.author}
										</span>
									</div>
								</div>
							</div>
						);
					})}
				</div>
			</Card>

			{stats.length > 0 && (
				<div className="grid gap-4 sm:grid-cols-3">
					{stats.map((stat) => (
						<Card key={stat.label} className="gap-2 py-5">
							<CardHeader>
								<CardDescription className="font-mono text-[10px] uppercase tracking-[0.25em]">
									{stat.label}
								</CardDescription>
								<CardTitle className="font-display text-2xl font-semibold tracking-tight">
									{stat.value}
								</CardTitle>
							</CardHeader>
							<CardContent className="text-xs text-muted-foreground">
								{stat.note}
							</CardContent>
						</Card>
					))}
				</div>
			)}
		</div>
	);
}
