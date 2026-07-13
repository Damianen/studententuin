import { test, expect } from '@playwright/test';
import {
	uniqueUser,
	uniqueSubdomain,
	register,
	login,
	pickOption,
} from './helpers';

test.describe.configure({ mode: 'serial' });

const user = uniqueUser('projects');
const subdomain = uniqueSubdomain('garden');
const renamed = `${subdomain}-renamed`;

test.describe('projects', () => {
	test('a fresh account starts with an empty garden', async ({ page }) => {
		await register(page, user);
		await expect(page.getByText('Empty plot')).toBeVisible();
	});

	test('create a database project', async ({ page }) => {
		await login(page, user);
		await page.goto('/projects/new/database');

		await page.getByLabel('Name', { exact: true }).fill('my-database');
		await page.getByLabel('Subdomain').fill(subdomain);
		// Engine (postgres) and version (16) keep their defaults; credentials
		// are generated server-side since phase 5, so there is no password field.
		await pickOption(page, 'PostgreSQL', 'PostgreSQL');
		await page.getByLabel('Database name').fill('app_db');
		await page.getByRole('button', { name: 'Create database' }).click();

		// ensureSubdomain claims the subdomain, then the database is planted on it.
		await expect(page).toHaveURL(/\/projects\/[0-9a-f-]{36}$/);
		await expect(
			page.getByRole('heading', { name: subdomain, exact: true })
		).toBeVisible();
		await expect(page.getByText('my-database').first()).toBeVisible();
		await expect(page.getByText('provisioning').first()).toBeVisible();
	});

	test('add an application to the same project', async ({ page }) => {
		await login(page, user);
		await page.goto(
			`/projects/new/app?subdomain=${encodeURIComponent(subdomain)}`
		);

		await expect(page.getByLabel('Subdomain')).toHaveValue(subdomain);
		await page.getByLabel('Name', { exact: true }).fill('my-app');
		await page
			.getByLabel('Repository URL')
			.fill('https://github.com/example/my-app');
		await page.getByLabel('Branch').fill('main');
		await page.getByLabel('Build command').fill('npm run build');
		await page.getByLabel('Start command').fill('npm start');
		await page.getByRole('button', { name: 'Create application' }).click();

		await expect(page).toHaveURL(/\/projects\/[0-9a-f-]{36}$/);
		await expect(page.getByText('my-app').first()).toBeVisible();
	});

	test('the garden lists the project with both resources', async ({
		page,
	}) => {
		await login(page, user);

		const card = page
			.getByRole('link', { name: `Open project ${subdomain}` })
			.locator('..');
		await expect(card).toBeVisible();
		await expect(card.getByText('my-app')).toBeVisible();
		await expect(card.getByText('my-database')).toBeVisible();
		await expect(
			card.getByText(`${subdomain}.studententuin.com`)
		).toBeVisible();
	});

	test('rename the project', async ({ page }) => {
		await login(page, user);
		await page.getByRole('link', { name: `Open project ${subdomain}` }).click();

		await page.getByRole('button', { name: 'Rename project' }).click();
		// The inline rename input has no label; it grabs focus on open.
		const nameInput = page.locator('input:focus');
		await nameInput.fill(renamed);
		await nameInput.press('Enter');

		await expect(
			page.getByRole('heading', { name: renamed, exact: true })
		).toBeVisible();
	});

	test('delete the project with type-to-confirm', async ({ page }) => {
		// Deleting is blocked (409) while the database provision is in flight,
		// so wait for the badge to settle first (cold hosts pull the image).
		test.setTimeout(300_000);
		await login(page, user);
		await page.getByRole('link', { name: `Open project ${renamed}` }).click();
		await expect(page).toHaveURL(/\/projects\/[0-9a-f-]{36}$/);
		await expect(page.getByText(/^(running|failed)$/).first()).toBeVisible({
			timeout: 240_000,
		});

		await page.getByRole('button', { name: 'Delete project' }).click();
		const dialog = page.getByRole('dialog');
		await expect(dialog).toBeVisible();

		const confirm = dialog.getByRole('button', { name: 'Delete', exact: true });
		await expect(confirm).toBeDisabled();
		await dialog.locator('#confirm-text').fill(renamed);
		await expect(confirm).toBeEnabled();
		await confirm.click();

		await expect(page).toHaveURL(/\/projects$/);
		await expect(page.getByText('Empty plot')).toBeVisible();
	});
});
