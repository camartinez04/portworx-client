FROM docker.io/golang:1.19-alpine as builder

RUN mkdir /app

COPY ./cmd /app/cmd

COPY ./pkg /app/pkg

WORKDIR /app

RUN go mod init github.com/camartinez04/portworx-client/broker

RUN go get github.com/go-chi/chi/v5 && go get github.com/go-chi/cors && go get google.golang.org/grpc && go get google.golang.org/protobuf && go get github.com/alexedwards/scs/v2 && go get github.com/libopenstorage/openstorage-sdk-clients/sdk/golang

RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

RUN chmod +x /app/brokerApp

FROM alpine:latest 

RUN mkdir /app

COPY --from=builder /app/brokerApp /app

CMD [ "/app/brokerApp"]
