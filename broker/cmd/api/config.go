package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/Nerzal/gocloak/v11"
	"github.com/alexedwards/scs/v2"
	"github.com/camartinez04/portworx-client/broker/pkg/config"
	api "github.com/libopenstorage/openstorage-sdk-clients/sdk/golang"
	"google.golang.org/grpc"
)

var KeycloakURL = os.Getenv("KEYCLOAK_URL")
var KeycloakClientID = os.Getenv("KEYCLOAK_CLIENT_ID")
var KeycloakSecret = os.Getenv("KEYCLOAK_SECRET")
var KeycloakRealm = os.Getenv("KEYCLOAK_REALM")

var (
	UseTls  = flag.Bool("usetls", false, "Connect to server using TLS. Loads CA from the system")
	Token   = flag.String("token", os.Getenv("PORTWORX_TOKEN"), "Authorization token if any")
	Address = flag.String("address", os.Getenv("PORTWORX_GRPC_URL"), "Address to server as <address>:<port>")
)

var Application *AppConfig

var KeycloakToken string

var KeycloakRefreshToken string

var Session *scs.SessionManager

var App AppConfig

const (
	Bytes   = uint64(1)
	KB      = Bytes * uint64(1024)
	MB      = KB * uint64(1024)
	GB      = MB * uint64(1024)
	WebPort = ":8080"
)

// AppConfig holds the application configuration
type AppConfig struct {
	Session      *scs.SessionManager
	Conn         *grpc.ClientConn
	InfoLog      *log.Logger
	ErrorLog     *log.Logger
	InProduction bool
	Models       Models
	NewKeycloak  *Keycloak
}

type Keycloak struct {
	gocloak      gocloak.GoCloak // keycloak client
	clientId     string          // clientId specified in Keycloak
	clientSecret string          // client secret specified in Keycloak
	realm        string          // realm specified in Keycloak
}

