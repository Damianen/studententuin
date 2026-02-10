import { useState } from 'react';
import { useNavigate } from 'react-router';
import { useAuth } from '../contexts/auth_context';

export function LoginPage() {
	const navigate = useNavigate();
	const [email, setEmail] = useState('');
	const [password, setPassword] = useState('');
	const [error, setError] = useState('');
	const [loading, setLoading] = useState(false);
	const { login } = useAuth();

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		setLoading(true);

		try {
			await login({ email, password });
			navigate('/');
		} catch (err) {
			setError(err instanceof Error ? err.message : 'inloggen mislukt');
		} finally {
			setLoading(false);
		}
	};

	return (
		<div>
			<h1>Login</h1>
			{error && <p>{error}</p>}
			<form onSubmit={handleSubmit}>
				<div>
					<label htmlFor="email">Email</label>
					<br />
					<input
						id="email"
						type="email"
						value={email}
						onChange={(e) => setEmail(e.target.value)}
						required
					/>
				</div>
				<div>
					<label htmlFor="password">Wachtwoord</label>
					<br />
					<input
						id="password"
						type="password"
						value={password}
						onChange={(e) => setPassword(e.target.value)}
						required
					/>
				</div>
				<button type="submit" disabled={loading}>
					{loading ? 'Laden...' : 'Inloggen'}
				</button>
			</form>
		</div>
	);
}
