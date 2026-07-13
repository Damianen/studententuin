import { test, expect } from '@playwright/test';
import { uniqueUser, uniqueSubdomain, register, login } from './helpers';

// Phase 5 acceptance: create database in UI → running postgres container on
// the hosting server; the connection string lands on the record; the logs tab
// shows the real container output. Needs the full compose stack (manager with
// the Docker socket) — the manager pulls postgres:16 on first run, so the
// running-state budget is generous.
test.describe.configure({ mode: 'serial' });

const user = uniqueUser('db');
const subdomain = uniqueSubdomain('db');

// The provision pipeline pulls the image on a cold host; four minutes covers
// a slow CI network, while a warm host settles in seconds.
const PROVISION_BUDGET_MS = 240_000;

// Captured after create; later tests navigate directly — the project card's
// overlay link is covered by the "add an application" link on a db-only
// project, so role-based clicks are flaky here.
let projectUrl = '';

test.describe('database provisioning', () => {
	test('create a postgres database and watch it reach running', async ({
		page,
	}) => {
		test.setTimeout(PROVISION_BUDGET_MS + 60_000);

		await register(page, user);
		await page.goto('/projects/new/database');

		await page.getByLabel('Name', { exact: true }).fill('e2e-postgres');
		await page.getByLabel('Subdomain').fill(subdomain);
		// Engine and version selects keep their defaults (PostgreSQL 16);
		// credentials are generated server-side.
		await page.getByLabel('Database name').fill('app_db');
		await page.getByRole('button', { name: 'Create database' }).click();

		await expect(page).toHaveURL(/\/projects\/[0-9a-f-]{36}$/);
		projectUrl = new URL(page.url()).pathname;
		await expect(page.getByText('provisioning').first()).toBeVisible();

		// The page polls while provisioning; the badge flips on its own.
		await expect(page.getByText('running').first()).toBeVisible({
			timeout: PROVISION_BUDGET_MS,
		});
	});

	test('the connection string is on the card', async ({ page }) => {
		await login(page, user);
		await page.goto(projectUrl);

		await page
			.getByRole('button', { name: 'Show connection string' })
			.click();
		await expect(page.getByText(/postgres:\/\/app:[0-9a-f]{32}@stt-db-/)).toBeVisible();
		await expect(
			page.getByRole('button', { name: 'Copy connection string' })
		).toBeVisible();
	});

	test('the logs tab shows real postgres output', async ({ page }) => {
		await login(page, user);
		await page.goto(`${projectUrl}/db`);
		await page.getByRole('link', { name: 'Logs' }).click();

		// Real container output, not the seeded sample data of phases 2-4.
		await expect(
			page.getByText('database system is ready to accept connections').first()
		).toBeVisible({ timeout: 60_000 });
		await expect(
			page.locator('[data-slot="card"]', {
				has: page.getByText('Runtime logs'),
			}).getByText('sample data')
		).toHaveCount(0);
	});

	test('deleting the project sweeps the database', async ({ page }) => {
		await login(page, user);
		await page.goto(projectUrl);

		await page.getByRole('button', { name: 'Delete project' }).click();
		const dialog = page.getByRole('dialog');
		await dialog.locator('#confirm-text').fill(subdomain);
		await dialog.getByRole('button', { name: 'Delete', exact: true }).click();

		await expect(page).toHaveURL(/\/projects$/);
		// The container, network, and data volume are removed server-side —
		// the manager integration test asserts that half; here the project is
		// simply gone.
		await expect(page.getByText('Empty plot')).toBeVisible();
	});
});
