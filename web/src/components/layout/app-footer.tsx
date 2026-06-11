import { LogoMark } from '@/components/logo';

export function AppFooter() {
	return (
		<footer className="relative z-10 border-t border-border">
			<div className="mx-auto flex w-full max-w-6xl flex-wrap items-center justify-between gap-4 px-6 py-10 text-sm text-muted-foreground">
				<div className="flex items-center gap-2.5">
					<LogoMark className="size-6" />
					<span className="font-display text-base text-foreground">
						studententuin
					</span>
				</div>
				<span className="font-mono text-xs">
					© {new Date().getFullYear()} — grown with care, for students
				</span>
			</div>
		</footer>
	);
}
