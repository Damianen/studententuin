const MINUTE = 60;
const HOUR = 3600;
const DAY = 86400;

// Coarse buckets on purpose — the deployments list refreshes on every
// deploy, not on a timer, so second-level precision would just go stale.
export function formatTimeAgo(iso: string): string {
	const seconds = Math.floor((Date.now() - new Date(iso).getTime()) / 1000);
	if (seconds < MINUTE) return 'just now';
	if (seconds < HOUR) {
		const minutes = Math.floor(seconds / MINUTE);
		return minutes === 1 ? '1 minute ago' : `${minutes} minutes ago`;
	}
	if (seconds < DAY) {
		const hours = Math.floor(seconds / HOUR);
		return hours === 1 ? '1 hour ago' : `${hours} hours ago`;
	}
	const days = Math.floor(seconds / DAY);
	if (days === 1) return 'yesterday';
	if (days < 7) return `${days} days ago`;
	// Beyond a week an absolute date reads better than "53 days ago".
	return new Date(iso).toLocaleDateString([], {
		day: 'numeric',
		month: 'short',
		year: 'numeric',
	});
}

export function formatDuration(seconds: number): string {
	const whole = Math.round(seconds);
	if (whole < MINUTE) return `${whole}s`;
	return `${Math.floor(whole / MINUTE)}m ${String(whole % MINUTE).padStart(2, '0')}s`;
}
