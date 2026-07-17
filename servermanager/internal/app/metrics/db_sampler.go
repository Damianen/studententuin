package metrics

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"servermanager/internal/domain"
	"servermanager/internal/ports"
)

// dbExecTimeout bounds one psql exec so a wedged database can't stall the
// whole tick.
const dbExecTimeout = 5 * time.Second

// pgStatsSQL returns "conn|xacts|bytes" in one statement: client backends
// (minus our own), total committed+rolled-back transactions, and total size
// of the non-template databases.
const pgStatsSQL = `SELECT (SELECT count(*) FROM pg_stat_activity WHERE backend_type = 'client backend' AND pid <> pg_backend_pid()), COALESCE(sum(xact_commit+xact_rollback),0), (SELECT COALESCE(sum(pg_database_size(datname)),0) FROM pg_database WHERE NOT datistemplate) FROM pg_stat_database`

// pgStatsArgv builds the fixed psql command. It works credential-free by the
// official image's invariants: initdb writes `local all all trust`, the exec
// runs as the container's postgres OS user (spec User), the "postgres"
// maintenance DB always exists, and POSTGRES_USER is a superuser. user is the
// only variable part and is re-validated before it gets here.
func pgStatsArgv(user string) []string {
	return []string{"psql", "-X", "-A", "-t", "-F", "|", "-U", user, "-d", "postgres", "-c", pgStatsSQL}
}

// dbSampler adds postgres vitals (conn/qps/disk) to a database container's
// tick values via one psql exec. Owned by the collector goroutine — no locks.
type dbSampler struct {
	exec ports.ContainerExec

	// idents caches the validated POSTGRES_USER per container ID. ok=false is
	// a permanent verdict: identifiers are immutable for a container's
	// lifetime, so an invalid one skips psql until the container is replaced.
	idents      map[string]dbIdentity
	identWarned map[string]bool
	// xactBaselines back the qps delta, keyed by container ID like the cpu
	// baselines.
	xactBaselines map[string]xactBaseline
}

type dbIdentity struct {
	user string
	ok   bool
}

type xactBaseline struct {
	xacts float64
	at    time.Time
}

func newDBSampler(exec ports.ContainerExec) *dbSampler {
	return &dbSampler{
		exec:          exec,
		idents:        map[string]dbIdentity{},
		identWarned:   map[string]bool{},
		xactBaselines: map[string]xactBaseline{},
	}
}

// sample adds conn/qps/disk to values. Failures warn, drop the qps baseline
// (the counters may have reset), and leave values untouched — cpu/mem still
// flow for the container.
func (s *dbSampler) sample(ctx context.Context, ctr domain.ManagedContainer, now time.Time, values map[string]float64) {
	user, ok := s.identity(ctx, ctr)
	if !ok {
		return
	}

	execCtx, cancel := context.WithTimeout(ctx, dbExecTimeout)
	out, err := s.exec.ExecCapture(execCtx, ctr.ID, pgStatsArgv(user))
	cancel()
	if err != nil {
		// Stopped/restarting containers land here (409) — routine, not fatal.
		slog.Warn("metrics: db stats exec", "container", ctr.Name, "error", err)
		delete(s.xactBaselines, ctr.ID)
		return
	}

	conn, xacts, bytes, err := parsePGStats(out)
	if err != nil {
		slog.Warn("metrics: db stats parse", "container", ctr.Name, "error", err)
		delete(s.xactBaselines, ctr.ID)
		return
	}

	values["conn"] = conn
	values["disk"] = bytes / (1024 * 1024)

	// qps is transactions/sec — a tps proxy, pinned as such: true
	// queries-per-second needs pg_stat_statements the images don't load.
	prev, ok := s.xactBaselines[ctr.ID]
	s.xactBaselines[ctr.ID] = xactBaseline{xacts: xacts, at: now}
	if ok && xacts >= prev.xacts && now.After(prev.at) {
		values["qps"] = (xacts - prev.xacts) / now.Sub(prev.at).Seconds()
	}
	// A negative delta means postgres restarted and reset its counters: the
	// fresh baseline above already replaced the stale one; skip the point.
}

// identity resolves and caches the container's POSTGRES_USER. Both identity
// envs are re-validated with the provision-path rule before the user lands in
// an argv — the manager treats them as request-provided even when it set them
// itself. Validation failure warns once (key name only, the value could be
// anything) and skips psql for this container permanently.
func (s *dbSampler) identity(ctx context.Context, ctr domain.ManagedContainer) (string, bool) {
	if ident, ok := s.idents[ctr.ID]; ok {
		return ident.user, ident.ok
	}

	env, err := s.exec.ContainerEnv(ctx, ctr.ID, "POSTGRES_USER", "POSTGRES_DB")
	if err != nil {
		// Transient inspect failure: not cached, retried next tick.
		slog.Warn("metrics: reading db identity env", "container", ctr.Name, "error", err)
		return "", false
	}

	for _, key := range []string{"POSTGRES_USER", "POSTGRES_DB"} {
		if err := domain.ValidateDBIdent(key, env[key]); err != nil {
			if !s.identWarned[ctr.ID] {
				slog.Warn("metrics: db identity failed validation — skipping db stats",
					"container", ctr.Name, "key", key)
				s.identWarned[ctr.ID] = true
			}
			s.idents[ctr.ID] = dbIdentity{ok: false}
			return "", false
		}
	}

	ident := dbIdentity{user: env["POSTGRES_USER"], ok: true}
	s.idents[ctr.ID] = ident
	return ident.user, true
}

// prune forgets per-ID state for containers not seen this tick.
func (s *dbSampler) prune(seen map[string]bool) {
	for id := range s.idents {
		if !seen[id] {
			delete(s.idents, id)
		}
	}
	for id := range s.identWarned {
		if !seen[id] {
			delete(s.identWarned, id)
		}
	}
	for id := range s.xactBaselines {
		if !seen[id] {
			delete(s.xactBaselines, id)
		}
	}
}

// parsePGStats parses the "conn|xacts|bytes" line psql -A -t -F '|' prints.
// Anything else — NULLs, error text, extra fields — is a parse error, never a
// zero point.
func parsePGStats(out string) (conn, xacts, bytes float64, err error) {
	fields := strings.Split(strings.TrimSpace(out), "|")
	if len(fields) != 3 {
		return 0, 0, 0, fmt.Errorf("expected conn|xacts|bytes, got %d fields", len(fields))
	}
	parsed := make([]float64, 3)
	for i, f := range fields {
		v, err := strconv.ParseFloat(strings.TrimSpace(f), 64)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("field %d is not a number", i)
		}
		parsed[i] = v
	}
	return parsed[0], parsed[1], parsed[2], nil
}
