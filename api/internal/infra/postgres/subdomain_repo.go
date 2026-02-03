package postgres

import (
	"api/internal/domain"
	"context"

	"gorm.io/gorm"
)

type GormSubdomainRepo struct {
	DB *gorm.DB
}

func (repo *GormSubdomainRepo) FindByID(id string, context context.Context) (*domain.Subdomain, error) {
	var subdomain domain.Subdomain
	err := repo.DB.WithContext(context).First(&subdomain, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &subdomain, nil
}

func (repo *GormSubdomainRepo) Create(subdomain *domain.Subdomain, context context.Context) error {
	err := repo.DB.WithContext(context).Create(subdomain).Error

	if err != nil {
		return err
	}

	return nil
}

func (repo *GormSubdomainRepo) Update(subdomain *domain.Subdomain, context context.Context) error {
	return repo.DB.WithContext(context).Save(subdomain).Error
}

func (repo *GormSubdomainRepo) Delete(subdomain *domain.Subdomain, context context.Context) error {
	return repo.DB.WithContext(context).Delete(subdomain).Error
}
