import { test, expect } from '@playwright/test';

test.describe('smoke', () => {
	test('landing page renders', async ({ page }) => {
		await page.goto('/');
		await expect(page.getByText('Free hosting for students').first()).toBeVisible();
		await expect(
			page.getByRole('heading', { level: 1 }).filter({ hasText: 'Plant' })
		).toBeVisible();
	});

	test('API health is reachable through the web proxy', async ({ request }) => {
		const response = await request.get('/api/health');
		expect(response.ok()).toBeTruthy();
		expect(await response.json()).toEqual({ server: 'running' });
	});

	test('anonymous visit to /projects redirects to login', async ({ page }) => {
		await page.goto('/projects');
		await expect(page).toHaveURL(/\/login$/);
		await expect(page.getByRole('button', { name: 'Sign in' })).toBeVisible();
	});

	test('login with wrong credentials shows an error', async ({ page }) => {
		await page.goto('/login');
		await page.getByLabel('Email').fill('nobody@test.local');
		await page.getByLabel('Password', { exact: true }).fill('wrong-password');
		await page.getByRole('button', { name: 'Sign in' }).click();

		await expect(page.getByRole('alert')).toHaveText(
			'email or password not correct!'
		);
		await expect(page).toHaveURL(/\/login$/);
	});
});
