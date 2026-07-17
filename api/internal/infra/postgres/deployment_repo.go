package postgres

import (
	"api/internal/domain"
	"context"

	"gorm.io/gorm"
)

type GormDeploymentRepo struct {
	DB *gorm.DB
}

func (repo *GormDeploymentRepo) Create(deployment *domain.Deployment, context context.Context) error {
	return repo.DB.WithContext(context).Create(deployment).Error
}

func (repo *GormDeploymentRepo) UpdateByManagerID(managerID string, updates map[string]any, context context.Context) error {
	// Zero rows affected (a pre-history deploy) is fine: nothing to finalize.
	return repo.DB.WithContext(context).Model(&domain.Deployment{}).
		Where("manager_deployment_id = ?", managerID).Updates(updates).Error
}

func (repo *GormDeploymentRepo) ListByApplication(appID string, limit int, context context.Context) ([]domain.Deployment, error) {
	var deployments []domain.Deployment
	err := repo.DB.WithContext(context).
		Where("application_id = ?", appID).
		Order("started_at DESC").
		Limit(limit).
		Find(&deployments).Error
	if err != nil {
		return nil, err
	}
	return deployments, nil
}

func (repo *GormDeploymentRepo) FailInFlight(reason string, context context.Context) error {
	return repo.DB.WithContext(context).Model(&domain.Deployment{}).
		Where("status = ?", domain.DeploymentStatusInFlight).
		Updates(map[string]any{
			"status": domain.DeploymentStatusFailed,
			"error":  reason,
		}).Error
}
