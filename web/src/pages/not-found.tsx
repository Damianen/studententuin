import { Link } from 'react-router';
import { Sprout } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useAuth } from '@/contexts/auth_context';

export default function NotFound() {
	const { isAuthenticated } = useAuth();

	return (
		<div className="relative flex flex-1 flex-col items-center justify-center gap-6 px-6 py-24 text-center">
			{/* Ghost number in the soil, like the landing page's step numerals. */}
			<span
				aria-hidden
				className="pointer-events-none absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 select-none font-display text-[clamp(12rem,40vw,24rem)] font-semibold leading-none text-foreground/[0.045]"
			>
				404
			</span>
			<span className="relative flex size-14 items-center justify-center rounded-full border border-primary/40 bg-primary/10 text-primary">
				<Sprout className="size-6" />
			</span>
			<div className="relative space-y-3">
				<p className="font-mono text-xs uppercase tracking-[0.22em] text-primary">
					404 — fallow ground
				</p>
				<h1 className="font-display text-4xl font-semibold tracking-tight sm:text-5xl">
					Nothing growing{' '}
					<em className="glow-text italic text-primary">here</em>
				</h1>
				<p className="text-muted-foreground">
					The page you are looking for does not exist.
				</p>
			</div>
			<Button asChild className="relative">
				<Link to={isAuthenticated ? '/projects' : '/'}>
					{isAuthenticated ? 'Back to your garden' : 'Back home'}
				</Link>
			</Button>
		</div>
	);
}
