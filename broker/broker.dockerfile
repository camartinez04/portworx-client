FROM docker.io/golang:1.20 as builder

RUN mkdir /app

COPY ./cmd /app/cmd

COPY ./pkg /app/pkg

WORKDIR /app

RUN go mod init github.com/camartinez04/portworx-client/broker

RUN go get github.com/go-chi/chi/v5 && go get github.com/go-chi/cors && go get google.golang.org/grpc && go get google.golang.org/protobuf && go get github.com/alexedwards/scs/v2 && go get github.com/libopenstorage/openstorage-sdk-clients/sdk/golang && go get github.com/Nerzal/gocloak/v11 && go get golang.org/x/sys && go get golang.org/x/net/idna && go mod tidy

RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

RUN chmod +x /app/brokerApp

# ========================================================================================================================

FROM busybox:latest 

ENV APP_HOME /app

RUN mkdir /app

RUN adduser 1000 -D -h $APP_HOME && mkdir -p $APP_HOME && chown 1000:1000 $APP_HOME

USER 1000

WORKDIR /app

COPY --chown=0:0 --from=builder /app/brokerApp /app

CMD [ "/app/brokerApp"]
