package ports

import (
	"api/internal/domain"
	"context"
)

type DatabaseRepo interface {
	FindByID(id string, context context.Context) (*domain.Database, error)
	Create(db *domain.Database, context context.Context) error
	Update(db *domain.Database, context context.Context) error
	Delete(db *domain.Database, context context.Context) error
}
