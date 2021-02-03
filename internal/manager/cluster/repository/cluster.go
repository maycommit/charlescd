package repository

import (
	"github.com/google/uuid"
	"github.com/maycommit/circlerr/internal/manager/cluster"
	"github.com/maycommit/circlerr/internal/manager/customerror"
	"github.com/maycommit/circlerr/internal/manager/models"
	"gorm.io/gorm"
)

type clusterRepository struct {
	db *gorm.DB
}

func NewClusterRepository(db *gorm.DB) cluster.Repository {
	return clusterRepository{db: db}
}

func (r clusterRepository) FindAll(workspaceId uuid.UUID) ([]models.Cluster, error) {
	var clusters []models.Cluster

	if res := r.db.Where("workspace_id = ?", workspaceId).Find(&clusters); res.Error != nil {
		return nil, customerror.New("Find all clusters failed", res.Error, nil, "repository.FindAll.Find")
	}

	return clusters, nil
}

func (r clusterRepository) Save(workspaceId uuid.UUID, cluster models.Cluster) (models.Cluster, error) {
	cluster.WorkspaceID = workspaceId
	if res := r.db.Save(&cluster); res.Error != nil {
		return models.Cluster{}, customerror.New("Save cluster failed", res.Error, nil, "repository.Save.Save")
	}

	return cluster, nil
}

func (r clusterRepository) GetByID(id uuid.UUID) (models.Cluster, error) {
	var cluster models.Cluster

	if res := r.db.Model(models.Cluster{}).Where("id = ?", id).First(&cluster); res.Error != nil {
		return models.Cluster{}, customerror.New("Find cluster failed", res.Error, nil, "repository.Save.First")
	}

	return cluster, nil
}

func (r clusterRepository) Update(id uuid.UUID, cluster models.Cluster) (models.Cluster, error) {
	if res := r.db.Model(models.Cluster{}).Where("id = ?", id).Updates(&cluster); res.Error != nil {
		return models.Cluster{}, customerror.New("Update cluster failed", res.Error, nil, "repository.Update.Updates")
	}

	return cluster, nil
}

func (r clusterRepository) Delete(id uuid.UUID) error {
	if res := r.db.Delete(models.Cluster{}, id); res.Error != nil {
		return customerror.New("Delete cluster failed", res.Error, nil, "repository.Delete.Delete")
	}

	return nil
}
