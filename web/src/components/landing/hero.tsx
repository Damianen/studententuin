import { useRef } from 'react';
import { Link } from 'react-router';
import { Sprout } from 'lucide-react';
import { useAuth } from '@/contexts/auth_context';
import { SUBDOMAIN_SUFFIX } from '@/lib/subdomain_lookup';
import { gsap, useGSAP, SplitText, MOTION_OK } from './animation';
import { VineAnchor } from './scroll-vine';

/*
 * Hero as a botanical specimen sheet: terraced display type stepping
 * down-right, a mono specimen-label rail, field coordinates in the
 * corners — and the CTA is the seed the vine grows from.
 */

const SPECIMEN = [
	{ label: 'Specimen', value: 'student code, wild variety' },
	{ label: 'Soil', value: `your-name.${SUBDOMAIN_SUFFIX}` },
	{ label: 'Season', value: 'free while you study' },
	{ label: 'Yield', value: 'apps · databases · ideas' },
];

export function Hero() {
	const { isAuthenticated } = useAuth();
	const root = useRef<HTMLElement>(null);
	const seedZoneRef = useRef<HTMLDivElement>(null);
	const seedRef = useRef<HTMLDivElement>(null);

	useGSAP(
		() => {
			const mm = gsap.matchMedia();
			mm.add(MOTION_OK, () => {
				// Only the cream lines are split into chars. The glowing
				// italic line must stay ONE element — per-char text-shadows
				// render as visible boxes around the letters.
				const split = SplitText.create('.hero-line', {
					type: 'chars,words',
					autoSplit: true,
					onSplit: (self) =>
						gsap.from(self.chars, {
							yPercent: 80,
							autoAlpha: 0,
							rotation: () => gsap.utils.random(-6, 6),
							duration: 1.1,
							stagger: { each: 0.03, from: 'start' },
							ease: 'power4.out',
							delay: 0.2,
						}),
				});
				gsap.from('.hero-title-glow', {
					y: 56,
					autoAlpha: 0,
					scale: 0.96,
					transformOrigin: '0% 100%',
					duration: 1.2,
					ease: 'power3.out',
					delay: 0.6,
				});
				gsap.from('.hero-fade', {
					y: 24,
					autoAlpha: 0,
					duration: 0.9,
					stagger: 0.1,
					ease: 'power3.out',
					delay: 0.8,
				});
				gsap.from('.spec-rule', {
					scaleY: 0,
					transformOrigin: 'top',
					duration: 1.1,
					ease: 'power3.inOut',
					delay: 0.7,
				});
				gsap.from('.spec-item', {
					x: 26,
					autoAlpha: 0,
					duration: 0.8,
					stagger: 0.12,
					ease: 'power3.out',
					delay: 0.9,
				});
				gsap.from('.hero-seed', {
					scale: 0.5,
					autoAlpha: 0,
					duration: 1,
					ease: 'back.out(1.8)',
					delay: 1.05,
				});
				gsap.to('.seed-ring', {
					scale: 1.22,
					opacity: 0,
					duration: 2.2,
					stagger: 0.55,
					repeat: -1,
					ease: 'sine.out',
				});
				return () => split.revert();
			});

			// Magnetic seed — mouse-only devices, motion allowed.
			mm.add(`${MOTION_OK} and (pointer: fine)`, () => {
				const zone = seedZoneRef.current;
				const seed = seedRef.current;
				if (!zone || !seed) return;

				const xTo = gsap.quickTo(seed, 'x', {
					duration: 0.45,
					ease: 'power3.out',
				});
				const yTo = gsap.quickTo(seed, 'y', {
					duration: 0.45,
					ease: 'power3.out',
				});
				const onMove = (event: PointerEvent) => {
					const rect = zone.getBoundingClientRect();
					xTo((event.clientX - rect.left - rect.width / 2) * 0.3);
					yTo((event.clientY - rect.top - rect.height / 2) * 0.3);
				};
				const onLeave = () => {
					gsap.to(seed, {
						x: 0,
						y: 0,
						duration: 0.7,
						ease: 'elastic.out(1, 0.4)',
					});
				};
				zone.addEventListener('pointermove', onMove);
				zone.addEventListener('pointerleave', onLeave);
				return () => {
					zone.removeEventListener('pointermove', onMove);
					zone.removeEventListener('pointerleave', onLeave);
				};
			});
		},
		{ scope: root }
	);

	return (
		<section
			ref={root}
			className="relative flex min-h-svh flex-col px-6 pt-20 sm:px-10 lg:px-16 xl:pt-24"
		>
			{/* Moonlight pooling at the top, a faint glow where the seed waits. */}
			<div
				aria-hidden
				className="pointer-events-none absolute inset-0 bg-[radial-gradient(70rem_34rem_at_30%_-14%,oklch(0.38_0.07_145/0.4),transparent)]"
			/>
			<div
				aria-hidden
				className="pointer-events-none absolute inset-x-0 bottom-0 h-64 bg-[radial-gradient(36rem_14rem_at_50%_100%,oklch(0.6_0.14_130/0.16),transparent)]"
			/>

			<div className="relative z-10 mx-auto flex w-full max-w-7xl flex-1 flex-col">
				<div className="flex flex-1 items-center py-6">
					<div className="grid w-full items-center gap-12 lg:grid-cols-[1fr_auto]">
						{/* Terraced headline, stepping down-right like garden beds. */}
						<div>
							<p className="hero-fade mb-6 flex items-center gap-2.5 font-mono text-[11px] uppercase tracking-[0.25em] text-muted-foreground">
								<span className="size-1.5 animate-pulse rounded-full bg-primary" />
								Free hosting for students
							</p>
							<h1 className="font-display font-semibold leading-[0.92] tracking-tight">
								<span className="hero-line block text-[clamp(3rem,min(10.5vw,15svh),8.5rem)]">
									Plant
								</span>
								<span className="hero-line ml-[clamp(2.5rem,9vw,8rem)] block text-[clamp(3rem,min(10.5vw,15svh),8.5rem)]">
									an idea.
								</span>
								<em className="hero-title-glow glow-text ml-[clamp(1rem,4vw,3.5rem)] mt-3 block italic text-[clamp(2.2rem,min(6.5vw,9.5svh),5.25rem)] text-primary">
									Grow the web.
								</em>
							</h1>
							<p className="hero-fade ml-[clamp(1rem,4vw,3.5rem)] mt-6 max-w-sm text-balance text-muted-foreground">
								A living patch of internet for every student — your
								subdomain, your apps, your databases. No bills, no ops.
							</p>
						</div>

						{/* Specimen label — the botanical record card. */}
						<aside className="relative hidden gap-7 self-center pl-8 lg:flex lg:flex-col">
							<span
								aria-hidden
								className="spec-rule absolute inset-y-1 left-0 w-px bg-linear-to-b from-primary/60 via-border to-transparent"
							/>
							<p className="spec-item flex items-center gap-2 font-mono text-[10px] uppercase tracking-[0.25em] text-primary">
								<Sprout className="size-3.5" />
								Field notes
							</p>
							{SPECIMEN.map((row) => (
								<div key={row.label} className="spec-item">
									<dt className="font-mono text-[10px] uppercase tracking-[0.25em] text-muted-foreground">
										{row.label}
									</dt>
									<dd className="mt-1.5 font-mono text-sm text-foreground/90">
										{row.value}
									</dd>
								</div>
							))}
						</aside>
					</div>
				</div>

				{/* Bottom band: coordinates · the seed · the invitation down. */}
				<div className="grid items-end gap-6 pb-6 sm:grid-cols-3">
					<p className="hero-fade hidden font-mono text-[10px] uppercase tracking-[0.25em] text-muted-foreground/70 sm:block">
						52.2°N 5.5°E — grown in NL
					</p>

					<div className="flex flex-col items-center">
						<div ref={seedZoneRef} className="p-4">
							<div ref={seedRef}>
								<Link
									to={isAuthenticated ? '/projects' : '/register'}
									className="hero-seed group relative flex size-28 items-center justify-center rounded-full border border-primary/40 bg-primary/10 backdrop-blur-md transition-colors duration-300 hover:bg-primary/20 sm:size-32 lg:size-36"
								>
									<span
										aria-hidden
										className="seed-ring absolute inset-0 rounded-full border border-primary/35"
									/>
									<span
										aria-hidden
										className="seed-ring absolute inset-0 rounded-full border border-primary/25"
									/>
									<span className="px-5 text-center font-mono text-[11px] uppercase leading-relaxed tracking-[0.2em] text-primary">
										{isAuthenticated ? 'Your garden →' : 'Claim your plot →'}
									</span>
								</Link>
							</div>
						</div>
						{!isAuthenticated && (
							<Link
								to="/login"
								className="hero-fade font-mono text-xs text-muted-foreground underline-offset-4 transition-colors hover:text-primary hover:underline"
							>
								or sign in
							</Link>
						)}
					</div>

					<p className="hero-fade hidden text-right font-mono text-[10px] uppercase tracking-[0.25em] text-muted-foreground/70 sm:block">
						Follow the vine ↓
					</p>
				</div>
			</div>

			{/* The seed — where the vine breaks ground. */}
			<div className="absolute inset-x-0 bottom-0">
				<VineAnchor x={0.5} />
			</div>
		</section>
	);
}
