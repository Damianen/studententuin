import { useState } from 'react';
import type { FormEvent } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router';
import { ArrowLeft } from 'lucide-react';
import { toast } from 'sonner';
import DatabaseController from '@/controllers/database_controller';
import { ensureSubdomain, SUBDOMAIN_SUFFIX } from '@/lib/subdomain_lookup';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select';
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';

const VERSION_PLACEHOLDERS: Record<string, string> = {
	postgres: '16',
	mysql: '8.0',
	mongodb: '7',
};

export default function NewDatabase() {
	const navigate = useNavigate();
	const [searchParams] = useSearchParams();
	const [subdomain, setSubdomain] = useState(
		searchParams.get('subdomain') ?? ''
	);
	const [type, setType] = useState('postgres');
	const [error, setError] = useState<string | null>(null);
	const [pending, setPending] = useState(false);

	const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
		event.preventDefault();
		const formData = new FormData(event.currentTarget);
		setError(null);
		setPending(true);

		try {
			const project = await ensureSubdomain(subdomain);
			await DatabaseController.create(project.id, {
				name: String(formData.get('name') ?? '').trim(),
				type,
				version: String(formData.get('version') ?? '').trim(),
				db_name: String(formData.get('db_name') ?? '').trim(),
				db_password: String(formData.get('db_password') ?? '').trim(),
			});
			toast.success('Database planted');
			navigate(`/projects/${project.id}`);
		} catch (err) {
			setError(
				err instanceof Error ? err.message : 'Failed to create database'
			);
		} finally {
			setPending(false);
		}
	};

	return (
		<div className="mx-auto w-full max-w-2xl flex-1 px-6 py-10">
			<Link
				to="/projects"
				className="mb-6 inline-flex items-center gap-1.5 font-mono text-[11px] uppercase tracking-[0.2em] text-muted-foreground transition-colors hover:text-primary"
			>
				<ArrowLeft className="size-3.5" />
				Your garden
			</Link>

			<div className="mb-8">
				<p className="mb-3 font-mono text-xs uppercase tracking-[0.22em] text-primary">
					New planting — database
				</p>
				<h1 className="font-display text-4xl font-semibold tracking-tight">
					Plant a <em className="italic text-primary">database</em>
				</h1>
				<p className="mt-3 text-muted-foreground">
					Provision a managed database on your own subdomain.
				</p>
			</div>

			<Card>
				<CardHeader>
					<CardTitle className="font-mono text-xs uppercase tracking-[0.22em] text-primary">
						Database details
					</CardTitle>
					<CardDescription>
						Pick an engine and we take care of running it.
					</CardDescription>
				</CardHeader>
				<CardContent>
					<form onSubmit={handleSubmit} className="space-y-6">
						<div className="space-y-2">
							<Label htmlFor="name">Name</Label>
							<Input id="name" name="name" placeholder="my-database" required />
						</div>

						<div className="space-y-2">
							<Label htmlFor="subdomain">Subdomain</Label>
							<div className="flex items-center gap-2">
								<Input
									id="subdomain"
									value={subdomain}
									onChange={(event) => setSubdomain(event.target.value)}
									placeholder="mydb"
									required
									pattern="^[a-z0-9\-]+$"
									title="Only lowercase letters, numbers, and hyphens allowed"
									className="flex-1 font-mono"
								/>
								<span className="text-sm text-muted-foreground">
									.{SUBDOMAIN_SUFFIX}
								</span>
							</div>
						</div>

						<div className="grid gap-6 sm:grid-cols-2">
							<div className="space-y-2">
								<Label>Engine</Label>
								<Select value={type} onValueChange={setType}>
									<SelectTrigger className="w-full">
										<SelectValue placeholder="Select an engine…" />
									</SelectTrigger>
									<SelectContent>
										<SelectItem value="postgres">PostgreSQL</SelectItem>
										<SelectItem value="mysql">MySQL</SelectItem>
										<SelectItem value="mongodb">MongoDB</SelectItem>
									</SelectContent>
								</Select>
							</div>
							<div className="space-y-2">
								<Label htmlFor="version">Version</Label>
								<Input
									id="version"
									name="version"
									placeholder={VERSION_PLACEHOLDERS[type]}
									required
								/>
							</div>
						</div>

						<div className="space-y-2">
							<Label htmlFor="db_name">Database name</Label>
							<Input id="db_name" name="db_name" placeholder="app_db" required />
						</div>

						<div className="space-y-2">
							<Label htmlFor="db_password">Database password</Label>
							<Input
								id="db_password"
								name="db_password"
								type="password"
								autoComplete="new-password"
								required
							/>
						</div>

						{error && (
							<p className="text-sm text-destructive" role="alert">
								{error}
							</p>
						)}

						<div className="flex justify-end gap-2">
							<Button type="button" variant="ghost" asChild>
								<Link to="/projects">Cancel</Link>
							</Button>
							<Button type="submit" disabled={pending}>
								{pending ? 'Planting…' : 'Create database'}
							</Button>
						</div>
					</form>
				</CardContent>
			</Card>
		</div>
	);
}
