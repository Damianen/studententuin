import { ReactNode } from 'react';

export function MainLayout({ children }: { children: ReactNode }) {
	return (
		<div className="">
			<main className="">{children}</main>
		</div>
	);
}
