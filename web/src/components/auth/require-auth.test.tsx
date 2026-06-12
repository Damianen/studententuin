import { render, screen } from '@testing-library/react';
import { MemoryRouter, Route, Routes, useLocation } from 'react-router';
import { RequireAuth } from './require-auth';
import { useAuth } from '@/contexts/auth_context';

vi.mock('@/contexts/auth_context', () => ({
	useAuth: vi.fn(),
}));

function LoginProbe() {
	const location = useLocation();
	const from = (location.state as { from?: string } | null)?.from;
	return <div>login page, from: {from}</div>;
}

function renderGuard() {
	return render(
		<MemoryRouter initialEntries={['/projects']}>
			<Routes>
				<Route element={<RequireAuth />}>
					<Route path="/projects" element={<div>protected content</div>} />
				</Route>
				<Route path="/login" element={<LoginProbe />} />
			</Routes>
		</MemoryRouter>
	);
}

const baseAuth = {
	user: null,
	loading: false,
	isAuthenticated: false,
	login: vi.fn(),
	logout: vi.fn(),
	register: vi.fn(),
	updateUser: vi.fn(),
	deleteUser: vi.fn(),
};

describe('RequireAuth', () => {
	it('renders the outlet when authenticated', () => {
		vi.mocked(useAuth).mockReturnValue({
			...baseAuth,
			isAuthenticated: true,
			user: { email: 'a@b.com', name: 'Alice', status: 'active' },
		});

		renderGuard();

		expect(screen.getByText('protected content')).toBeInTheDocument();
	});

	it('redirects to /login with the origin in state when anonymous', () => {
		vi.mocked(useAuth).mockReturnValue(baseAuth);

		renderGuard();

		expect(screen.queryByText('protected content')).not.toBeInTheDocument();
		expect(screen.getByText('login page, from: /projects')).toBeInTheDocument();
	});

	it('shows skeletons (not a redirect) while auth state is loading', () => {
		vi.mocked(useAuth).mockReturnValue({ ...baseAuth, loading: true });

		renderGuard();

		expect(screen.queryByText('protected content')).not.toBeInTheDocument();
		expect(screen.queryByText(/login page/)).not.toBeInTheDocument();
	});
});
