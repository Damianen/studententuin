import { lazy, Suspense, useEffect } from 'react';
import { Outlet, Route, Routes, useLocation } from 'react-router';
import { AuthProvider } from '@/contexts/auth_context';
import { AppHeader } from '@/components/layout/app-header';
import { AppFooter } from '@/components/layout/app-footer';
import { RequireAuth } from '@/components/auth/require-auth';
import { PublicOnly } from '@/components/auth/public-only';
import { Toaster } from '@/components/ui/sonner';

const Landing = lazy(() => import('@/pages/landing'));
const Login = lazy(() => import('@/pages/login'));
const Register = lazy(() => import('@/pages/register'));
const Account = lazy(() => import('@/pages/account'));
const NotFound = lazy(() => import('@/pages/not-found'));
const Projects = lazy(() => import('@/pages/projects/projects'));
const ProjectDetail = lazy(() => import('@/pages/projects/project-detail'));
const NewApplication = lazy(() => import('@/pages/projects/new-application'));
const NewDatabase = lazy(() => import('@/pages/projects/new-database'));

function ScrollToTop() {
	const { pathname } = useLocation();
	useEffect(() => {
		window.scrollTo(0, 0);
	}, [pathname]);
	return null;
}

function AppLayout() {
	return (
		<div className="relative flex min-h-svh flex-col overflow-x-clip bg-background text-foreground">
			{/* The same moonlit greenhouse the landing page lives in. */}
			<div
				aria-hidden
				className="pointer-events-none absolute inset-x-0 top-0 h-[40rem] bg-[radial-gradient(70rem_34rem_at_30%_-14%,oklch(0.38_0.07_145/0.35),transparent)]"
			/>
			<div aria-hidden className="grain-overlay" />
			<AppHeader />
			<main className="relative z-10 flex flex-1 flex-col">
				<Suspense fallback={<div className="flex-1" />}>
					<Outlet />
				</Suspense>
			</main>
			<AppFooter />
		</div>
	);
}

export default function App() {
	return (
		<AuthProvider>
			<ScrollToTop />
			<Routes>
				{/* The landing page is its own full-bleed world — no shared chrome. */}
				<Route
					index
					element={
						<Suspense fallback={null}>
							<Landing />
						</Suspense>
					}
				/>
				<Route element={<AppLayout />}>
					<Route element={<PublicOnly />}>
						<Route path="login" element={<Login />} />
						<Route path="register" element={<Register />} />
					</Route>
					<Route element={<RequireAuth />}>
						<Route path="projects" element={<Projects />} />
						<Route path="projects/new/app" element={<NewApplication />} />
						<Route path="projects/new/database" element={<NewDatabase />} />
						<Route path="projects/:id" element={<ProjectDetail />} />
						<Route path="account" element={<Account />} />
					</Route>
					<Route path="*" element={<NotFound />} />
				</Route>
			</Routes>
			<Toaster position="bottom-right" />
		</AuthProvider>
	);
}
