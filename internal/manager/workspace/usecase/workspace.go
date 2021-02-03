package usecase

import (
	"github.com/google/uuid"
	"github.com/maycommit/circlerr/internal/manager/models"
	"github.com/maycommit/circlerr/internal/manager/workspace"
)

type workspaceUsecase struct {
	workspaceRepository workspace.Repository
}

func NewWorkspaceUsecase(r workspace.Repository) workspace.UseCase {
	return workspaceUsecase{
		workspaceRepository: r,
	}
}

func (u workspaceUsecase) FindAll() ([]models.Workspace, error) {
	return u.workspaceRepository.FindAll()
}

func (u workspaceUsecase) Save(workspace models.Workspace) (models.Workspace, error) {
	return u.workspaceRepository.Save(workspace)
}

func (u workspaceUsecase) GetByID(id uuid.UUID) (models.Workspace, error) {
	return u.workspaceRepository.GetByID(id)
}

func (u workspaceUsecase) Update(id uuid.UUID, workspace models.Workspace) (models.Workspace, error) {
	return u.workspaceRepository.Update(id, workspace)
}

func (u workspaceUsecase) Delete(id uuid.UUID) error {
	return u.workspaceRepository.Delete(id)
}
