package db

import (
	"api/internal/app/ports"
	"api/internal/domain"
	"context"
	"fmt"
	"net/url"
	"time"
)

const (
	defaultProvisionPollInterval = 2 * time.Second
	// defaultProvisionPollBudget covers the manager's pull (5m) + health
	// (60s) budgets plus slack.
	defaultProvisionPollBudget = 7 * time.Minute
	// maxLostStrikes: consecutive "no container, no provision job" polls
	// before the provision is declared lost (manager restart mid-provision).
	maxLostStrikes = 3
)

// ProvisionCreds carries the generated credentials from create to the poller,
// which composes the connection string once the database is actually up. They
// live only in this process — the record stores them inside ConnectionString.
type ProvisionCreds struct {
	User     string
	Password string
	DBName   string
}

// ProvisionPoller watches in-flight database provisions and writes the
// outcome onto the database row: the api owns all Status transitions (§4).
// Polling state is in-memory — an api restart drops the watch and the row
// stays "provisioning" until the user deletes and recreates (MVP limit).
type ProvisionPoller struct {
	databaseRepo  ports.DatabaseRepo
	serverManager ports.ServerManagerClient
	clock         ports.Clock

	// Interval and Budget are overridable for tests.
	Interval time.Duration
	Budget   time.Duration
}

func NewProvisionPoller(repo ports.DatabaseRepo, sm ports.ServerManagerClient, clock ports.Clock) *ProvisionPoller {
	return &ProvisionPoller{
		databaseRepo:  repo,
		serverManager: sm,
		clock:         clock,
		Interval:      defaultProvisionPollInterval,
		Budget:        defaultProvisionPollBudget,
	}
}

// Watch follows one provision to its terminal state in the background.
func (p *ProvisionPoller) Watch(dbID string, creds ProvisionCreds) {
	go p.poll(dbID, creds)
}

func (p *ProvisionPoller) poll(dbID string, creds ProvisionCreds) {
	ctx := context.Background()
	deadline := p.clock.Now().Add(p.Budget)
	strikes := 0

	for {
		time.Sleep(p.Interval)

		status, err := p.serverManager.DatabaseStatus(ctx, dbID)
		switch {
		case err != nil:
			// Transient manager trouble: keep polling within the budget.
		case status.Provision != nil && status.Provision.Status == "running":
			p.markRunning(ctx, dbID, creds, status)
			return
		case status.Provision != nil && status.Provision.Status == "failed":
			p.markFailed(ctx, dbID)
			return
		case status.Provision == nil && status.Exists && status.Running && status.Health == "healthy":
			// The manager lost the job (restart) but the container made it —
			// Docker state is the recoverable truth.
			p.markRunning(ctx, dbID, creds, status)
			return
		case status.Provision == nil:
			if strikes++; strikes >= maxLostStrikes {
				p.markFailed(ctx, dbID)
				return
			}
		default:
			strikes = 0
		}

		if p.clock.Now().After(deadline) {
			p.markFailed(ctx, dbID)
			return
		}
	}
}

func (p *ProvisionPoller) markRunning(ctx context.Context, dbID string, creds ProvisionCreds, status *ports.DBStatus) {
	containerName := domainContainerName(dbID)
	image := ""
	containerID := ""
	if prov := status.Provision; prov != nil {
		if prov.ContainerName != "" {
			containerName = prov.ContainerName
		}
		image = prov.Image
		containerID = prov.ContainerID
	}

	updates := map[string]any{
		"status":                domain.DatabaseStatusRunning,
		"connection_string":     connectionString(creds, containerName),
		"docker_container_name": containerName,
		"updated_at":            p.clock.Now(),
	}
	if image != "" {
		updates["docker_image"] = image
	}
	if containerID != "" {
		updates["docker_container_id"] = containerID
	}
	if err := p.databaseRepo.Update(dbID, updates, ctx); err != nil {
		fmt.Println("recording successful provision:", err.Error())
	}
}

func (p *ProvisionPoller) markFailed(ctx context.Context, dbID string) {
	if err := p.databaseRepo.Update(dbID, map[string]any{
		"status":     domain.DatabaseStatusFailed,
		"updated_at": p.clock.Now(),
	}, ctx); err != nil {
		fmt.Println("recording failed provision:", err.Error())
	}
}

// connectionString composes the postgres URL the linked app receives as
// DATABASE_URL. The host is the container name — resolvable only on the
// database's private network, never from the internet. sslmode=disable is an
// accepted MVP tradeoff on that isolated bridge (revisit with phase 7).
func connectionString(creds ProvisionCreds, containerName string) string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(creds.User, creds.Password),
		Host:     containerName + ":5432",
		Path:     "/" + creds.DBName,
		RawQuery: "sslmode=disable",
	}
	return u.String()
}

// domainContainerName mirrors the manager's deterministic naming rule
// (stt-db-<uuid>) as the fallback when the provision state is absent.
func domainContainerName(dbID string) string { return "stt-db-" + dbID }
