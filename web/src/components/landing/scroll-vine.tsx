import { useEffect, useRef, useState } from 'react';
import type { RefObject } from 'react';
import { gsap, useGSAP, ScrollTrigger, MOTION_OK } from './animation';

/*
 * The bioluminescent vine that grows down the whole landing page as you
 * scroll. Sections drop <VineAnchor x={0..1}> markers; after layout we
 * measure them, thread a smooth bezier through the waypoints, scatter
 * leaves along the path and scrub the whole thing with one ScrollTrigger.
 */

interface Leaf {
	x: number;
	y: number;
	angle: number;
	size: number;
	flip: boolean;
	fill: string;
	/** position along the stem as a fraction of its total length */
	frac: number;
}

interface VineLayout {
	width: number;
	height: number;
	d: string;
	leaves: Leaf[];
}

const LEAF_BLADE = 'M0 0C9 -3 15 -13 13 -27C3.5 -23 -4.5 -12 0 0Z';
const LEAF_RIB = 'M0 0C4.5 -7 8 -14 11.5 -23';
const LEAF_FILLS = [
	'oklch(0.62 0.11 148)',
	'oklch(0.72 0.14 135)',
	'oklch(0.55 0.09 152)',
];
const BERRY_FILL = 'oklch(0.78 0.12 85)';

// Deterministic jitter so resize rebuilds keep the same plant.
function jitter(seed: number): number {
	const value = Math.sin(seed * 127.1 + 311.7) * 43758.5453;
	return value - Math.floor(value);
}

function buildPath(points: { x: number; y: number }[]): string {
	let d = `M ${points[0].x.toFixed(1)} ${points[0].y.toFixed(1)}`;
	for (let i = 1; i < points.length; i++) {
		const from = points[i - 1];
		const to = points[i];
		// Vertical control handles keep every join tangent-vertical, so the
		// stem winds down the page in smooth S-curves whatever the anchors do.
		const handle = (to.y - from.y) * 0.42;
		d +=
			` C ${from.x.toFixed(1)} ${(from.y + handle).toFixed(1)},` +
			` ${to.x.toFixed(1)} ${(to.y - handle).toFixed(1)},` +
			` ${to.x.toFixed(1)} ${to.y.toFixed(1)}`;
	}
	return d;
}

function measureVine(container: HTMLElement): VineLayout | null {
	const rect = container.getBoundingClientRect();
	const width = rect.width;
	const height = container.scrollHeight;
	if (width < 10 || height < 10) return null;

	const anchors = Array.from(
		container.querySelectorAll<HTMLElement>('[data-vine]')
	);
	if (anchors.length < 2) return null;

	const points = anchors.map((anchor) => {
		const anchorRect = anchor.getBoundingClientRect();
		const fraction = Number.parseFloat(anchor.dataset.vine ?? '0.5');
		// Keep the stem on-screen on narrow viewports.
		const x = gsap.utils.clamp(24, width - 24, width * fraction);
		return { x, y: anchorRect.top - rect.top };
	});
	const d = buildPath(points);

	// Probe path in a hidden in-DOM svg — getTotalLength on detached
	// elements is unreliable across browsers.
	const svgNS = 'http://www.w3.org/2000/svg';
	const probeSvg = document.createElementNS(svgNS, 'svg');
	probeSvg.setAttribute('aria-hidden', 'true');
	probeSvg.style.cssText =
		'position:absolute;width:0;height:0;overflow:hidden;visibility:hidden';
	const probe = document.createElementNS(svgNS, 'path');
	probe.setAttribute('d', d);
	probeSvg.appendChild(probe);
	container.appendChild(probeSvg);

	const total = probe.getTotalLength();
	const leaves: Leaf[] = [];
	let distance = 140;
	let index = 0;
	while (distance < total - 160) {
		const point = probe.getPointAtLength(distance);
		const ahead = probe.getPointAtLength(Math.min(distance + 1, total));
		const tangent =
			(Math.atan2(ahead.y - point.y, ahead.x - point.x) * 180) / Math.PI;
		const flip = index % 2 === 0;
		const lean = (jitter(index) - 0.5) * 26;
		leaves.push({
			x: point.x,
			y: point.y,
			// Leaves sprout up-and-outward, alternating sides, biased by the
			// stem's local direction (tangent is ~90° when growing straight down).
			angle: (flip ? 60 : -60) + (tangent - 90) + lean,
			size: 1.05 + jitter(index + 57) * 0.65,
			flip,
			fill:
				index % 7 === 3
					? BERRY_FILL
					: LEAF_FILLS[index % LEAF_FILLS.length],
			frac: distance / total,
		});
		distance += 115 + jitter(index + 13) * 95;
		index++;
	}

	probeSvg.remove();
	return { width, height, d, leaves };
}

