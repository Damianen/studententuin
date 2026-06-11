import { useState } from 'react';
import type { FormEvent } from 'react';
import { Link, useLocation, useNavigate } from 'react-router';
import { useAuth } from '@/contexts/auth_context';
import { LogoMark } from '@/components/logo';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';

export default function Login() {
	const { login } = useAuth();
	const navigate = useNavigate();
	const location = useLocation();
	const [email, setEmail] = useState('');
	const [password, setPassword] = useState('');
	const [error, setError] = useState<string | null>(null);
	const [pending, setPending] = useState(false);

	const from = (location.state as { from?: string } | null)?.from;

	const handleSubmit = async (event: FormEvent) => {
		event.preventDefault();
		setError(null);
		setPending(true);
		try {
			await login({ email, password });
			navigate(from ?? '/projects', { replace: true });
		} catch (err) {
			setError(
				err instanceof Error ? err.message : 'Signing in failed. Try again.'
			);
		} finally {
			setPending(false);
		}
	};

	return (
		<div className="flex flex-1 items-center justify-center px-6 py-16">
			<Card className="w-full max-w-md shadow-lifted">
				<CardHeader className="text-center">
					<LogoMark className="mx-auto mb-4 size-11" />
					<p className="mb-2 font-mono text-xs uppercase tracking-[0.22em] text-primary">
						The garden gate
					</p>
					<CardTitle className="font-display text-3xl font-semibold tracking-tight">
						Welcome <em className="italic text-primary">back</em>
					</CardTitle>
					<CardDescription>Sign in to tend to your projects.</CardDescription>
				</CardHeader>
				<CardContent>
					<form onSubmit={handleSubmit} className="space-y-5">
						<div className="space-y-2">
							<Label htmlFor="email">Email</Label>
							<Input
								id="email"
								type="email"
								autoComplete="email"
								required
								value={email}
								onChange={(event) => setEmail(event.target.value)}
								placeholder="you@university.nl"
							/>
						</div>
						<div className="space-y-2">
							<Label htmlFor="password">Password</Label>
							<Input
								id="password"
								type="password"
								autoComplete="current-password"
								required
								value={password}
								onChange={(event) => setPassword(event.target.value)}
							/>
						</div>
						{error && (
							<p className="text-sm text-destructive" role="alert">
								{error}
							</p>
						)}
						<Button type="submit" className="w-full" disabled={pending}>
							{pending ? 'Signing in…' : 'Sign in'}
						</Button>
					</form>
					<p className="mt-6 text-center text-sm text-muted-foreground">
						New here?{' '}
						<Link
							to="/register"
							className="font-medium text-primary underline-offset-4 hover:underline"
						>
							Create an account
						</Link>
					</p>
				</CardContent>
			</Card>
		</div>
	);
}
