import { useMemo, useState } from 'react';
import type { FormEvent } from 'react';
import { Link, useNavigate, useParams } from 'react-router';
import {
	AppWindow,
	ArrowLeft,
	ArrowUpRight,
	Database,
	ExternalLink,
	Eye,
	EyeOff,
	GitBranch,
	Play,
	Plus,
	Square,
} from 'lucide-react';
import { toast } from 'sonner';
import SubdomainController from '@/controllers/subdomain_controller';
import ApplicationController from '@/controllers/application_controller';
import DatabaseController from '@/controllers/database_controller';
import type { SubdomainListItemDto } from '@/dtos/subdomain_dtos';
import type { ApplicationDto } from '@/dtos/application_dtos';
import type { DatabaseDto } from '@/dtos/database_dtos';
import { useAsync } from '@/hooks/use_async';
import { cn } from '@/lib/utils';
import { defaultEnvVars, makeSeries, seededInt } from '@/lib/mock_telemetry';
import type { ResourceKind } from '@/lib/mock_telemetry';
import { headlineStats, overviewSpec } from '@/lib/metric_specs';
import { StatusBadge } from '@/components/status-badge';
import { CopyButton } from '@/components/copy-button';
import { ConfirmDialog } from '@/components/confirm-dialog';
import { AreaChart } from '@/components/resource/area-chart';
import { DemoBadge } from '@/components/resource/demo-badge';
import { LogsTerminal } from '@/components/resource/logs-terminal';
import { MetricsPanel } from '@/components/resource/metrics-panel';
import { DeploymentsList } from '@/components/resource/deployments-list';
import { EnvVarsCard } from '@/components/resource/env-vars';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Skeleton } from '@/components/ui/skeleton';
import {
	Card,
	CardAction,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';

const TAB_DEFS: Record<ResourceKind, Array<{ slug: string; label: string }>> = {
	application: [
		{ slug: '', label: 'Overview' },
		{ slug: 'logs', label: 'Logs' },
		{ slug: 'metrics', label: 'Metrics' },
		{ slug: 'deployments', label: 'Deployments' },
		{ slug: 'settings', label: 'Settings' },
	],
	database: [
		{ slug: '', label: 'Overview' },
		{ slug: 'logs', label: 'Logs' },
		{ slug: 'metrics', label: 'Metrics' },
		{ slug: 'settings', label: 'Settings' },
	],
};

function errorMessage(err: unknown): string {
	return err instanceof Error ? err.message : 'Something went wrong';
}

// Stop/Start for the app's container. Hidden while a deploy is pending —
// there is nothing sensible to drive until it settles.
function LifecycleControls({
	subdomainId,
	appId,
	status,
	onChanged,
}: {
	subdomainId: string;
	appId: string;
	status: string;
	onChanged: () => void;
}) {
	const [pending, setPending] = useState(false);
	const running = status.toLowerCase() === 'running';
	const startable = ['stopped', 'failed'].includes(status.toLowerCase());
	if (!running && !startable) {
		return null;
	}

	const act = async () => {
		setPending(true);
		try {
			if (running) {
				await ApplicationController.stop(subdomainId, appId);
				toast.success('Application stopped');
			} else {
				await ApplicationController.start(subdomainId, appId);
				toast.success('Application started');
			}
			onChanged();
		} catch (err) {
			toast.error(errorMessage(err));
		} finally {
			setPending(false);
		}
	};

	return (
		<Button
			variant="outline"
			size="sm"
			className="ml-auto"
			disabled={pending}
			onClick={act}
		>
			{running ? (
				<Square className="size-3.5" />
			) : (
				<Play className="size-3.5" />
			)}
			{running ? 'Stop' : 'Start'}
		</Button>
	);
}

function InfoTerm({ children }: { children: React.ReactNode }) {
	return (
		<dt className="font-mono text-[10px] uppercase tracking-[0.25em] text-muted-foreground">
			{children}
		</dt>
	);
}

function ConnectionString({ value }: { value: string }) {
	const [revealed, setRevealed] = useState(false);
	return (
		<div className="flex items-center gap-1 rounded-lg bg-muted/60 px-3 py-2">
			<code className="min-w-0 flex-1 truncate font-mono text-xs">
				{revealed ? value : '••••••••••••••••••••••••••••••••'}
			</code>
			<Button
				variant="ghost"
				size="icon"
				className="size-7 text-muted-foreground"
				onClick={() => setRevealed((current) => !current)}
				aria-label={revealed ? 'Hide connection string' : 'Show connection string'}
			>
				{revealed ? <EyeOff className="size-3.5" /> : <Eye className="size-3.5" />}
			</Button>
			<CopyButton value={value} label="Copy connection string" />
		</div>
	);
}

function OverviewTab({
	kind,
	subdomain,
	application,
	database,
	basePath,
}: {
	kind: ResourceKind;
	subdomain: SubdomainListItemDto;
	application?: ApplicationDto;
	database?: DatabaseDto;
	basePath: string;
}) {
	const resource = (kind === 'application' ? application : database)!;
	const chartSpec = overviewSpec(kind);
	const chartData = useMemo(
		() => makeSeries(`${resource.id}:${chartSpec.key}`, chartSpec),
		[resource.id, chartSpec],
	);
	const stats = useMemo(() => {
		const uptimeDays = seededInt(`${resource.id}:uptime`, 3, 41);
		return [
			{ label: 'Uptime', value: `${uptimeDays} days` },
			...headlineStats(kind, resource.id).slice(0, 3),
		];
	}, [kind, resource.id]);

	return (
		<div className="space-y-5">
			<Card>
				<CardHeader>
					<CardTitle className="font-mono text-xs uppercase tracking-[0.22em] text-primary">
						Planting details
					</CardTitle>
				</CardHeader>
				<CardContent className="space-y-4">
					<dl className="grid gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
						{kind === 'application' && application ? (
							<>
								<div>
									<InfoTerm>Repository</InfoTerm>
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
									<InfoTerm>Branch</InfoTerm>
									<dd className="mt-1.5 flex items-center gap-1.5 font-mono text-xs">
										<GitBranch className="size-3.5 text-muted-foreground" />
										{application.branch}
									</dd>
								</div>
								<div>
									<InfoTerm>Runtime</InfoTerm>
									<dd className="mt-1.5 capitalize">{application.type}</dd>
								</div>
							</>
						) : (
							<div>
								<InfoTerm>Engine</InfoTerm>
								<dd className="mt-1.5">
									{database?.type} {database?.version}
								</dd>
							</div>
						)}
						<div>
							<InfoTerm>Domain</InfoTerm>
							<dd className="mt-1.5 flex items-center gap-1 font-mono text-xs">
								<a
									href={`https://${subdomain.fullDomain}`}
									target="_blank"
									rel="noreferrer"
									className="inline-flex items-center gap-1 text-primary underline-offset-4 hover:underline"
								>
									{subdomain.fullDomain}
									<ExternalLink className="size-3" />
								</a>
								<CopyButton value={subdomain.fullDomain} label="Copy domain" />
							</dd>
						</div>
					</dl>
					{kind === 'database' &&
						(database?.connection_string ? (
							<div className="space-y-1.5">
								<span className="font-mono text-[10px] uppercase tracking-[0.25em] text-muted-foreground">
									Connection string
								</span>
								<ConnectionString value={database.connection_string} />
							</div>
						) : (
							<p className="text-sm text-muted-foreground">
								The connection string will appear here once the database is
								provisioned.
							</p>
						))}
				</CardContent>
			</Card>

			<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
				{stats.map((stat) => (
					<Card key={stat.label} className="gap-1 py-5">
						<CardHeader>
							<CardDescription className="font-mono text-[10px] uppercase tracking-[0.25em]">
								{stat.label}
							</CardDescription>
							<CardTitle className="font-display text-2xl font-semibold tracking-tight">
								{stat.value}
							</CardTitle>
						</CardHeader>
					</Card>
				))}
			</div>

			<Card className="gap-4">
				<CardHeader>
					<CardTitle className="font-mono text-xs uppercase tracking-[0.22em] text-primary">
						{chartSpec.title} — last 24 hours
					</CardTitle>
					<CardAction className="flex items-center gap-2">
						<DemoBadge />
						<Button variant="outline" size="sm" asChild>
							<Link to={`${basePath}/metrics`}>
								All vitals
								<ArrowUpRight className="size-3.5" />
							</Link>
						</Button>
					</CardAction>
				</CardHeader>
				<CardContent>
					<AreaChart
						data={chartData}
						color={chartSpec.color}
						unit={chartSpec.unit}
					/>
				</CardContent>
			</Card>
		</div>
	);
}

function ApplicationSettingsCard({
	subdomain,
	application,
	onChanged,
}: {
	subdomain: SubdomainListItemDto;
	application: ApplicationDto;
	onChanged: () => Promise<void>;
}) {
	const [pending, setPending] = useState(false);
	const [form, setForm] = useState({
		name: application.name,
		repo_url: application.repo_url,
		branch: application.branch,
		build_command: '',
		start_command: '',
	});

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
		} catch (err) {
			toast.error(errorMessage(err));
		} finally {
			setPending(false);
		}
	};

	return (
		<Card>
			<CardHeader>
				<CardTitle className="font-mono text-xs uppercase tracking-[0.22em] text-primary">
					Configuration
				</CardTitle>
				<CardDescription>
					How {application.name} is built and started.
				</CardDescription>
			</CardHeader>
			<CardContent>
				<form onSubmit={handleSave} className="grid gap-4 sm:grid-cols-2">
					<div className="space-y-2">
						<Label htmlFor="app-name">Name</Label>
						<Input
							id="app-name"
							value={form.name}
							onChange={(event) => setForm({ ...form, name: event.target.value })}
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
					<div className="flex items-end justify-end">
						<Button type="submit" disabled={pending}>
							{pending ? 'Saving…' : 'Save changes'}
						</Button>
					</div>
				</form>
			</CardContent>
		</Card>
	);
}

