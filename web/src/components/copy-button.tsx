import { useState } from 'react';
import { Check, Copy } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';

export function CopyButton({
	value,
	label = 'Copy',
	className,
}: {
	value: string;
	label?: string;
	className?: string;
}) {
	const [copied, setCopied] = useState(false);

	const handleCopy = async (event: React.MouseEvent) => {
		event.preventDefault();
		event.stopPropagation();
		await navigator.clipboard.writeText(value);
		setCopied(true);
		setTimeout(() => setCopied(false), 1500);
	};

	return (
		<Button
			type="button"
			variant="ghost"
			size="icon"
			className={cn('size-7 text-muted-foreground', className)}
			onClick={handleCopy}
			aria-label={label}
		>
			{copied ? (
				<Check className="size-3.5 text-status-running" />
			) : (
				<Copy className="size-3.5" />
			)}
		</Button>
	);
}
