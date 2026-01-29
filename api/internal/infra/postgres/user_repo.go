package postgres

import (
	"api/internal/domain"
	"context"

	"gorm.io/gorm"
)


type GormUserRepo struct {
	DB *gorm.DB
}

func (repo *GormUserRepo) FindByID(id string, context context.Context) (*domain.User, error) {
	var user domain.User
	err := repo.DB.WithContext(context).First(&user, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *GormUserRepo) FindByEmail(id string, context context.Context) (*domain.User, error) {
	var user domain.User
	err := repo.DB.WithContext(context).First(&user, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *GormUserRepo) Create(user *domain.User, context context.Context) error {
	err := repo.DB.WithContext(context).Create(user).Error

	if err != nil {
		return err
	}

	return nil
}

func (repo *GormUserRepo) Update(user *domain.User, context context.Context) error {
	return repo.DB.WithContext(context).Save(user).Error
}

func (repo *GormUserRepo) Delete(user *domain.User, context context.Context) error {
	return repo.DB.WithContext(context).Delete(user).Error
}

