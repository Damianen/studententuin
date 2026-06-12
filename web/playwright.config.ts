import { defineConfig, devices } from '@playwright/test';

// E2E tests run against the docker-compose stack (web + api + postgres):
//   docker compose up -d --build --wait
//   npm run test:e2e
// Override the target with BASE_URL or WEB_PORT when 3000 is taken locally.
const baseURL =
	process.env.BASE_URL ?? `http://localhost:${process.env.WEB_PORT ?? 3000}`;

export default defineConfig({
	testDir: './e2e',
	// The stack shares one Postgres; parallel workers would interleave data.
	workers: 1,
	fullyParallel: false,
	retries: process.env.CI ? 1 : 0,
	timeout: 30_000,
	reporter: [['list'], ['html', { open: 'never' }]],
	use: {
		baseURL,
		trace: 'on-first-retry',
	},
	projects: [
		{
			name: 'chromium',
			use: {
				...devices['Desktop Chrome'],
				// Use a system browser locally (e.g. PW_CHROMIUM_PATH=/usr/bin/chromium)
				// instead of downloading one; CI runs `playwright install chromium`.
				launchOptions: {
					executablePath: process.env.PW_CHROMIUM_PATH || undefined,
				},
			},
		},
	],
});
