package postgres

import (
	"api/internal/domain"
	"context"
	"errors"

	"gorm.io/gorm"
)

type GormDatabaseRepo struct {
	DB *gorm.DB
}

func (repo *GormDatabaseRepo) FindByID(id string, context context.Context) (*domain.Database, error) {
	var database domain.Database
	err := repo.DB.WithContext(context).First(&database, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &database, nil
}

func (repo *GormDatabaseRepo) FindBySubdomainID(subdomainID string, context context.Context) (*domain.Database, error) {
	var database domain.Database
	err := repo.DB.WithContext(context).First(&database, "subdomain_id = ?", subdomainID).Error
	if err != nil {
		return nil, err
	}
	return &database, nil
}

func (repo *GormDatabaseRepo) Create(database *domain.Database, context context.Context) error {
	return repo.DB.WithContext(context).Create(database).Error
}

func (repo *GormDatabaseRepo) Update(id string, updates map[string]any, context context.Context) error {
	if len(updates) == 0 {
		return errors.New("nothing to update")
	}
	return repo.DB.WithContext(context).Model(&domain.Database{}).Where("id = ?", id).Updates(updates).Error
}

func (repo *GormDatabaseRepo) Delete(database *domain.Database, context context.Context) error {
	return repo.DB.WithContext(context).Delete(database).Error
}
