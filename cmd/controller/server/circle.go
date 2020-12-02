package server

import (
	"context"

	circlepb "charlescd/pkg/grpc/circle"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/health"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *GRPCServer) CircleTree(ctx context.Context, in *circlepb.Circle) (*circlepb.CircleTreeResponse, error) {
	circle, err := s.CircleClientset.CircleV1alpha1().Circles(in.Namespace).Get(context.TODO(), in.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	circleTree := &circlepb.CircleTreeResponse{}
	for _, proj := range circle.Status.Projects {

		resources := []*circlepb.ResourceNode{}
		for _, res := range proj.Resources {
			resKey := kube.NewResourceKey(
				res.Group,
				res.Kind,
				in.Namespace,
				res.Name,
			)

			cl := *s.ClusterCache

			cl.IterateHierarchy(resKey, func(resource *cache.Resource, namespaceResources map[kube.ResourceKey]*cache.Resource) {
				node := circlepb.ResourceNode{}

				node.Ref = &circlepb.ResourceStatus{
					Kind:              resource.Ref.Kind,
					Name:              resource.Ref.Name,
					CreationTimestamp: resource.CreationTimestamp.String(),
				}

				var status *health.HealthStatus
				if resource.Resource != nil {
					status, _ = health.GetResourceHealth(resource.Resource, nil)
				}

				if status != nil {
					statusStr := string(status.Status)
					node.Ref.Health = &circlepb.ResourceHealth{
						Status:  &statusStr,
						Message: &status.Message,
					}
				}

				for _, parent := range resource.OwnerRefs {
					newParent := &circlepb.ResourceParent{
						Kind: parent.Kind,
						Name: parent.Name,
					}

					if parent.Controller != nil {
						newParent.Controller = *parent.Controller
					}

					node.Parents = append(node.Parents, newParent)
				}

				resources = append(resources, &node)
			})
		}

		circleTree.Nodes = append(circleTree.Nodes, &circlepb.ProjectNode{
			Name:      proj.Name,
			Resources: resources,
		})
	}

	return circleTree, nil
}
