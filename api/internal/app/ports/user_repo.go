package ports

import (
	"api/internal/domain"
	"context"
)

type UserRepo interface {
	FindByID(id string, context context.Context) (*domain.User, error)
	FindByEmail(email string, context context.Context) (*domain.User, error)
	Create(user *domain.User, context context.Context) error
	Delete(user *domain.User, context context.Context) error
	Update(id string, updates map[string]any, context context.Context) error
}
