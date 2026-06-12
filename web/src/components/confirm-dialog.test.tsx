import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { ConfirmDialog } from './confirm-dialog';
import { Button } from '@/components/ui/button';

function renderDialog(props: Partial<Parameters<typeof ConfirmDialog>[0]> = {}) {
	const onConfirm = vi.fn();
	render(
		<ConfirmDialog
			trigger={<Button>Delete project</Button>}
			title="Delete project"
			description="This cannot be undone."
			onConfirm={onConfirm}
			{...props}
		/>
	);
	return onConfirm;
}

describe('ConfirmDialog', () => {
	it('confirms immediately when no confirmText is required', async () => {
		const onConfirm = renderDialog();

		await userEvent.click(screen.getByRole('button', { name: 'Delete project' }));
		const confirm = await screen.findByRole('button', { name: 'Delete' });
		expect(confirm).toBeEnabled();

		await userEvent.click(confirm);
		expect(onConfirm).toHaveBeenCalledTimes(1);
	});

	it('keeps confirm disabled until the exact phrase is typed', async () => {
		const onConfirm = renderDialog({ confirmText: 'myproject' });

		await userEvent.click(screen.getByRole('button', { name: 'Delete project' }));
		const confirm = await screen.findByRole('button', { name: 'Delete' });
		const input = screen.getByLabelText(/to confirm/);

		expect(confirm).toBeDisabled();

		await userEvent.type(input, 'wrong');
		expect(confirm).toBeDisabled();

		await userEvent.clear(input);
		await userEvent.type(input, 'myproject');
		expect(confirm).toBeEnabled();

		await userEvent.click(confirm);
		expect(onConfirm).toHaveBeenCalledTimes(1);
	});

	it('cancel closes the dialog without confirming', async () => {
		const onConfirm = renderDialog();

		await userEvent.click(screen.getByRole('button', { name: 'Delete project' }));
		await userEvent.click(await screen.findByRole('button', { name: 'Cancel' }));

		await waitFor(() =>
			expect(screen.queryByRole('dialog')).not.toBeInTheDocument()
		);
		expect(onConfirm).not.toHaveBeenCalled();
	});

	it('clears the typed phrase when the dialog is dismissed and reopened', async () => {
		renderDialog({ confirmText: 'myproject' });

		await userEvent.click(screen.getByRole('button', { name: 'Delete project' }));
		await userEvent.type(screen.getByLabelText(/to confirm/), 'myproject');
		// Dismissing (Escape/overlay) goes through onOpenChange, which resets
		// the typed phrase. The Cancel button only flips `open` and does not.
		await userEvent.keyboard('{Escape}');
		await waitFor(() =>
			expect(screen.queryByRole('dialog')).not.toBeInTheDocument()
		);

		await userEvent.click(screen.getByRole('button', { name: 'Delete project' }));
		expect(screen.getByLabelText(/to confirm/)).toHaveValue('');
		expect(screen.getByRole('button', { name: 'Delete' })).toBeDisabled();
	});
});
