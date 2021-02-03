package workspace

import (
	"github.com/google/uuid"
	"github.com/maycommit/circlerr/internal/manager/models"
)

type Repository interface {
	FindAll() ([]models.Workspace, error)
	Save(user models.Workspace) (models.Workspace, error)
	GetByID(id uuid.UUID) (models.Workspace, error)
	Update(id uuid.UUID, user models.Workspace) (models.Workspace, error)
	Delete(id uuid.UUID) error
}
