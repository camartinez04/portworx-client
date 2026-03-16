# Portworx Client - Project Guide

## Project Overview

A web-based UI client for managing **Portworx** storage clusters. The app is split into two Go microservices:

- **Broker** (`broker/`) — REST API middleware that wraps Portworx's gRPC API
- **Frontend** (`frontend/`) — Server-side rendered HTML web UI using Go templates

## Quick Start

```bash
cd project
make up_build     # Build images and start all services
make down         # Stop all services
```

Access the app at `http://localhost:8082/portworx/`
Default credentials: `pxadmin` / `pxAdmin123$`

## Repository Layout

```
portworx-client/
├── broker/              # REST API service (port 8081)
│   ├── cmd/api/         # main.go, routes.go, handlers, middleware, keycloak
│   ├── pkg/
│   │   ├── cluster/     # Cluster info/capacity/alarms
│   │   ├── volumes/     # Volume CRUD + update operations
│   │   ├── nodes/       # Node enumeration and inspection
│   │   ├── snapshots/   # Cloud and local snapshot management
│   │   └── config/      # Shared data models/structs
│   ├── vendor/          # Vendored Go dependencies
│   └── broker.dockerfile
├── frontend/            # Web UI service (port 8082)
│   ├── cmd/web/         # main.go, routes.go, handlers, repository, render
│   ├── static/
│   │   ├── templates/   # HTML page templates
│   │   ├── css/         # Stylesheets (compiled)
│   │   ├── scss/        # SCSS source
│   │   ├── js/          # JavaScript
│   │   └── vendors/     # Third-party JS/CSS
│   ├── vendor/          # Vendored Go dependencies
│   └── frontend.dockerfile
├── project/             # Docker Compose + Makefile
├── kubernetes/          # K8s deployment manifests
├── openshift/           # OpenShift configs
└── docs/                # Architecture and developer docs
```

## Technology Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.22.2 |
| HTTP Router | go-chi/chi v5 |
| Auth | Keycloak 24.0+ via gocloak v11 |
| Sessions | alexedwards/scs v2.8.0 |
| CSRF | justinas/nosurf |
| Portworx SDK | libopenstorage/openstorage-sdk-clients v0.109.0 |
| gRPC | google.golang.org/grpc v1.63.2 |
| Validation | asaskevich/govalidator |
| Templates | Go html/template (server-side rendered) |
| Container | Docker / Docker Compose |
| Orchestration | Kubernetes / OpenShift |

## Build Commands

```bash
# Full Docker Compose build and start
cd project && make up_build

# Build Linux binaries manually
cd broker && env GOOS=linux CGO_ENABLED=0 go build -o brokerApp ./cmd/api
cd frontend && env GOOS=linux CGO_ENABLED=0 go build -o frontEndApp ./cmd/web

# Push Docker images to Docker Hub
cd project && make push_images
```

## Environment Variables

### Broker
| Variable | Description | Default |
|----------|-------------|---------|
| `PORTWORX_GRPC_URL` | Portworx gRPC endpoint | `localhost:9020` |
| `KEYCLOAK_URL` | Keycloak server URL | required |
| `KEYCLOAK_REALM` | Keycloak realm | `portworx` |
| `KEYCLOAK_CLIENT_ID` | OAuth client ID | required |
| `KEYCLOAK_SECRET` | OAuth client secret | required |
| `PORTWORX_TOKEN` | Optional Portworx auth token | `""` |

### Frontend
| Variable | Description | Default |
|----------|-------------|---------|
| `BROKER_URL` | Broker service URL | `http://localhost:8081` |
| `KEYCLOAK_URL` | Keycloak server URL | required |
| `KEYCLOAK_REALM` | Keycloak realm | `portworx` |
| `KEYCLOAK_CLIENT_ID` | OAuth client ID | required |
| `KEYCLOAK_SECRET` | OAuth client secret | required |

## Authentication Flow

1. User hits frontend login page (`/portworx/login`)
2. Credentials are validated against Keycloak (`postLoginHTTP`)
3. Keycloak issues a JWT access token
4. Token stored in server-side SCS session
5. All protected routes (`/portworx/client/*` and `/broker/*`) use `AuthKeycloak` middleware
6. Middleware validates the token via Keycloak introspection on each request
7. Token expiration triggers automatic redirect to login

## Service Ports

| Service | Port |
|---------|------|
| Frontend | 8082 |
| Broker | 8081 |
| Portworx mock gRPC | 9100 |
| Portworx mock REST | 9110 |

## Docker Images

- `calvarado2004/portworx-client-frontend:latest`
- `calvarado2004/portworx-client-broker:latest`

## Key Conventions

- All broker routes are under `/broker/*` and protected by `AuthKeycloak` middleware
- All frontend protected routes are under `/portworx/client/*`
- Static files served from `./static/` under `/portworx/static/*` and `/portworx/client/static/*`
- Go dependencies are vendored (`vendor/` directories) — always use `go mod vendor` after adding deps
- Session token key in SCS: `"token"` (access token) and `"username"` (logged-in user)
- Broker returns JSON; frontend renders HTML from templates

## Related Documentation

- `broker/CLAUDE.md` — Broker service details, API reference
- `frontend/CLAUDE.md` — Frontend service details, template guide
- `docs/ARCHITECTURE.md` — Full architecture diagram and data flow
- `README.md` — Quick-start and Keycloak setup guide
- Postman API Docs: https://documenter.getpostman.com/view/17794050/VUqpsxJW
