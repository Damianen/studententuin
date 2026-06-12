import { useMemo } from 'react';
import type { LucideIcon } from 'lucide-react';
import { Check, GitBranch, RefreshCw, Rocket, User, X } from 'lucide-react';
import { toast } from 'sonner';
import { makeDeployments, seededInt } from '@/lib/mock_telemetry';
import type { DeploymentStatus } from '@/lib/mock_telemetry';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import {
	Card,
	CardAction,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';
import { DemoBadge } from '@/components/resource/demo-badge';

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
	seedId,
	branch,
	resourceName,
}: {
	seedId: string;
	branch: string;
	resourceName: string;
}) {
	const deployments = useMemo(
		() => makeDeployments(seedId, branch),
		[seedId, branch],
	);

	const stats = useMemo(() => {
		const failures = seededInt(`${seedId}:deploy-fail`, 1, 8);
		return [
			{
				label: 'Total deployments',
				value: String(seededInt(`${seedId}:deploy-total`, 24, 140)),
				note: 'Last 90 days',
			},
			{
				label: 'Success rate',
				value: `${100 - failures}%`,
				note: `${failures} failed in the last 90 days`,
			},
			{
				label: 'Avg deploy time',
				value: `${seededInt(`${seedId}:deploy-min`, 1, 4)}m ${seededInt(`${seedId}:deploy-sec`, 0, 59)}s`,
				note: 'From push to live',
			},
			{
				label: 'Last deployment',
				value: deployments[0].timeAgo,
				note: STATUS_META[deployments[0].status].label,
			},
		];
	}, [seedId, deployments]);

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
					<CardAction className="flex items-center gap-2">
						<DemoBadge />
						<Button
							size="sm"
							onClick={() =>
								toast.info(
									`Deploys from the dashboard are coming soon — push to ${branch} to redeploy.`,
								)
							}
						>
							<Rocket className="size-3.5" />
							Deploy now
						</Button>
					</CardAction>
				</CardHeader>
				<div className="divide-y border-t">
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

			<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
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
		</div>
	);
}
