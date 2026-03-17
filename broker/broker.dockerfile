FROM docker.io/golang:1.25 as builder

RUN mkdir /app

COPY ./cmd    /app/cmd
COPY ./pkg    /app/pkg
COPY ./vendor /app/vendor
COPY go.mod   /app/go.mod
COPY go.sum   /app/go.sum

WORKDIR /app

RUN CGO_ENABLED=0 go build -mod=vendor -o brokerApp ./cmd/api

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
