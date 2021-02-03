package workspace

import (
	"github.com/google/uuid"
	"github.com/maycommit/circlerr/internal/manager/models"
)

type UseCase interface {
	FindAll() ([]models.Workspace, error)
	Save(workspace models.Workspace) (models.Workspace, error)
	GetByID(id uuid.UUID) (models.Workspace, error)
	Update(id uuid.UUID, workspace models.Workspace) (models.Workspace, error)
	Delete(id uuid.UUID) error
}
