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

type ManifestHandler struct {
	usecase cluster.UseCase
}

func NewManifestHandler(e *echo.Group, u cluster.UseCase) {
	handler := ManifestHandler{
		usecase: u,
	}

	path := "/workspaces/:workspaceId/clusters/:clusterId/manifests/*"
	e.Any(path, handler.listHandler)
}

func (h ManifestHandler) listHandler(c echo.Context) error {

	clusterId, err := uuid.Parse(c.Param("clusterId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	cluster, err := h.usecase.GetByID(clusterId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	splitPath := strings.Split(c.Request().URL.Path, "manifests/")

	c.QueryParams().Add("namespace", "default")

	client := &http.Client{}
	req, err := http.NewRequest(c.Request().Method, fmt.Sprintf("%s/api/v1/manifests/%s?%s", cluster.Address, splitPath[1], c.QueryParams().Encode()), c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	var manifests interface{}
	err = json.NewDecoder(resp.Body).Decode(&manifests)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(resp.StatusCode, manifests)
}
