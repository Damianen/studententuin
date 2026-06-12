import { test, expect } from '@playwright/test';
import { uniqueUser, register, login, attemptLogin } from './helpers';

test.describe.configure({ mode: 'serial' });

const user = uniqueUser('account');
const newName = 'E2E Account Renamed';

test.describe('account', () => {
	test('update the display name', async ({ page }) => {
		await register(page, user);

		await page.getByRole('button', { name: 'Account menu' }).click();
		await page.getByRole('menuitem', { name: 'Account' }).click();
		await expect(page).toHaveURL(/\/account$/);
		await expect(page.getByLabel('Name')).toHaveValue(user.name);

		await page.getByLabel('Name').fill(newName);
		await page.getByRole('button', { name: 'Save changes' }).click();
		await expect(page.getByText('Account updated')).toBeVisible();

		// The change survives a full reload, so it really hit the API.
		await page.reload();
		await expect(page.getByLabel('Name')).toHaveValue(newName);
		await expect(
			page.getByRole('button', { name: 'Account menu' })
		).toContainText(newName);
	});

	test('delete the account and verify the credentials are gone', async ({
		page,
	}) => {
		await login(page, user);
		await page.goto('/account');

		await page.getByRole('button', { name: 'Delete account' }).click();
		const dialog = page.getByRole('dialog');
		const confirm = dialog.getByRole('button', { name: 'Delete my account' });
		await expect(confirm).toBeDisabled();
		await dialog.locator('#confirm-text').fill(user.email);
		await expect(confirm).toBeEnabled();
		await confirm.click();

		// navigate('/') races the route guard's redirect to /login; either
		// destination means the session ended.
		await expect(page).toHaveURL(/\/(login)?$/);

		// The account is gone: logging in with the old credentials fails.
		await attemptLogin(page, user);
		await expect(page.getByRole('alert')).toHaveText(
			'email or password not correct!'
		);
	});
});
