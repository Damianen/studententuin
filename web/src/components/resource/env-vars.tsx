import { useRef, useState } from 'react';
import { Eye, EyeOff, Plus, Trash2 } from 'lucide-react';
import { toast } from 'sonner';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
	Card,
	CardAction,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';
import { DemoBadge } from '@/components/resource/demo-badge';

interface EnvVariable {
	id: number;
	key: string;
	value: string;
	visible: boolean;
}

export function EnvVarsCard({
	initial,
}: {
	initial: Array<{ key: string; value: string }>;
}) {
	const nextId = useRef(initial.length);
	const [variables, setVariables] = useState<EnvVariable[]>(() =>
		initial.map((variable, index) => ({ id: index, visible: false, ...variable })),
	);
	const [draft, setDraft] = useState({ key: '', value: '' });

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

	return (
		<Card>
			<CardHeader>
				<CardTitle className="font-mono text-xs uppercase tracking-[0.22em] text-primary">
					Environment
				</CardTitle>
				<CardDescription>
					Variables available to the application at runtime.
				</CardDescription>
				<CardAction>
					<DemoBadge />
				</CardAction>
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

				<div className="flex justify-end border-t pt-4">
					<Button
						onClick={() =>
							toast.info(
								"Environment variables are a preview — the API doesn't store them yet.",
							)
						}
					>
						Save changes
					</Button>
				</div>
			</CardContent>
		</Card>
	);
}
