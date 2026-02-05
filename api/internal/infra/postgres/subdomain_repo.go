package postgres

import (
	"api/internal/domain"
	"context"
	"errors"

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

func (repo *GormSubdomainRepo) FindByFullDomain(fullDomain string, context context.Context) (*domain.Subdomain, error) {
	var subdomain domain.Subdomain
	err := repo.DB.WithContext(context).First(&subdomain, "full_domain = ?", fullDomain).Error

	if err != nil {
		return nil, err
	}

	return &subdomain, nil
}

func (repo *GormSubdomainRepo) Create(subdomain *domain.Subdomain, context context.Context) error {
	err := repo.DB.WithContext(context).Create(subdomain).Error

	if isUniqueViolation(err) {
		return ErrFullDomainAlreadyInUse
	}
	return err
}

func (repo *GormSubdomainRepo) Update(id string, updates map[string]any, context context.Context) error {
	if len(updates) == 0 {
		return errors.New("nothing to update")
	}
	return repo.DB.WithContext(context).Model(&domain.Subdomain{}).Where("id = ?", id).Updates(updates).Error
}

func (repo *GormSubdomainRepo) Delete(subdomain *domain.Subdomain, context context.Context) error {
	return repo.DB.WithContext(context).Delete(subdomain).Error
}
