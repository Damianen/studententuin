import { useState } from 'react';
import type { FormEvent } from 'react';
import { Link, useNavigate } from 'react-router';
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

export default function Register() {
	const { register } = useAuth();
	const navigate = useNavigate();
	const [name, setName] = useState('');
	const [email, setEmail] = useState('');
	const [password, setPassword] = useState('');
	const [confirm, setConfirm] = useState('');
	const [error, setError] = useState<string | null>(null);
	const [pending, setPending] = useState(false);

	const handleSubmit = async (event: FormEvent) => {
		event.preventDefault();
		setError(null);

		if (password !== confirm) {
			setError('Passwords do not match.');
			return;
		}
		if (password.length < 8) {
			setError('Password must be at least 8 characters.');
			return;
		}

		setPending(true);
		try {
			await register({ name, email, password });
			navigate('/projects', { replace: true });
		} catch (err) {
			setError(
				err instanceof Error ? err.message : 'Registration failed. Try again.'
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
						Claim your plot
					</p>
					<CardTitle className="font-display text-3xl font-semibold tracking-tight">
						Plant your first <em className="italic text-primary">project</em>
					</CardTitle>
					<CardDescription>
						Free hosting for students — claim your own corner of the garden.
					</CardDescription>
				</CardHeader>
				<CardContent>
					<form onSubmit={handleSubmit} className="space-y-5">
						<div className="space-y-2">
							<Label htmlFor="name">Name</Label>
							<Input
								id="name"
								autoComplete="name"
								required
								value={name}
								onChange={(event) => setName(event.target.value)}
								placeholder="Jan Jansen"
							/>
						</div>
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
								autoComplete="new-password"
								required
								minLength={8}
								value={password}
								onChange={(event) => setPassword(event.target.value)}
							/>
						</div>
						<div className="space-y-2">
							<Label htmlFor="confirm">Confirm password</Label>
							<Input
								id="confirm"
								type="password"
								autoComplete="new-password"
								required
								value={confirm}
								onChange={(event) => setConfirm(event.target.value)}
							/>
						</div>
						{error && (
							<p className="text-sm text-destructive" role="alert">
								{error}
							</p>
						)}
						<Button type="submit" className="w-full" disabled={pending}>
							{pending ? 'Creating account…' : 'Create account'}
						</Button>
					</form>
					<p className="mt-6 text-center text-sm text-muted-foreground">
						Already have an account?{' '}
						<Link
							to="/login"
							className="font-medium text-primary underline-offset-4 hover:underline"
						>
							Sign in
						</Link>
					</p>
				</CardContent>
			</Card>
		</div>
	);
}
