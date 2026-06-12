import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { AuthProvider, useAuth } from './auth_context';
import UserController from '@/controllers/user_controller';
import AuthController from '@/controllers/auth_controller';

vi.mock('@/controllers/user_controller', () => ({
	default: {
		get: vi.fn(),
		register: vi.fn(),
		patch: vi.fn(),
		delete: vi.fn(),
	},
}));
vi.mock('@/controllers/auth_controller', () => ({
	default: { login: vi.fn(), logout: vi.fn() },
}));

const user = { email: 'a@b.com', name: 'Alice', status: 'active' };

function Probe() {
	const auth = useAuth();
	return (
		<div>
			<span data-testid="loading">{String(auth.loading)}</span>
			<span data-testid="user">{auth.user?.name ?? 'none'}</span>
			<button onClick={() => auth.login({ email: 'a@b.com', password: 'pw' })}>
				login
			</button>
			<button onClick={() => auth.logout()}>logout</button>
			<button
				onClick={() =>
					auth.register({ email: 'a@b.com', password: 'pw123456', name: 'Alice' })
				}
			>
				register
			</button>
			<button onClick={() => auth.deleteUser()}>delete</button>
		</div>
	);
}

function renderProbe() {
	return render(
		<AuthProvider>
			<Probe />
		</AuthProvider>
	);
}

beforeEach(() => {
	vi.mocked(UserController.get).mockReset();
	vi.mocked(AuthController.login).mockReset();
	vi.mocked(AuthController.logout).mockReset();
});

describe('AuthProvider', () => {
	it('loads the current user on mount', async () => {
		vi.mocked(UserController.get).mockResolvedValue(user);

		renderProbe();

		expect(screen.getByTestId('loading')).toHaveTextContent('true');
		await waitFor(() =>
			expect(screen.getByTestId('loading')).toHaveTextContent('false')
		);
		expect(screen.getByTestId('user')).toHaveTextContent('Alice');
	});

	it('stays unauthenticated when the user fetch fails', async () => {
		vi.mocked(UserController.get).mockRejectedValue(
			new Error('authentication required')
		);

		renderProbe();

		await waitFor(() =>
			expect(screen.getByTestId('loading')).toHaveTextContent('false')
		);
		expect(screen.getByTestId('user')).toHaveTextContent('none');
	});

	it('login fetches and sets the user', async () => {
		vi.mocked(UserController.get)
			.mockRejectedValueOnce(new Error('authentication required'))
			.mockResolvedValueOnce(user);
		vi.mocked(AuthController.login).mockResolvedValue(undefined);

		renderProbe();
		await waitFor(() =>
			expect(screen.getByTestId('loading')).toHaveTextContent('false')
		);

		await userEvent.click(screen.getByRole('button', { name: 'login' }));

		expect(AuthController.login).toHaveBeenCalledWith({
			email: 'a@b.com',
			password: 'pw',
		});
		await waitFor(() =>
			expect(screen.getByTestId('user')).toHaveTextContent('Alice')
		);
	});

	it('logout clears the user', async () => {
		vi.mocked(UserController.get).mockResolvedValue(user);
		vi.mocked(AuthController.logout).mockResolvedValue(undefined);

		renderProbe();
		await waitFor(() =>
			expect(screen.getByTestId('user')).toHaveTextContent('Alice')
		);

		await userEvent.click(screen.getByRole('button', { name: 'logout' }));

		await waitFor(() =>
			expect(screen.getByTestId('user')).toHaveTextContent('none')
		);
	});

	it('register creates the account and chains into login', async () => {
		vi.mocked(UserController.get)
			.mockRejectedValueOnce(new Error('authentication required'))
			.mockResolvedValueOnce(user);
		vi.mocked(UserController.register).mockResolvedValue(undefined);
		vi.mocked(AuthController.login).mockResolvedValue(undefined);

		renderProbe();
		await waitFor(() =>
			expect(screen.getByTestId('loading')).toHaveTextContent('false')
		);

		await userEvent.click(screen.getByRole('button', { name: 'register' }));

		expect(UserController.register).toHaveBeenCalledWith({
			email: 'a@b.com',
			password: 'pw123456',
			name: 'Alice',
		});
		expect(AuthController.login).toHaveBeenCalledWith({
			email: 'a@b.com',
			password: 'pw123456',
		});
		await waitFor(() =>
			expect(screen.getByTestId('user')).toHaveTextContent('Alice')
		);
	});

	it('deleteUser clears the user', async () => {
		vi.mocked(UserController.get).mockResolvedValue(user);
		vi.mocked(UserController.delete).mockResolvedValue(undefined);

		renderProbe();
		await waitFor(() =>
			expect(screen.getByTestId('user')).toHaveTextContent('Alice')
		);

		await userEvent.click(screen.getByRole('button', { name: 'delete' }));

		await waitFor(() =>
			expect(screen.getByTestId('user')).toHaveTextContent('none')
		);
	});

	it('useAuth throws outside the provider', () => {
		const spy = vi.spyOn(console, 'error').mockImplementation(() => {});
		expect(() => render(<Probe />)).toThrow(
			'useAuth can only be used within AuthProvider'
		);
		spy.mockRestore();
	});
});
