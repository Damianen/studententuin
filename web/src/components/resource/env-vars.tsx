import { useRef, useState } from 'react';
import { Eye, EyeOff, Plus, RefreshCw, Rocket, Trash2 } from 'lucide-react';
import { toast } from 'sonner';
import ApplicationController from '@/controllers/application_controller';
import type { ApplicationDto } from '@/dtos/application_dtos';
import type { SubdomainListItemDto } from '@/dtos/subdomain_dtos';
import { STAGE_LABELS, useDeployment } from '@/hooks/use_deployment';
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

// Mirrors the key rule the api and servermanager enforce.
const ENV_KEY_PATTERN = /^[A-Za-z_][A-Za-z0-9_]*$/;

interface EnvVariable {
	id: number;
	key: string;
	value: string;
	visible: boolean;
}

function errorMessage(err: unknown): string {
	return err instanceof Error ? err.message : 'Something went wrong';
}

// Sorted for a stable order — jsonb does not keep insertion order.
function sortedEntries(env: Record<string, string>): Array<[string, string]> {
	return Object.entries(env).sort(([a], [b]) => a.localeCompare(b));
}

export function EnvVarsCard({
	subdomain,
	application,
	onChanged,
}: {
	subdomain: SubdomainListItemDto;
	application: ApplicationDto;
	onChanged: () => Promise<void>;
}) {
	const storedEnv = application.environment_variables ?? {};
	const nextId = useRef(Object.keys(storedEnv).length);
	const [variables, setVariables] = useState<EnvVariable[]>(() =>
		sortedEntries(storedEnv).map(([key, value], index) => ({
			id: index,
			key,
			value,
			visible: false,
		})),
	);
	const [draft, setDraft] = useState({ key: '', value: '' });
	const [pending, setPending] = useState(false);
	// Set after a successful save; drives the "takes effect on next deploy" offer.
	const [saved, setSaved] = useState(false);

	const {
		deploying,
		stage,
		error: deployError,
		deploy,
	} = useDeployment(subdomain.id, application.id, (ok) => {
		if (ok) {
			setSaved(false);
			toast.success(`${application.name} is live`);
		} else {
			toast.error('Deployment failed');
		}
		void onChanged();
	});

	const update = (id: number, patch: Partial<EnvVariable>) => {
		setVariables((current) =>
			current.map((variable) =>
				variable.id === id ? { ...variable, ...patch } : variable,
			),
		);
	};

	const addDraft = () => {
		if (!draft.key.trim() || !draft.value.trim()) return;
		setVariables((current) => [
			...current,
			{
				id: nextId.current++,
				key: draft.key.trim(),
				value: draft.value.trim(),
				visible: false,
			},
		]);
		setDraft({ key: '', value: '' });
	};

	const handleSave = async () => {
		const env: Record<string, string> = {};
		for (const variable of variables) {
			const key = variable.key.trim();
			if (!ENV_KEY_PATTERN.test(key)) {
				toast.error(
					key === ''
						? 'Variable keys cannot be empty'
						: `Invalid variable key "${key}" — use letters, digits and underscores, not starting with a digit`,
				);
				return;
			}
			if (key in env) {
				toast.error(`Duplicate variable key "${key}"`);
				return;
			}
			env[key] = variable.value;
		}

		setPending(true);
		try {
			await ApplicationController.patch(subdomain.id, application.id, {
				environment_variables: env,
			});
			await onChanged();
			setVariables(
				sortedEntries(env).map(([key, value]) => ({
					id: nextId.current++,
					key,
					value,
					visible: false,
				})),
			);
			setSaved(true);
			toast.success('Environment variables saved');
		} catch (err) {
			toast.error(errorMessage(err));
		} finally {
			setPending(false);
		}
	};

	const offerDeploy = saved && Boolean(application.repo_url);

	return (
		<Card>
			<CardHeader>
				<CardTitle className="font-mono text-xs uppercase tracking-[0.22em] text-primary">
					Environment
				</CardTitle>
				<CardDescription>
					Variables available to the application at runtime.
				</CardDescription>
			</CardHeader>
			<CardContent className="space-y-4">
				{variables.map((variable) => (
					<div key={variable.id} className="flex items-end gap-2">
						<div className="flex-1 space-y-2">
							<Label htmlFor={`env-key-${variable.id}`}>Key</Label>
							<Input
								id={`env-key-${variable.id}`}
								value={variable.key}
								onChange={(event) => update(variable.id, { key: event.target.value })}
								placeholder="VARIABLE_NAME"
								className="font-mono text-xs"
							/>
						</div>
						<div className="flex-1 space-y-2">
							<Label htmlFor={`env-value-${variable.id}`}>Value</Label>
							<div className="relative">
								<Input
									id={`env-value-${variable.id}`}
									type={variable.visible ? 'text' : 'password'}
									value={variable.value}
									onChange={(event) =>
										update(variable.id, { value: event.target.value })
									}
									placeholder="value"
									className="pr-9 font-mono text-xs"
								/>
								<Button
									type="button"
									variant="ghost"
									size="icon"
									className="absolute right-0 top-0 h-full text-muted-foreground"
									onClick={() => update(variable.id, { visible: !variable.visible })}
									aria-label={variable.visible ? 'Hide value' : 'Show value'}
								>
									{variable.visible ? (
										<EyeOff className="size-3.5" />
									) : (
										<Eye className="size-3.5" />
									)}
								</Button>
							</div>
						</div>
						<Button
							type="button"
							variant="ghost"
							size="icon"
							className="text-muted-foreground hover:text-destructive"
							onClick={() =>
								setVariables((current) =>
									current.filter((item) => item.id !== variable.id),
								)
							}
							aria-label={`Remove ${variable.key}`}
						>
							<Trash2 className="size-4" />
						</Button>
					</div>
				))}

				<div className="space-y-4 border-t pt-4">
					<p className="font-mono text-[10px] uppercase tracking-[0.25em] text-muted-foreground">
						Add new variable
					</p>
					<div className="flex items-end gap-2">
						<div className="flex-1 space-y-2">
							<Label htmlFor="env-new-key">Key</Label>
							<Input
								id="env-new-key"
								value={draft.key}
								onChange={(event) => setDraft({ ...draft, key: event.target.value })}
								placeholder="VARIABLE_NAME"
								className="font-mono text-xs"
							/>
						</div>
						<div className="flex-1 space-y-2">
							<Label htmlFor="env-new-value">Value</Label>
							<Input
								id="env-new-value"
								value={draft.value}
								onChange={(event) =>
									setDraft({ ...draft, value: event.target.value })
								}
								placeholder="value"
								className="font-mono text-xs"
							/>
						</div>
						<Button
							type="button"
							variant="outline"
							onClick={addDraft}
							disabled={!draft.key.trim() || !draft.value.trim()}
						>
							<Plus className="size-4" />
							Add
						</Button>
					</div>
				</div>

				<div className="space-y-3 border-t pt-4">
					{offerDeploy && !deploying && (
						<p className="text-xs text-muted-foreground">
							Saved — changes take effect on the next deploy.
						</p>
					)}
					{deploying && (
						<p className="text-xs text-status-pending">
							Deploying — {STAGE_LABELS[stage ?? ''] ?? stage}
						</p>
					)}
					{!deploying && deployError && (
						<p className="text-xs text-status-failed">
							{deployError} — the deployments tab has the build log.
						</p>
					)}
					<div className="flex justify-end gap-2">
						{offerDeploy && (
							<Button
								variant="outline"
								onClick={deploy}
								disabled={deploying || pending}
							>
								{deploying ? (
									<RefreshCw className="size-3.5 animate-spin" />
								) : (
									<Rocket className="size-3.5" />
								)}
								{deploying ? 'Deploying…' : 'Deploy now'}
							</Button>
						)}
						<Button onClick={handleSave} disabled={pending || deploying}>
							{pending ? 'Saving…' : 'Save changes'}
						</Button>
					</div>
				</div>
			</CardContent>
		</Card>
	);
}
