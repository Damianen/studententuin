# Studententuin — Improvement Roadmap

Assessment + suggestions from a full repo/PDF review (June 2026). The two PDFs in the
root (project plan with epics, and the runc/gVisor/Kata performance research) are the
source of the goals referenced below.

## Where the project stands

- `api/` — solid: clean architecture (controllers → services → ports → GORM repos),
  Gin + GORM + Postgres, JWT cookie auth with bcrypt, ownership checks on every
  resource, good unit tests, CI with tests + govulncheck + gosec.
- `web/` — well-layered (router → UI → controllers → DTOs → services → API), route
  guards, 401 handling, polished resource-detail UI. Logs/metrics/deployments/env-vars
  tabs run on **seeded mock data** (`web/src/lib/mock_telemetry.ts`).
- `servermanager/` — **a stub**. Epic 4 (the actual hosting) doesn't exist yet. The
  runc/gVisor/Kata research is done but nothing uses it.

The gap isn't code quality — it's that the riskiest, most differentiating component
(the hosting layer) is still on paper. Resist polishing the dashboard further until one
real container can be deployed, started, and show real logs end to end.

## 1. Build the servermanager — this is the platform

- Use the **Docker Go SDK** (`github.com/docker/docker/client`), never
  `exec.Command("docker", ...)` — string-built shell commands + user input (repo URLs,
  env values, app names) = command injection.
- Make the **runtime configurable per container** (`HostConfig.Runtime`). Research
  concluded runc for MVP, but Kata was strong under load — config-driven runtime makes
  a "hardened tier = Kata" trivial later.
- **Enforce the limits already modeled**: `domain.Application` has `MemoryLimit`,
  `CpuLimit`, `EnvironmentVariables`, `Volumes` — all dead fields today. Wire them to
  `HostConfig.Resources`. Never run user containers without memory/CPU/pids limits.
- Lock down user containers: no docker socket mount ever, `no-new-privileges`, drop
  capabilities, isolated network per app (no host network, no access to the manager),
  read-only root FS where possible, disk quotas.
- **Authenticate the API ↔ servermanager link** (the manager is root-equivalent).
  Minimum: long random bearer token over a private network; better: mTLS. Bind to the
  private interface only.
- Deployment flow: prefer **Cloud Native Buildpacks or Nixpacks** over hand-rolled
  base images — language detection + OCI image without per-runtime Dockerfiles. Build
  in a throwaway container, not on the host.
- Routing/TLS for subdomains: **Traefik** (docker-label-driven) or **Caddy**
  (on-demand TLS) for `app.domain` → container + wildcard HTTPS.
- Metrics: the research itself recommends host-level observation (cAdvisor /
  systemd-cgtop) because `docker stats` undercounts Kata. The frontend's
  `web/src/lib/metric_specs.ts` already defines the wanted series (CPU, memory,
  response time, req/min) — treat that as the API contract and replace
  `mock_telemetry.ts` tab by tab. **Logs are the easy first win**: stream
  `GET /apps/:id/logs` from the Docker logs API into the existing `logs-terminal.tsx`.

## 2. Finish the auth epic (Epic 1's own "done" criteria aren't met)

- **Email verification** (`User.EmailVerifiedAt` exists but is never set) and
  **password reset** — both in the epic's definition of done, both missing.
- **Rate limiting on login/register** — nothing stops brute force today.
- **Token revocation / sessions**: architecture diagram has a `sessions` table but the
  implementation is a stateless 24h JWT — logout can't invalidate anything.
  Server-side sessions (simpler, fits single-VPS) or short-lived JWT + refresh rotation.
- `Secure` cookie flag is hardcoded `false`
  (`api/internal/api/http/auth/auth_controller.go:48`) — make it config-driven. Fail
  fast at startup if `JWT_TOKEN` is missing or equals the dev default.
- **Security headers middleware** (HSTS, X-Content-Type-Options, X-Frame-Options, CSP,
  Permissions-Policy) — ZAP baseline flagged exactly these (8 warnings). This is the
  prerequisite for flipping the DAST job to blocking (`fail_action: true` in
  `.github/workflows/dast.yml` + a `.zap/rules.tsv` for accepted findings).
- Config-driven CORS origin instead of hardcoded `localhost:3000` in
  `api/cmd/api/main.go:33-38`.
- Verify with a test that an app/database ID is validated as belonging to the
  subdomain in the URL (not just that the subdomain belongs to the user).

## 3. ~~CI/CD gaps (Epic 5)~~ — DONE (June 2026)

Delivered: frontend + servermanager workflows, ZAP baseline DAST (report-only until
security headers land), gitleaks, Dependabot (gomod/npm/docker/actions), all actions
SHA-pinned. See `.github/workflows/`.

## 4. Operational basics

- ~~Dockerfiles + docker-compose~~ — **DONE** (June 2026): `api/Dockerfile`,
  root `docker-compose.yml` (db/api/web, healthchecks, only `WEB_PORT` published,
  default 3000 — locally use `WEB_PORT=3001`, 3000 is taken by Grafana).
- **Versioned migrations** (golang-migrate or Atlas) instead of GORM AutoMigrate —
  deliberately deferred from the docker pass. Note: `db.go:40-45` overwrites `err`
  per AutoMigrate call, only the last error is checked.
- **Structured logging** (`log/slog`) instead of `fmt.Println` in controllers.
- **Postgres backups** — nightly `pg_dump` to off-server storage, minimum.
- Frontend tests: zero today. `playwright-core` already in devDeps — one
  login → create project → create app smoke test covers the critical path. Vitest for
  `useAsync`/controllers next.

## 5. Product-level (later epics)

- Deploy pipeline: GitHub webhook → redeploy on push (the UI already has a disabled
  "Deploy now" button in `deployments-list.tsx`).
- Quotas per user level (the "level" concept from Epic 2).
- Free tier / resource saving (Epic 8): scale-to-zero — connects directly to the cold
  start research (runc's fast cold starts matter most there).

## Suggested order

1. ~~Dockerfiles + compose~~ ✅ (migrations still open → item 4)
2. Servermanager MVP: authenticated API, deploy-from-git via buildpacks, start/stop,
   logs endpoint, resource limits from day one
3. Wire frontend logs tab to the real endpoint (first mock replaced)
4. Auth hardening: rate limiting, sessions/revocation, email verification + reset,
   security headers (→ then flip ZAP to blocking)
5. Reverse proxy + wildcard TLS for subdomains
6. Metrics (cAdvisor or docker stats), then env-var persistence, then deployment history
7. ~~CI expansion + ZAP~~ ✅
