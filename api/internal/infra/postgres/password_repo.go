package postgres

import (
	"api/internal/domain"
	"context"
	"gorm.io/gorm"
)

type GormPasswordRepo struct {
	DB *gorm.DB
}

func (repo *GormPasswordRepo) FindById(id string, context context.Context) (*domain.PasswordCredential, error) {
	var passwordCred domain.PasswordCredential
	err := repo.DB.WithContext(context).First(&passwordCred, "user_id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &passwordCred, nil
}

func (repo *GormPasswordRepo) Create(pc *domain.PasswordCredential, context context.Context) error {
	err := repo.DB.WithContext(context).Create(pc).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *GormPasswordRepo) Update(pc *domain.PasswordCredential, context context.Context) error {
	return repo.DB.WithContext(context).Save(pc).Error
}

func (repo *GormPasswordRepo) Delete(pc *domain.PasswordCredential, context context.Context) error {
	return repo.DB.WithContext(context).Delete(pc).Error
}
