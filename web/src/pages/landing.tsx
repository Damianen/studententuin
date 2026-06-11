import { lazy, Suspense, useRef } from 'react';
import { useGSAP, ScrollTrigger } from '@/components/landing/animation';
import { LandingNav } from '@/components/landing/landing-nav';
import { LandingFooter } from '@/components/landing/landing-footer';
import { ScrollVine } from '@/components/landing/scroll-vine';
import { Hero } from '@/components/landing/hero';
import { SeedTicker } from '@/components/landing/seed-ticker';
import { GrowthSteps } from '@/components/landing/growth-steps';
import { DeployDemo } from '@/components/landing/deploy-demo';
import { FeatureGrid } from '@/components/landing/feature-grid';
import { CtaFinale } from '@/components/landing/cta-finale';

const FireflyCanvas = lazy(
	() => import('@/components/landing/firefly-canvas')
);

export default function Landing() {
	const pageRef = useRef<HTMLDivElement>(null);

	useGSAP(() => {
		// Fraunces swaps in late; re-measure scroll positions once it lands.
		let cancelled = false;
		document.fonts.ready.then(() => {
			if (!cancelled) ScrollTrigger.refresh();
		});
		return () => {
			cancelled = true;
		};
	});

	return (
		<div className="relative min-h-svh overflow-x-clip bg-background text-foreground">
			<Suspense fallback={null}>
				<FireflyCanvas />
			</Suspense>
			<div aria-hidden className="grain-overlay" />
			<LandingNav />

			<div ref={pageRef} className="relative">
				<ScrollVine container={pageRef} />
				<main className="relative z-10">
					<Hero />
					<SeedTicker />
					<GrowthSteps />
					<DeployDemo />
					<FeatureGrid />
					<CtaFinale />
				</main>
				<LandingFooter />
			</div>
		</div>
	);
}
