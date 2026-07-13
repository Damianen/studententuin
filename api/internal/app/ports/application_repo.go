package ports

import (
	"api/internal/domain"
	"context"
)

type ApplicationRepo interface {
	FindByID(id string, context context.Context) (*domain.Application, error)
	// FindBySubdomainID returns the subdomain's application (1:1 today);
	// gorm.ErrRecordNotFound when the subdomain has none.
	FindBySubdomainID(subdomainID string, context context.Context) (*domain.Application, error)
	Create(app *domain.Application, context context.Context) error
	Update(id string, updates map[string]any, context context.Context) error
	Delete(app *domain.Application, context context.Context) error
}
