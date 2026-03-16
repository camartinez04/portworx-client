# Architecture

## System Overview

```
┌──────────────────────────────────────────────────────────────────┐
│                          User Browser                            │
└───────────────────────────────┬──────────────────────────────────┘
                                │ HTTP :8082
                                ▼
┌──────────────────────────────────────────────────────────────────┐
│                      Frontend Service                            │
│                     (Go 1.22, chi v5)                            │
│                                                                  │
│  ┌─────────────┐  ┌──────────────┐  ┌──────────────────────┐    │
│  │   routes.go  │  │  handlers.go  │  │    repository.go     │    │
│  │             │  │              │  │  (broker HTTP calls) │    │
│  └─────────────┘  └──────────────┘  └──────────────────────┘    │
│  ┌─────────────┐  ┌──────────────┐  ┌──────────────────────┐    │
│  │  render.go   │  │  keycloak.go  │  │    middleware.go     │    │
│  │ (templates) │  │  (auth flow) │  │  (CSRF, session)    │    │
│  └─────────────┘  └──────────────┘  └──────────────────────┘    │
└───────────────────────────────┬──────────────────────────────────┘
          │ token validate      │ HTTP :8081
          │                     ▼
          │         ┌──────────────────────────────────────────────┐
          │         │               Broker Service                 │
          │         │              (Go 1.22, chi v5)               │
          │         │                                              │
          │         │  ┌─────────────┐  ┌──────────────────────┐  │
          │         │  │  routes.go   │  │     handlers.go       │  │
          │         │  └─────────────┘  └──────────────────────┘  │
          │         │  ┌──────────────────────────────────────┐    │
          │         │  │           pkg/ packages               │    │
          │         │  │  cluster | volumes | nodes | snaps   │    │
          │         │  └──────────────────────────────────────┘    │
          │         └───────────────────────┬──────────────────────┘
          │                   gRPC :9020    │
          ▼                                 ▼
┌─────────────────────┐       ┌─────────────────────────────┐
│   Keycloak Server   │       │    Portworx Storage Cluster  │
│  (Identity Mgmt)    │       │    (OpenStorage gRPC API)   │
│  realm: portworx    │       │    or mock-sdk-server        │
└─────────────────────┘       └─────────────────────────────┘
```

## Request Flow

### Unauthenticated Request (Login)

```
Browser → GET /portworx/login
        → Frontend: renders login page (login.page.html)

Browser → POST /portworx/login {username, password}
        → Frontend.PostLoginHTTP()
        → Keycloak.Login(username, password)
        → Keycloak returns JWT access token
        → Frontend stores token in SCS session
        → Redirect to /portworx/client/cluster
```

### Authenticated Request (e.g., View Volumes)

```
Browser → GET /portworx/client/volumes
        → Frontend: AuthKeycloak middleware
          → SCS: load session → get "token"
          → Keycloak: introspect token (still valid?)
          → Valid → proceed; Invalid → redirect /portworx/login

        → Frontend.VolumesHTTP handler
        → Repository.GetAllVolumes(token)
          → HTTP GET http://broker:8081/broker/getallvolumescomplete
          → Header: Authorization: Bearer <token>

        → Broker: AuthKeycloak middleware
          → Validates Bearer token with Keycloak
          → Valid → proceed

        → Broker.getAllVolumesCompleteHTTP handler
        → pkg/volumes.GetAllVolumesComplete(pxConn)
          → gRPC: api.OpenStorageVolumeClient.Enumerate()
          → Returns []VolumeData JSON

        → Frontend: receives JSON, passes to render.RenderTemplate()
        → Template: volumes.page.html rendered with data
        → Browser: HTML response
```

## Authentication Architecture

### Keycloak Configuration

- **Realm**: `portworx`
- **Client ID**: `portworx-frontend` (configurable)
- **Realm export**: `kubernetes/realm-export.json`
- **Default admin user**: `pxadmin` in realm users
- **Token lifetime**: Configured in Keycloak realm settings

### Session Management (SCS)

Both services use `alexedwards/scs/v2`:
- Cookie-based session store
- Session lifetime: 24 hours
- `SameSite=Lax`, `HttpOnly=true`
- Session loaded via `SessionLoad` middleware on every request

### Middleware Chain (Frontend)

