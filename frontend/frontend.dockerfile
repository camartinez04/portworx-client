FROM docker.io/golang:1.25 AS builder

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

ENV APP_HOME=/app

RUN mkdir /app && \
    adduser -D -h /app -u 1000 appuser && \
    chown 1000:1000 /app

USER 1000

WORKDIR /app

COPY --chown=1000:1000 --from=builder /app/frontendApp /app/
COPY --chown=1000:1000 --from=builder /app/static      /app/static

CMD ["/app/frontendApp"]
