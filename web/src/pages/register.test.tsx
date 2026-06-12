import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter, Route, Routes } from 'react-router';
import Register from './register';
import { useAuth } from '@/contexts/auth_context';

vi.mock('@/contexts/auth_context', () => ({
	useAuth: vi.fn(),
}));

const register = vi.fn();

function renderRegister() {
	return render(
		<MemoryRouter initialEntries={['/register']}>
			<Routes>
				<Route path="/register" element={<Register />} />
				<Route path="/projects" element={<div>projects page</div>} />
			</Routes>
		</MemoryRouter>
	);
}

beforeEach(() => {
	register.mockReset();
	vi.mocked(useAuth).mockReturnValue({
		user: null,
		loading: false,
		isAuthenticated: false,
		login: vi.fn(),
		logout: vi.fn(),
		register,
		updateUser: vi.fn(),
		deleteUser: vi.fn(),
	});
});

async function fillForm({
	name = 'Alice',
	email = 'a@b.com',
	password = 'secret123',
	confirm = 'secret123',
} = {}) {
	await userEvent.type(screen.getByLabelText('Name'), name);
	await userEvent.type(screen.getByLabelText('Email'), email);
	await userEvent.type(screen.getByLabelText('Password'), password);
	await userEvent.type(screen.getByLabelText('Confirm password'), confirm);
	await userEvent.click(
		screen.getByRole('button', { name: 'Create account' })
	);
}

describe('Register page', () => {
	it('registers and navigates to /projects', async () => {
		register.mockResolvedValue(undefined);
		renderRegister();

		await fillForm();

		expect(register).toHaveBeenCalledWith({
			name: 'Alice',
			email: 'a@b.com',
			password: 'secret123',
		});
		expect(await screen.findByText('projects page')).toBeInTheDocument();
	});

	it('rejects mismatched passwords without calling register', async () => {
		renderRegister();

		await fillForm({ confirm: 'different123' });

		expect(await screen.findByRole('alert')).toHaveTextContent(
			'Passwords do not match.'
		);
		expect(register).not.toHaveBeenCalled();
	});

	it('shows the API error when registration fails', async () => {
		register.mockRejectedValue(new Error('email already in use'));
		renderRegister();

		await fillForm();

		expect(await screen.findByRole('alert')).toHaveTextContent(
			'email already in use'
		);
		expect(screen.queryByText('projects page')).not.toBeInTheDocument();
	});
});
