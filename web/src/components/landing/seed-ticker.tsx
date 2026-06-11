import { Sprout } from 'lucide-react';
import { SUBDOMAIN_SUFFIX } from '@/lib/subdomain_lookup';

const ITEMS = [
	`your-name.${SUBDOMAIN_SUFFIX}`,
	'git push → live',
	'PostgreSQL',
	'MySQL',
	'MongoDB',
	'Node.js',
	'free for students',
	'no credit card',
];

function Row({ hidden }: { hidden?: boolean }) {
	return (
		<div
			aria-hidden={hidden}
			className="flex shrink-0 items-center gap-10 pr-10"
		>
			{ITEMS.map((item) => (
				<span key={item} className="flex items-center gap-10">
					<span className="font-mono text-sm uppercase tracking-[0.18em] text-muted-foreground">
						{item}
					</span>
					<Sprout className="size-3.5 shrink-0 text-primary/60" />
				</span>
			))}
		</div>
	);
}

export function SeedTicker() {
	return (
		<section className="relative z-10 -mx-2 -rotate-[1.2deg] overflow-hidden border-y border-border bg-card/60 py-4 backdrop-blur-sm">
			<div className="marquee-track flex">
				<Row />
				<Row hidden />
			</div>
		</section>
	);
}
