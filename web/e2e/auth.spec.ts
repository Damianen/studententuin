import { test, expect } from '@playwright/test';
import { uniqueUser, register, login, logout } from './helpers';

test.describe.configure({ mode: 'serial' });

const user = uniqueUser('auth');

test.describe('authentication', () => {
	test('register auto-logs-in and lands on the empty garden', async ({
		page,
	}) => {
		await register(page, user);

		await expect(page.getByText('Empty plot')).toBeVisible();
		await expect(
			page.getByRole('button', { name: 'Account menu' })
		).toContainText(user.name);
	});

	test('register with an email that is already taken fails', async ({
		page,
	}) => {
		await page.goto('/register');
		await page.getByLabel('Name').fill(user.name);
		await page.getByLabel('Email').fill(user.email);
		await page.getByLabel('Password', { exact: true }).fill(user.password);
		await page.getByLabel('Confirm password').fill(user.password);
		await page.getByRole('button', { name: 'Create account' }).click();

		await expect(page.getByRole('alert')).toHaveText('email already in use');
		await expect(page).toHaveURL(/\/register$/);
	});

	test('logout ends the session and the guard kicks back in', async ({
		page,
	}) => {
		await login(page, user);
		await expect(page).toHaveURL(/\/projects$/);

		await logout(page);

		await page.goto('/projects');
		await expect(page).toHaveURL(/\/login$/);
	});

	test('login works again and an authenticated user cannot see /login', async ({
		page,
	}) => {
		await login(page, user);
		await expect(page).toHaveURL(/\/projects$/);

		// PublicOnly bounces authenticated users back to their garden.
		await page.goto('/login');
		await expect(page).toHaveURL(/\/projects$/);
	});
});
