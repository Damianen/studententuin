import { render, screen } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router';
import { PublicOnly } from './public-only';
import { useAuth } from '@/contexts/auth_context';

vi.mock('@/contexts/auth_context', () => ({
	useAuth: vi.fn(),
}));

function renderGuard() {
	return render(
		<MemoryRouter initialEntries={['/login']}>
			<Routes>
				<Route element={<PublicOnly />}>
					<Route path="/login" element={<div>login form</div>} />
				</Route>
				<Route path="/projects" element={<div>projects page</div>} />
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

describe('PublicOnly', () => {
	it('renders the outlet when anonymous', () => {
		vi.mocked(useAuth).mockReturnValue(baseAuth);

		renderGuard();

		expect(screen.getByText('login form')).toBeInTheDocument();
	});

	it('redirects to /projects when authenticated', () => {
		vi.mocked(useAuth).mockReturnValue({
			...baseAuth,
			isAuthenticated: true,
			user: { email: 'a@b.com', name: 'Alice', status: 'active' },
		});

		renderGuard();

		expect(screen.queryByText('login form')).not.toBeInTheDocument();
		expect(screen.getByText('projects page')).toBeInTheDocument();
	});

	it('renders nothing while auth state is loading', () => {
		vi.mocked(useAuth).mockReturnValue({ ...baseAuth, loading: true });

		renderGuard();

		expect(screen.queryByText('login form')).not.toBeInTheDocument();
		expect(screen.queryByText('projects page')).not.toBeInTheDocument();
	});
});
