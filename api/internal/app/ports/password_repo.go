package ports

import (
	"api/internal/domain"
	"context"
)

type PasswordRepo interface {
	FindById(id string, context context.Context) (*domain.PasswordCredential, error)
	Create(pc *domain.PasswordCredential, context context.Context) error
	Update(pc *domain.PasswordCredential, context context.Context) error
	Delete(pc *domain.PasswordCredential, context context.Context) error
}
