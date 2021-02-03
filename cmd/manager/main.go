package main

import (
	"github.com/labstack/echo"
	_circleHttp "github.com/maycommit/circlerr/internal/manager/circle/web/http"
	_manifestHttp "github.com/maycommit/circlerr/internal/manager/manifest/web/http"
	_projectHttp "github.com/maycommit/circlerr/internal/manager/project/web/http"
	_workspaceRepository "github.com/maycommit/circlerr/internal/manager/workspace/repository"
	_workspaceUsecase "github.com/maycommit/circlerr/internal/manager/workspace/usecase"
	_workspaceHttp "github.com/maycommit/circlerr/internal/manager/workspace/web/http"

	_clusterRepository "github.com/maycommit/circlerr/internal/manager/cluster/repository"
	_clusterUsecase "github.com/maycommit/circlerr/internal/manager/cluster/usecase"
	_clusterHttp "github.com/maycommit/circlerr/internal/manager/cluster/web/http"
)

func main() {

	sqlDB, gormDB, err := ConnectDatabase()
	if err != nil {
		panic(err)
	}

	err = RunMigrations(sqlDB)
	if err != nil {
		panic(err)
	}

	workspaceRepository := _workspaceRepository.NewWorkspaceRepository(gormDB)
	workspaceUsecase := _workspaceUsecase.NewWorkspaceUsecase(workspaceRepository)

	clusterRepository := _clusterRepository.NewClusterRepository(gormDB)
	clusterUsecase := _clusterUsecase.NewClusterUsecase(clusterRepository)

	e := echo.New()
	api := e.Group("/api")
	v1 := api.Group("/v1")
	{
		{
			_circleHttp.NewCircleHandler(v1, clusterUsecase)
			_projectHttp.NewProjectHandler(v1, clusterUsecase)
			_manifestHttp.NewManifestHandler(v1, clusterUsecase)
			_clusterHttp.NewClusterHandler(v1, clusterUsecase)
			_workspaceHttp.NewWorkspaceHandler(v1, workspaceUsecase)

		}
	}

	e.Logger.Fatal(e.Start(":8080"))
}
