# Frontend Service

Server-side rendered web UI. Handles user sessions, renders HTML templates, and calls the broker REST API for all data.

## Entry Points

| File | Purpose |
|------|---------|
| `cmd/web/main.go` | Bootstrap, session setup, template cache, HTTP server |
| `cmd/web/routes.go` | chi router + all route registrations |
| `cmd/web/handlers.go` | HTTP handlers for all pages |
| `cmd/web/repository.go` | Data access layer — all broker API calls |
| `cmd/web/render.go` | Template rendering helper |
| `cmd/web/middleware.go` | `NoSurf` (CSRF) and `SessionLoad` middleware |
| `cmd/web/keycloak.go` | Keycloak authentication logic |
| `cmd/web/models.go` | Frontend-specific data structs and template data |
| `cmd/web/forms.go` | Form struct, validation helpers |
| `cmd/web/helpers.go` | Utility functions (byte→GB, date formatting, etc.) |

## AppConfig & Repository Pattern

```go
// AppConfig holds application-wide dependencies
type AppConfig struct {
    TemplateCache   map[string]*template.Template
    Session         *scs.SessionManager
    NewKeycloak     *gocloak.GoCloak
    KeycloakRealm   string
    KeycloakClientID string
    KeycloakSecret  string
    BrokerURL       string
}

// DBRepo is the interface for all data operations
type DBRepo interface {
    GetLoginHTTP(w, r)
    PostLoginHTTP(w, r)
    ClusterHTTP(w, r)
    VolumesHTTP(w, r)
    // ... all handler methods
}

// Repository is the global Repo variable, assigned at startup
var Repo *Repository
```

All handlers live on `*Repository`. The `App` field inside `Repository` provides access to session, keycloak, templates, and config.

## Routes

### Public (no auth)
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET/POST | `/` | `GetLoginHTTP` / `PostLoginHTTP` | Root login |
| GET/POST | `/portworx/` | same | Login page |
| GET/POST | `/portworx/login` | same | Login page |
| GET | `/portworx/static/*` | FileServer | Static assets |

### Protected (`/portworx/client/*`) — requires `AuthKeycloak`
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/` | `ClusterHTTP` | Cluster dashboard |
| GET | `/cluster` | `ClusterHTTP` | Cluster dashboard |
| GET/POST | `/logout` | `LogoutHTTP` | Logout and clear session |
| GET | `/volumes` | `VolumesHTTP` | All volumes list |
| GET | `/volume/{volume_id}` | `VolumeInformationHTTP` | Single volume detail |
| GET | `/nodes` | `NodesHTTP` | All nodes list |
| GET | `/node/{node_id}` | `NodeInformationHTTP` | Single node detail |
| GET | `/snapshots` | `GetAllSnapsHTTP` | All cloud snapshots |
| GET | `/snapshot/{cred_id}/{bucket}/{snap_id}` | `SpecificSpapInformationHTTP` | Snapshot detail |
| GET | `/cloud-credentials` | `CloudCredentialsHTTP` | Credentials list |
| GET | `/cloud-credential/{cloud_cred_id}` | `CloudCredentialsInformationHTTP` | Credential detail |
| GET/POST | `/create-credentials` | `CreateCloudCredentialsHTTP` / `PostCreateCloudCredentialsHTTP` | New AWS creds form |
| GET/POST | `/create-volume` | `CreateVolumeHTTP` / `PostCreateVolumeHTTP` | New volume form |
| GET/POST | `/create-cloudsnap` | `CreateCloudSnapHTTP` / `PostCreateCloudSnapHTTP` | New cloud snap form |
| GET | `/delete-volume/{volume_id}` | `DeleteVolumeHTTP` | Delete volume (confirm) |
| GET | `/delete-cloudsnap/{bucket}/{snap_id}` | `DeleteCloudSnapHTTP` | Delete snapshot |
| GET | `/update-volume-halevel/{volume_id}/{ha-level}` | `UpdateVolumeHALevelHTTP` | Update HA level |
| GET | `/update-volume-size/{volume_id}/{size}` | `UpdateVolumeSizeHTTP` | Resize volume |
| GET | `/update-volume-ioprofile/{volume_id}/{ioprofile}` | `UpdateVolumeIOProfileHTTP` | Change IO profile |
| GET | `/portworx/client/static/*` | FileServer | Static assets (authenticated) |

## Templates

All templates are in `static/templates/`. The template cache is built at startup by `render.go`.

| Template File | Page |
|--------------|------|
| `base-template.layout.html` | Main layout (nav, head, scripts) |
| `login.page.html` | Login form |
| `index.page.html` | Cluster dashboard |
| `volumes.page.html` | Volume list table |
| `volume-specific.page.html` | Volume detail + update controls |
| `nodes.page.html` | Node list table |
| `node-specific.page.html` | Node detail view |
| `snapshots.page.html` | Cloud snapshots list |
| `snap-specific.page.html` | Snapshot detail |
| `create-volume.page.html` | Volume creation form |
| `create-credentials.page.html` | AWS credential form |
| `create-cloudsnap.page.html` | Cloud snapshot creation form |
| `cloud-credentials.page.html` | Credentials list |
| `cloud-credential-specific.page.html` | Credential detail |

### Template Data (`TemplateData`)

```go
type TemplateData struct {
    StringMap       map[string]string
    IntMap          map[string]int
    FloatMap        map[string]float32
    Data            map[string]interface{}
    CSRFToken       string
    Flash           string
    Warning         string
    Error           string
    Form            *forms.Form
    IsAuthenticated int
}
```

Pass data to templates via `render.RenderTemplate(w, r, "page-name.page.html", &models.TemplateData{...})`.

## Session Keys

| Key | Type | Description |
|-----|------|-------------|
| `"token"` | string | Keycloak access token |
| `"username"` | string | Logged-in username |
| `"refresh_token"` | string | Keycloak refresh token |

## Static Assets

```
static/
├── css/          # Compiled CSS (Bootstrap-based)
├── scss/         # SCSS source files
├── js/           # Application JavaScript (jQuery, custom)
├── vendors/      # AdminLTE, Bootstrap, Font Awesome, etc.
└── templates/    # HTML Go templates
```

Static files are served at:
- `/portworx/static/*` (public, pre-auth)
- `/portworx/client/static/*` (authenticated)

## Broker API Calls

All broker calls are in `cmd/web/repository.go`. Pattern:

```go
func (m *Repository) callBrokerEndpoint(token string, ...) (*SomeStruct, error) {
    req, _ := http.NewRequest("GET", m.App.BrokerURL + "/broker/endpoint", nil)
    req.Header.Set("Authorization", "Bearer " + token)
    client := &http.Client{}
    resp, err := client.Do(req)
    // decode JSON response
}
```

The token is retrieved from the SCS session:
```go
token := m.App.Session.GetString(r.Context(), "token")
```

## Adding a New Page

1. Create handler in `cmd/web/handlers.go`
2. Add broker call in `cmd/web/repository.go`
3. Create template in `static/templates/<name>.page.html`
4. Register route in `cmd/web/routes.go`
5. Add nav link in `base-template.layout.html`

## Dependencies (key)

```
github.com/go-chi/chi/v5          - HTTP router
github.com/alexedwards/scs/v2     - Session management
github.com/justinas/nosurf        - CSRF protection
github.com/Nerzal/gocloak/v11     - Keycloak client
github.com/asaskevich/govalidator  - Form validation
```
