import { useAuth } from '../contexts/auth_context';
import { useNavigate } from 'react-router';

export function HomePage() {
	const { user, loading, isAuthenticated, logout } = useAuth();
	const navigate = useNavigate();

	const handleLogout = async () => {
		try {
			await logout();
			navigate('/login');
		} catch (err) {
			console.error(err);
		}
	};

	if (loading) {
		return <p>Laden...</p>;
	}

	if (!isAuthenticated) {
		return (
			<div>
				<h1>Home</h1>
				<p>Je bent niet ingelogd.</p>
				<button onClick={() => navigate('/login')}>Inloggen</button>
			</div>
		);
	}

	return (
		<div>
			<h1>Home</h1>
			<p>Welkom, {user!.email}!</p>
			<button onClick={handleLogout}>Uitloggen</button>
		</div>
	);
}
