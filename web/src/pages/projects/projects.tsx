import { Link } from 'react-router';
import { AppWindow, Database, Plus, Sprout } from 'lucide-react';
import SubdomainController from '@/controllers/subdomain_controller';
import type { SubdomainListItemDto } from '@/dtos/subdomain_dtos';
import { useAsync } from '@/hooks/use_async';
import { StatusBadge } from '@/components/status-badge';
import { CopyButton } from '@/components/copy-button';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import {
	Card,
	CardAction,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';

function ProjectCard({ subdomain }: { subdomain: SubdomainListItemDto }) {
	const { application, database } = subdomain;

	return (
		<Card className="relative gap-4 transition-[translate,border-color,box-shadow] duration-300 hover:-translate-y-1 hover:border-primary/40 hover:shadow-glow">
			<Link
				to={`/projects/${subdomain.id}`}
				className="absolute inset-0 z-0 rounded-2xl"
				aria-label={`Open project ${subdomain.name}`}
			/>
			<CardHeader>
				<CardTitle className="font-display text-xl font-semibold tracking-tight">
					{subdomain.name}
				</CardTitle>
				<CardDescription className="flex items-center gap-1">
					<span className="truncate font-mono text-xs">
						{subdomain.fullDomain}
					</span>
					<CopyButton
						value={subdomain.fullDomain}
						label="Copy domain"
						className="relative z-10"
					/>
				</CardDescription>
				<CardAction>
					<span
						className={`mt-1.5 block size-2 rounded-full ${
							subdomain.isActive ? 'bg-status-running' : 'bg-status-stopped'
						}`}
						title={subdomain.isActive ? 'Active' : 'Inactive'}
					/>
				</CardAction>
			</CardHeader>
			<CardContent className="space-y-2.5">
				{application ? (
					<div className="flex items-center justify-between gap-2 rounded-lg bg-muted/60 px-3 py-2.5">
						<span className="flex min-w-0 items-center gap-2 text-sm">
							<AppWindow className="size-4 shrink-0 text-muted-foreground" />
							<span className="truncate font-medium">{application.name}</span>
							<span className="shrink-0 text-xs text-muted-foreground">
								{application.type}
							</span>
						</span>
						<StatusBadge status={application.status} />
					</div>
				) : (
					<Link
						to={`/projects/new/app?subdomain=${encodeURIComponent(subdomain.name)}`}
						className="relative z-10 flex items-center gap-2 rounded-lg border border-dashed px-3 py-2.5 text-sm text-muted-foreground transition-colors hover:border-primary/40 hover:text-primary"
					>
						<Plus className="size-4" />
						Add application
					</Link>
				)}
				{database ? (
					<div className="flex items-center justify-between gap-2 rounded-lg bg-muted/60 px-3 py-2.5">
						<span className="flex min-w-0 items-center gap-2 text-sm">
							<Database className="size-4 shrink-0 text-muted-foreground" />
							<span className="truncate font-medium">{database.name}</span>
							<span className="shrink-0 text-xs text-muted-foreground">
								{database.type} {database.version}
							</span>
						</span>
						<StatusBadge status={database.status} />
					</div>
				) : (
					<Link
						to={`/projects/new/database?subdomain=${encodeURIComponent(subdomain.name)}`}
						className="relative z-10 flex items-center gap-2 rounded-lg border border-dashed px-3 py-2.5 text-sm text-muted-foreground transition-colors hover:border-primary/40 hover:text-primary"
					>
						<Plus className="size-4" />
						Add database
					</Link>
				)}
			</CardContent>
		</Card>
	);
}

export default function Projects() {
	const {
		data: subdomains,
		loading,
		error,
		reload,
	} = useAsync(SubdomainController.getAll);

	return (
		<div className="mx-auto w-full max-w-6xl flex-1 px-6 py-10">
			<div className="mb-10 flex flex-wrap items-end justify-between gap-4">
				<div>
					<p className="mb-3 flex items-center gap-2.5 font-mono text-xs uppercase tracking-[0.22em] text-primary">
						<span className="size-1.5 animate-pulse rounded-full bg-primary" />
						Field overview
					</p>
					<h1 className="font-display text-4xl font-semibold tracking-tight">
						Your <em className="italic text-primary">garden</em>
					</h1>
				</div>
				<div className="flex gap-2">
					<Button variant="outline" asChild>
						<Link to="/projects/new/database">
							<Database className="size-4" />
							New database
						</Link>
					</Button>
					<Button asChild>
						<Link to="/projects/new/app">
							<Plus className="size-4" />
							New application
						</Link>
					</Button>
				</div>
			</div>

			{loading && (
				<div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
					{Array.from({ length: 6 }, (_, index) => (
						<Skeleton key={index} className="h-48 rounded-2xl" />
					))}
				</div>
			)}

			{!loading && error && (
				<div className="flex flex-col items-center gap-4 py-24 text-center">
					<p className="text-muted-foreground">
						Could not load your projects: {error}
					</p>
					<Button variant="outline" onClick={reload}>
						Try again
					</Button>
				</div>
			)}

			{!loading && !error && subdomains && subdomains.length === 0 && (
				<Card className="mx-auto max-w-lg py-12 text-center">
					<CardContent className="flex flex-col items-center gap-5">
						<span className="flex size-14 items-center justify-center rounded-full border border-primary/40 bg-primary/10 text-primary">
							<Sprout className="size-6" />
						</span>
						<div className="space-y-2">
							<p className="font-mono text-xs uppercase tracking-[0.22em] text-primary">
								Empty plot
							</p>
							<h2 className="font-display text-2xl font-semibold tracking-tight">
								Plant your first{' '}
								<em className="italic text-primary">project</em>
							</h2>
							<p className="text-sm text-muted-foreground">
								Claim a subdomain and put an application or database on it.
								It only takes a minute.
							</p>
						</div>
						<div className="flex gap-2">
							<Button asChild>
								<Link to="/projects/new/app">New application</Link>
							</Button>
							<Button variant="outline" asChild>
								<Link to="/projects/new/database">New database</Link>
							</Button>
						</div>
					</CardContent>
				</Card>
			)}

			{!loading && !error && subdomains && subdomains.length > 0 && (
				<div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
					{subdomains.map((subdomain) => (
						<ProjectCard key={subdomain.id} subdomain={subdomain} />
					))}
				</div>
			)}
		</div>
	);
}