function DatabaseSettingsCard({
	subdomain,
	database,
	onChanged,
}: {
	subdomain: SubdomainListItemDto;
	database: DatabaseDto;
	onChanged: () => Promise<void>;
}) {
	const [pending, setPending] = useState(false);
	const [form, setForm] = useState({
		name: database.name,
		version: database.version,
	});

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
		} catch (err) {
			toast.error(errorMessage(err));
		} finally {
			setPending(false);
		}
	};

	return (
		<Card>
			<CardHeader>
				<CardTitle className="font-mono text-xs uppercase tracking-[0.22em] text-primary">
					Configuration
				</CardTitle>
				<CardDescription>
					The basics of how {database.name} runs.
				</CardDescription>
			</CardHeader>
			<CardContent>
				<form onSubmit={handleSave} className="grid gap-4 sm:grid-cols-2">
					<div className="space-y-2">
						<Label htmlFor="db-name">Name</Label>
						<Input
							id="db-name"
							value={form.name}
							onChange={(event) => setForm({ ...form, name: event.target.value })}
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
					<div className="flex items-end justify-end sm:col-span-2">
						<Button type="submit" disabled={pending}>
							{pending ? 'Saving…' : 'Save changes'}
						</Button>
					</div>
				</form>
			</CardContent>
		</Card>
	);
}

