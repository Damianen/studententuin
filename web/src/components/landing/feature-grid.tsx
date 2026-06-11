import { useRef } from 'react';
import { Database, GitBranch, Globe, GraduationCap } from 'lucide-react';
import { gsap, useGSAP, MOTION_OK } from './animation';
import { VineAnchor } from './scroll-vine';

const FEATURES = [
	{
		icon: Globe,
		title: 'Your own subdomain',
		description: 'A real domain that is yours for as long as you study.',
	},
	{
		icon: GitBranch,
		title: 'Git-based deploys',
		description: 'Connect a repository and every push becomes a deploy.',
	},
	{
		icon: Database,
		title: 'Managed databases',
		description: 'PostgreSQL, MySQL or MongoDB — provisioned in one click.',
	},
	{
		icon: GraduationCap,
		title: 'Made for students',
		description: 'Free while you learn. No credit card, no surprises.',
	},
];

export function FeatureGrid() {
	const root = useRef<HTMLElement>(null);

	useGSAP(
		() => {
			const mm = gsap.matchMedia();
			mm.add(MOTION_OK, () => {
				gsap.from('.feature-head', {
					y: 30,
					autoAlpha: 0,
					duration: 0.9,
					ease: 'power3.out',
					scrollTrigger: { trigger: root.current, start: 'top 78%' },
				});
				gsap.from('.feature-card', {
					y: 44,
					autoAlpha: 0,
					duration: 0.85,
					stagger: 0.1,
					ease: 'power3.out',
					scrollTrigger: { trigger: '.feature-row', start: 'top 82%' },
				});
			});
		},
		{ scope: root }
	);

	return (
		<section ref={root} id="features" className="relative scroll-mt-20">
			<VineAnchor x={0.3} />
			<div className="mx-auto w-full max-w-6xl px-6 py-24 sm:py-32">
				<div className="feature-head mb-12 max-w-xl">
					<p className="mb-3 font-mono text-xs uppercase tracking-[0.22em] text-primary">
						03 — What you get
					</p>
					<h2 className="font-display text-4xl font-semibold tracking-tight sm:text-5xl">
						Everything a garden{' '}
						<em className="italic text-primary">needs</em>
					</h2>
				</div>

				<div className="feature-row grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
					{FEATURES.map((feature) => (
						<div
							key={feature.title}
							className="feature-card group rounded-2xl border border-border bg-card/40 p-6 backdrop-blur-sm transition-[translate,border-color,box-shadow] duration-300 hover:-translate-y-1 hover:border-primary/40 hover:shadow-glow"
						>
							<span className="mb-4 flex size-10 items-center justify-center rounded-lg bg-primary/12 text-primary transition-transform duration-300 group-hover:-rotate-6 group-hover:scale-110">
								<feature.icon className="size-5" />
							</span>
							<h3 className="text-sm font-semibold">{feature.title}</h3>
							<p className="mt-1.5 text-sm leading-relaxed text-muted-foreground">
								{feature.description}
							</p>
						</div>
					))}
				</div>
			</div>
		</section>
	);
}
