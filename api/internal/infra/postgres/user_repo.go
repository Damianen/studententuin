package postgres

import (
	"api/internal/domain"
	"context"
	"errors"

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
	err := repo.DB.WithContext(context).First(&user, "email = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *GormUserRepo) Create(user *domain.User, context context.Context) error {
	err := repo.DB.WithContext(context).Create(user).Error
	if isUniqueViolation(err) {
		return ErrEmailAlreadyInUse
	}
	return err
}

func (repo *GormUserRepo) Update(id string, updates map[string]any, context context.Context) error {
	if len(updates) == 0 {
		return errors.New("nothing to update")
	}
	return repo.DB.WithContext(context).Model(&domain.User{}).Where("id = ?", id).Updates(updates).Error
}

func (repo *GormUserRepo) Delete(user *domain.User, context context.Context) error {
	return repo.DB.WithContext(context).Delete(user).Error
}

