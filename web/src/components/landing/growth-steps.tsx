import { Fragment, useRef } from 'react';
import type { ReactNode } from 'react';
import { Sprout } from 'lucide-react';
import { gsap, useGSAP, MOTION_OK } from './animation';
import { SUBDOMAIN_SUFFIX } from '@/lib/subdomain_lookup';
import { VineAnchor } from './scroll-vine';

interface Step {
	season: string;
	title: string;
	description: string;
	detail: ReactNode;
}

const STEPS: Step[] = [
	{
		season: 'Sow',
		title: 'Claim your plot',
		description:
			'Pick a name and the plot is yours — a real subdomain, free for as long as you study.',
		detail: (
			<code className="inline-flex items-center gap-2 rounded-full border border-primary/25 bg-primary/10 px-4 py-2 font-mono text-xs text-primary">
				<Sprout className="size-3.5" />
				your-name.{SUBDOMAIN_SUFFIX}
			</code>
		),
	},
	{
		season: 'Tend',
		title: 'Plant something',
		description:
			'Point us at a Git repository or provision a managed database — PostgreSQL, MySQL or MongoDB. We build, water, and keep the lights on.',
		detail: (
			<code className="inline-block rounded-lg border border-border bg-background/70 px-4 py-2 font-mono text-xs text-foreground/90">
				<span className="text-primary">$</span> git push tuin main
			</code>
		),
	},
	{
		season: 'Harvest',
		title: 'Watch it grow',
		description:
			'Your project goes live on your own domain. Update it, share it, or plant the next idea right beside it.',
		detail: (
			<span className="inline-flex items-center gap-2 rounded-full border border-status-running/30 bg-status-running/10 px-4 py-2 font-mono text-xs text-status-running">
				<span className="size-1.5 animate-pulse rounded-full bg-status-running" />
				live — and growing
			</span>
		),
	},
];

export function GrowthSteps() {
	const root = useRef<HTMLElement>(null);

	useGSAP(
		() => {
			const mm = gsap.matchMedia();
			mm.add(MOTION_OK, () => {
				gsap.from('.step-head', {
					y: 32,
					autoAlpha: 0,
					duration: 0.9,
					ease: 'power3.out',
					scrollTrigger: { trigger: root.current, start: 'top 75%' },
				});
				gsap.utils
					.toArray<HTMLElement>('.step-card', root.current!)
					.forEach((card, index) => {
						gsap.from(card, {
							y: 64,
							autoAlpha: 0,
							rotate: index % 2 === 0 ? -1.5 : 1.5,
							duration: 1,
							ease: 'power3.out',
							scrollTrigger: { trigger: card, start: 'top 80%' },
						});
					});
				gsap.utils
					.toArray<HTMLElement>('.step-num', root.current!)
					.forEach((num) => {
						gsap.to(num, {
							yPercent: -30,
							ease: 'none',
							scrollTrigger: {
								trigger: num,
								start: 'top bottom',
								end: 'bottom top',
								scrub: 1,
							},
						});
					});
			});
		},
		{ scope: root }
	);

	return (
		<section ref={root} id="grow" className="relative scroll-mt-20">
			<div className="mx-auto w-full max-w-6xl px-6 py-28 sm:py-36">
				<div className="step-head max-w-xl">
					<p className="mb-3 font-mono text-xs uppercase tracking-[0.22em] text-primary">
						01 — How it grows
					</p>
					<h2 className="font-display text-4xl font-semibold tracking-tight sm:text-5xl">
						Three seasons from seed to{' '}
						<em className="italic text-primary">shipped</em>
					</h2>
				</div>

				<div className="mt-12 space-y-28 sm:space-y-40">
					{STEPS.map((step, index) => {
						const onLeft = index % 2 === 0;
						return (
							<Fragment key={step.season}>
								{/* The vine swings to the empty side of each row. */}
								<VineAnchor x={onLeft ? 0.76 : 0.24} />
								<div
									className={`relative flex ${
										onLeft ? 'justify-start' : 'justify-end'
									}`}
								>
									<span
										aria-hidden
										className={`step-num pointer-events-none absolute -top-20 select-none font-display text-[clamp(9rem,22vw,17rem)] font-semibold leading-none text-foreground/[0.045] ${
											onLeft ? '-right-2 sm:right-8' : '-left-2 sm:left-8'
										}`}
									>
										0{index + 1}
									</span>
									<div className="step-card relative w-full max-w-lg rounded-3xl border border-border bg-card/50 p-8 shadow-lifted backdrop-blur-md sm:p-10">
										<p className="font-mono text-xs uppercase tracking-[0.22em] text-primary">
											0{index + 1} — {step.season}
										</p>
										<h3 className="mt-3 font-display text-3xl font-semibold tracking-tight">
											{step.title}
										</h3>
										<p className="mt-3 leading-relaxed text-muted-foreground">
											{step.description}
										</p>
										<div className="mt-7">{step.detail}</div>
									</div>
								</div>
							</Fragment>
						);
					})}
				</div>
			</div>
		</section>
	);
}
