import { cn } from '@/lib/utils';

export const LogoMark = ({ className }: { className?: string }) => (
	<svg
		viewBox="0 0 32 32"
		fill="none"
		xmlns="http://www.w3.org/2000/svg"
		className={cn('size-8', className)}
	>
		<rect width="32" height="32" rx="9" className="fill-primary" />
		<path
			d="M16 24c0-5 0-9 0-11"
			className="stroke-primary-foreground"
			strokeWidth="2"
			strokeLinecap="round"
		/>
		<path
			d="M16 15c0-4 2.5-6.5 6-7 0 4-2.5 6.5-6 7z"
			className="fill-accent"
		/>
		<path
			d="M16 18c0-3.5-2.2-5.7-5.5-6.2 0 3.5 2.2 5.7 5.5 6.2z"
			className="fill-primary-foreground"
		/>
	</svg>
);

export const Logo = ({ className }: { className?: string }) => {
	return (
		<div className={cn('flex items-center gap-2.5', className)}>
			<LogoMark />
			<span className="font-display text-xl font-semibold tracking-tight text-foreground">
				studententuin
			</span>
		</div>
	);
};
