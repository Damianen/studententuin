import gsap from 'gsap';
import { useGSAP } from '@gsap/react';
import { ScrollTrigger } from 'gsap/ScrollTrigger';
import { SplitText } from 'gsap/SplitText';
import { DrawSVGPlugin } from 'gsap/DrawSVGPlugin';
import { ScrollToPlugin } from 'gsap/ScrollToPlugin';

gsap.registerPlugin(useGSAP, ScrollTrigger, SplitText, DrawSVGPlugin, ScrollToPlugin);
ScrollTrigger.config({ ignoreMobileResize: true });

if (import.meta.env.DEV) {
	// Expose for headless QA scripts (scripts/landing-qa.mjs) and devtools.
	(window as unknown as Record<string, unknown>).__ScrollTrigger = ScrollTrigger;
}

// All landing animations are gated behind this media query; with reduced
// motion every element simply keeps its final (CSS) state.
export const MOTION_OK = '(prefers-reduced-motion: no-preference)';

export { gsap, useGSAP, ScrollTrigger, SplitText };
