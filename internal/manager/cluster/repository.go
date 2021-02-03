package cluster

import (
	"github.com/google/uuid"
	"github.com/maycommit/circlerr/internal/manager/models"
)

type Repository interface {
	FindAll(workspaceId uuid.UUID) ([]models.Cluster, error)
	Save(workspaceId uuid.UUID, cluster models.Cluster) (models.Cluster, error)
	GetByID(id uuid.UUID) (models.Cluster, error)
	Update(id uuid.UUID, cluster models.Cluster) (models.Cluster, error)
	Delete(id uuid.UUID) error
}
