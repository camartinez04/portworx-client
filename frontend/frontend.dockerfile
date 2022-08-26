FROM docker.io/golang:1.19-alpine as builder

RUN mkdir /app

COPY ./cmd /app/cmd

COPY ./static /app/static

WORKDIR /app

RUN go mod init github.com/camartinez04/portworx-client/frontend

RUN go get github.com/go-chi/chi/v5 && go get github.com/go-chi/cors && go get github.com/alexedwards/scs/v2 && go get github.com/justinas/nosurf

RUN CGO_ENABLED=0 go build -o frontendApp ./cmd/web

RUN chmod +x /app/frontendApp

FROM alpine:latest 

RUN mkdir /app

COPY --from=builder /app/frontendApp /app

COPY --from=builder /app/static /app/static

WORKDIR /app

CMD [ "/app/frontendApp"]
