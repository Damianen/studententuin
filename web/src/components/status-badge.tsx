import { cn } from '@/lib/utils';

const STATUS_DOTS: Record<string, string> = {
	running: 'bg-status-running',
	provisioning: 'bg-status-pending animate-pulse',
	pending: 'bg-status-pending animate-pulse',
	stopped: 'bg-status-stopped',
	failed: 'bg-status-failed',
};

export function StatusBadge({
	status,
	className,
}: {
	status: string;
	className?: string;
}) {
	const dot = STATUS_DOTS[status.toLowerCase()] ?? 'bg-status-stopped';

	return (
		<span
			className={cn(
				'inline-flex items-center gap-1.5 rounded-full border border-border bg-card px-2.5 py-0.5 text-xs font-medium capitalize',
				className
			)}
		>
			<span className={cn('size-1.5 rounded-full', dot)} />
			{status.toLowerCase()}
		</span>
	);
}
