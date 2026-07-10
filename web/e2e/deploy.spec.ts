import { test, expect } from '@playwright/test';
import { uniqueUser, uniqueSubdomain, register, login } from './helpers';

// The Epic 4 acceptance path: paste a GitHub URL, deploy, watch it go
// pending→running, see real logs, stop it from the UI. Needs a real public
// repository (the manager clones over https from an allowlisted host), so
// the suite only runs when E2E_DEPLOY_REPO points at one — see
// e2e-fixtures/node-hello/README.md.
const deployRepo = process.env.E2E_DEPLOY_REPO;

test.describe.configure({ mode: 'serial' });

const user = uniqueUser('deploy');
const subdomain = uniqueSubdomain('deploy');

test.describe('deploy from git', () => {
	test.skip(!deployRepo, 'set E2E_DEPLOY_REPO to a public repo to enable');

	test('create the application', async ({ page }) => {
		await register(page, user);
		await page.goto(
			`/projects/new/app?subdomain=${encodeURIComponent(subdomain)}`
		);

		await page.getByLabel('Name', { exact: true }).fill('deploy-me');
		await page.getByLabel('Repository URL').fill(deployRepo!);
		await page.getByLabel('Branch').fill('main');
		await page.getByLabel('Build command').fill('npm run build');
		await page.getByLabel('Start command').fill('npm start');
		await page.getByRole('button', { name: 'Create application' }).click();
		await expect(page).toHaveURL(/\/projects\/[0-9a-f-]{36}$/);
	});

	test('deploy and watch it go live', async ({ page }) => {
		// Clone + nixpacks build + start: give the whole journey ten minutes.
		test.setTimeout(10 * 60 * 1000);

		await login(page, user);
		await page
			.getByRole('link', { name: `Open project ${subdomain}` })
			.click();
		await page.getByRole('link', { name: 'Open', exact: true }).click();
		await page.getByRole('link', { name: 'Deployments' }).click();

		await page.getByRole('button', { name: 'Deploy now' }).click();
		await expect(page.getByText('Deploying', { exact: true })).toBeVisible();

		// The page reloads the record when the deployment settles.
		await expect(page.getByText('running').first()).toBeVisible({
			timeout: 9 * 60 * 1000,
		});

		// Real container output lands in the logs tab.
		await page.getByRole('link', { name: 'Logs' }).click();
		await expect(page.getByText('listening on', { exact: false })).toBeVisible(
			{ timeout: 60 * 1000 }
		);
	});

	test('save env vars and redeploy from the settings tab', async ({ page }) => {
		// A redeploy repeats the whole clone + build journey.
		test.setTimeout(10 * 60 * 1000);

		await login(page, user);
		await page
			.getByRole('link', { name: `Open project ${subdomain}` })
			.click();
		await page.getByRole('link', { name: 'Open', exact: true }).click();
		await page.getByRole('link', { name: 'Settings' }).click();

		const envCard = page.locator('[data-slot="card"]', {
			has: page.getByText('Environment', { exact: true }),
		});
		await envCard.getByLabel('Key').last().fill('E2E_FLAG');
		await envCard.getByLabel('Value').last().fill('from-settings');
		await envCard.getByRole('button', { name: 'Add' }).click();
		await envCard.getByRole('button', { name: 'Save changes' }).click();
		await expect(
			page.getByText('Environment variables saved').first()
		).toBeVisible();

		// Saving offers a redeploy in place, with live stage feedback.
		await envCard.getByRole('button', { name: 'Deploy now' }).click();
		await expect(envCard.getByText('Deploying —')).toBeVisible();
		await expect(page.getByText('deploy-me is live').first()).toBeVisible({
			timeout: 9 * 60 * 1000,
		});
	});

	test('stop it from the UI', async ({ page }) => {
		await login(page, user);
		await page
			.getByRole('link', { name: `Open project ${subdomain}` })
			.click();
		await page.getByRole('link', { name: 'Open', exact: true }).click();

		await page.getByRole('button', { name: 'Stop' }).click();
		await expect(page.getByText('stopped').first()).toBeVisible({
			timeout: 60 * 1000,
		});
		await expect(page.getByRole('button', { name: 'Start' })).toBeVisible();
	});
});
