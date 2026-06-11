import { useState } from 'react';
import type { ReactNode } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog';

interface ConfirmDialogProps {
	trigger: ReactNode;
	title: string;
	description: string;
	confirmLabel?: string;
	// When set, the user must type this exact value to enable the confirm button.
	confirmText?: string;
	onConfirm: () => Promise<void> | void;
}

export function ConfirmDialog({
	trigger,
	title,
	description,
	confirmLabel = 'Delete',
	confirmText,
	onConfirm,
}: ConfirmDialogProps) {
	const [open, setOpen] = useState(false);
	const [typed, setTyped] = useState('');
	const [pending, setPending] = useState(false);

	const blocked = confirmText !== undefined && typed !== confirmText;

	const handleConfirm = async () => {
		setPending(true);
		try {
			await onConfirm();
			setOpen(false);
		} finally {
			setPending(false);
		}
	};

	return (
		<Dialog
			open={open}
			onOpenChange={(next) => {
				setOpen(next);
				if (!next) setTyped('');
			}}
		>
			<DialogTrigger asChild>{trigger}</DialogTrigger>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>{title}</DialogTitle>
					<DialogDescription>{description}</DialogDescription>
				</DialogHeader>
				{confirmText !== undefined && (
					<div className="space-y-2">
						<Label htmlFor="confirm-text">
							Type <span className="font-mono text-xs">{confirmText}</span> to
							confirm
						</Label>
						<Input
							id="confirm-text"
							value={typed}
							onChange={(event) => setTyped(event.target.value)}
							autoComplete="off"
						/>
					</div>
				)}
				<DialogFooter>
					<Button
						variant="ghost"
						onClick={() => setOpen(false)}
						disabled={pending}
					>
						Cancel
					</Button>
					<Button
						variant="destructive"
						onClick={handleConfirm}
						disabled={blocked || pending}
					>
						{pending ? 'Working…' : confirmLabel}
					</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
}
