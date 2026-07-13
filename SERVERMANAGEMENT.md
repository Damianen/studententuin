# Servermanagement — Build Plan (Epic 4)

Working document for building the `servermanager/` hosting layer — item 1 of
`IMPROVEMENTS.md`. Grounded in the two root PDFs: the project plan
(`Studententuin.pdf`, Epic 4 "Applicatie management" + architecture diagram) and the
performance research (`Studententuin_eindverslag_refs_opgeschoond.pdf`, runc vs gVisor
vs Kata). We build this together, phase by phase — check off items as they land.

**Epic 4 definition of done** (from the project plan): an application can be hosted,
taken offline, and the user can see and edit everything about their application
(logs, env vars, metrics, manage lifecycle).

---

## 1. What the research decided — and how the design applies it

The research compared three hosting architectures on Debian 13 / Docker 28.4.0
(12 threads, 31.2 GiB RAM):

| Measurement | runc | gVisor | Kata Containers |
|---|---|---|---|
| Cold start, mean of 20 (ms) | **1006** | 1272 | 1960 |
| Idle memory via `docker stats` (MiB) | 2.72 | 23.73 | 1.35 ⚠ undercounted |
| Load test throughput (req/s, AB c=20) | 38.8 | 109.3 | **157.8** |
| Load test mean response (ms) | 521 | 183 | **127** |
| P95 (ms) | 938 | 392 | **223** |
| Failed requests | 0 | 0 | 0 |

Design decisions taken directly from the research conclusions:

1. **runc is the default runtime for the MVP** — fastest cold starts, simplest
   operations, standard workflow. This is the research's explicit recommendation for
   a first version of Studententuin.
2. **Runtime is configurable per container** (`HostConfig.Runtime`), not hardcoded.
   Kata won the load test convincingly (4× runc throughput, lowest P95), so a future
   "hardened tier" (stronger isolation for long-running/paid workloads) is a config
   value, not a rewrite. gVisor earned no role: slower than runc at start, behind
   Kata under load, highest measured overhead.
3. **Metrics must be host-level, not `docker stats`** — the research showed
   `docker stats` attributes only 1.35 MiB to Kata containers because the microVM
   overhead lives outside the cgroup. The metrics phase uses cAdvisor (or cgroup
   reads) so numbers stay honest regardless of runtime.
4. **Scale-to-zero (Epic 8) is viable on runc** — ~1.0s cold start is an acceptable
   first-request wake latency. The deploy flow below keeps images around after stop,
   precisely so a stopped app can be restarted in ~1s later.
