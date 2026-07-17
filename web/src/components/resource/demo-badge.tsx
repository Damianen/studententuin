export function DemoBadge({
	title = 'Response and request series arrive with edge routing (phase 7) — until then this is seeded sample data.',
}: {
	title?: string;
}) {
	return (
		<span
			className="inline-flex items-center rounded-full border border-dashed px-2.5 py-0.5 font-mono text-[10px] uppercase tracking-[0.18em] text-muted-foreground"
			title={title}
		>
			sample data
		</span>
	);
}
