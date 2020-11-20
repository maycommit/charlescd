GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOTOOL=$(GOCMD) tool

CMD_MANAGER=cmd/manager/*.go
CMD_OPERATOR=cmd/operator/main.go
CMD_CONTROLLER=cmd/controller/main.go

DIST_PATH=dist
CMD_PATH=cmd/*.go
BINARY_NAME=octopipe

start-controller:
				$(GORUN) $(CMD_CONTROLLER)
# build: 
# 				$(GOBUILD) -o $(DIST_PATH)/$(BINARY_NAME) $(CMD_PATH)
# test:
# 				$(GOTEST) ./...
# cover:
# 				$(GOTEST) -coverprofile cover.out ./...
# 				$(GOTOOL) cover -func=cover.out
# cover-browser:
# 				$(GOTEST) -coverprofile cover.out ./...
# 				$(GOTOOL) cover -html=cover.out -o cover.html
# 				open cover.html