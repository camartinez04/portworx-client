FROM docker.io/golang:1.22 as builder

RUN mkdir /app

COPY ./cmd /app/cmd

COPY ./static /app/static

WORKDIR /app

RUN go mod init github.com/camartinez04/portworx-client/portworx

RUN go get github.com/go-chi/chi/v5 && go get github.com/go-chi/cors && go get github.com/alexedwards/scs/v2 && go get github.com/justinas/nosurf && go get github.com/libopenstorage/openstorage-sdk-clients/sdk/golang && go get github.com/asaskevich/govalidator && go get github.com/go-chi/cors && go get github.com/Nerzal/gocloak/v11

RUN CGO_ENABLED=0 go build -o frontendApp ./cmd/web

RUN chmod +x /app/frontendApp

# ========================================================================================================================

FROM busybox:latest 

ENV APP_HOME /app

RUN mkdir /app

RUN adduser 1000 -D -h $APP_HOME && mkdir -p $APP_HOME && chown 1000:1000 $APP_HOME

USER 1000

WORKDIR /app

COPY --chown=0:0 --from=builder /app/frontendApp /app

COPY --chown=0:0 --from=builder /app/static /app/static


CMD [ "/app/frontendApp"]
