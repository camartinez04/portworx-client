# portworx-cli
A Portworx Client with Go

- Built in Go version 1.19
- Uses [libopenstorage](https://github.com/libopenstorage/openstorage-sdk-clients)
- Uses [gRPC](https://pkg.go.dev/google.golang.org/grpc)
- Uses [chi router](https://github.com/go-chi/chi)
- Uses [nosurf](https://github.com/justinas/nosurf)

# API Reference

[api documentation](https://documenter.getpostman.com/view/17794050/VUqpsxJW)

# Docker Compose way for developing

Have Docker running on your laptop and docker-compose installed as well.

```
cd project

make up_build

docker-compose ps

```

Create a test volume once the project is up, we will be using the mock service from OpenStorage as gRPC endpoint.

```
curl --location --request POST 'localhost:8080/postcreatevolume' \
--header 'Volume-Name: postman-volume' \
--header 'Volume-Size: 10' \
--header 'Volume-Ha-Level: 2' \
--header 'Volume-Encryption-Enabled: true' \
--header 'Volume-Sharedv4-Enabled: false' \
--header 'Volume-No-Discard: true'

```

Check on the frontend that the mock volume was created:

[http://localhost:8081/frontend/volume/postman-volume](http://localhost:8081/frontend/volume/postman-volume) 

# Test the Broker on Kubernetes

You need Portworx running on your Kubernetes cluster

Forward the portworx-api service that usually will be located on kube-system namespace.

```
kubectl port-forward svc/portworx-api -n kube-system 9020:9020

export PORTWORX_GRPC_URL=localhost:9020


./broker/brokerApp
```

Open a Web Brower and try the routes included on broker/cmd/api/routes.go
