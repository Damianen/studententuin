import { render, screen } from '@testing-library/react';
import { StatusBadge } from './status-badge';

function dotOf(container: HTMLElement) {
	return container.querySelector('span > span') as HTMLElement;
}

describe('StatusBadge', () => {
	it('renders the status text lowercased', () => {
		render(<StatusBadge status="Running" />);
		expect(screen.getByText('running')).toBeInTheDocument();
	});

	it.each([
		['running', 'bg-status-running'],
		['provisioning', 'bg-status-pending'],
		['pending', 'bg-status-pending'],
		['stopped', 'bg-status-stopped'],
		['failed', 'bg-status-failed'],
	])('shows the %s dot', (status, expectedClass) => {
		const { container } = render(<StatusBadge status={status} />);
		expect(dotOf(container)).toHaveClass(expectedClass);
	});

	it('pulses while provisioning', () => {
		const { container } = render(<StatusBadge status="provisioning" />);
		expect(dotOf(container)).toHaveClass('animate-pulse');
	});

	it('falls back to the stopped dot for unknown statuses', () => {
		const { container } = render(<StatusBadge status="mystery" />);
		expect(dotOf(container)).toHaveClass('bg-status-stopped');
	});
});
