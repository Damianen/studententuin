import { Link, NavLink, useNavigate } from 'react-router';
import { ChevronDown, LogOut, Settings, Sprout } from 'lucide-react';
import { toast } from 'sonner';
import { useAuth } from '@/contexts/auth_context';
import { LogoMark } from '@/components/logo';
import { Button } from '@/components/ui/button';
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

/*
 * The landing page's floating glass capsule, carried into the app —
 * same moonlit glass, minus the scroll theatrics.
 */
export function AppHeader() {
	const { user, isAuthenticated, logout } = useAuth();
	const navigate = useNavigate();

	const handleLogout = async () => {
		try {
			await logout();
			navigate('/');
		} catch {
			toast.error('Signing out failed. Please try again.');
		}
	};

	return (
		<header className="sticky top-0 z-40 flex justify-center px-4 pt-4">
			<nav
				aria-label="Main"
				className="relative flex items-center gap-1 rounded-full border border-border bg-[oklch(0.19_0.023_152/0.72)] py-1.5 pl-2 pr-1.5 shadow-lifted backdrop-blur-xl"
			>
				<Link
					to={isAuthenticated ? '/projects' : '/'}
					aria-label="Studententuin home"
					className="flex items-center gap-2.5 rounded-full px-2 py-1 transition-colors hover:bg-foreground/5"
				>
					<LogoMark className="size-7" />
					<span className="hidden font-display text-base font-semibold tracking-tight sm:inline">
						studententuin
					</span>
				</Link>

				<span aria-hidden className="mx-1 h-5 w-px bg-border" />

				{isAuthenticated ? (
					<>
						<NavLink
							to="/projects"
							className={({ isActive }) =>
								`rounded-full px-3.5 py-2 font-mono text-[10px] uppercase tracking-[0.2em] transition-colors ${
									isActive
										? 'text-primary'
										: 'text-muted-foreground hover:bg-foreground/5 hover:text-primary'
								}`
							}
						>
							Garden
						</NavLink>

						<span aria-hidden className="mx-1 h-5 w-px bg-border" />

						<DropdownMenu>
							<DropdownMenuTrigger asChild>
								<Button
									variant="ghost"
									className="gap-2 rounded-full px-2"
									aria-label="Account menu"
								>
									<span className="flex size-7 items-center justify-center rounded-full border border-primary/25 bg-primary/10 text-primary">
										<Sprout className="size-4" />
									</span>
									<span className="hidden max-w-32 truncate text-sm font-medium sm:inline">
										{user?.name}
									</span>
									<ChevronDown className="hidden size-3.5 text-muted-foreground sm:block" />
								</Button>
							</DropdownMenuTrigger>
							<DropdownMenuContent align="end" className="w-52 rounded-xl">
								<DropdownMenuLabel className="font-normal">
									<div className="text-sm font-medium">{user?.name}</div>
									<div className="truncate text-xs text-muted-foreground">
										{user?.email}
									</div>
								</DropdownMenuLabel>
								<DropdownMenuSeparator />
								<DropdownMenuItem onClick={() => navigate('/account')}>
									<Settings className="size-4" />
									Account
								</DropdownMenuItem>
								<DropdownMenuSeparator />
								<DropdownMenuItem onClick={handleLogout}>
									<LogOut className="size-4" />
									Sign out
								</DropdownMenuItem>
							</DropdownMenuContent>
						</DropdownMenu>
					</>
				) : (
					<div className="flex items-center gap-1">
						<Button size="sm" variant="ghost" asChild>
							<Link to="/login">Sign in</Link>
						</Button>
						<Button size="sm" asChild>
							<Link to="/register">Claim your plot</Link>
						</Button>
					</div>
				)}
			</nav>
		</header>
	);
}
