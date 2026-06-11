import { useRef } from 'react';
import { Link } from 'react-router';
import { ArrowRight } from 'lucide-react';
import { useAuth } from '@/contexts/auth_context';
import { gsap, useGSAP, MOTION_OK } from './animation';
import { VineAnchor } from './scroll-vine';

const PETAL = 'M0 0C16 -12 18 -40 0 -54C-18 -40 -16 -12 0 0Z';
const PETAL_COUNT = 6;

/*
 * The payoff: the vine that grew down the whole page ends here, blooming
 * into a flower as the section scrolls into view.
 */
function Bloom() {
	return (
		<div className="relative mx-auto flex h-44 items-center justify-center">
			{/* The vine's last waypoint — dead center of the flower. */}
			<div className="absolute inset-x-0 top-1/2">
				<VineAnchor x={0.5} />
			</div>
			<svg
				aria-hidden
				width="176"
				height="176"
				viewBox="-88 -88 176 176"
				fill="none"
				className="bloom relative"
			>
				{Array.from({ length: PETAL_COUNT }, (_, i) => (
					<g key={`outer-${i}`} transform={`rotate(${i * (360 / PETAL_COUNT)})`}>
						<path
							className="bloom-petal"
							d={PETAL}
							fill={
								i % 2 === 0
									? 'oklch(0.87 0.18 125 / 0.85)'
									: 'oklch(0.74 0.14 138 / 0.9)'
							}
						/>
					</g>
				))}
				{Array.from({ length: PETAL_COUNT }, (_, i) => (
					<g
						key={`inner-${i}`}
						transform={`rotate(${i * (360 / PETAL_COUNT) + 30}) scale(0.55)`}
					>
						<path
							className="bloom-petal"
							d={PETAL}
							fill="oklch(0.93 0.1 115 / 0.9)"
						/>
					</g>
				))}
				<circle className="bloom-heart" r="11" fill="oklch(0.82 0.13 85)" />
				{Array.from({ length: 8 }, (_, i) => {
					const angle = (i / 8) * Math.PI * 2;
					return (
						<circle
							key={`seed-${i}`}
							className="bloom-heart"
							cx={Math.cos(angle) * 5.5}
							cy={Math.sin(angle) * 5.5}
							r="1.4"
							fill="oklch(0.55 0.1 70)"
						/>
					);
				})}
			</svg>
		</div>
	);
}

export function CtaFinale() {
	const root = useRef<HTMLElement>(null);
	const magnetZoneRef = useRef<HTMLDivElement>(null);
	const magnetRef = useRef<HTMLDivElement>(null);
	const { isAuthenticated } = useAuth();

	useGSAP(
		() => {
			const mm = gsap.matchMedia();

			mm.add(MOTION_OK, () => {
				gsap.from('.bloom-petal', {
					scale: 0,
					transformOrigin: '0px 0px',
					stagger: 0.05,
					ease: 'back.out(1.6)',
					scrollTrigger: {
						trigger: root.current,
						start: 'top 78%',
						end: 'top 32%',
						scrub: 0.9,
					},
				});
				gsap.from('.bloom-heart', {
					scale: 0,
					transformOrigin: 'center',
					ease: 'back.out(2)',
					scrollTrigger: {
						trigger: root.current,
						start: 'top 55%',
						end: 'top 32%',
						scrub: 0.9,
					},
				});
				// The whole flower turns, imperceptibly slowly.
				gsap.to('.bloom', {
					rotation: 360,
					duration: 240,
					repeat: -1,
					ease: 'none',
					transformOrigin: 'center',
				});
				gsap.from('.cta-inner > *:not(.bloom-wrap)', {
					y: 32,
					autoAlpha: 0,
					duration: 0.9,
					stagger: 0.1,
					ease: 'power3.out',
					scrollTrigger: { trigger: root.current, start: 'top 55%' },
				});
			});

			// Magnetic button — mouse-only devices, motion allowed.
			mm.add(`${MOTION_OK} and (pointer: fine)`, () => {
				const zone = magnetZoneRef.current;
				const magnet = magnetRef.current;
				if (!zone || !magnet) return;

				const xTo = gsap.quickTo(magnet, 'x', {
					duration: 0.4,
					ease: 'power3.out',
				});
				const yTo = gsap.quickTo(magnet, 'y', {
					duration: 0.4,
					ease: 'power3.out',
				});

				const onMove = (event: PointerEvent) => {
					const rect = zone.getBoundingClientRect();
					xTo((event.clientX - rect.left - rect.width / 2) * 0.35);
					yTo((event.clientY - rect.top - rect.height / 2) * 0.35);
				};
				const onLeave = () => {
					gsap.to(magnet, {
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
		<section ref={root} className="relative overflow-hidden">
			<div
				aria-hidden
				className="pointer-events-none absolute inset-x-0 bottom-0 h-[32rem] bg-[radial-gradient(50rem_24rem_at_50%_100%,oklch(0.45_0.1_135/0.22),transparent)]"
			/>
			<div className="cta-inner relative z-10 mx-auto flex w-full max-w-4xl flex-col items-center px-6 pb-32 pt-32 text-center sm:pb-40 sm:pt-36">
				<div className="bloom-wrap">
					<Bloom />
				</div>
				<p className="mt-6 font-mono text-xs uppercase tracking-[0.22em] text-primary">
					04 — Your turn
				</p>
				<h2 className="mt-5 font-display text-[clamp(2.6rem,7.5vw,5.5rem)] font-semibold leading-[1.04] tracking-tight">
					Ready to plant
					<br />
					<em className="glow-text italic text-primary">something?</em>
				</h2>
				<p className="mt-6 max-w-md text-balance text-muted-foreground">
					Claim your subdomain and have your first project growing before
					your next lecture.
				</p>
				<div ref={magnetZoneRef} className="mt-8 p-6">
					<div ref={magnetRef}>
						<Link
							to={isAuthenticated ? '/projects' : '/register'}
							className="inline-flex h-14 items-center gap-2.5 rounded-full bg-primary px-9 text-base font-semibold text-primary-foreground shadow-glow transition-shadow hover:shadow-lifted"
						>
							{isAuthenticated ? 'Open your garden' : 'Start growing — free'}
							<ArrowRight className="size-4.5" />
						</Link>
					</div>
				</div>
				<p className="font-mono text-[11px] uppercase tracking-[0.22em] text-muted-foreground">
					free for students · no credit card
				</p>
			</div>
		</section>
	);
}
