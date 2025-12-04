import { Hono } from 'hono';
import { logger } from 'hono/logger';
import { serve } from '@hono/node-server';

const app = new Hono();

app.use(logger());

app.get('/', (c) => c.json({ message: 'Hello from Hono!' }));

serve({
	fetch: app.fetch,
	port: 3000,
});

console.log('Server running at http://localhost:3000');
