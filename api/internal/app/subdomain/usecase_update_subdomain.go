package subdomain

import (
	"api/internal/app/ports"
	"context"
	"errors"
	"strings"
)

type UpdateSubdomain struct {
	subdomainRepo ports.SubdomainRepo
	clock ports.Clock
}

type SubdomainUpdateInput struct {
	ID string
	Name *string
	FullDomain *string
}

func (u *UpdateSubdomain) Execute(ctx context.Context, si SubdomainUpdateInput) error {
	now := u.clock.Now()
	updates := map[string]any{
		"updated_at": now,
	}

	err := ports.SetIfNotNil(updates, "full_domain", si.FullDomain, func(s string) (any, error) {
	    s = strings.ToLower(strings.TrimSpace(s))
	    if s == "" {
	        return nil, errors.New("domain cannot be empty")
	    }
	    return s, nil
	})
	if err != nil { return err }

	err = ports.SetIfNotNil(updates, "name", si.Name, func(s string) (any, error) {
	    if s == "" {
	        return nil, errors.New("name cannot be empty")
	    }
	    return s, nil
	})
	if err != nil { return err }

	if len(updates) == 1 {
		return errors.New("no fields to update")
	}

	return u.subdomainRepo.Update(si.ID, updates, ctx)
}
