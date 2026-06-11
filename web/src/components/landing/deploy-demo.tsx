import { useRef } from 'react';
import { Check } from 'lucide-react';
import { gsap, useGSAP, MOTION_OK } from './animation';
import { VineAnchor } from './scroll-vine';

const CHECKS = [
	'Automatic builds from your repository',
	'PostgreSQL, MySQL or MongoDB on the side',
	'Live on your own subdomain in minutes',
];

export function DeployDemo() {
	const root = useRef<HTMLElement>(null);

	useGSAP(
		() => {
			const mm = gsap.matchMedia();
			mm.add(MOTION_OK, () => {
				gsap.from('.demo-copy > *', {
					y: 30,
					autoAlpha: 0,
					duration: 0.9,
					stagger: 0.1,
					ease: 'power3.out',
					scrollTrigger: { trigger: root.current, start: 'top 72%' },
				});
				gsap.from('.demo-terminal', {
					y: 44,
					autoAlpha: 0,
					rotate: 1.5,
					duration: 1,
					ease: 'power3.out',
					scrollTrigger: { trigger: root.current, start: 'top 72%' },
				});
				gsap.from('.demo-line', {
					autoAlpha: 0,
					y: 8,
					duration: 0.4,
					stagger: 0.5,
					ease: 'power2.out',
					scrollTrigger: { trigger: '.demo-terminal', start: 'top 65%' },
				});
				gsap.from('.demo-live', {
					scale: 0.85,
					duration: 0.5,
					ease: 'back.out(2.5)',
					delay: 3,
					scrollTrigger: { trigger: '.demo-terminal', start: 'top 65%' },
				});
			});
		},
		{ scope: root }
	);

	return (
		<section ref={root} id="deploy" className="relative scroll-mt-20">
			<VineAnchor x={0.68} />
			<div className="mx-auto w-full max-w-6xl px-6 py-28 sm:py-36">
				<div className="grid items-center gap-12 lg:grid-cols-2 lg:gap-16">
					<div className="demo-copy max-w-lg">
						<p className="mb-3 font-mono text-xs uppercase tracking-[0.22em] text-primary">
							02 — From push to planted
						</p>
						<h2 className="font-display text-4xl font-semibold tracking-tight sm:text-5xl">
							Deploys that feel like{' '}
							<em className="italic text-primary">gardening</em>, not ops
						</h2>
						<p className="mt-5 text-muted-foreground">
							Push your branch and we take it from there — building your
							app, provisioning your database, and serving it on your own
							domain.
						</p>
						<ul className="mt-7 space-y-3">
							{CHECKS.map((check) => (
								<li key={check} className="flex items-start gap-3 text-sm">
									<span className="mt-0.5 flex size-5 shrink-0 items-center justify-center rounded-full bg-primary/15 text-primary">
										<Check className="size-3" />
									</span>
									{check}
								</li>
							))}
						</ul>
					</div>

					{/* Glass terminal — the vine glows through the blur behind it. */}
					<div className="demo-terminal overflow-hidden rounded-2xl border border-border bg-[oklch(0.13_0.018_160/0.75)] shadow-lifted backdrop-blur-md">
						<div className="flex items-center gap-1.5 border-b border-border px-4 py-3">
							<span className="size-2.5 rounded-full bg-[oklch(0.6_0.16_25)]" />
							<span className="size-2.5 rounded-full bg-[oklch(0.78_0.13_85)]" />
							<span className="size-2.5 rounded-full bg-primary/70" />
							<span className="ml-3 font-mono text-xs text-muted-foreground">
								studententuin — deploy
							</span>
						</div>
						<div className="space-y-2.5 px-5 py-6 font-mono text-[13px] leading-relaxed">
							<p className="demo-line text-foreground/90">
								<span className="text-primary">$</span> git push tuin main
							</p>
							<p className="demo-line text-muted-foreground">
								→ building mijn-app…
							</p>
							<p className="demo-line text-primary/90">
								✓ build complete · 12s
							</p>
							<p className="demo-line text-muted-foreground">
								→ provisioning postgres 16…
							</p>
							<p className="demo-line text-primary/90">✓ database ready</p>
							<p className="demo-line demo-live mt-4 inline-flex items-center gap-2 rounded-lg border border-primary/25 bg-primary/10 px-3 py-2 text-foreground">
								<span className="size-1.5 animate-pulse rounded-full bg-primary" />
								live at{' '}
								<span className="text-primary">
									mijntuin.studententuin.com
								</span>
							</p>
							<p className="demo-line text-foreground/80">
								<span className="text-primary">$</span>{' '}
								<span className="terminal-cursor inline-block h-3.5 w-1.5 translate-y-0.5 bg-foreground/80" />
							</p>
						</div>
					</div>
				</div>
			</div>
		</section>
	);
}
