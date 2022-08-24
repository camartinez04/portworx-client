# portworx-cli
A Portworx Client with Go

- Built in Go version 1.19
- Uses [libopenstorage](https://github.com/libopenstorage/openstorage-sdk-clients)
- Uses [gRPC](https://pkg.go.dev/google.golang.org/grpc)
- Uses [chi router](https://github.com/go-chi/chi)
- Uses [nosurf](https://github.com/justinas/nosurf)

# API Reference

[api documentation](https://documenter.getpostman.com/view/17794050/VUqpsxJW)

# How to try it

You need Portworx running on your Kubernetes cluster

Forward the portworx-api service that usually will be located on kube-system namespace.

```
kubectl port-forward svc/portworx-api -n kube-system 9020:9020

export PORTWORX_GRPC_URL=localhost:9020


./broker/brokerApp
```

Open a Web Brower and try the routes included on broker/cmd/api/routes.go
