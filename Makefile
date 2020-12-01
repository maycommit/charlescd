GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOTOOL=$(GOCMD) tool

CMD_OPERATOR=cmd/operator/*.go
CMD_CONTROLLER=cmd/controller/*.go

DIST_PATH=dist
CMD_PATH=cmd/*.go
BINARY_NAME=octopipe

pre-config:
				sh hack/prepare-development.sh

start-controller:
				$(GORUN) $(CMD_CONTROLLER)

start-manager:
				$(GORUN) $(CMD_CONTROLLER)
