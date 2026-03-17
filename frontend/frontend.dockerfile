FROM docker.io/golang:1.25 as builder

RUN mkdir /app

COPY ./cmd    /app/cmd
COPY ./static /app/static
COPY ./vendor /app/vendor
COPY go.mod   /app/go.mod
COPY go.sum   /app/go.sum

WORKDIR /app

RUN CGO_ENABLED=0 go build -mod=vendor -o frontendApp ./cmd/web

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
