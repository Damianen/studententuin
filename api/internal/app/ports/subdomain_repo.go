package ports

import (
	"api/internal/domain"
	"context"
)

type SubdomainRepo interface {
	FindByID(id string, context context.Context) (*domain.Subdomain, error)
	Create(subdomain *domain.Subdomain, context context.Context) error
	Update(subdomain *domain.Subdomain, context context.Context) error
	Delete(subdomain *domain.Subdomain, context context.Context) error
}
