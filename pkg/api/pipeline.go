package api

import (
	"octopipe/pkg/deployment"
	"octopipe/pkg/mozart"

	"github.com/gin-gonic/gin"
)

type PipelineApi struct {
	mozart mozart.UseCases
}

func (api *Api) NewPipelineApi(mozart mozart.UseCases) {
	path := "/deployments"
	controller := PipelineApi{mozart}

	api.v1.POST(path, controller.startPipeline)
}

func (api *PipelineApi) startPipeline(ctx *gin.Context) {
	var deployment *deployment.Deployment
	ctx.Bind(&deployment)

	api.mozart.Start(deployment)

	ctx.JSON(201, nil)
}
