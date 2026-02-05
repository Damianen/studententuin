package postgres

import (
	"errors"

	"github.com/lib/pq"
)

var ErrEmailAlreadyInUse = errors.New("email already in use")
var ErrFullDomainAlreadyInUse = errors.New("Domain already in use")

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return string(pqErr.Code) == "23505"
	}
	return false
}


