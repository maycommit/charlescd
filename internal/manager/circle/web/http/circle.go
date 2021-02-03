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

type CircleHandler struct {
	usecase cluster.UseCase
}

func NewCircleHandler(e *echo.Group, u cluster.UseCase) {
	handler := CircleHandler{
		usecase: u,
	}

	path := "/workspaces/:workspaceId/clusters/:clusterId/circles"
	e.Any(path, handler.listHandler)
	e.Any(fmt.Sprintf("%s/*", path), handler.singleHandler)
}

func (h CircleHandler) listHandler(c echo.Context) error {

	clusterId, err := uuid.Parse(c.Param("clusterId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	cluster, err := h.usecase.GetByID(clusterId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(c.Request().Method, fmt.Sprintf("%s/api/v1/circles?namespace=default", cluster.Address), c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	var circles []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&circles)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(resp.StatusCode, circles)
}

func (h CircleHandler) singleHandler(c echo.Context) error {

	clusterId, err := uuid.Parse(c.Param("clusterId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	cluster, err := h.usecase.GetByID(clusterId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	splitPath := strings.Split(c.Request().URL.Path, "circles")

	client := &http.Client{}
	req, err := http.NewRequest(c.Request().Method, fmt.Sprintf("%s/api/v1/circles%s?namespace=default", cluster.Address, splitPath[1]), c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	var circle interface{}
	err = json.NewDecoder(resp.Body).Decode(&circle)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(resp.StatusCode, circle)
}