type KeyCloakMiddleware struct {
	keycloak *Keycloak
	Session  *scs.SessionManager
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

type Controller struct {
	keycloak *Keycloak
}

// Models holds the models
type Models struct {
	LogEntry LogEntry
}

// LogEntry holds the log entry model
type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type VolumeInspect []string

// jsonResponse is the response format for JSON
type JsonResponse struct {
	Error    bool   `json:"error,omitempty"`
	Message  string `json:"message,omitempty"`
	VolumeID string `json:"volume_id,omitempty"`
	NodeID   string `json:"node_id,omitempty"`
	SnapID   string `json:"snap_id,omitempty"`
	CredID   string `json:"cred_id,omitempty"`
}

// JsonClusterInfo is the response format for JSON for ClusterInfo
type JsonClusterInfo struct {
	Error         bool   `json:"error,omitempty"`
	ClusterUUID   string `json:"cluster_uuid,omitempty"`
	ClusterStatus string `json:"cluster_status,omitempty"`
	ClusterName   string `json:"cluster_name,omitempty"`
}

// JsonClusterCapacity is the response format for JSON for ClusterCapacity
type JsonClusterCapacity struct {
	Error                   bool    `json:"error,omitempty"`
	ClusterCapacity         uint64  `json:"cluster_capacity,omitempty"`
	ClusterUsed             uint64  `json:"cluster_used,omitempty"`
	ClusterAvailable        uint64  `json:"cluster_available,omitempty"`
	ClusterPercentUsed      float64 `json:"cluster_percent_used,omitempty"`
	ClusterPercentAvailable float64 `json:"cluster_percent_available,omitempty"`
}

// JsonGetAllVolumesInfo is the response format for JSON for GetAllVolumesInfo
type JsonGetAllVolumesInfo struct {
	Error          bool                `json:"error,omitempty"`
	AllVolumesInfo []config.VolumeInfo `json:"all_volumes_info,omitempty"`
}

// JsonGetAllNodesInfo is the response format for JSON for GetAllNodesInfo
type JsonGetAllNodesInfo struct {
	Error        bool              `json:"error,omitempty"`
	AllNodesInfo []config.NodeInfo `json:"all_nodes_info,omitempty"`
}

// JsonGetVolumeInfo is the response format for JSON for GetVolumeInfo
type JsonGetVolumeInfo struct {
	Error      bool              `json:"error,omitempty"`
	VolumeInfo config.VolumeInfo `json:"volume_info,omitempty"`
}

// JsonGetNodeInfo is the response format for JSON for GetNodeInfo
type JsonGetNodeInfo struct {
	Error    bool            `json:"error,omitempty"`
	NodeInfo config.NodeInfo `json:"node_info,omitempty"`
}

// JsonVolumeUsage is the response format for JSON for VolumeUsage
type JsonVolumeUsage struct {
	Error                  bool    `json:"error,omitempty"`
	VolumeUsage            float64 `json:"volume_usage,omitempty"`
	AvailableSpace         float64 `json:"available_space,omitempty"`
	TotalSize              float64 `json:"total_size,omitempty"`
	VolumeUsagePercent     float32 `json:"volume_usage_percent,omitempty"`
	VolumeAvailablePercent float32 `json:"volume_available_percent,omitempty"`
}

// JsonVolumeInspect is the response format for JSON for VolumeInspect
type JsonVolumeInspect struct {
	Error              bool     `json:"error,omitempty"`
	VolumeInspect      any      `json:"volume_inspect,omitempty"`
	ReplicasInfo       []string `json:"replicas_info,omitempty"`
	VolumeNodes        []string `json:"volume_nodes,omitempty"`
	VolumeStatusString string   `json:"volume_status_string,omitempty"`
	IoProfileString    string   `json:"io_profile_string,omitempty"`
}

// JsonVolumeList is the response format for JSON for VolumeList
type JsonVolumeList struct {
	Error      bool                         `json:"error,omitempty"`
	VolumeList map[string]config.VolumeInfo `json:"volume_list,omitempty"`
}

// JsonNodeList is the response format for JSON for NodeList
type JsonNodeList struct {
	Error    bool                `json:"error,omitempty"`
	NodeList map[string][]string `json:"node_list,omitempty"`
}

// JsonNodesOfVolume is the response format for JSON for NodesOfVolume
type JsonNodesOfVolume struct {
	Error         bool     `json:"error,omitempty"`
	NodesOfVolume []string `json:"nodes_of_volume,omitempty"`
}

// JsonAllVolumesList is the response format for JSON for AllVolumesList
type JsonAllVolumesList struct {
	Error          bool     `json:"error,omitempty"`
	AllVolumesList []string `json:"all_volumes_list,omitempty"`
}

// JsonApiVolumesList is the response format for JSON for ApiVolumesList
type JsonApiVolumesList struct {
	Error          bool                                     `json:"error,omitempty"`
	ApiVolumesList map[string]*api.SdkVolumeInspectResponse `json:"all_volumes_list,omitempty"`
}

// JsonCloudSnapList is the response format for JSON for CloudSnapList
type JsonCloudSnapList struct {
	Error         bool                                 `json:"error,omitempty"`
	CloudSnapList map[string][]*api.SdkCloudBackupInfo `json:"cloud_snap_list,omitempty"`
}

// JsonCredentialInspect is the response format for JSON for CredentialInspect
type JsonCredentialInspect struct {
	Error             bool                             `json:"error,omitempty"`
	CredentialInspect api.SdkCredentialInspectResponse `json:"credential_inspect,omitempty"`
}

// JsonAlarmList is the response format for JSON for AlarmList
type JsonAlarmList struct {
	Error     bool         `json:"error,omitempty"`
	AlarmList []*api.Alert `json:"alarm_list,omitempty"`
}

// JsonCloudSnap is the response format for JSON for CloudSnap
type JsonCloudSnap struct {
	Error   bool   `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	TaskID  string `json:"task_id,omitempty"`
}

// JsonAllCloudSnapList is the response format for JSON for AllCloudSnapList
type JsonAllCloudSnapList struct {
	Error          bool                                            `json:"error,omitempty"`
	CloudSnapsList map[string]map[string][]*api.SdkCloudBackupInfo `json:"cloud_snaps_list,omitempty"`
}

// JsonSpecificCloudSnap is the response format for JSON for SpecificCloudSnap
type JsonSpecificCloudSnap struct {
	Error       bool                    `json:"error,omitempty"`
	CloudSnap   *api.SdkCloudBackupInfo `json:"cloud_snap,omitempty"`
	CloudSnapId string                  `json:"cloud_snap_id,omitempty"`
}

type OpenStorageSdkToken struct{}

type JsonAllCloudCredsList struct {
	Error          bool     `json:"error,omitempty"`
	CloudCredsList []string `json:"cloud_creds_list,omitempty"`
}
