GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOTOOL=$(GOCMD) tool

KUBECTLCMD=kubectl

DIST_PATH=dist
CMD_CLI=cli/*.go
CMD_CONTROLLER_PATH=cmd/controller/*.go
CMD_GITOPS_PATH=cmd/gitops/*.go
BINARY_NAME=circlerr-controller

CMD_K8SCONTROLLER_PATH=cmd/k8s/controller/*.go
CMD_K8SCONTROLLER_GITOPS_PATH=cmd/k8s/gitops/*.go
CMD_MANAGER_PATH=cmd/manager/*.go


# === K8S Controller ===
k8s-controller-config:
				sh hack/k8s/controller/prepare-development.sh

k8s-controller-start:
				$(GORUN) $(CMD_K8SCONTROLLER_PATH) -k8sconntype=out

k8s-controller-tests:
				$(GOTEST) ./k8scontroller

k8s-controller-e2e-deps-up:
				docker-compose -f k8scontroller/test/docker-compose.test.yaml up -d server node
				$(KUBECTLCMD) --kubeconfig=./k8scontroller/test/kubeconfig.yaml apply -f ./manifests/crds/circle-crd.yaml -f ./manifests/crds/project-crd.yaml
				docker-compose -f k8scontroller/test/docker-compose.test.yaml up -d app

k8s-controller-e2e-deps-down:
				docker-compose -f test/docker-compose.test.yaml down

k8s-controller-e2e:
				make k8scontroller-e2e-deps-up
				$(GOTEST) ./test/e2e
				make k8scontroller-e2e-deps-down

k8s-gitops-start:
				$(GORUN) $(CMD_K8SCONTROLLER_GITOPS_PATH) -k8sconntype=out
# ===

# === Manager ===
manager-deps-up:
				docker-compose -f manager/resources/docker-compose.yaml up -d

manager-deps-down:
				docker-compose -f manager/resources/docker-compose.yaml down

manager-start:
				$(GORUN) $(CMD_MANAGER_PATH)

manager-new-migration:
				sh manager/hack/generate-migration.sh $(filter-out $@,$(MAKECMDGOALS))

# ===
