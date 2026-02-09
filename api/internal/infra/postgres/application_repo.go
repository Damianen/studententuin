package postgres

import (
	"api/internal/domain"
	"context"
	"errors"

	"gorm.io/gorm"
)

type GormApplicationRepo struct {
	DB *gorm.DB
}

func (repo *GormApplicationRepo) FindByID(id string, context context.Context) (*domain.Application, error) {
	var application domain.Application
	err := repo.DB.WithContext(context).First(&application, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &application, nil
}

func (repo *GormApplicationRepo) Create(application *domain.Application, context context.Context) error {
	return repo.DB.WithContext(context).Create(application).Error
}

func (repo *GormApplicationRepo) Update(id string, updates map[string]any, context context.Context) error {
	if len(updates) == 0 {
		return errors.New("nothing to update")
	}
	return repo.DB.WithContext(context).Model(&domain.Application{}).Where("id = ?", id).Updates(updates).Error
}

func (repo *GormApplicationRepo) Delete(application *domain.Application, context context.Context) error {
	return repo.DB.WithContext(context).Delete(application).Error
}
