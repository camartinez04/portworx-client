package cluster

import (
	"context"
	"os"

	"log"

	api "github.com/libopenstorage/openstorage-sdk-clients/sdk/golang"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	Bytes = uint64(1)
	KB    = Bytes * uint64(1024)
	MB    = KB * uint64(1024)
	GB    = MB * uint64(1024)
)

// clusterInfo prints the Portworx cluster information
func ClusterInfo(conn *grpc.ClientConn) (clusterUUID string, clusterStatus string, clusterName string, erroFound error) {

	// Create a cluster client
	cluster := api.NewOpenStorageClusterClient(conn)

	// Print the cluster information
	clusterInfo, erroFound := cluster.InspectCurrent(
		context.Background(),
		&api.SdkClusterInspectCurrentRequest{})
	if erroFound != nil {
		return "", "", "", erroFound
	}

	clusterUUID = clusterInfo.GetCluster().GetId()

	clusterStatus = clusterInfo.GetCluster().GetStatus().String()

	clusterName = clusterInfo.GetCluster().GetName()

	return clusterUUID, clusterStatus, clusterName, nil
}

// clusterCapacity prints the Portworx cluster total capacity
func ClusterCapacity(conn *grpc.ClientConn) (uint64, uint64, error) {

	// --- Get Cluster capacity ---
	// First, get all node node IDs in this cluster
	nodeclient := api.NewOpenStorageNodeClient(conn)
	nodeEnumResp, err := nodeclient.Enumerate(
		context.Background(),
		&api.SdkNodeEnumerateRequest{})
	if err != nil {
		gerr, _ := status.FromError(err)
		log.Printf("Error Code[%d] Message[%s]\n",
			gerr.Code(), gerr.Message())
		os.Exit(1)
	}

	// Initialize the variables
	totalCapacity := uint64(0)
	totalUsed := uint64(0)

	// For each node ID, get its information
	for _, nodeID := range nodeEnumResp.GetNodeIds() {
		node, err := nodeclient.Inspect(
			context.Background(),
			&api.SdkNodeInspectRequest{
				NodeId: nodeID,
			},
		)
		if err != nil {
			gerr, _ := status.FromError(err)
			log.Printf("Error Code[%d] Message[%s]\n",
				gerr.Code(), gerr.Message())
			os.Exit(1)
		}

		// Get size from the pools
		// Use Pool instead of the disks, because disks could be in a RAID
		// configuration. The Pool returns the usable size.
		for _, pool := range node.GetNode().GetPools() {
			totalCapacity += pool.GetTotalSize()
			totalUsed += pool.GetUsed()
		}
	}

	gbCapacity := totalCapacity / 1024 / 1024 / 1024
	gbUsed := totalUsed / 1024 / 1024 / 1024

	return gbCapacity, gbUsed, nil

}
