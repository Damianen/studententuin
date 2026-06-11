import { useState } from 'react';
import type { FormEvent } from 'react';
import { useNavigate } from 'react-router';
import { toast } from 'sonner';
import { useAuth } from '@/contexts/auth_context';
import { ConfirmDialog } from '@/components/confirm-dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';

export default function Account() {
	const { user, updateUser, deleteUser } = useAuth();
	const navigate = useNavigate();
	const [name, setName] = useState(user?.name ?? '');
	const [email, setEmail] = useState(user?.email ?? '');
	const [pending, setPending] = useState(false);

	const handleSave = async (event: FormEvent) => {
		event.preventDefault();
		setPending(true);
		try {
			await updateUser({
				name: name.trim() !== user?.name ? name.trim() : undefined,
				email: email.trim() !== user?.email ? email.trim() : undefined,
			});
			toast.success('Account updated');
		} catch (err) {
			toast.error(
				err instanceof Error ? err.message : 'Updating your account failed'
			);
		} finally {
			setPending(false);
		}
	};

	const handleDelete = async () => {
		try {
			await deleteUser();
			toast.success('Your account has been deleted');
			navigate('/');
		} catch (err) {
			toast.error(
				err instanceof Error ? err.message : 'Deleting your account failed'
			);
		}
	};

	return (
		<div className="mx-auto w-full max-w-2xl flex-1 space-y-8 px-6 py-10">
			<div>
				<p className="mb-3 font-mono text-xs uppercase tracking-[0.22em] text-primary">
					The gardener
				</p>
				<h1 className="font-display text-4xl font-semibold tracking-tight">
					Your <em className="italic text-primary">account</em>
				</h1>
				<p className="mt-3 text-muted-foreground">
					Manage your profile and account settings.
				</p>
			</div>

			<Card>
				<CardHeader>
					<CardTitle className="font-mono text-xs uppercase tracking-[0.22em] text-primary">
						Profile
					</CardTitle>
					<CardDescription className="flex items-center gap-2">
						Status:
						<Badge variant="secondary" className="capitalize">
							{user?.status ?? 'unknown'}
						</Badge>
					</CardDescription>
				</CardHeader>
				<CardContent>
					<form onSubmit={handleSave} className="space-y-5">
						<div className="space-y-2">
							<Label htmlFor="name">Name</Label>
							<Input
								id="name"
								value={name}
								onChange={(event) => setName(event.target.value)}
								required
							/>
						</div>
						<div className="space-y-2">
							<Label htmlFor="email">Email</Label>
							<Input
								id="email"
								type="email"
								value={email}
								onChange={(event) => setEmail(event.target.value)}
								required
							/>
						</div>
						<div className="flex justify-end">
							<Button type="submit" disabled={pending}>
								{pending ? 'Saving…' : 'Save changes'}
							</Button>
						</div>
					</form>
				</CardContent>
			</Card>

			<div className="flex flex-wrap items-center justify-between gap-4 rounded-2xl border border-destructive/30 bg-destructive/5 px-5 py-4">
				<div>
					<h2 className="font-mono text-xs uppercase tracking-[0.2em] text-destructive">
						Delete account
					</h2>
					<p className="text-sm text-muted-foreground">
						Permanently removes your account and all of your projects.
					</p>
				</div>
				<ConfirmDialog
					trigger={<Button variant="destructive">Delete account</Button>}
					title="Delete account"
					description="This permanently removes your account, subdomains, applications, and databases. This cannot be undone."
					confirmText={user?.email}
					confirmLabel="Delete my account"
					onConfirm={handleDelete}
				/>
			</div>
		</div>
	);
}
