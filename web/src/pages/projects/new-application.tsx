import { useState } from 'react';
import type { FormEvent } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router';
import { ArrowLeft } from 'lucide-react';
import { toast } from 'sonner';
import ApplicationController from '@/controllers/application_controller';
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

export default function NewApplication() {
	const navigate = useNavigate();
	const [searchParams] = useSearchParams();
	const [subdomain, setSubdomain] = useState(
		searchParams.get('subdomain') ?? ''
	);
	const [type, setType] = useState('Nodejs');
	const [error, setError] = useState<string | null>(null);
	const [pending, setPending] = useState(false);

	const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
		event.preventDefault();
		const formData = new FormData(event.currentTarget);
		setError(null);
		setPending(true);

		try {
			const project = await ensureSubdomain(subdomain);
			await ApplicationController.create(project.id, {
				name: String(formData.get('name') ?? '').trim(),
				type,
				repo_url: String(formData.get('repo_url') ?? '').trim(),
				branch: String(formData.get('branch') ?? 'main').trim() || 'main',
				build_command: String(formData.get('build_command') ?? '').trim(),
				start_command: String(formData.get('start_command') ?? '').trim(),
			});
			toast.success('Application planted');
			navigate(`/projects/${project.id}`);
		} catch (err) {
			setError(
				err instanceof Error ? err.message : 'Failed to create application'
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
					New planting — application
				</p>
				<h1 className="font-display text-4xl font-semibold tracking-tight">
					Plant an <em className="italic text-primary">application</em>
				</h1>
				<p className="mt-3 text-muted-foreground">
					Deploy an application from a Git repository onto your own subdomain.
				</p>
			</div>

			<Card>
				<CardHeader>
					<CardTitle className="font-mono text-xs uppercase tracking-[0.22em] text-primary">
						Application details
					</CardTitle>
					<CardDescription>
						We pull the repository, build it, and serve it on your domain.
					</CardDescription>
				</CardHeader>
				<CardContent>
					<form onSubmit={handleSubmit} className="space-y-6">
						<div className="space-y-2">
							<Label htmlFor="name">Name</Label>
							<Input id="name" name="name" placeholder="my-app" required />
						</div>

						<div className="space-y-2">
							<Label htmlFor="subdomain">Subdomain</Label>
							<div className="flex items-center gap-2">
								<Input
									id="subdomain"
									value={subdomain}
									onChange={(event) => setSubdomain(event.target.value)}
									placeholder="myapp"
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

						<div className="space-y-2">
							<Label>Runtime</Label>
							<Select value={type} onValueChange={setType}>
								<SelectTrigger className="w-full">
									<SelectValue placeholder="Select a runtime…" />
								</SelectTrigger>
								<SelectContent>
									<SelectItem value="Nodejs">Node.js</SelectItem>
									<SelectItem value="dotnet" disabled>
										.NET (coming soon)
									</SelectItem>
								</SelectContent>
							</Select>
						</div>

						<div className="space-y-2">
							<Label htmlFor="repo_url">Repository URL</Label>
							<Input
								id="repo_url"
								name="repo_url"
								type="url"
								placeholder="https://github.com/you/my-app"
								required
							/>
						</div>

						<div className="grid gap-6 sm:grid-cols-2">
							<div className="space-y-2">
								<Label htmlFor="branch">Branch</Label>
								<Input id="branch" name="branch" placeholder="main" />
							</div>
							<div className="space-y-2">
								<Label htmlFor="build_command">Build command</Label>
								<Input
									id="build_command"
									name="build_command"
									placeholder="npm run build"
								/>
							</div>
						</div>

						<div className="space-y-2">
							<Label htmlFor="start_command">Start command</Label>
							<Input
								id="start_command"
								name="start_command"
								placeholder="npm start"
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
								{pending ? 'Planting…' : 'Create application'}
							</Button>
						</div>
					</form>
				</CardContent>
			</Card>
		</div>
	);
}
