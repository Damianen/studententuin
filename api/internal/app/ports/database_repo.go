package ports

import (
	"api/internal/domain"
	"context"
)

type DatabaseRepo interface {
	FindByID(id string, context context.Context) (*domain.Database, error)
	// FindBySubdomainID returns the subdomain's database (1:1 today);
	// gorm.ErrRecordNotFound when the subdomain has none.
	FindBySubdomainID(subdomainID string, context context.Context) (*domain.Database, error)
	Create(db *domain.Database, context context.Context) error
	Update(id string, updates map[string]any, context context.Context) error
	Delete(db *domain.Database, context context.Context) error
}
