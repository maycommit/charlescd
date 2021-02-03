# Circlerr Controller installation


## Controller installation
```
kubectl create namespace circlerr
kubectl apply -n circlerr -f https://raw.githubusercontent.com/circlerr/circlerr-k8s-controller/master/manifests/install.yaml
```

## Gitops installation
```
kubectl apply -n circlerr -f https://raw.githubusercontent.com/circlerr/circlerr-k8s-controller/master/manifests/gitops-install.yaml
```