```
Request
  └─ cors.Handler          (allow cross-origin)
  └─ middleware.Recoverer   (panic recovery)
  └─ NoSurf                 (CSRF token injection/validation)
  └─ SessionLoad            (SCS session hydration)
  └─ middleware.Heartbeat   (/ping endpoint)
  └─ [route-specific]
       └─ AuthKeycloak      (token validation — protected routes only)
            └─ Handler
```

### Middleware Chain (Broker)

```
Request
  └─ cors.Handler
  └─ middleware.Heartbeat
  └─ middleware.Recoverer
  └─ SessionLoad
  └─ [route-specific]
       └─ AuthKeycloak      (validates Bearer token — /broker/* routes)
            └─ Handler
```

## Data Models

Defined in `broker/pkg/config/config.go` (broker) and `frontend/cmd/web/models.go` (frontend).

### Core Structs

```go
// Cluster
ClusterData { UUID, Status, Name, TotalCapacity, UsedCapacity, ... }
ClusterAlarmData { Id, Severity, Message, Timestamp }

// Volume
VolumeData {
    Id, Name, Size, Format, Status,
    HALevel, IOProfile, Sharedv4, NoDiscard,
    ReplicaSets []ReplicaSet
}
VolumeUsageData { VolumeId, AvailableBytes, TotalBytes, UsedBytes }

// Node
NodeData { Id, SchedulerNodeName, Hostname, Status, ... }
NodeInfoData { pools, disks, network interfaces, ... }

// Snapshots / Cloud
CloudSnapData { Id, VolumeId, Timestamp, Status, CredentialId }
AWSCloudCredData { Id, Name, AccessKey, Endpoint, Region, BucketName }
```

## Volume Operations Detail

The volume update operations follow a consistent pattern:

```
Frontend GET /update-volume-<field>/{volume_id}/{value}
  → Repo.UpdateVolume<Field>HTTP
  → HTTP PATCH broker/patch<field>/{volume_id}  {body: {"<field>": value}}
  → Broker.patchUpdate<Field>HTTP
  → pkg/volumes.Update<Field>(pxConn, volumeId, value)
  → gRPC: api.OpenStorageVolumeClient.Update(volumeSpec)
```

Supported update operations:
- `size` — resize in GB
- `ha_level` — replication factor (1–3)
- `io_profile` — `db`, `sequential`, `auto`, `db_remote`, `cms`
- `sharedv4` — enable/disable NFS-style sharing
- `sharedv4_service` — enable/disable sharedv4 service
- `no_discard` — enable/disable nodiscard mount option
- `replica_set` — update node affinity for replicas

## Kubernetes Deployment

```
kubernetes/
├── keycloak.yaml           # Keycloak StatefulSet + Service
├── postgres.yaml           # PostgreSQL for Keycloak backend
├── pxBrokerDeploy.yaml     # Broker Deployment + Service
├── pxFrontendDeploy.yaml   # Frontend Deployment + Service
└── realm-export.json       # Keycloak realm config (import on setup)
```

### Kubernetes Environment Variables

Broker deployment sets:
- `PORTWORX_GRPC_URL` — actual Portworx cluster gRPC address

Frontend deployment sets:
- `BROKER_URL` — internal service URL (e.g., `http://broker-service:8081`)
- `KEYCLOAK_URL`, `KEYCLOAK_REALM`, `KEYCLOAK_CLIENT_ID`, `KEYCLOAK_SECRET`

## Docker Images

Built via multi-stage Dockerfiles:
1. Stage 1: `golang:1.22` — compile Go binary with `CGO_ENABLED=0`
2. Stage 2: `busybox:latest` — minimal runtime image

Images pushed to Docker Hub:
- `calvarado2004/portworx-client-broker:latest`
- `calvarado2004/portworx-client-frontend:latest`

## Local Development with Mock Portworx

The `project/docker-compose.yaml` includes `openstorage/mock-sdk-server` which simulates the Portworx gRPC API at port `9100`. This allows full local development without a real Portworx cluster.

```bash
cd project
make up_build
# Frontend: http://localhost:8082/portworx/
# Broker:   http://localhost:8081/broker/
# Mock PX:  localhost:9100 (gRPC)
```

**Note**: The mock server does not require Keycloak. Set `KEYCLOAK_SECRET` in the docker-compose if using Keycloak locally.
