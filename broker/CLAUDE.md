# Broker Service

REST API middleware service. Translates HTTP requests from the frontend into gRPC calls against the Portworx OpenStorage API.

## Entry Points

| File | Purpose |
|------|---------|
| `cmd/api/main.go` | Application bootstrap, flag parsing, gRPC dial, SCS session setup |
| `cmd/api/routes.go` | All HTTP route registrations |
| `cmd/api/handlers.go` | HTTP handler functions |
| `cmd/api/keycloak.go` | Keycloak client initialization and token operations |
| `cmd/api/middleware.go` | `AuthKeycloak` middleware (token validation) |
| `cmd/api/helpers.go` | Shared utility functions |

## Package Structure

```
broker/
├── cmd/api/
│   ├── main.go          # AppConfig struct, startup, gRPC connection
│   ├── routes.go        # chi router setup + all routes
│   ├── handlers.go      # HTTP handler implementations
│   ├── middleware.go    # AuthKeycloak middleware
│   ├── keycloak.go      # Keycloak client + token validation
│   └── helpers.go       # Utilities
└── pkg/
    ├── cluster/
    │   └── cluster.go   # GetClusterInfo, GetClusterCapacity, GetClusterAlarms
    ├── volumes/
    │   └── volumes.go   # Volume CRUD, resize, HA level, IO profile, sharedv4, no-discard, replica set
    ├── nodes/
    │   └── nodes.go     # GetListOfNodes, GetNodeInfo, GetAllNodesInfo
    ├── snapshots/
    │   └── snapshots.go # Cloud/local snapshot creation, listing, deletion
    └── config/
        └── config.go    # All shared structs (ClusterData, VolumeData, NodeData, SnapData, etc.)
```

## API Routes

All protected routes require `Authorization: Bearer <token>` header.

### Auth
| Method | Path | Description |
|--------|------|-------------|
| POST/GET | `/login` | Authenticate with Keycloak, returns session token |
| GET | `/logout` | Invalidate session |
| GET | `/ping` | Health check |

### Cluster
| Method | Path | Description |
|--------|------|-------------|
| GET | `/broker/getpxcluster` | Cluster UUID, name, status |
| GET | `/broker/getpxclustercapacity` | Total/used/available capacity |
| GET | `/broker/getpxclusteralarms` | Active cluster alarms |

### Volumes
| Method | Path | Description |
|--------|------|-------------|
| POST | `/broker/postcreatevolume` | Create new volume |
| GET | `/broker/getallvolumes` | List all volume IDs |
| GET | `/broker/getallvolumesinfo` | List volumes with basic info |
| GET | `/broker/getallvolumescomplete` | Full volume details for all volumes |
| GET | `/broker/getvolumeinfo/{volume_id}` | Single volume details |
| GET | `/broker/getvolumeid/{volume_name}` | Get ID by name |
| GET | `/broker/getnodesofvolume/{volume_name}` | Nodes hosting a volume |
| GET | `/broker/getinspectvolume/{volume_name}` | Deep inspect by name |
| GET | `/broker/getvolumeusage/{volume_name}` | Volume usage stats |
| GET | `/broker/getreplicaspernode/{node_id}` | Replicas on a node |
| PATCH | `/broker/patchvolumesize/{volume_id}` | Resize volume |
| PATCH | `/broker/patchvolumeioprofile/{volume_id}` | Change IO profile |
| PATCH | `/broker/patchvolumehalevel/{volume_id}` | Change HA level (replication factor) |
| PATCH | `/broker/patchvolumesharedv4/{volume_id}` | Toggle sharedv4 |
| PATCH | `/broker/patchvolumesharedv4service/{volume_id}` | Toggle sharedv4 service |
| PATCH | `/broker/patchvolumenodiscard/{volume_id}` | Toggle no-discard flag |
| PATCH | `/broker/patchvolumereplicaset/{volume_id}` | Update replica set nodes |
| DELETE | `/broker/deletevolume/{volume_id}` | Delete volume |

### Nodes
| Method | Path | Description |
|--------|------|-------------|
| GET | `/broker/getlistofnodes` | List node IDs |
| GET | `/broker/getallnodesinfo` | All nodes with details |
| GET | `/broker/getnodeinfo/{node_id}` | Single node details |

### Snapshots & Cloud
| Method | Path | Description |
|--------|------|-------------|
| POST | `/broker/postcreatecloudsnap` | Create cloud snapshot |
| POST | `/broker/postcreatelocalsnap` | Create local snapshot |
| GET | `/broker/getcloudsnaps/{volume_id}` | Cloud snaps for a volume |
| GET | `/broker/getallcloudsnaps` | All cloud snapshots |
| GET | `/broker/getspecificcloudsnapshot` | Specific snapshot details |
| DELETE | `/broker/deletecloudsnap` | Delete cloud snapshot |
| POST | `/broker/postcreateawscloudcreds` | Create AWS cloud credentials |
| GET | `/broker/getinspectawscloudcreds` | Inspect AWS credentials |
| GET | `/broker/getallcloudcredsids` | List all credential IDs |
| DELETE | `/broker/deleteawscloudcreds` | Delete AWS credentials |

## AppConfig Struct

```go
type AppConfig struct {
    Session        *scs.SessionManager
    NewKeycloak    *gocloak.GoCloak
    KeycloakRealm  string
    KeycloakClientID string
    KeycloakSecret string
    PortworxToken  string
    PXConn         *grpc.ClientConn  // gRPC connection to Portworx
}
```

## gRPC Connection

The broker dials Portworx at startup using `grpc.Dial` with `grpc.WithTransportCredentials(insecure.NewCredentials())`. The connection is shared via `AppConfig.PXConn` and passed to service functions in `pkg/`.

- **Dev/mock**: `portworx-openstorage:9100` (mock-sdk-server Docker image)
- **Production**: Real Portworx cluster gRPC endpoint, typically `<node-ip>:9020`

## Adding a New API Endpoint

1. Add handler function to `cmd/api/handlers.go`
2. Add service logic to the appropriate `pkg/<domain>/` package
3. Register route in `cmd/api/routes.go` under the `/broker` route group
4. Add required structs to `pkg/config/config.go`

## Dependencies (key)

```
github.com/go-chi/chi/v5                              - HTTP router
github.com/go-chi/cors                                - CORS middleware
github.com/Nerzal/gocloak/v11                         - Keycloak client
github.com/alexedwards/scs/v2                         - Session manager
github.com/libopenstorage/openstorage-sdk-clients      - Portworx SDK
google.golang.org/grpc                                - gRPC
```

Run `go mod tidy && go mod vendor` after adding new dependencies.
