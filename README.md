# CharlesCD

[![GoReport Widget]][GoReport Status]

[GoReport Status]: https://goreportcard.com/report/github.com/maycommit/charlescd
[GoReport Widget]: https://goreportcard.com/badge/github.com/maycommit/charlescd

## Whats is CharlesCD Fork?
This fork uses argo's [gitops-engine](https://github.com/argoproj/gitops-engine) library to build a "circleops" model. Based on kubernetes only.

## TODO

- Create Circle CRD ✅
- Implementing argo gitops-engine ✅
- Continuous resync by interval ✅
- Create API types with kubernetes/code-generator
- Create controller based on kubernetes/sample-controller
- Resync by circle modification unsing k8s informers
- Add support helm charts
- Update resource status in circle CRD
