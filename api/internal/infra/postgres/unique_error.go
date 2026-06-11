package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var ErrEmailAlreadyInUse = errors.New("email already in use")
var ErrFullDomainAlreadyInUse = errors.New("Domain already in use")

func isUniqueViolation(err error) bool {
	// GORM's postgres driver is backed by pgx, so unique violations
	// surface as *pgconn.PgError, not lib/pq errors.
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
