import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter, Route, Routes } from 'react-router';
import Login from './login';
import { useAuth } from '@/contexts/auth_context';

vi.mock('@/contexts/auth_context', () => ({
	useAuth: vi.fn(),
}));

const login = vi.fn();

function renderLogin(initialEntries: (string | object)[] = ['/login']) {
	return render(
		<MemoryRouter initialEntries={initialEntries}>
			<Routes>
				<Route path="/login" element={<Login />} />
				<Route path="/projects" element={<div>projects page</div>} />
				<Route path="/account" element={<div>account page</div>} />
			</Routes>
		</MemoryRouter>
	);
}

beforeEach(() => {
	login.mockReset();
	vi.mocked(useAuth).mockReturnValue({
		user: null,
		loading: false,
		isAuthenticated: false,
		login,
		logout: vi.fn(),
		register: vi.fn(),
		updateUser: vi.fn(),
		deleteUser: vi.fn(),
	});
});

async function submit(email: string, password: string) {
	await userEvent.type(screen.getByLabelText('Email'), email);
	await userEvent.type(screen.getByLabelText('Password'), password);
	await userEvent.click(screen.getByRole('button', { name: 'Sign in' }));
}

describe('Login page', () => {
	it('submits credentials and navigates to /projects', async () => {
		login.mockResolvedValue(undefined);
		renderLogin();

		await submit('a@b.com', 'secret123');

		expect(login).toHaveBeenCalledWith({
			email: 'a@b.com',
			password: 'secret123',
		});
		expect(await screen.findByText('projects page')).toBeInTheDocument();
	});

	it('navigates back to the page the user came from', async () => {
		login.mockResolvedValue(undefined);
		renderLogin([{ pathname: '/login', state: { from: '/account' } }]);

		await submit('a@b.com', 'secret123');

		expect(await screen.findByText('account page')).toBeInTheDocument();
	});

	it('shows the error message when login fails', async () => {
		login.mockRejectedValue(new Error('email or password not correct!'));
		renderLogin();

		await submit('a@b.com', 'wrong-pass');

		expect(await screen.findByRole('alert')).toHaveTextContent(
			'email or password not correct!'
		);
		expect(screen.queryByText('projects page')).not.toBeInTheDocument();
	});
});