function DangerZone({
	kind,
	subdomain,
	resource,
}: {
	kind: ResourceKind;
	subdomain: SubdomainListItemDto;
	resource: ApplicationDto | DatabaseDto;
}) {
	const navigate = useNavigate();
	const noun = kind === 'application' ? 'application' : 'database';

	const handleDelete = async () => {
		try {
			if (kind === 'application') {
				await ApplicationController.delete(subdomain.id, resource.id);
			} else {
				await DatabaseController.delete(subdomain.id, resource.id);
			}
			toast.success(
				`${kind === 'application' ? 'Application' : 'Database'} deleted`,
			);
			navigate(`/projects/${subdomain.id}`);
		} catch (err) {
			toast.error(errorMessage(err));
		}
	};

	return (
		<div className="flex flex-wrap items-center justify-between gap-4 rounded-2xl border border-destructive/30 bg-destructive/5 px-5 py-4">
			<div>
				<h2 className="font-mono text-xs uppercase tracking-[0.2em] text-destructive">
					Delete this {noun}
				</h2>
				<p className="text-sm text-muted-foreground">
					Uproots {resource.name} from {subdomain.fullDomain} for good.
				</p>
			</div>
			<ConfirmDialog
				trigger={<Button variant="destructive">Delete {noun}</Button>}
				title={`Delete ${noun}`}
				description={
					kind === 'application'
						? `This permanently removes "${resource.name}" and its deployment from ${subdomain.fullDomain}.`
						: `This permanently removes "${resource.name}" and all of its data. This cannot be undone.`
				}
				confirmText={resource.name}
				onConfirm={handleDelete}
			/>
		</div>
	);
}

