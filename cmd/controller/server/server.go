package server

import (
	circleclientset "charlescd/pkg/client/clientset/versioned"
	circlepb "charlescd/pkg/grpc/circle"

	"github.com/argoproj/gitops-engine/pkg/cache"
)

type GRPCServer struct {
	ClusterCache    *cache.ClusterCache
	CircleClientset *circleclientset.Clientset
	circlepb.UnimplementedCircleServiceServer
}
