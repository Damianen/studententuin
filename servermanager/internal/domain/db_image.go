package domain

import (
	"fmt"
	"regexp"
)

// DBImage is one allowlisted database image. DataDir is per entry because
// upstream images move it between majors (postgres 18 uses
// /var/lib/postgresql instead of .../data).
type DBImage struct {
	Ref     string
	DataDir string
}

// dbImages is the server-side image allowlist (SERVERMANAGEMENT.md §7):
// requests carry type+version, never an image reference. Postgres-only
// through phase 5 (§10.3); debian variants — musl's stubbed locales make
// alpine risky for collation-sensitive data.
var dbImages = map[string]map[string]DBImage{
	"postgres": {
		"16": {Ref: "postgres:16", DataDir: "/var/lib/postgresql/data"},
		"17": {Ref: "postgres:17", DataDir: "/var/lib/postgresql/data"},
	},
}

// DBImageFor resolves a database type+version against the allowlist.
func DBImageFor(dbType, version string) (DBImage, error) {
	versions, ok := dbImages[dbType]
	if !ok {
		return DBImage{}, fmt.Errorf("database type %q is not supported: %w", dbType, ErrInvalid)
	}
	img, ok := versions[version]
	if !ok {
		return DBImage{}, fmt.Errorf("%s version %q is not supported: %w", dbType, version, ErrInvalid)
	}
	return img, nil
}

// Postgres identifiers used as POSTGRES_USER/POSTGRES_DB and in the
// healthcheck argv. Conservative on purpose: they also appear inside the
// connection string the api composes.
var dbIdentRe = regexp.MustCompile(`^[a-z_][a-z0-9_]{0,62}$`)

func ValidateDBIdent(field, v string) error {
	if !dbIdentRe.MatchString(v) {
		return fmt.Errorf("%s must match %s: %w", field, dbIdentRe.String(), ErrInvalid)
	}
	return nil
}

// ValidateDBPassword bounds the generated secret: it travels only via the
// container env, but a broken caller shouldn't get to smuggle arbitrary blobs.
func ValidateDBPassword(v string) error {
	if len(v) < 16 || len(v) > 128 {
		return fmt.Errorf("db_password must be 16-128 characters: %w", ErrInvalid)
	}
	for _, r := range v {
		if r <= 0x20 || r >= 0x7f {
			return fmt.Errorf("db_password contains non-printable characters: %w", ErrInvalid)
		}
	}
	return nil
}
