package repository

import (
	"github.com/google/uuid"
	"github.com/maycommit/circlerr/internal/manager/customerror"
	"github.com/maycommit/circlerr/internal/manager/models"
	"github.com/maycommit/circlerr/internal/manager/workspace"
	"gorm.io/gorm"
)

type workspaceRepository struct {
	db *gorm.DB
}

func NewWorkspaceRepository(db *gorm.DB) workspace.Repository {
	return workspaceRepository{db: db}
}

func (r workspaceRepository) FindAll() ([]models.Workspace, error) {
	var workspaces []models.Workspace

	if res := r.db.Preload("Clusters").Find(&workspaces); res.Error != nil {
		return nil, customerror.New("Find all workspaces failed", res.Error, nil, "repository.FindAll.Find")
	}

	return workspaces, nil
}

func (r workspaceRepository) Save(workspace models.Workspace) (models.Workspace, error) {
	if res := r.db.Save(&workspace); res.Error != nil {
		return models.Workspace{}, customerror.New("Save workspace failed", res.Error, nil, "repository.Save.Save")
	}

	return workspace, nil
}

func (r workspaceRepository) GetByID(id uuid.UUID) (models.Workspace, error) {
	var workspace models.Workspace

	if res := r.db.Model(models.Workspace{}).Where("id = ?", id).First(&workspace); res.Error != nil {
		return models.Workspace{}, customerror.New("Find workspace failed", res.Error, nil, "repository.Save.First")
	}

	return workspace, nil
}

func (r workspaceRepository) Update(id uuid.UUID, workspace models.Workspace) (models.Workspace, error) {
	if res := r.db.Model(models.Workspace{}).Where("id = ?", id).Updates(&workspace); res.Error != nil {
		return models.Workspace{}, customerror.New("Update workspace failed", res.Error, nil, "repository.Update.Updates")
	}

	return workspace, nil
}

func (r workspaceRepository) Delete(id uuid.UUID) error {
	if res := r.db.Delete(models.Workspace{}, id); res.Error != nil {
		return customerror.New("Delete workspace failed", res.Error, nil, "repository.Delete.Delete")
	}

	return nil
}
