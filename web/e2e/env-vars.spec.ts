import { test, expect } from '@playwright/test';
import type { Page } from '@playwright/test';
import { uniqueUser, uniqueSubdomain, register, login } from './helpers';

// Phase 4 round-trip: env vars edited in the settings tab persist through the
// api into postgres and come back after a full reload. No deploy involved, so
// unlike deploy.spec.ts this runs without E2E_DEPLOY_REPO — the repository URL
// is never cloned.
test.describe.configure({ mode: 'serial' });

const user = uniqueUser('env');
const subdomain = uniqueSubdomain('env');

// The settings tab has two "Save changes" buttons (configuration + env), so
// every interaction is scoped to the Environment card.
function envCard(page: Page) {
	return page.locator('[data-slot="card"]', {
		has: page.getByText('Environment', { exact: true }),
	});
}

async function openSettings(page: Page) {
	await login(page, user);
	// The project card is covered by its overlay link, so navigate by role.
	await page
		.getByRole('link', { name: `Open project ${subdomain}` })
		.click();
	await page.getByRole('link', { name: 'Open', exact: true }).click();
	await page.getByRole('link', { name: 'Settings' }).click();
}

test.describe('environment variables round-trip', () => {
	test('create the application', async ({ page }) => {
		await register(page, user);
		await page.goto(
			`/projects/new/app?subdomain=${encodeURIComponent(subdomain)}`
		);

		await page.getByLabel('Name', { exact: true }).fill('env-app');
		await page
			.getByLabel('Repository URL')
			.fill('https://github.com/e2e/placeholder');
		await page.getByLabel('Branch').fill('main');
		await page.getByLabel('Build command').fill('npm run build');
		await page.getByLabel('Start command').fill('npm start');
		await page.getByRole('button', { name: 'Create application' }).click();
		await expect(page).toHaveURL(/\/projects\/[0-9a-f-]{36}$/);
	});

	test('save an env var and read it back after reload', async ({ page }) => {
		await openSettings(page);
		const card = envCard(page);

		// A fresh application starts with no variables — only the draft row.
		await expect(card.getByLabel('Key')).toHaveCount(1);

		await card.getByLabel('Key').fill('E2E_FLAG');
		await card.getByLabel('Value').fill('hello');
		await card.getByRole('button', { name: 'Add' }).click();
		await card.getByRole('button', { name: 'Save changes' }).click();

		await expect(
			page.getByText('Environment variables saved').first()
		).toBeVisible();
		// Saving surfaces the redeploy offer.
		await expect(
			card.getByText('take effect on the next deploy')
		).toBeVisible();
		await expect(card.getByRole('button', { name: 'Deploy now' })).toBeVisible();

		await page.reload();
		await expect(envCard(page).getByLabel('Key').first()).toHaveValue(
			'E2E_FLAG'
		);
	});

	test('reject an invalid key without losing the stored one', async ({
		page,
	}) => {
		await openSettings(page);
		const card = envCard(page);

		await expect(card.getByLabel('Key').first()).toHaveValue('E2E_FLAG');

		// Keys must not start with a digit; the save is blocked client-side.
		await card.getByLabel('Key').last().fill('1BAD');
		await card.getByLabel('Value').last().fill('nope');
		await card.getByRole('button', { name: 'Add' }).click();
		await card.getByRole('button', { name: 'Save changes' }).click();
		await expect(
			page.getByText('Invalid variable key "1BAD"').first()
		).toBeVisible();

		// Dropping the bad row makes the save go through again.
		await card.getByRole('button', { name: 'Remove 1BAD' }).click();
		await card.getByRole('button', { name: 'Save changes' }).click();
		await expect(
			page.getByText('Environment variables saved').first()
		).toBeVisible();

		await page.reload();
		const reloaded = envCard(page);
		await expect(reloaded.getByLabel('Key').first()).toHaveValue('E2E_FLAG');
		await expect(reloaded.getByLabel('Key')).toHaveCount(2); // stored row + draft
	});

	// This app is never deployed, so it covers the phase-6 no-data path
	// without needing E2E_DEPLOY_REPO.
	test('metrics and deployments show their empty states before a deploy', async ({
		page,
	}) => {
		await login(page, user);
		await page
			.getByRole('link', { name: `Open project ${subdomain}` })
			.click();
		await page.getByRole('link', { name: 'Open', exact: true }).click();

		await page.getByRole('link', { name: 'Metrics' }).click();
		await expect(page.getByText('No metrics yet').first()).toBeVisible();

		await page.getByRole('link', { name: 'Deployments' }).click();
		await expect(
			page.getByText('No deployments yet — plant the first one with Deploy now.'),
		).toBeVisible();
	});
});
