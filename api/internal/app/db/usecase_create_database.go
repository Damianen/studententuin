package db

import (
	"api/internal/app/ports"
	"api/internal/domain"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
)

// dbMemoryDefault is what the api asks for when the user set no limit. Higher
// than the app default on purpose: the manager carves a 256m /dev/shm out of
// the memory cgroup for postgres, and 256m total would let the database OOM
// itself.
const dbMemoryDefault = "512m"

// dbUser is the generated role every provisioned database gets. Fixed: the
// credential that matters is the password, and one recognizable role keeps
// connection strings predictable.
const dbUser = "app"

type CreateDatabase struct {
	databaseRepo  ports.DatabaseRepo
	serverManager ports.ServerManagerClient
	poller        *ProvisionPoller
}

type DatabaseInput struct {
	SubdomainID string
	Name        string
	Version     string
	Type        domain.DatabaseType
	DbName      string
	MemoryLimit *string
	CpuLimit    *string
}

// Execute creates the record and fires the async provision on the hosting
// server. Credentials are generated here (§6 phase 5) — the user never picks
// a database password; the composed connection string lands on the record
// when the poller sees the provision succeed. A provision that cannot even be
// requested marks the record failed rather than failing the create: the row
// is the user's handle for retrying via delete + recreate.
func (c *CreateDatabase) Execute(ctx context.Context, di DatabaseInput) error {
	id, err := uuid.Parse(di.SubdomainID)
	if err != nil {
		return err
	}

	password, err := generatePassword()
	if err != nil {
		return fmt.Errorf("generating database password: %w", err)
	}

	dbName := di.DbName
	database := &domain.Database{
		ID:          uuid.New(),
		SubdomainID: id,
		Name:        di.Name,
		Version:     di.Version,
		Type:        di.Type,
		Status:      domain.DatabaseStatusProvisioning,
		DbName:      &dbName,
		DbUser:      strPtr(dbUser),
		MemoryLimit: di.MemoryLimit,
		CpuLimit:    di.CpuLimit,
	}
	if err := c.databaseRepo.Create(database, ctx); err != nil {
		return err
	}

	memoryLimit := dbMemoryDefault
	if di.MemoryLimit != nil && *di.MemoryLimit != "" {
		memoryLimit = *di.MemoryLimit
	}
	cpuLimit := ""
	if di.CpuLimit != nil {
		cpuLimit = *di.CpuLimit
	}

	dbID := database.ID.String()
	_, err = c.serverManager.ProvisionDatabase(ctx, dbID, ports.DBProvisionSpec{
		Type:        string(di.Type),
		Version:     di.Version,
		DBName:      di.DbName,
		DBUser:      dbUser,
		DBPassword:  password,
		MemoryLimit: memoryLimit,
		CpuLimit:    cpuLimit,
	})
	if err != nil {
		fmt.Println("provision request failed:", err.Error())
		if uerr := c.databaseRepo.Update(dbID, map[string]any{"status": domain.DatabaseStatusFailed}, ctx); uerr != nil {
			fmt.Println("marking database failed:", uerr.Error())
		}
		return nil
	}

	c.poller.Watch(dbID, ProvisionCreds{User: dbUser, Password: password, DBName: di.DbName})
	return nil
}

// generatePassword returns 32 hex chars from crypto/rand — URL-safe by
// construction, so the connection string needs no escaping surprises.
func generatePassword() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func strPtr(s string) *string { return &s }
