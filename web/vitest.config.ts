import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';
import path from 'node:path';

// Separate from vite.config.ts on purpose: tests don't need the Tailwind
// plugin or the dev-server /api proxy.
export default defineConfig({
	plugins: [react()],
	resolve: {
		alias: {
			'@': path.resolve(__dirname, 'src'),
		},
	},
	test: {
		environment: 'jsdom',
		// globals is load-bearing: Testing Library only auto-cleans up
		// between tests when vitest exposes afterEach globally.
		globals: true,
		setupFiles: ['./src/test/setup.ts'],
		include: ['src/**/*.test.{ts,tsx}'],
	},
});