export function VineAnchor({ x }: { x: number }) {
	return (
		<div
			data-vine={x}
			aria-hidden
			className="pointer-events-none relative h-0 w-full"
		/>
	);
}

export function ScrollVine({
	container,
}: {
	container: RefObject<HTMLDivElement | null>;
}) {
	const svgRef = useRef<SVGSVGElement>(null);
	const stemRef = useRef<SVGPathElement>(null);
	const tipRef = useRef<SVGGElement>(null);
	const tipHaloRef = useRef<SVGCircleElement>(null);
	const [layout, setLayout] = useState<VineLayout | null>(null);

	// Passive effect on purpose: the parent's ref is attached bottom-up, so a
	// layout effect here would run before `container.current` exists.
	useEffect(() => {
		const el = container.current;
		if (!el) return;
		let frame = 0;
		const measure = () => {
			const next = measureVine(el);
			if (!next) return;
			setLayout((prev) =>
				prev && prev.d === next.d && prev.height === next.height
					? prev
					: next
			);
		};
		measure();
		const observer = new ResizeObserver(() => {
			cancelAnimationFrame(frame);
			frame = requestAnimationFrame(measure);
		});
		observer.observe(el);
		return () => {
			cancelAnimationFrame(frame);
			observer.disconnect();
		};
	}, [container]);

	useGSAP(
		() => {
			const stem = stemRef.current;
			const tip = tipRef.current;
			const svg = svgRef.current;
			if (!layout || !stem || !tip || !svg) return;

			const mm = gsap.matchMedia();
			mm.add(MOTION_OK, () => {
				const total = stem.getTotalLength();

				// y → arc-length lookup. The stem only ever descends, so y is
				// monotonic and we can binary-search it. This lets the drawn
				// tip track the viewport instead of raw scroll fraction —
				// otherwise the vine drifts off screen on long pages.
				const SAMPLES = 240;
				const lookup: { len: number; y: number }[] = [];
				for (let i = 0; i <= SAMPLES; i++) {
					const len = (i / SAMPLES) * total;
					lookup.push({ len, y: stem.getPointAtLength(len).y });
				}
				const fracForY = (y: number) => {
					if (y <= lookup[0].y) return 0;
					if (y >= lookup[SAMPLES].y) return 1;
					let lo = 0;
					let hi = SAMPLES;
					while (lo + 1 < hi) {
						const mid = (lo + hi) >> 1;
						if (lookup[mid].y <= y) lo = mid;
						else hi = mid;
					}
					const a = lookup[lo];
					const b = lookup[hi];
					const t = (y - a.y) / Math.max(b.y - a.y, 0.001);
					return (a.len + t * (b.len - a.len)) / total;
				};

				const timeline = gsap.timeline({
					paused: true,
					defaults: { ease: 'none' },
				});

				timeline.from('.vine-draw', { drawSVG: 0, duration: 1 }, 0);

				gsap.utils
					.toArray<SVGGElement>('.vine-leaf', svg)
					.forEach((leaf) => {
						const frac = Number(leaf.dataset.frac ?? 0);
						timeline.from(
							leaf,
							{
								scale: 0,
								duration: 0.035,
								ease: 'back.out(2.2)',
								transformOrigin: '0px 0px',
							},
							gsap.utils.clamp(0, 0.96, frac - 0.01)
						);
					});

				// The glowing tip rides the freshly drawn end of the stem. The
				// draw tween spans the whole timeline (duration 1, every other
				// tween ends earlier), so timeline progress == draw progress.
				const setX = gsap.quickSetter(tip, 'x', 'px');
				const setY = gsap.quickSetter(tip, 'y', 'px');
				const placeTip = () => {
					const point = stem.getPointAtLength(
						timeline.progress() * total
					);
					setX(point.x);
					setY(point.y);
				};
				timeline.eventCallback('onUpdate', placeTip);
				timeline.to(tip, { autoAlpha: 0, duration: 0.025 }, 0.965);

				// Drive the timeline so the freshly drawn tip rides at ~62% of
				// the viewport height; the eased quickTo gives the organic
				// trail-and-catch-up feel a raw scrub can't do here.
				const progressTo = gsap.quickTo(timeline, 'progress', {
					duration: 0.7,
					ease: 'power3.out',
				});
				const el = container.current;
				const targetFrac = () => {
					if (!el) return 0;
					const rect = el.getBoundingClientRect();
					// rect.top is the container's offset from the viewport top,
					// so -rect.top is how far we've scrolled into it.
					return fracForY(-rect.top + window.innerHeight * 0.62);
				};
				ScrollTrigger.create({
					trigger: el,
					start: 'top top',
					end: 'bottom bottom',
					onUpdate: () => progressTo(targetFrac()),
					onRefresh: () => progressTo(targetFrac()),
				});
				timeline.progress(targetFrac());
				placeTip();

				// Soft breathing on the halo, independent of scroll.
				if (tipHaloRef.current) {
					gsap.to(tipHaloRef.current, {
						opacity: 0.18,
						scale: 1.35,
						transformOrigin: 'center',
						duration: 1.6,
						repeat: -1,
						yoyo: true,
						ease: 'sine.inOut',
					});
				}
			});
		},
		{ dependencies: [layout], revertOnUpdate: true, scope: svgRef }
	);

	if (!layout) return null;

	return (
		<svg
			ref={svgRef}
			aria-hidden
			fill="none"
			width={layout.width}
			height={layout.height}
			viewBox={`0 0 ${layout.width} ${layout.height}`}
			className="pointer-events-none absolute inset-0 z-[1] h-full w-full"
		>
			<defs>
				<linearGradient
					id="vine-stem-gradient"
					gradientUnits="userSpaceOnUse"
					x1="0"
					y1="0"
					x2="0"
					y2={layout.height}
				>
					<stop offset="0" stopColor="oklch(0.68 0.12 145)" />
					<stop offset="0.55" stopColor="oklch(0.78 0.15 135)" />
					<stop offset="1" stopColor="oklch(0.87 0.18 125)" />
				</linearGradient>
				<radialGradient id="vine-tip-gradient">
					<stop offset="0" stopColor="oklch(0.92 0.16 122)" />
					<stop
						offset="1"
						stopColor="oklch(0.92 0.16 122)"
						stopOpacity="0"
					/>
				</radialGradient>
			</defs>

			{/* Wide faint strokes under the core = cheap glow, no SVG filters. */}
			<path
				className="vine-draw"
				d={layout.d}
				stroke="oklch(0.87 0.18 125)"
				strokeOpacity="0.08"
				strokeWidth="22"
				strokeLinecap="round"
			/>
			<path
				className="vine-draw"
				d={layout.d}
				stroke="oklch(0.87 0.18 125)"
				strokeOpacity="0.22"
				strokeWidth="9"
				strokeLinecap="round"
			/>
			<path
				ref={stemRef}
				className="vine-draw"
				d={layout.d}
				stroke="url(#vine-stem-gradient)"
				strokeWidth="3.5"
				strokeLinecap="round"
			/>

			{layout.leaves.map((leaf, i) => (
				<g
					key={`${i}-${leaf.x.toFixed(0)}-${leaf.y.toFixed(0)}`}
					transform={`translate(${leaf.x} ${leaf.y}) rotate(${leaf.angle}) scale(${
						leaf.flip ? -leaf.size : leaf.size
					} ${leaf.size})`}
				>
					{leaf.fill === BERRY_FILL ? (
						<g className="vine-leaf" data-frac={leaf.frac}>
							<circle cx="3" cy="-9" r="3.4" fill={BERRY_FILL} />
							<circle
								cx="9"
								cy="-15"
								r="2.6"
								fill={BERRY_FILL}
								opacity="0.8"
							/>
							<circle
								cx="-2"
								cy="-16"
								r="2.1"
								fill={BERRY_FILL}
								opacity="0.65"
							/>
						</g>
					) : (
						<g className="vine-leaf" data-frac={leaf.frac}>
							<path d={LEAF_BLADE} fill={leaf.fill} opacity="0.92" />
							<path
								d={LEAF_RIB}
								stroke="oklch(0.16 0.021 155 / 0.55)"
								strokeWidth="1"
							/>
						</g>
					)}
				</g>
			))}

			<g ref={tipRef} className="vine-tip">
				<circle
					ref={tipHaloRef}
					r="24"
					fill="url(#vine-tip-gradient)"
					opacity="0.5"
				/>
				<circle r="5" fill="oklch(0.87 0.18 125)" />
				<circle r="2.2" fill="oklch(0.98 0.03 110)" />
			</g>
		</svg>
	);
}
