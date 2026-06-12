import { useState } from 'react';
import type { FormEvent } from 'react';
import { Link, useNavigate, useParams } from 'react-router';
import {
	AppWindow,
	ArrowLeft,
	ArrowUpRight,
	Check,
	Database,
	ExternalLink,
	Eye,
	EyeOff,
	GitBranch,
	Pencil,
	Plus,
	Trash2,
	X,
} from 'lucide-react';
import { toast } from 'sonner';
import SubdomainController from '@/controllers/subdomain_controller';
import ApplicationController from '@/controllers/application_controller';
import DatabaseController from '@/controllers/database_controller';
import type { SubdomainListItemDto } from '@/dtos/subdomain_dtos';
import type { ApplicationDto } from '@/dtos/application_dtos';
import type { DatabaseDto } from '@/dtos/database_dtos';
import { useAsync } from '@/hooks/use_async';
import { StatusBadge } from '@/components/status-badge';
import { CopyButton } from '@/components/copy-button';
import { ConfirmDialog } from '@/components/confirm-dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Skeleton } from '@/components/ui/skeleton';
import { Separator } from '@/components/ui/separator';
import {
	Card,
	CardAction,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';

function errorMessage(err: unknown): string {
	return err instanceof Error ? err.message : 'Something went wrong';
}

function ProjectHeader({
	subdomain,
	onChanged,
}: {
	subdomain: SubdomainListItemDto;
	onChanged: () => Promise<void>;
}) {
	const [editing, setEditing] = useState(false);
	const [name, setName] = useState(subdomain.name);
	const [pending, setPending] = useState(false);

	const handleRename = async (event: FormEvent) => {
		event.preventDefault();
		if (!name.trim() || name.trim() === subdomain.name) {
			setEditing(false);
			return;
		}
		setPending(true);
		try {
			await SubdomainController.patch(subdomain.id, { name: name.trim() });
			await onChanged();
			toast.success('Project renamed');
			setEditing(false);
		} catch (err) {
			toast.error(errorMessage(err));
		} finally {
			setPending(false);
		}
	};

	return (
		<div className="space-y-3">
			<p className="font-mono text-xs uppercase tracking-[0.22em] text-primary">
				Plot record
			</p>
			{editing ? (
				<form onSubmit={handleRename} className="flex items-center gap-2">
					<Input
						value={name}
						onChange={(event) => setName(event.target.value)}
						className="max-w-xs text-lg font-semibold"
						autoFocus
					/>
					<Button type="submit" size="icon" disabled={pending}>
						<Check className="size-4" />
					</Button>
					<Button
						type="button"
						size="icon"
						variant="ghost"
						onClick={() => {
							setName(subdomain.name);
							setEditing(false);
						}}
					>
						<X className="size-4" />
					</Button>
				</form>
			) : (
				<div className="flex items-center gap-2">
					<h1 className="font-display text-4xl font-semibold tracking-tight">
						{subdomain.name}
					</h1>
					<Button
						variant="ghost"
						size="icon"
						className="text-muted-foreground"
						onClick={() => setEditing(true)}
						aria-label="Rename project"
					>
						<Pencil className="size-4" />
					</Button>
				</div>
			)}
			<div className="flex flex-wrap items-center gap-2 text-sm text-muted-foreground">
				<span className="font-mono text-xs">{subdomain.fullDomain}</span>
				<CopyButton value={subdomain.fullDomain} label="Copy domain" />
				<a
					href={`https://${subdomain.fullDomain}`}
					target="_blank"
					rel="noreferrer"
					className="inline-flex items-center gap-1 text-xs text-primary underline-offset-4 hover:underline"
				>
					Visit <ExternalLink className="size-3" />
				</a>
				<span
					className={`ml-1 inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-medium ${
						subdomain.isActive ? '' : 'text-muted-foreground'
					}`}
				>
					<span
						className={`size-1.5 rounded-full ${
							subdomain.isActive ? 'bg-status-running' : 'bg-status-stopped'
						}`}
					/>
					{subdomain.isActive ? 'Active' : 'Inactive'}
				</span>
			</div>
		</div>
	);
}

function ApplicationSection({
	subdomain,
	application,
	onChanged,
}: {
	subdomain: SubdomainListItemDto;
	application?: ApplicationDto;
	onChanged: () => Promise<void>;
}) {
	const [editing, setEditing] = useState(false);
	const [pending, setPending] = useState(false);
	const [form, setForm] = useState({
		name: application?.name ?? '',
		repo_url: application?.repo_url ?? '',
		branch: application?.branch ?? '',
		build_command: '',
		start_command: '',
	});

	if (!application) {
		return (
			<Link
				to={`/projects/new/app?subdomain=${encodeURIComponent(subdomain.name)}`}
				className="flex items-center justify-center gap-2 rounded-2xl border border-dashed px-6 py-10 text-sm text-muted-foreground transition-colors hover:border-primary/40 hover:text-primary"
			>
				<Plus className="size-4" />
				Add an application to this project
			</Link>
		);
	}

	const handleSave = async (event: FormEvent) => {
		event.preventDefault();
		setPending(true);
		try {
			await ApplicationController.patch(subdomain.id, application.id, {
				name: form.name.trim() || undefined,
				repo_url: form.repo_url.trim() || undefined,
				branch: form.branch.trim() || undefined,
				build_command: form.build_command.trim() || undefined,
				start_command: form.start_command.trim() || undefined,
			});
			await onChanged();
			toast.success('Application updated');
			setEditing(false);
		} catch (err) {
			toast.error(errorMessage(err));
		} finally {
			setPending(false);
		}
	};

	const handleDelete = async () => {
		try {
			await ApplicationController.delete(subdomain.id, application.id);
			await onChanged();
			toast.success('Application deleted');
		} catch (err) {
			toast.error(errorMessage(err));
		}
	};

	return (
		<Card>
			<CardHeader>
				<CardTitle className="flex items-center gap-2 font-mono text-xs uppercase tracking-[0.22em] text-primary">
					<AppWindow className="size-3.5" />
					Application
				</CardTitle>
				<CardDescription className="font-display text-lg tracking-tight text-foreground">
					{application.name}{' '}
					<span className="font-sans text-sm text-muted-foreground">
						· {application.type}
					</span>
				</CardDescription>
				<CardAction className="flex items-center gap-1">
					<StatusBadge status={application.status} />
					<Button variant="outline" size="sm" asChild className="ml-1">
						<Link to={`/projects/${subdomain.id}/app`}>
							Open
							<ArrowUpRight className="size-3.5" />
						</Link>
					</Button>
					<Button
						variant="ghost"
						size="icon"
						className="text-muted-foreground"
						onClick={() => setEditing((value) => !value)}
						aria-label="Edit application"
					>
						<Pencil className="size-4" />
					</Button>
					<ConfirmDialog
						trigger={
							<Button
								variant="ghost"
								size="icon"
								className="text-muted-foreground hover:text-destructive"
								aria-label="Delete application"
							>
								<Trash2 className="size-4" />
							</Button>
						}
						title="Delete application"
						description={`This permanently removes "${application.name}" and its deployment from ${subdomain.fullDomain}.`}
						onConfirm={handleDelete}
					/>
				</CardAction>
			</CardHeader>
			<CardContent>
				{editing ? (
					<form onSubmit={handleSave} className="grid gap-4 sm:grid-cols-2">
						<div className="space-y-2">
							<Label htmlFor="app-name">Name</Label>
							<Input
								id="app-name"
								value={form.name}
								onChange={(event) =>
									setForm({ ...form, name: event.target.value })
								}
							/>
						</div>
						<div className="space-y-2">
							<Label htmlFor="app-repo">Repository URL</Label>
							<Input
								id="app-repo"
								value={form.repo_url}
								onChange={(event) =>
									setForm({ ...form, repo_url: event.target.value })
								}
							/>
						</div>
						<div className="space-y-2">
							<Label htmlFor="app-branch">Branch</Label>
							<Input
								id="app-branch"
								value={form.branch}
								onChange={(event) =>
									setForm({ ...form, branch: event.target.value })
								}
							/>
						</div>
						<div className="space-y-2">
							<Label htmlFor="app-build">Build command</Label>
							<Input
								id="app-build"
								value={form.build_command}
								placeholder="leave blank to keep current"
								onChange={(event) =>
									setForm({ ...form, build_command: event.target.value })
								}
							/>
						</div>
						<div className="space-y-2">
							<Label htmlFor="app-start">Start command</Label>
							<Input
								id="app-start"
								value={form.start_command}
								placeholder="leave blank to keep current"
								onChange={(event) =>
									setForm({ ...form, start_command: event.target.value })
								}
							/>
						</div>
						<div className="flex items-end justify-end gap-2">
							<Button
								type="button"
								variant="ghost"
								onClick={() => setEditing(false)}
							>
								Cancel
							</Button>
							<Button type="submit" disabled={pending}>
								{pending ? 'Saving…' : 'Save changes'}
							</Button>
						</div>
					</form>
				) : (
					<dl className="grid gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
						<div>
							<dt className="font-mono text-[10px] uppercase tracking-[0.25em] text-muted-foreground">
								Repository
							</dt>
							<dd className="mt-1.5 break-all">
								<a
									href={application.repo_url}
									target="_blank"
									rel="noreferrer"
									className="text-primary underline-offset-4 hover:underline"
								>
									{application.repo_url}
								</a>
							</dd>
						</div>
						<div>
							<dt className="font-mono text-[10px] uppercase tracking-[0.25em] text-muted-foreground">
								Branch
							</dt>
							<dd className="mt-1.5 flex items-center gap-1.5 font-mono text-xs">
								<GitBranch className="size-3.5 text-muted-foreground" />
								{application.branch}
							</dd>
						</div>
					</dl>
				)}
			</CardContent>
		</Card>
	);
}

function DatabaseSection({
	subdomain,
	database,
	onChanged,
}: {
	subdomain: SubdomainListItemDto;
	database?: DatabaseDto;
	onChanged: () => Promise<void>;
}) {
	const [editing, setEditing] = useState(false);
	const [pending, setPending] = useState(false);
	const [revealed, setRevealed] = useState(false);
	const [form, setForm] = useState({
		name: database?.name ?? '',
		version: database?.version ?? '',
	});

	if (!database) {
		return (
			<Link
				to={`/projects/new/database?subdomain=${encodeURIComponent(subdomain.name)}`}
				className="flex items-center justify-center gap-2 rounded-2xl border border-dashed px-6 py-10 text-sm text-muted-foreground transition-colors hover:border-primary/40 hover:text-primary"
			>
				<Plus className="size-4" />
				Add a database to this project
			</Link>
		);
	}

	const handleSave = async (event: FormEvent) => {
		event.preventDefault();
		setPending(true);
		try {
			await DatabaseController.patch(subdomain.id, database.id, {
				name: form.name.trim() || undefined,
				version: form.version.trim() || undefined,
			});
			await onChanged();
			toast.success('Database updated');
			setEditing(false);
		} catch (err) {
			toast.error(errorMessage(err));
		} finally {
			setPending(false);
		}
	};

	const handleDelete = async () => {
		try {
			await DatabaseController.delete(subdomain.id, database.id);
			await onChanged();
			toast.success('Database deleted');
		} catch (err) {
			toast.error(errorMessage(err));
		}
	};

	return (
		<Card>
			<CardHeader>
				<CardTitle className="flex items-center gap-2 font-mono text-xs uppercase tracking-[0.22em] text-primary">
					<Database className="size-3.5" />
					Database
				</CardTitle>
				<CardDescription className="font-display text-lg tracking-tight text-foreground">
					{database.name}{' '}
					<span className="font-sans text-sm text-muted-foreground">
						· {database.type} {database.version}
					</span>
				</CardDescription>
				<CardAction className="flex items-center gap-1">
					<StatusBadge status={database.status} />
					<Button variant="outline" size="sm" asChild className="ml-1">
						<Link to={`/projects/${subdomain.id}/db`}>
							Open
							<ArrowUpRight className="size-3.5" />
						</Link>
					</Button>
					<Button
						variant="ghost"
						size="icon"
						className="text-muted-foreground"
						onClick={() => setEditing((value) => !value)}
						aria-label="Edit database"
					>
						<Pencil className="size-4" />
					</Button>
					<ConfirmDialog
						trigger={
							<Button
								variant="ghost"
								size="icon"
								className="text-muted-foreground hover:text-destructive"
								aria-label="Delete database"
							>
								<Trash2 className="size-4" />
							</Button>
						}
						title="Delete database"
						description={`This permanently removes "${database.name}" and all of its data. This cannot be undone.`}
						onConfirm={handleDelete}
					/>
				</CardAction>
			</CardHeader>
			<CardContent>
				{editing ? (
					<form onSubmit={handleSave} className="grid gap-4 sm:grid-cols-2">
						<div className="space-y-2">
							<Label htmlFor="db-name">Name</Label>
							<Input
								id="db-name"
								value={form.name}
								onChange={(event) =>
									setForm({ ...form, name: event.target.value })
								}
							/>
						</div>
						<div className="space-y-2">
							<Label htmlFor="db-version">Version</Label>
							<Input
								id="db-version"
								value={form.version}
								onChange={(event) =>
									setForm({ ...form, version: event.target.value })
								}
							/>
						</div>
						<div className="flex items-end justify-end gap-2 sm:col-span-2">
							<Button
								type="button"
								variant="ghost"
								onClick={() => setEditing(false)}
							>
								Cancel
							</Button>
							<Button type="submit" disabled={pending}>
								{pending ? 'Saving…' : 'Save changes'}
							</Button>
						</div>
					</form>
				) : database.connection_string ? (
					<div className="space-y-1.5">
						<span className="font-mono text-[10px] uppercase tracking-[0.25em] text-muted-foreground">
							Connection string
						</span>
						<div className="flex items-center gap-1 rounded-lg bg-muted/60 px-3 py-2">
							<code className="min-w-0 flex-1 truncate font-mono text-xs">
								{revealed
									? database.connection_string
									: '••••••••••••••••••••••••••••••••'}
							</code>
							<Button
								variant="ghost"
								size="icon"
								className="size-7 text-muted-foreground"
								onClick={() => setRevealed((value) => !value)}
								aria-label={
									revealed ? 'Hide connection string' : 'Show connection string'
								}
							>
								{revealed ? (
									<EyeOff className="size-3.5" />
								) : (
									<Eye className="size-3.5" />
								)}
							</Button>
							<CopyButton
								value={database.connection_string}
								label="Copy connection string"
							/>
						</div>
					</div>
				) : (
					<p className="text-sm text-muted-foreground">
						The connection string will appear here once the database is
						provisioned.
					</p>
				)}
			</CardContent>
		</Card>
	);
}

export default function ProjectDetail() {
	const { id } = useParams<{ id: string }>();
	const navigate = useNavigate();
	const {
		data: subdomains,
		loading,
		error,
		reload,
	} = useAsync(SubdomainController.getAll);

	const subdomain = subdomains?.find((item) => item.id === id);

	const handleDeleteProject = async () => {
		if (!subdomain) return;
		try {
			await SubdomainController.delete(subdomain.id);
			toast.success('Project deleted');
			navigate('/projects');
		} catch (err) {
			toast.error(errorMessage(err));
		}
	};

	if (loading) {
		return (
			<div className="mx-auto w-full max-w-4xl flex-1 space-y-6 px-6 py-10">
				<Skeleton className="h-9 w-56" />
				<Skeleton className="h-5 w-72" />
				<Skeleton className="h-44 w-full rounded-2xl" />
				<Skeleton className="h-44 w-full rounded-2xl" />
			</div>
		);
	}

	if (error) {
		return (
			<div className="flex flex-1 flex-col items-center justify-center gap-4 px-6 py-24 text-center">
				<p className="text-muted-foreground">
					Could not load this project: {error}
				</p>
				<Button variant="outline" onClick={reload}>
					Try again
				</Button>
			</div>
		);
	}

	if (!subdomain) {
		return (
			<div className="flex flex-1 flex-col items-center justify-center gap-4 px-6 py-24 text-center">
				<p className="font-mono text-xs uppercase tracking-[0.22em] text-primary">
					Fallow ground
				</p>
				<h1 className="font-display text-3xl font-semibold tracking-tight">
					Project not <em className="italic text-primary">found</em>
				</h1>
				<p className="text-muted-foreground">
					It may have been deleted, or the link is wrong.
				</p>
				<Button variant="outline" asChild>
					<Link to="/projects">
						<ArrowLeft className="size-4" />
						Back to projects
					</Link>
				</Button>
			</div>
		);
	}

	return (
		<div className="mx-auto w-full max-w-4xl flex-1 space-y-8 px-6 py-10">
			<Link
				to="/projects"
				className="inline-flex items-center gap-1.5 font-mono text-[11px] uppercase tracking-[0.2em] text-muted-foreground transition-colors hover:text-primary"
			>
				<ArrowLeft className="size-3.5" />
				Your garden
			</Link>

			<ProjectHeader subdomain={subdomain} onChanged={reload} />

			<div className="space-y-5">
				<ApplicationSection
					key={`app-${subdomain.application?.id ?? 'none'}`}
					subdomain={subdomain}
					application={subdomain.application}
					onChanged={reload}
				/>
				<DatabaseSection
					key={`db-${subdomain.database?.id ?? 'none'}`}
					subdomain={subdomain}
					database={subdomain.database}
					onChanged={reload}
				/>
			</div>

			<Separator />

			<div className="flex flex-wrap items-center justify-between gap-4 rounded-2xl border border-destructive/30 bg-destructive/5 px-5 py-4">
				<div>
					<h2 className="font-mono text-xs uppercase tracking-[0.2em] text-destructive">
						Delete this project
					</h2>
					<p className="text-sm text-muted-foreground">
						Removes the subdomain and everything growing on it.
					</p>
				</div>
				<ConfirmDialog
					trigger={<Button variant="destructive">Delete project</Button>}
					title="Delete project"
					description={`This permanently removes ${subdomain.fullDomain}, including any application or database on it. This cannot be undone.`}
					confirmText={subdomain.name}
					onConfirm={handleDeleteProject}
				/>
			</div>
		</div>
	);
}
