import { useRef } from 'react';
import { Link } from 'react-router';
import { ArrowRight } from 'lucide-react';
import { useAuth } from '@/contexts/auth_context';
import { LogoMark } from '@/components/logo';
import { Button } from '@/components/ui/button';
import { gsap, useGSAP, ScrollTrigger, MOTION_OK } from './animation';

const LINKS = [
	{ label: 'Grow', href: '#grow' },
	{ label: 'Deploy', href: '#deploy' },
	{ label: 'Features', href: '#features' },
];

/*
 * Floating glass capsule instead of a classic top bar. The chartreuse
 * line along its base fills as you scroll — the garden growing.
 */
export function LandingNav() {
	const { isAuthenticated } = useAuth();
	const capsuleRef = useRef<HTMLElement>(null);
	const growLineRef = useRef<HTMLSpanElement>(null);

	useGSAP(
		() => {
			const mm = gsap.matchMedia();
			mm.add(MOTION_OK, () => {
				gsap.from(capsuleRef.current, {
					y: -28,
					autoAlpha: 0,
					duration: 0.9,
					ease: 'power3.out',
					delay: 0.4,
				});
				gsap.fromTo(
					growLineRef.current,
					{ scaleX: 0 },
					{
						scaleX: 1,
						ease: 'none',
						scrollTrigger: { start: 0, end: 'max', scrub: 0.4 },
					}
				);
				// Slip away while scrolling down, return on the first scroll
				// up — keeps the finale bloom unobstructed. The 90px ride
				// clears the viewport entirely, so no fade is needed.
				const yTo = gsap.quickTo(capsuleRef.current, 'y', {
					duration: 0.55,
					ease: 'power3.out',
				});
				ScrollTrigger.create({
					start: 160,
					end: 'max',
					onUpdate: (self) => yTo(self.direction === 1 ? -90 : 0),
					onLeaveBack: () => yTo(0),
				});
			});
		},
		{ scope: capsuleRef }
	);

	const scrollTo = (href: string) => {
		gsap.to(window, {
			scrollTo: { y: href, offsetY: 88 },
			duration: 1.1,
			ease: 'power3.inOut',
		});
	};

	return (
		<header className="fixed inset-x-0 top-4 z-50 flex justify-center px-4">
			<nav
				ref={capsuleRef}
				aria-label="Main"
				className="relative flex items-center gap-1 overflow-hidden rounded-full border border-border bg-[oklch(0.19_0.023_152/0.72)] py-1.5 pl-2 pr-1.5 shadow-lifted backdrop-blur-xl"
			>
				<Link
					to="/"
					aria-label="Studententuin home"
					className="flex items-center gap-2.5 rounded-full px-2 py-1 transition-colors hover:bg-foreground/5"
				>
					<LogoMark className="size-7" />
					<span className="hidden font-display text-base font-semibold tracking-tight lg:inline">
						studententuin
					</span>
				</Link>

				<span aria-hidden className="mx-1 hidden h-5 w-px bg-border md:inline-block" />

				<div className="hidden items-center md:flex">
					{LINKS.map((link) => (
						<a
							key={link.href}
							href={link.href}
							onClick={(event) => {
								event.preventDefault();
								scrollTo(link.href);
							}}
							className="rounded-full px-3.5 py-2 font-mono text-[10px] uppercase tracking-[0.2em] text-muted-foreground transition-colors hover:bg-foreground/5 hover:text-primary"
						>
							{link.label}
						</a>
					))}
				</div>

				<span aria-hidden className="mx-1 hidden h-5 w-px bg-border md:inline-block" />

				{isAuthenticated ? (
					<Button size="sm" className="rounded-full" asChild>
						<Link to="/projects">
							Your garden
							<ArrowRight className="size-3.5" />
						</Link>
					</Button>
				) : (
					<div className="flex items-center gap-1">
						<Button size="sm" variant="ghost" className="rounded-full" asChild>
							<Link to="/login">Sign in</Link>
						</Button>
						<Button size="sm" className="rounded-full" asChild>
							<Link to="/register">Claim your plot</Link>
						</Button>
					</div>
				)}

				{/* Scroll progress — the capsule's base slowly grows green. */}
				<span
					ref={growLineRef}
					aria-hidden
					className="absolute inset-x-6 bottom-0 h-[2px] origin-left rounded-full bg-linear-to-r from-primary/10 via-primary/60 to-primary"
				/>
			</nav>
		</header>
	);
}