function SettingsTab({
	kind,
	subdomain,
	application,
	database,
	onChanged,
}: {
	kind: ResourceKind;
	subdomain: SubdomainListItemDto;
	application?: ApplicationDto;
	database?: DatabaseDto;
	onChanged: () => Promise<void>;
}) {
	const resource = (kind === 'application' ? application : database)!;
	return (
		<div className="space-y-5">
			{kind === 'application' && application ? (
				<>
					<ApplicationSettingsCard
						key={application.id}
						subdomain={subdomain}
						application={application}
						onChanged={onChanged}
					/>
					<EnvVarsCard initial={defaultEnvVars(subdomain.database)} />
				</>
			) : database ? (
				<DatabaseSettingsCard
					key={database.id}
					subdomain={subdomain}
					database={database}
					onChanged={onChanged}
				/>
			) : null}
			<DangerZone kind={kind} subdomain={subdomain} resource={resource} />
		</div>
	);
}

export default function ResourceDetail({ kind }: { kind: ResourceKind }) {
	const { id, tab } = useParams<{ id: string; tab?: string }>();
	const {
		data: subdomains,
		loading,
		error,
		reload,
	} = useAsync(SubdomainController.getAll);

	const subdomain = subdomains?.find((item) => item.id === id);
	const resource =
		kind === 'application' ? subdomain?.application : subdomain?.database;
	const basePath = `/projects/${id}/${kind === 'application' ? 'app' : 'db'}`;
	const tabs = TAB_DEFS[kind];
	const activeTab = tabs.some((item) => item.slug === (tab ?? ''))
		? (tab ?? '')
		: '';

	if (loading) {
		return (
			<div className="mx-auto w-full max-w-5xl flex-1 space-y-6 px-6 py-10">
				<Skeleton className="h-9 w-56" />
				<Skeleton className="h-5 w-72" />
				<Skeleton className="h-64 w-full rounded-2xl" />
				<Skeleton className="h-64 w-full rounded-2xl" />
			</div>
		);
	}

	if (error) {
		return (
			<div className="flex flex-1 flex-col items-center justify-center gap-4 px-6 py-24 text-center">
				<p className="text-muted-foreground">
					Could not load this {kind}: {error}
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

	if (!resource) {
		const createPath =
			kind === 'application'
				? `/projects/new/app?subdomain=${encodeURIComponent(subdomain.name)}`
				: `/projects/new/database?subdomain=${encodeURIComponent(subdomain.name)}`;
		return (
			<div className="flex flex-1 flex-col items-center justify-center gap-4 px-6 py-24 text-center">
				<p className="font-mono text-xs uppercase tracking-[0.22em] text-primary">
					Fallow ground
				</p>
				<h1 className="font-display text-3xl font-semibold tracking-tight">
					Nothing planted <em className="italic text-primary">here</em> yet
				</h1>
				<p className="text-muted-foreground">
					{subdomain.name} has no {kind} growing on it.
				</p>
				<div className="flex gap-2">
					<Button asChild>
						<Link to={createPath}>
							<Plus className="size-4" />
							Plant one now
						</Link>
					</Button>
					<Button variant="outline" asChild>
						<Link to={`/projects/${subdomain.id}`}>
							<ArrowLeft className="size-4" />
							Back to plot
						</Link>
					</Button>
				</div>
			</div>
		);
	}

	return (
		<div className="mx-auto w-full max-w-5xl flex-1 space-y-8 px-6 py-10">
			<Link
				to={`/projects/${subdomain.id}`}
				className="inline-flex items-center gap-1.5 font-mono text-[11px] uppercase tracking-[0.2em] text-muted-foreground transition-colors hover:text-primary"
			>
				<ArrowLeft className="size-3.5" />
				{subdomain.name}
			</Link>

			<div className="space-y-3">
				<p className="font-mono text-xs uppercase tracking-[0.22em] text-primary">
					{kind === 'application' ? 'Application record' : 'Database record'}
				</p>
				<div className="flex flex-wrap items-center gap-3">
					<span className="flex size-10 items-center justify-center rounded-xl border border-primary/30 bg-primary/10 text-primary">
						{kind === 'application' ? (
							<AppWindow className="size-5" />
						) : (
							<Database className="size-5" />
						)}
					</span>
					<h1 className="font-display text-4xl font-semibold tracking-tight">
						{resource.name}
					</h1>
					<StatusBadge status={resource.status} />
					{kind === 'application' && (
						<LifecycleControls
							subdomainId={subdomain.id}
							appId={resource.id}
							status={resource.status}
							onChanged={reload}
						/>
					)}
				</div>
				<div className="flex flex-wrap items-center gap-2 text-sm text-muted-foreground">
					<span className="font-mono text-xs">{subdomain.fullDomain}</span>
					<CopyButton value={subdomain.fullDomain} label="Copy domain" />
					{kind === 'application' && (
						<a
							href={`https://${subdomain.fullDomain}`}
							target="_blank"
							rel="noreferrer"
							className="inline-flex items-center gap-1 text-xs text-primary underline-offset-4 hover:underline"
						>
							Visit <ExternalLink className="size-3" />
						</a>
					)}
					<span className="text-xs">
						·{' '}
						{kind === 'application'
							? resource.type
							: `${resource.type} ${(resource as DatabaseDto).version}`}
					</span>
				</div>
			</div>

			<nav className="flex gap-6 overflow-x-auto border-b">
				{tabs.map((item) => (
					<Link
						key={item.slug}
						to={item.slug ? `${basePath}/${item.slug}` : basePath}
						className={cn(
							'-mb-px shrink-0 border-b-2 pb-3 font-mono text-[11px] uppercase tracking-[0.2em] transition-colors',
							activeTab === item.slug
								? 'border-primary text-primary'
								: 'border-transparent text-muted-foreground hover:text-foreground',
						)}
					>
						{item.label}
					</Link>
				))}
			</nav>

			{activeTab === '' && (
				<OverviewTab
					kind={kind}
					subdomain={subdomain}
					application={subdomain.application}
					database={subdomain.database}
					basePath={basePath}
				/>
			)}
			{activeTab === 'logs' && (
				<LogsTerminal
					key={resource.id}
					kind={kind}
					subdomainId={subdomain.id}
					resourceId={resource.id}
					resourceName={resource.name}
				/>
			)}
			{activeTab === 'metrics' && (
				<MetricsPanel kind={kind} seedId={resource.id} />
			)}
			{activeTab === 'deployments' && kind === 'application' && (
				<DeploymentsList
					seedId={resource.id}
					branch={(resource as ApplicationDto).branch}
					resourceName={resource.name}
					subdomainId={subdomain.id}
					appId={resource.id}
					repoUrl={(resource as ApplicationDto).repo_url}
					onChanged={reload}
				/>
			)}
			{activeTab === 'settings' && (
				<SettingsTab
					kind={kind}
					subdomain={subdomain}
					application={subdomain.application}
					database={subdomain.database}
					onChanged={reload}
				/>
			)}
		</div>
	);
}
