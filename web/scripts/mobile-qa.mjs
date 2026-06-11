// Mobile QA sweep: screenshots every route at a phone viewport, checks for
// horizontal overflow and console errors. Run with the dev server up:
//   node scripts/mobile-qa.mjs
import { chromium } from 'playwright-core';
import { mkdirSync } from 'node:fs';

const BASE = process.env.QA_BASE_URL ?? 'http://localhost:5173';
const EMAIL = process.env.QA_EMAIL ?? 'test@student.nl';
const PASSWORD = process.env.QA_PASSWORD ?? 'hunter2hunter2';
const OUT = new URL('./qa-screens/', import.meta.url).pathname;
mkdirSync(OUT, { recursive: true });

const browser = await chromium.launch({
	executablePath: '/usr/bin/chromium',
	headless: true,
});

const failures = [];

async function checkPage(page, label, path, { fullPage = true } = {}) {
	const consoleErrors = [];
	const onConsole = (msg) => {
		if (msg.type() !== 'error') return;
		// The auth probe 401s by design when logged out; browsers always log it.
		if (/Failed to load resource.*401/.test(msg.text())) return;
		consoleErrors.push(msg.text());
	};
	const onPageError = (err) => consoleErrors.push(String(err));
	page.on('console', onConsole);
	page.on('pageerror', onPageError);

	await page.goto(`${BASE}${path}`, { waitUntil: 'networkidle' });
	await page.waitForTimeout(600);

	const overflow = await page.evaluate(() => {
		const doc = document.documentElement;
		if (doc.scrollWidth <= doc.clientWidth + 1) return null;
		const offenders = [];
		for (const el of document.querySelectorAll('*')) {
			const rect = el.getBoundingClientRect();
			if (rect.right > doc.clientWidth + 1 || rect.left < -1) {
				offenders.push(
					`${el.tagName.toLowerCase()}.${[...el.classList].slice(0, 3).join('.')} [${Math.round(rect.left)}..${Math.round(rect.right)}]`
				);
				if (offenders.length >= 5) break;
			}
		}
		return {
			scrollWidth: doc.scrollWidth,
			clientWidth: doc.clientWidth,
			offenders,
		};
	});

	await page.screenshot({ path: `${OUT}${label}.png`, fullPage });

	page.off('console', onConsole);
	page.off('pageerror', onPageError);

	if (overflow) {
		failures.push(
			`${label}: horizontal overflow ${overflow.scrollWidth}px > ${overflow.clientWidth}px\n    ${overflow.offenders.join('\n    ')}`
		);
	}
	if (consoleErrors.length) {
		failures.push(`${label}: console errors:\n    ${consoleErrors.join('\n    ')}`);
	}
	console.log(
		`✓ ${label}${overflow ? ' [OVERFLOW]' : ''}${consoleErrors.length ? ' [CONSOLE ERRORS]' : ''}`
	);
}

// ---- Pass 1: mobile viewport, reduced motion (deterministic screenshots)
const mobile = await browser.newContext({
	baseURL: BASE,
	viewport: { width: 390, height: 844 },
	deviceScaleFactor: 2,
	isMobile: true,
	hasTouch: true,
	reducedMotion: 'reduce',
});

// Authenticate via the shared cookie jar.
const login = await mobile.request.post('/api/auth/login', {
	data: { email: EMAIL, password: PASSWORD },
});
let authed = login.ok() && (await login.json()).code === 200;
if (!authed) {
	await mobile.request.post('/api/user/register', {
		data: { email: EMAIL, password: PASSWORD, name: 'QA Bot' },
	});
	const retry = await mobile.request.post('/api/auth/login', {
		data: { email: EMAIL, password: PASSWORD },
	});
	authed = retry.ok() && (await retry.json()).code === 200;
}
if (!authed) {
	console.error('Could not authenticate QA user — aborting.');
	await browser.close();
	process.exit(1);
}

const subdomains = (await (await mobile.request.get('/api/subdomain')).json())
	.data ?? [];
const projectId = subdomains[0]?.id;

const page = await mobile.newPage();
await checkPage(page, 'mobile-landing', '/');
await checkPage(page, 'mobile-projects', '/projects');
if (projectId) {
	await checkPage(page, 'mobile-project-detail', `/projects/${projectId}`);
} else {
	console.warn('! no subdomain found, skipping /projects/:id');
}
await checkPage(page, 'mobile-new-app', '/projects/new/app');
await checkPage(page, 'mobile-new-database', '/projects/new/database');
await checkPage(page, 'mobile-account', '/account');
await mobile.clearCookies();
await checkPage(page, 'mobile-login', '/login');
await checkPage(page, 'mobile-register', '/register');
await mobile.close();

// ---- Pass 2: desktop, normal motion — stepwise scroll through the landing
const desktop = await browser.newContext({
	baseURL: BASE,
	viewport: { width: 1440, height: 900 },
});
const dpage = await desktop.newPage();
await dpage.goto(`${BASE}/`, { waitUntil: 'networkidle' });
await dpage.waitForTimeout(1800); // hero intro
for (let step = 0; step <= 5; step++) {
	await dpage.evaluate(
		(s) =>
			window.scrollTo({
				top: (document.body.scrollHeight - window.innerHeight) * (s / 5),
			}),
		step
	);
	await dpage.waitForTimeout(900);
	await dpage.screenshot({ path: `${OUT}desktop-scroll-${step}.png` });
}
console.log('✓ desktop scroll sequence captured');
await desktop.close();

await browser.close();

if (failures.length) {
	console.error('\nFAILURES:\n' + failures.map((f) => `  ✗ ${f}`).join('\n'));
	process.exit(1);
}
console.log('\nAll routes clean: no overflow, no console errors.');
