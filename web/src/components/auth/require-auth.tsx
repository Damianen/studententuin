import { Navigate, Outlet, useLocation } from 'react-router';
import { useAuth } from '@/contexts/auth_context';
import { Skeleton } from '@/components/ui/skeleton';

export function RequireAuth() {
	const { isAuthenticated, loading } = useAuth();
	const location = useLocation();

	if (loading) {
		return (
			<div className="mx-auto w-full max-w-5xl space-y-4 px-6 py-12">
				<Skeleton className="h-9 w-48" />
				<Skeleton className="h-40 w-full" />
				<Skeleton className="h-40 w-full" />
			</div>
		);
	}

	if (!isAuthenticated) {
		return <Navigate to="/login" state={{ from: location.pathname }} replace />;
	}

	return <Outlet />;
}
