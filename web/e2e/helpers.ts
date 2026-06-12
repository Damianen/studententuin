import type { Page } from '@playwright/test';
import { expect } from '@playwright/test';

export interface TestUser {
	name: string;
	email: string;
	password: string;
}

// Unique per run so e2e never depends on a clean database.
export function uniqueUser(prefix: string): TestUser {
	const id = `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`;
	return {
		name: `E2E ${prefix}`,
		email: `e2e-${prefix}-${id}@test.local`,
		password: 'e2e-password-123',
	};
}

export function uniqueSubdomain(prefix: string): string {
	return `e2e-${prefix}-${Date.now().toString(36)}`;
}

export async function register(page: Page, user: TestUser) {
	await page.goto('/register');
	await page.getByLabel('Name').fill(user.name);
	await page.getByLabel('Email').fill(user.email);
	await page.getByLabel('Password', { exact: true }).fill(user.password);
	await page.getByLabel('Confirm password').fill(user.password);
	await page.getByRole('button', { name: 'Create account' }).click();
	await expect(page).toHaveURL(/\/projects$/);
}

// Submits the login form without waiting for the outcome.
export async function attemptLogin(page: Page, user: TestUser) {
	await page.goto('/login');
	await page.getByLabel('Email').fill(user.email);
	await page.getByLabel('Password', { exact: true }).fill(user.password);
	await page.getByRole('button', { name: 'Sign in' }).click();
}

export async function login(page: Page, user: TestUser) {
	await attemptLogin(page, user);
	await expect(page).toHaveURL(/\/projects$/);
}

export async function logout(page: Page) {
	await page.getByRole('button', { name: 'Account menu' }).click();
	await page.getByRole('menuitem', { name: 'Sign out' }).click();
	// The header navigates to / but the route guard races it to /login when
	// signing out from a protected page; either way the session is gone.
	await expect(page).toHaveURL(/\/(login)?$/);
	await expect(
		page.getByRole('button', { name: 'Account menu' })
	).toBeHidden();
}

// Radix Select doesn't respond to selectOption(); click trigger, then option.
export async function pickOption(page: Page, trigger: string, option: string) {
	await page.getByRole('combobox').filter({ hasText: trigger }).click();
	await page.getByRole('option', { name: option }).click();
}
