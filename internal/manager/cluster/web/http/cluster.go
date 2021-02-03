package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/maycommit/circlerr/internal/manager/cluster"
	"github.com/maycommit/circlerr/internal/manager/models"

	"github.com/google/uuid"
)

type ClusterHandler struct {
	usecase cluster.UseCase
}

func NewClusterHandler(e *echo.Group, u cluster.UseCase) {
	handler := ClusterHandler{
		usecase: u,
	}

	path := "/workspaces/:workspaceId/clusters"
	e.GET(path, handler.list)
	e.POST(path, handler.save)
	e.GET(fmt.Sprintf("%s/%s", path, ":clusterId"), handler.getById)
	e.PUT(fmt.Sprintf("%s/%s", path, ":clusterId"), handler.update)
	e.DELETE(fmt.Sprintf("%s/%s", path, ":clusterId"), handler.delete)
}

func (h ClusterHandler) list(c echo.Context) error {
	workspaceUUID, err := uuid.Parse(c.Param("workspaceId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	clusters, err := h.usecase.FindAll(workspaceUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, clusters)
}

func (h ClusterHandler) save(c echo.Context) error {
	workspaceUUID, err := uuid.Parse(c.Param("workspaceId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	var cluster models.Cluster
	err = c.Bind(&cluster)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	createdCluster, err := h.usecase.Save(workspaceUUID, cluster)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, createdCluster)
}

func (h ClusterHandler) getById(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("clusterId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	cluster, err := h.usecase.GetByID(uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, cluster)
}

func (h ClusterHandler) update(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("clusterId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	var cluster models.Cluster
	bindErr := c.Bind(&cluster)
	if bindErr != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	cluster, err = h.usecase.Update(uuid, cluster)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, cluster)
}

func (h ClusterHandler) delete(c echo.Context) error {
	uuid, err := uuid.Parse(c.Param("clusterId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = h.usecase.Delete(uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusNoContent, nil)
}
