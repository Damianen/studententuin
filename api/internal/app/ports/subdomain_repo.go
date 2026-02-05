package ports

import (
	"api/internal/domain"
	"context"
)

type SubdomainRepo interface {
	FindByID(id string, context context.Context) (*domain.Subdomain, error)
	FindByFullDomain(domain string, context context.Context) (*domain.Subdomain, error)
	Create(subdomain *domain.Subdomain, context context.Context) error
	Update(id string, updates map[string]any, context context.Context) error
	Delete(subdomain *domain.Subdomain, context context.Context) error
}
