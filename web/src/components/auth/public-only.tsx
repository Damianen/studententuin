import { Navigate, Outlet } from 'react-router';
import { useAuth } from '@/contexts/auth_context';

export function PublicOnly() {
	const { isAuthenticated, loading } = useAuth();

	if (loading) {
		return null;
	}

	if (isAuthenticated) {
		return <Navigate to="/projects" replace />;
	}

	return <Outlet />;
}
