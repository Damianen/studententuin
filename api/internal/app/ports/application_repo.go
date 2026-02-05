package ports

import (
	"api/internal/domain"
	"context"
)

type ApplicationRepo interface {
	FindByID(id string, context context.Context) (*domain.Application, error)
	Create(app *domain.Application, context context.Context) error
	Update(id string, updates map[string]any, context context.Context) error
	Delete(app *domain.Application, context context.Context) error
}
