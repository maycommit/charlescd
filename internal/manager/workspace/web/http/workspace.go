package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/maycommit/circlerr/internal/manager/models"
	"github.com/maycommit/circlerr/internal/manager/workspace"

	"github.com/google/uuid"
)

type WorkspaceHandler struct {
	usecase workspace.UseCase
}

func NewWorkspaceHandler(e *echo.Group, u workspace.UseCase) {
	handler := WorkspaceHandler{
		usecase: u,
	}

	path := "/workspaces"
	e.GET(path, handler.list)
	e.POST(path, handler.save)
	e.GET(fmt.Sprintf("%s/%s", path, ":workspaceId"), handler.getById)
	e.PUT(fmt.Sprintf("%s/%s", path, ":workspaceId"), handler.update)
	e.DELETE(fmt.Sprintf("%s/%s", path, ":workspaceId"), handler.delete)
}

func (h WorkspaceHandler) list(c echo.Context) error {
	workspaces, err := h.usecase.FindAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, workspaces)
}

func (h WorkspaceHandler) save(c echo.Context) error {
	workspace := new(models.Workspace)
	err := c.Bind(workspace)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	createdWorkspace, err := h.usecase.Save(*workspace)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, createdWorkspace)
}

func (h WorkspaceHandler) getById(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("workspaceId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	workspace, err := h.usecase.GetByID(uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, workspace)
}

func (h WorkspaceHandler) update(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("workspaceId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	var workspace models.Workspace
	bindErr := c.Bind(&workspace)
	if bindErr != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	workspace, err = h.usecase.Update(uuid, workspace)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, workspace)
}

func (h WorkspaceHandler) delete(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("workspaceId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = h.usecase.Delete(uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusNoContent, nil)
}