5. **Research follow-up recommendation** ("test more workload types, observe at host
   level, weigh operational aspects") maps to: structured logging in the manager,
   host-level metrics, and keeping the runtime abstraction thin so a second
   benchmark round can swap runtimes without code changes.

---

## 2. Architecture

From the project plan's architecture diagram: the main site (web + api + postgres)
lives on the VPS; the servermanager runs on a **separate hosting server** with full
Docker access and is the only thing that touches user containers.

```
┌─ VPS ──────────────────────────┐      ┌─ Hosting server ───────────────────┐
│  web (React) ── api (Gin) ── pg │      │  servermanager (Go + Docker SDK)   │
│                  │              │      │     │        │         │           │
│                  └── bearer ────┼──────┼──►  │        │         │           │
│              token, private net │      │  [app ctr] [app ctr] [db ctr] ...  │
└────────────────────────────────┘      │  (Traefik edge — later phase)      │
                                        └────────────────────────────────────┘
```

Trust model:

- The **browser never talks to the servermanager**. Everything goes through the api,
  which owns authentication (JWT cookie) and ownership checks (subdomain → user).
- The **api ↔ servermanager link is authenticated**: shared bearer token over the
  private network for the MVP (long random value, constant-time compare), mTLS as a
  later hardening step. The manager binds to the private interface only — it is
  root-equivalent (it holds the Docker socket), treat it like a database password.
- The servermanager is **stateless about ownership**: it identifies containers purely
  by the application/database UUID the api sends. Postgres on the VPS stays the
  single source of truth; the manager reconciles Docker state to what it's told.
- **User containers never get the Docker socket, host network, or the manager's
  network.** See the hardening defaults in §3.3.

Local dev: both run on one machine. Note the manager's stub currently listens on
`:8080`, which is taken locally (openmrs) — make the port env-driven (`SM_PORT`,
suggest 8095 locally) just like the api uses 8090 via `api/.env`.

---

## 3. The servermanager service

### 3.1 Module layout

Mirror the api's clean architecture so the codebase feels like one project:

```
servermanager/
  cmd/servermanager/main.go      config load, fail-fast on missing token, router
  internal/
    api/http/                    gin handlers + auth middleware (bearer token)
      apps/                      deploy/start/stop/remove/status/logs/metrics
      databases/                 provision/remove/status (phase 5)
      deployments/               deployment job status
    app/                         services + usecases (deploy orchestration, lifecycle)
    ports/                       ContainerRuntime, ImageBuilder, SourceFetcher, Clock
    infra/
      docker/                    Docker Go SDK adapter (implements ContainerRuntime)
      build/                     buildpacks/nixpacks adapter (implements ImageBuilder)
      git/                       go-git clone adapter (implements SourceFetcher)
  internal/domain/               ContainerSpec, DeploymentJob, value parsing (limits)
```

Key rules carried over from `IMPROVEMENTS.md`:

- **Docker Go SDK** (`github.com/moby/moby/client` + `github.com/moby/moby/api` —
  Engine v29 moved SDK releases there; the old `github.com/docker/docker` module is
  frozen at v28.5.x and no longer gets CVE fixes) everywhere. Never
  `exec.Command("docker", ...)` — repo URLs, env values and app names are user input;
  string-built shell commands are command injection.
- Ports (interfaces) + adapters, so services are unit-testable with a mocked
  `ContainerRuntime` — same pattern as `api/internal/app/ports/`.
- `log/slog` structured logging from day one (the api still needs this retrofit;
  don't repeat the mistake here).

### 3.2 Internal HTTP API (api → servermanager)

All routes under `/v1`, bearer token required, JSON in/out. The `:id` is always the
application/database UUID from the api's database; the manager derives container
names from it (`stt-app-<uuid>`, `stt-db-<uuid>`) — names are never user input.

| Route | Purpose |
|---|---|
| `POST /v1/apps/:id/deploy` | async: clone → build → (re)create container → start. Returns `202 {deployment_id}` |
| `GET  /v1/deployments/:id` | job status: `queued → cloning → building → starting → running \| failed` + error detail + build log tail |
| `POST /v1/apps/:id/run` | sync create+start from a prebuilt image (phase 1; phase 3's deploy ends in this same service path). 409 if the container exists — DELETE first |
| `POST /v1/apps/:id/start` / `:id/stop` | lifecycle on the existing container |
| `DELETE /v1/apps/:id` | remove container + network + anonymous volumes. Image removal deferred to phase 3 (manager doesn't own hand-pushed images); named-volume removal deferred to phase 3/5 |
| `GET  /v1/apps/:id` | actual Docker state (exists, running, started-at, restart count) |
| `GET  /v1/apps/:id/logs?tail=200&since=<rfc3339>` | structured log lines from the Docker logs API |
| `GET  /v1/apps/:id/logs/stream?tail=200` | live NDJSON log tail (`follow=true`); the api bridges it onto a browser WebSocket |
| `GET  /v1/apps/:id/metrics?range=24h` | series matching the frontend contract (phase 6) |
| `POST /v1/dbs/:id/provision` | async: pull → create → start → wait healthy. Returns `202 {provision_id}`; 409 while one is in flight or the container exists |
| `GET  /v1/dbs/:id` | Docker state (exists, running, health) + the latest provision job (`queued → pulling → starting → running \| failed` + error) |
| `DELETE /v1/dbs/:id` | idempotent sweep: container + `stt-dbnet-*` network + named data volume |
| `GET  /v1/dbs/:id/logs` (+`/stream`) | database container logs, same shapes as the app endpoints |
| `GET  /health` | unauthenticated liveness for compose/CI |

Deploy request body (everything the api already stores on `domain.Application`):

```json
{
  "repository_url": "https://github.com/user/repo",
  "branch": "main",
  "type": "Nodejs",
  "build_command": null,
  "start_command": null,  // on /run this is an argv ARRAY — the manager never shell-splits strings
  "env": {"NODE_ENV": "production"},
  "port": 3000,
  "memory_limit": "256m",
  "cpu_limit": "0.5",
  "runtime": "runc"
}
```

Log response shape — this is dictated by the frontend's existing
`LogEntry` (`web/src/lib/mock_telemetry.ts:12`): `{id, timestamp, level, message}`.
Level mapping for the MVP: stdout → `info`, stderr → `error` (Docker multiplexes the
two streams, so this is reliable without parsing); refine with line-content
heuristics later.

Deployments are **async jobs**: deploy returns immediately, the api polls
`GET /v1/deployments/:id` (background goroutine) and updates `Application.Status` +
`LastDeployedAt` + `DockerImage`/`DockerContainerID`/`DockerContainerName` on
completion. In-memory job store is fine for the MVP (single manager instance); jobs
are reconstructable from Docker state after a restart.

### 3.3 Container hardening defaults (non-negotiable, from day one)

Every user container is created with:

- `Resources.Memory` + `MemorySwap` (= same value: no swap), `NanoCPUs`,
  `PidsLimit` — parsed from the `MemoryLimit`/`CpuLimit` fields that already exist on
  `domain.Application`/`domain.Database` (today dead fields). Server-side caps with
  sane defaults (e.g. 256m / 0.5 CPU / 256 pids) when unset — **never unlimited**.
- `SecurityOpt: ["no-new-privileges:true"]`, `CapDrop: ["ALL"]` (add back only what's
  proven needed), non-root user where the image allows.
- **Isolated bridge network per app** (`stt-net-<uuid>`): app + its database only.
  No host network. No connectivity to the manager or other tenants. `ICC` off on any
  shared network.
  - Phase 5 refined the attachment direction: the shared app↔db network is
    **owned by the database** (`stt-dbnet-<db-uuid>`; the distinct prefix keeps
    the app-removal sweep of `stt-net-*` from touching it) and the app container
    is connected to it at deploy/run when the spec carries `database_id`. Same
    isolation, but the database can exist before/without an app and survives app
    deletion; every cutover re-attaches the new container.
- No bind mounts from the host; volumes are named Docker volumes only
  (`domain.Application.Volumes`). **The Docker socket is never mounted into any user
  container, under any circumstances.**
- `ReadonlyRootfs: true` where the app type allows, with a tmpfs for `/tmp`.
- `RestartPolicy: on-failure` with a max retry count (avoid crash-loop burning CPU).
- `HostConfig.Runtime` set from config: `runc` default, `kata` for the future
  hardened tier (§1.2). The value comes from a server-side allowlist — never free
  text from a request.
- Log driver `json-file` with rotation (`max-size=10m`, `max-file=3`) so user apps
  can't fill the disk via stdout.

### 3.4 Build pipeline (deploy-from-git)

Per the project plan: "containerize applications with a base, pull the application
code, build and start in Docker containers." Concretely:

1. **Fetch**: clone with `go-git` (pure Go — no `git` binary, no argv injection
   surface). Validate the URL first: scheme `https`, host allowlist
   (`github.com`, `gitlab.com` to start), depth-1 clone of the requested branch,
   size cap, timeout. Clone into a throwaway temp dir.
2. **Build**: **Nixpacks or Cloud Native Buildpacks** (`pack`) rather than hand-rolled
   per-language Dockerfiles — auto-detection gives us more than the current
   `ApplicationType` enum (only `Nodejs`) for free, and produces a normal OCI image.
   The builder itself runs **in a container, never on the host**, with no network
   access beyond what the build needs. Decision on Nixpacks vs CNB is open — see §10.
3. **Run**: create the container from the built image with the §3.3 hardening +
   `EnvironmentVariables` from the request (SDK `Env` slice — no shell, no injection)
   + the app's `StartCommand` if set.
4. **Cutover**: start new container, health-check it (TCP/HTTP on the app port, with
   timeout), then stop+remove the old one. Failed deploys leave the previous
   container running and mark the job `failed` with the build log attached.
5. **Cleanup**: prune dangling images beyond the last successful one per app
   (disk is the scarce resource on a budget host).

---

## 4. Changes in `api/`

Follow the existing clean-architecture pattern exactly:

- **New port** `api/internal/app/ports/server_manager.go`:

  ```go
  type ServerManagerClient interface {
      Deploy(ctx context.Context, app *domain.Application) (deploymentID string, err error)
      DeploymentStatus(ctx context.Context, deploymentID string) (*DeploymentStatus, error)
      Start(ctx context.Context, appID string) error
      Stop(ctx context.Context, appID string) error
      Remove(ctx context.Context, appID string) error
      Logs(ctx context.Context, appID string, opts LogOptions) ([]LogEntry, error)
      Metrics(ctx context.Context, appID string, rng string) (*MetricsResponse, error)
  }
  ```

  Implemented in `api/internal/infra/servermanager/` as a plain HTTP client (bearer
  token from env). Mock it in `api/internal/mocks/` like the other ports.

- **New usecases** in `api/internal/app/app/`: `usecase_deploy_application.go`,
  `usecase_start_application.go`, `usecase_stop_application.go`,
  `usecase_get_logs.go` — each re-using the existing ownership check
  (subdomain belongs to user **and** app belongs to subdomain — the check
  `IMPROVEMENTS.md` §2 wants a test for).
- **New routes** on the existing group (`application_router.go`):

  ```
  POST /subdomain/:id/application/:appId/deploy
  POST /subdomain/:id/application/:appId/start
  POST /subdomain/:id/application/:appId/stop
  GET  /subdomain/:id/application/:appId/logs
  GET  /subdomain/:id/application/:appId/metrics      (phase 6)
  GET  /subdomain/:id/application/:appId/deployments  (phase 6+)
  ```

- **Status lifecycle** stays on the existing `ApplicationStatus` values:
  `pending` (deploy requested/building) → `running` → `stopped`/`failed`. The api owns
  the transition writes via the existing `ApplicationRepo.Update`; a background
  poller watches in-flight deployments.
- Config: `SERVERMANAGER_URL`, `SERVERMANAGER_TOKEN` env vars; fail fast at startup
  if unset (same treatment IMPROVEMENTS.md wants for `JWT_TOKEN`).
- Delete flows: deleting an application/subdomain must also call `Remove` so
  containers don't orphan on the hosting server.

---

## 5. Frontend wiring (replace mocks tab by tab)

The contracts are already defined — treat them as the API spec:

- `web/src/lib/metric_specs.ts` — application series: `cpu` (%), `mem` (MB),
  `resp` (ms), `req` (/min); database series: `conn`, `qps` (/s), `cpu` (%),
  `disk` (MB). Points are `{time, value}` (`MetricPoint`).
- `web/src/lib/mock_telemetry.ts` — `LogEntry {id, timestamp, level, message}`,
  `DeploymentRecord {id, status, commit, branch, message, author, timeAgo, duration}`.

**Logs are the first win** (per IMPROVEMENTS.md): add `getLogs(subdomainId, appId)`
to `web/src/services/application_service.ts`, swap `makeLogs(...)` in
`logs-terminal.tsx:35` for a `useAsync` fetch, drop the `DemoBadge`, keep the
existing search/filter/download UI untouched. Then, in later phases: enable the
disabled "Deploy now" button in `deployments-list.tsx` (phase 3 ✅), wire env-var
persistence through the PATCH + GET round-trip (phase 4 ✅), real database
provisioning with live status/connection string/logs (phase 5 ✅ — the DemoBadge
now survives only on metrics), then metrics charts.

Source split for app metrics, planned now so the contract doesn't churn later:
`cpu`/`mem` come from container/host stats (cAdvisor, per §1.3); `resp`/`req` can
only come from the edge proxy (Traefik access metrics) — they land after phase 7.

---

## 6. Phased delivery (the build order)

### Phase 0 — Skeleton & contract ✦ start here
- [x] Restructure `servermanager/` into the §3.1 layout (`cmd/`, `internal/...`)
- [x] Config loading (`SM_PORT`, `SM_TOKEN`, `SM_BIND_ADDR`, `SM_DEFAULT_RUNTIME`,
      limits caps); fail fast if `SM_TOKEN` missing
- [x] Gin (or stdlib) router + bearer-token middleware (constant-time compare),
      `/health`, slog logging
- [x] Define `ports.ContainerRuntime` + domain types (`ContainerSpec`,
      `DeploymentJob`, limit parsing `"256m"`/`"0.5"` → bytes/NanoCPUs, with tests)
- [x] CI: make sure `servermanager-tests.yml` runs the new tests + gosec/govulncheck
      like the api workflow

### Phase 1 — Container lifecycle on a prebuilt image
- [x] Docker SDK adapter: create/start/stop/remove/inspect with the full §3.3
      hardening (limits, no-new-privileges, cap-drop, per-app network, log rotation,
      configurable runtime)
- [x] Routes: start/stop/remove/status against a hand-pushed test image (plus
      `POST /v1/apps/:id/run` to create+start from a prebuilt image — the path
      phase 3's deploy job will reuse)
- [x] Unit tests with mocked runtime; integration test behind a build tag that talks
      to real Docker (runs in CI on the ubuntu runner)
- [x] **Acceptance**: `curl` against the manager can run, stop, and inspect a
      hardened hello-world container with memory/CPU limits visible in `docker inspect`

### Phase 2 — Logs end-to-end (first mock replaced 🎉)
- [x] Manager: `GET /v1/apps/:id/logs` from the Docker logs API (tail, since,
      timestamps; stdout→info, stderr→error)
- [x] api: `ServerManagerClient` port + HTTP adapter + mock; logs usecase with
      ownership check (incl. app-belongs-to-subdomain); `GET .../logs` route
- [x] web: `getLogs` service + wire `logs-terminal.tsx` (5s poll), `DemoBadge`
      stays only on database logs (mock until phase 5)
- [x] **Acceptance**: the dashboard logs tab shows real output of a real container
- [x] Stretch: live tail — WebSocket end-to-end (manager `follow=true` NDJSON
      stream → api WS bridge via `coder/websocket` → web, poll fallback when
      the stream is unavailable; needed gin ≥ 1.12 in the api)

### Phase 3 — Deploy from git
- [x] `SourceFetcher` (go-git, URL validation, depth-1, size+time caps)
- [x] `ImageBuilder` (nixpacks — decided in §10; the binary only *generates*
      `.nixpacks/Dockerfile`, the daemon builds it via `ImageBuild`, so build
      steps run in daemon-managed containers, never on the host)
- [x] Async deploy job + `GET /v1/deployments/:id`; cutover replaces the
      container in place (network survives) with an inspect-after-grace health
      check; failed clone/build keeps the old container running, failed deploys
      surface the error + build-log tail; images GC'd to the last successful
- [x] api: deploy usecase, status poller, persist `DockerImage`/`DockerContainerID`/
      `DockerContainerName`/`LastDeployedAt`/`Status`; start/stop routes; delete
      now removes the container on the hosting server first (§4)
- [x] web: "Deploy now" → live stage feedback (cloning/building/starting),
      failure panel with build log, Stop/Start controls in the resource header
- [x] **Acceptance (= Epic 4 core)**: paste a GitHub URL of a Node app in the UI,
      deploy, watch status go pending→running, see its real logs; stop it from
      the UI. Covered end-to-end by `web/e2e/deploy.spec.ts` (needs
      `E2E_DEPLOY_REPO`, see `e2e-fixtures/node-hello/README.md`) plus the
      manager-side `TestIntegrationDeployPipeline`.
- MVP limits (documented): manager restart loses in-memory jobs (api poller
  fails after repeated 404s); no rollback-to-previous-image on failed start;
  no build cache (`--no-cache` keeps the generated Dockerfile legacy-builder
  compatible); deploys are single-flight per app.

### Phase 4 — Env vars & app settings round-trip
- [x] Deploy sends `EnvironmentVariables`; editing env vars in the UI persists via
      the existing PATCH and offers redeploy ("changing env" is an explicit
      servermanager duty in the architecture diagram)
- [x] Replace `defaultEnvVars` mock in the env tab with stored values
- [x] **Acceptance**: add a variable in the settings tab, save, hard-reload —
      it's still there; Deploy now from the same card and `docker inspect`
      shows it on the container. Covered by `web/e2e/env-vars.spec.ts` (no
      deploy repo needed) plus a redeploy leg in `deploy.spec.ts`.
- Landed with: the latent `enviroment_variables` column-key typo in the update
  usecase fixed (a unit test now pins the real column name); `serializer:json`
  on the jsonb map field (phase 4 is the first writer of that column); env
  keys validated in the api with the same `^[A-Za-z_][A-Za-z0-9_]*$` rule the
  manager enforces — a 400 names the key, never the value.

### Phase 5 — Databases
- [x] Manager provisions `domain.Database` containers from official images
      (postgres first; type/version from an **allowlist**, §3.3 hardening, named
      volume, generated credentials)
- [x] Connection string stored on the database record; injected as `DATABASE_URL`
      into the linked app (shared private network — never published to the
      internet; see the §3.3 note on the attachment direction)
- [x] **Acceptance**: create database in UI → running postgres container; its app
      can reach it; `ConnectionString` no longer null. Covered end-to-end by
      `web/e2e/databases.spec.ts` (no deploy repo needed) plus the manager-side
      `TestIntegrationProvisionPipeline` (provision → healthy → reachable from a
      linked container → delete sweeps container+network+volume → re-provision).
- Landed with: async provision job in the manager (`queued → pulling → starting
  → running | failed`, dbID-keyed single-flight store) polled by the api's
  `ProvisionPoller`; readiness via a TCP-forced `pg_isready` Docker healthcheck
  (the entrypoint's socket-only init server false-positives otherwise); postgres
  runs as the named `postgres` user under CapDrop ALL (fresh named volumes
  copy-up image ownership, so no chown is ever needed); `ShmSize` 256m against
  the postgres shared-memory footgun, so the api asks 512m memory for databases;
  db logs wired end-to-end (poll + WS stream), DemoBadge now metrics-only;
  db_password removed from the create flow — the api generates `app` + 32-hex
  creds and composes `postgres://…?sslmode=disable` (accepted MVP tradeoff on
  the isolated bridge); engine (type/version) immutable after create; the
  existing get/update/delete IDOR (no db-belongs-to-subdomain check — a foreign
  connection-string read and foreign volume destruction once containers exist)
  fixed alongside; subdomain delete now also sweeps its app container and its
  database (container + network + volume) on the hosting server.
- MVP limits (documented): manager restart mid-provision loses the job (stray
  container → 409; recover via delete + recreate); no re-provision endpoint and
  no credential rotation (delete + recreate); pull-if-missing pins images de
  facto — security patches need a manual pull; a running app picks up
  `DATABASE_URL` only on its next deploy (the UI offers that redeploy when the
  database turns running).

### Phase 6 — Metrics (research-informed)
- [ ] cAdvisor (or direct cgroup v2 reads) on the hosting server — host-level per
      §1.3, **not** `docker stats`
- [ ] Manager aggregates into the `metric_specs.ts` series contract (`cpu`, `mem`
      now; `conn`/`qps`/`disk` for databases from pg_stat + volume size)
- [ ] api proxy route + web charts off mock (`makeSeries` stays only for absent data)
- [ ] Real deployment history backing `deployments-list.tsx` (persist deployment
      records: commit, duration, status)

### Phase 7 — Edge routing & TLS
- [ ] Traefik on the hosting server, docker-label-driven: `app.<domain>` →
      container, wildcard TLS (or Caddy on-demand) — manager sets labels at create
- [ ] This unlocks `resp`/`req` metrics from Traefik for the remaining charts

### Phase 8 — Tiers & scale-to-zero (Epic 8 hook)
- [ ] `runtime: kata` hardened tier behind config + per-user level (Epic 2 "level")
- [ ] Idle detection → stop container; wake-on-request at the edge (runc ~1s cold
      start makes this acceptable, per the research)

---

## 7. Security checklist (gate before anything public)

- [x] No `exec.Command` with user-influenced strings anywhere in the manager (gosec
      rule in CI as backstop). The one exec — the nixpacks generator — uses a
      fixed argv (binary from config, fixed flags, manager-generated temp
      paths); user values travel via a config file, never argv.
- [ ] Bearer token: ≥32 random bytes, constant-time compare, private interface bind,
      never logged
- [x] Every container: limits + pids cap + no-new-privileges + cap-drop ALL +
      isolated network + no socket + log rotation (§3.3) — covered by a unit test
      that asserts the generated `HostConfig`
      (`internal/infra/docker/spec_mapper_test.go`, plus the daemon-side twin in
      the integration test)
- [x] Repo URL validation (scheme/host allowlist, https-only, no credentials/
      ports/queries) + clone size cap + per-stage timeouts
- [x] Image/version allowlists for databases (postgres 16/17, in manager code —
      the api mirrors it for clean 400s); runtime allowlist unchanged. Manager
      re-validates identifiers/password shape defensively; the provision job
      store never retains credentials.
- [x] Builds run containerized (docker daemon legacy builder; each RUN step is
      a container) on throwaway clones; build context tars never follow
      symlinks and exclude `.git`
- [x] Env var values treated as secrets in logs: neither service logs env
      maps, and the api's PATCH validation errors name the offending key only —
      values never reach a log line. (A user's own build log can still echo
      env their build prints; it is shown only to the owner.)
- [ ] Disk quotas / volume size caps; image GC
- [ ] Plays into Epic 5's "done": ZAP clean, SAST clean — manager endpoints are
      internal-only so the DAST surface stays the api

---

## 8. Testing & CI

- Unit: services against mocked `ContainerRuntime`/`ImageBuilder`/`SourceFetcher`
  (the api's testing style, `application_service_test.go` as the template).
- Integration: build-tagged tests against real Docker — lifecycle, limits actually
  applied, network isolation (container A cannot reach container B), log retrieval.
  GitHub's ubuntu runners have Docker; wire into `servermanager-tests.yml`.
- One e2e smoke once phase 3 lands: compose up everything + deploy a fixture repo →
  poll until running → assert logs. Doubles as the frontend's missing
  Playwright critical-path test (IMPROVEMENTS.md §4).

---

## 9. Configuration reference

| Var | Where | Notes |
|---|---|---|
| `SM_PORT` | servermanager | default 8080; locally use 8095 (8080 taken by openmrs) |
| `SM_BIND_ADDR` | servermanager | private interface IP in prod, 127.0.0.1 dev |
| `SM_TOKEN` | both | shared bearer secret; fail fast if missing |
| `SM_DEFAULT_RUNTIME` | servermanager | `runc` (research §1); `kata` allowlisted |
| `SM_MAX_MEMORY` / `SM_MAX_CPU` / `SM_DEFAULT_*` | servermanager | server-side caps for user limits |
| `SM_GIT_HOSTS` | servermanager | clone host allowlist; default `github.com,gitlab.com` |
| `SM_CLONE_TIMEOUT` / `SM_CLONE_MAX_SIZE` | servermanager | defaults `120s` / `200m` |
| `SM_BUILD_TIMEOUT` | servermanager | default `10m` |
| `SM_HEALTH_GRACE` | servermanager | post-start health wait; default `3s`, capped at 60s |
| `SM_NIXPACKS_BIN` | servermanager | default `nixpacks` (in the Docker image; pinned version) |
| `SM_DB_PULL_TIMEOUT` | servermanager | first pull of a database image; default `5m` |
| `SM_DB_HEALTH_BUDGET` | servermanager | db start → healthy (covers initdb); default `60s`, capped at 5m |
| `SERVERMANAGER_URL` | api | e.g. `http://10.0.0.2:8080` |
| `SERVERMANAGER_TOKEN` | api | same secret |

Local dev compose: add a `servermanager` service with `/var/run/docker.sock`
mounted (the manager legitimately needs it — this is exactly why **its** API needs
auth and user containers never get the socket). Sibling containers it creates run on
the host's Docker, outside the compose network — fine for dev.

---

## 10. Open decisions (input wanted)

1. **Nixpacks vs Cloud Native Buildpacks** — ✅ decided (phase 3): **Nixpacks**.
   Pinned release in the manager image + CI; generates the Dockerfile with
   `--no-cache` (no BuildKit-only directives) and the daemon's legacy builder
   builds it. Swapping to CNB (or BuildKit-over-the-port when the legacy
   builder disappears) stays behind the `ImageBuilder` port.
2. **Deploy status propagation**: api polls the manager (proposed, simplest) vs
   manager callbacks/webhooks to the api. Polling wins for MVP; revisit if we want
   live build logs streaming in the UI.
3. **Database engine scope** — ✅ decided (phase 5): **postgres-only** (16/17
   allowlisted). The domain enum keeps mysql/mongodb and the create dialog shows
   them disabled ("coming soon"); adding an engine is an allowlist entry plus
   engine-specific env/healthcheck/data-dir handling in the provision usecase.
4. **Job persistence**: in-memory (proposed for single-instance MVP, reconstructed
   from Docker state on restart) vs SQLite in the manager. Decide when deployment
   history (phase 6) lands — history likely belongs in the api's postgres anyway.
5. **One hosting server for now** — the manager API is already shaped so the api
   could fan out to several managers later (id-based, stateless), but nothing in the
   MVP should pretend to be multi-node.
