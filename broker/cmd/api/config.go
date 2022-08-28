package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
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
	Error              bool     `json:"error,omitempty"`
	Message            string   `json:"message,omitempty"`
	Data               any      `json:"data,omitempty"`
	VolumeID           string   `json:"volume_id,omitempty"`
	ClusterCapacity    string   `json:"cluster_capacity,omitempty"`
	ClusterUUID        string   `json:"cluster_uuid,omitempty"`
	VolumeInspect      any      `json:"volume_inspect,omitempty"`
	NodesOfVolume      []string `json:"nodes_of_volume,omitempty"`
	NodeList           []string `json:"node_list,omitempty"`
	VolumeList         []any    `json:"volume_list,omitempty"`
	AllVolumesList     []string `json:"all_volumes_list,omitempty"`
	ReplicasInfo       []string `json:"replicas_info,omitempty"`
	VolumeNodes        []string `json:"volume_nodes,omitempty"`
	VolumeStatusString string   `json:"volume_status_string,omitempty"`
	IoProfileString    string   `json:"io_profile_string,omitempty"`
}

type JsonVolumeUsage struct {
	Error                  bool    `json:"error,omitempty"`
	VolumeUsage            float64 `json:"volume_usage"`
	AvailableSpace         float64 `json:"available_space"`
	TotalSize              float64 `json:"total_size"`
	VolumeUsagePercent     float32 `json:"volume_usage_percent"`
	VolumeAvailablePercent float32 `json:"volume_available_percent"`
}

const (
	Bytes   = uint64(1)
	KB      = Bytes * uint64(1024)
	MB      = KB * uint64(1024)
	GB      = MB * uint64(1024)
	webPort = "8080"
)

var (
	useTls  = flag.Bool("usetls", false, "Connect to server using TLS. Loads CA from the system")
	token   = flag.String("token", "", "Authorization token if any")
	address = flag.String("address", os.Getenv("PORTWORX_GRPC_URL"), "Address to server as <address>:<port>")
)

type OpenStorageSdkToken struct{}
