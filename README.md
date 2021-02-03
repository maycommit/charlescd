# Circlerr Controller

[![GoReport Widget]][GoReport Status]

[GoReport Status]: https://goreportcard.com/report/github.com/maycommit/circlerr
[GoReport Widget]: https://goreportcard.com/badge/github.com/maycommit/circlerr

## Whats is Circlerr Controller?
This project uses argo's [gitops-engine](https://github.com/argoproj/gitops-engine) library to build a "circleops" model. Based on kubernetes only.

## Installation
[Instructions](https://github.com/circlerr/circlerr-k8s-controller/tree/master/manifests)

## Development
Start your kubernetes locally, example: minikube etc...

### Start controller
```
make config
make start-controller
```

### Start gitops
```
make start-gitops
```
