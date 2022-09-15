package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/camartinez04/portworx-client/broker/pkg/config"
	api "github.com/libopenstorage/openstorage-sdk-clients/sdk/golang"
	"google.golang.org/grpc"
)

type AppConfig struct {
	Session      *scs.SessionManager
	Conn         *grpc.ClientConn
	InfoLog      *log.Logger
	ErrorLog     *log.Logger
	InProduction bool
	Models       Models
}

type Models struct {
	LogEntry LogEntry
}

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
}

type JsonClusterInfo struct {
	Error         bool   `json:"error,omitempty"`
	ClusterUUID   string `json:"cluster_uuid,omitempty"`
	ClusterStatus string `json:"cluster_status,omitempty"`
	ClusterName   string `json:"cluster_name,omitempty"`
}

type JsonClusterCapacity struct {
	Error                   bool    `json:"error,omitempty"`
	ClusterCapacity         uint64  `json:"cluster_capacity,omitempty"`
	ClusterUsed             uint64  `json:"cluster_used,omitempty"`
	ClusterAvailable        uint64  `json:"cluster_available,omitempty"`
	ClusterPercentUsed      float64 `json:"cluster_percent_used,omitempty"`
	ClusterPercentAvailable float64 `json:"cluster_percent_available,omitempty"`
}

type JsonGetAllVolumesInfo struct {
	Error          bool                `json:"error,omitempty"`
	AllVolumesInfo []config.VolumeInfo `json:"all_volumes_info,omitempty"`
}

type JsonGetAllNodesInfo struct {
	Error        bool              `json:"error,omitempty"`
	AllNodesInfo []config.NodeInfo `json:"all_nodes_info,omitempty"`
}

type JsonGetVolumeInfo struct {
	Error      bool              `json:"error,omitempty"`
	VolumeInfo config.VolumeInfo `json:"volume_info,omitempty"`
}

type JsonGetNodeInfo struct {
	Error    bool            `json:"error,omitempty"`
	NodeInfo config.NodeInfo `json:"node_info,omitempty"`
}

type JsonVolumeUsage struct {
	Error                  bool    `json:"error,omitempty"`
	VolumeUsage            float64 `json:"volume_usage,omitempty"`
	AvailableSpace         float64 `json:"available_space,omitempty"`
	TotalSize              float64 `json:"total_size,omitempty"`
	VolumeUsagePercent     float32 `json:"volume_usage_percent,omitempty"`
	VolumeAvailablePercent float32 `json:"volume_available_percent,omitempty"`
}

type JsonVolumeInspect struct {
	Error              bool     `json:"error,omitempty"`
	VolumeInspect      any      `json:"volume_inspect,omitempty"`
	ReplicasInfo       []string `json:"replicas_info,omitempty"`
	VolumeNodes        []string `json:"volume_nodes,omitempty"`
	VolumeStatusString string   `json:"volume_status_string,omitempty"`
	IoProfileString    string   `json:"io_profile_string,omitempty"`
}

type JsonVolumeList struct {
	Error      bool                         `json:"error,omitempty"`
	VolumeList map[string]config.VolumeInfo `json:"volume_list,omitempty"`
}

type JsonNodeList struct {
	Error    bool                `json:"error,omitempty"`
	NodeList map[string][]string `json:"node_list,omitempty"`
}

type JsonNodesOfVolume struct {
	Error         bool     `json:"error,omitempty"`
	NodesOfVolume []string `json:"nodes_of_volume,omitempty"`
}

type JsonAllVolumesList struct {
	Error          bool     `json:"error,omitempty"`
	AllVolumesList []string `json:"all_volumes_list,omitempty"`
}

type JsonApiVolumesList struct {
	Error          bool                                     `json:"error,omitempty"`
	ApiVolumesList map[string]*api.SdkVolumeInspectResponse `json:"all_volumes_list,omitempty"`
}

type JsonCloudSnapList struct {
	Error         bool                      `json:"error,omitempty"`
	CloudSnapList []*api.SdkCloudBackupInfo `json:"cloud_snap_list,omitempty"`
}

type JsonCredentialInspect struct {
	Error             bool                             `json:"error,omitempty"`
	CredentialInspect api.SdkCredentialInspectResponse `json:"credential_inspect,omitempty"`
}

const (
	Bytes   = uint64(1)
	KB      = Bytes * uint64(1024)
	MB      = KB * uint64(1024)
	GB      = MB * uint64(1024)
	webPort = ":8080"
)

var (
	useTls  = flag.Bool("usetls", false, "Connect to server using TLS. Loads CA from the system")
	token   = flag.String("token", os.Getenv("PORTWORX_TOKEN"), "Authorization token if any")
	address = flag.String("address", os.Getenv("PORTWORX_GRPC_URL"), "Address to server as <address>:<port>")
)

type OpenStorageSdkToken struct{}
