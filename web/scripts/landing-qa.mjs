// Visual QA for the landing page: screenshots at several scroll depths,
// desktop + mobile, plus reduced-motion and app-route spot checks.
// Run with the dev server up: node scripts/landing-qa.mjs
import { chromium } from 'playwright-core';
import { mkdirSync } from 'node:fs';

const BASE = process.env.QA_BASE_URL ?? 'http://localhost:5199';
const OUT = new URL('./qa-landing/', import.meta.url).pathname;
mkdirSync(OUT, { recursive: true });

const browser = await chromium.launch({
	executablePath: '/usr/bin/chromium',
	headless: true,
});

const failures = [];

async function sweep(label, viewport, { reducedMotion = 'no-preference' } = {}) {
	const context = await browser.newContext({ viewport, reducedMotion });
	const page = await context.newPage();
	const consoleErrors = [];
	page.on('console', (msg) => {
		if (msg.type() !== 'error') return;
		// The auth probe 4xxs by design when logged out / API is elsewhere.
		if (/Failed to load resource.*40[14]/.test(msg.text())) return;
		consoleErrors.push(msg.text());
	});
	page.on('pageerror', (err) => consoleErrors.push(String(err)));

	await page.goto(`${BASE}/`, { waitUntil: 'networkidle' });
	await page.waitForTimeout(1500);

	const overflow = await page.evaluate(() => {
		const doc = document.documentElement;
		if (doc.scrollWidth <= doc.clientWidth + 1) return null;
		return { scrollWidth: doc.scrollWidth, clientWidth: doc.clientWidth };
	});
	if (overflow) failures.push(`${label}: horizontal overflow ${JSON.stringify(overflow)}`);

	const max = await page.evaluate(
		() => document.documentElement.scrollHeight - window.innerHeight
	);
	for (const pct of [0, 20, 40, 60, 80, 100]) {
		await page.evaluate((y) => window.scrollTo(0, y), (max * pct) / 100);
		// let the scrubbed animations catch up
		await page.waitForTimeout(1200);
		await page.screenshot({ path: `${OUT}${label}-${String(pct).padStart(3, '0')}.png` });
	}

	if (consoleErrors.length) {
		failures.push(`${label}: console errors:\n  ${consoleErrors.join('\n  ')}`);
	}
	await context.close();
}

async function route(label, path, viewport) {
	const context = await browser.newContext({ viewport });
	const page = await context.newPage();
	await page.goto(`${BASE}${path}`, { waitUntil: 'networkidle' });
	await page.waitForTimeout(600);
	await page.screenshot({ path: `${OUT}${label}.png`, fullPage: true });
	await context.close();
}

await sweep('desktop', { width: 1440, height: 900 });
await sweep('mobile', { width: 390, height: 844 });
await sweep('reduced-motion', { width: 1440, height: 900 }, { reducedMotion: 'reduce' });
await route('app-login', '/login', { width: 1440, height: 900 });
await route('app-register', '/register', { width: 1440, height: 900 });
await route('app-404', '/nope', { width: 1440, height: 900 });

await browser.close();

if (failures.length) {
	console.error('FAILURES:\n' + failures.join('\n'));
	process.exit(1);
}
console.log('QA sweep passed. Screenshots in', OUT);
