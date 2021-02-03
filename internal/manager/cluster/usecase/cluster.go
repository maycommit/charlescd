package usecase

import (
	"github.com/google/uuid"
	"github.com/maycommit/circlerr/internal/manager/cluster"
	"github.com/maycommit/circlerr/internal/manager/models"
)

type clusterUsecase struct {
	clusterRepository cluster.Repository
}

func NewClusterUsecase(r cluster.Repository) cluster.UseCase {
	return clusterUsecase{
		clusterRepository: r,
	}
}

func (u clusterUsecase) FindAll(workspaceId uuid.UUID) ([]models.Cluster, error) {
	return u.clusterRepository.FindAll(workspaceId)
}

func (u clusterUsecase) Save(workspaceId uuid.UUID, cluster models.Cluster) (models.Cluster, error) {
	return u.clusterRepository.Save(workspaceId, cluster)
}

func (u clusterUsecase) GetByID(id uuid.UUID) (models.Cluster, error) {
	return u.clusterRepository.GetByID(id)
}

func (u clusterUsecase) Update(id uuid.UUID, cluster models.Cluster) (models.Cluster, error) {
	return u.clusterRepository.Update(id, cluster)
}

func (u clusterUsecase) Delete(id uuid.UUID) error {
	return u.clusterRepository.Delete(id)
}
