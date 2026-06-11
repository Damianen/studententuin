# Studententuin — web

The frontend for Studententuin: a Vite + React single-page app, styled with
Tailwind CSS v4 and shadcn/ui primitives in a warm, garden-inspired theme.

## Development

```bash
npm install
npm run dev
```

The dev server runs on http://localhost:5173 and proxies `/api/*` to the Go
backend on `http://localhost:8080` (see `vite.config.ts`), so cookies work
same-origin and no CORS configuration is needed. Start the API first:

```bash
cd ../api && go run ./cmd/api
```

## Scripts

| Command           | What it does                            |
| ----------------- | --------------------------------------- |
| `npm run dev`     | Dev server with HMR and `/api` proxy    |
| `npm run build`   | Type-check + production build to `dist` |
| `npm run preview` | Serve the production build locally      |
| `npm run lint`    | ESLint over `src`                       |

## Production

`Dockerfile` builds the app and serves `dist/` with nginx (`nginx.conf`),
which handles the SPA fallback and proxies `/api/` to the `api` service.
