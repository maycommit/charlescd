package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/maycommit/circlerr/internal/manager/cluster"
)

type ProjectHandler struct {
	usecase cluster.UseCase
}

func NewProjectHandler(e *echo.Group, u cluster.UseCase) {
	handler := ProjectHandler{
		usecase: u,
	}

	path := "/workspaces/:workspaceId/clusters/:clusterId/projects"
	e.Any(path, handler.listHandler)
	e.Any(fmt.Sprintf("%s/*", path), handler.singleHandler)
}

func (h ProjectHandler) listHandler(c echo.Context) error {

	clusterId, err := uuid.Parse(c.Param("clusterId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	cluster, err := h.usecase.GetByID(clusterId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(c.Request().Method, fmt.Sprintf("%s/api/v1/projects?namespace=default", cluster.Address), c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	var projects []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&projects)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(resp.StatusCode, projects)
}

func (h ProjectHandler) singleHandler(c echo.Context) error {

	clusterId, err := uuid.Parse(c.Param("clusterId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	cluster, err := h.usecase.GetByID(clusterId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	splitPath := strings.Split(c.Request().URL.Path, "projects")

	fmt.Println(fmt.Sprintf("%s/api/v1/projects%s?namespace=default", cluster.Address, splitPath[1]))

	client := &http.Client{}
	req, err := http.NewRequest(c.Request().Method, fmt.Sprintf("%s/api/v1/projects%s?namespace=default", cluster.Address, splitPath[1]), c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	var project map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&project)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(resp.StatusCode, project)
}
