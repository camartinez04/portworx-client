package cluster

import (
	"context"
	"os"

	"log"

	"github.com/camartinez04/portworx-client/broker/pkg/helpers"
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
func ClusterCapacity(conn *grpc.ClientConn) (mbCapacity uint64, mbUsed uint64, mbAvailable uint64, percentUsed float64, percentAvailable float64, errorFound error) {

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

	mbCapacity = totalCapacity / 1024 / 1024
	mbUsed = totalUsed / 1024 / 1024
	mbAvailable = mbCapacity - mbUsed
	percentUsed = helpers.RoundFloat(((float64(mbUsed) / float64(mbCapacity)) * 100), 2)
	percentAvailable = 100 - percentUsed

	return mbCapacity, mbUsed, mbAvailable, percentUsed, percentAvailable, nil

}

// ClusterAlarms prints the Portworx cluster alarms
func ClusterAlarms(conn *grpc.ClientConn) (alarms []*api.Alert, errorFound error) {

	// Create a cluster client
	alertsClient := api.NewOpenStorageAlertsClient(conn)

	alertsToClient, erroFound := alertsClient.EnumerateWithFilters(context.Background(), &api.SdkAlertsEnumerateWithFiltersRequest{
		Queries: []*api.SdkAlertsQuery{{}},
	})
	if erroFound != nil {
		return nil, erroFound
	}

	// Get the context of the gRPC connection
	alertsToClient.Context()
	if erroFound != nil {
		log.Printf("Error value: %v", erroFound)
		return nil, erroFound
	}

	// Declare the context as done
	alertsToClient.Context().Done()
	if erroFound != nil {
		log.Printf("Error value: %v", erroFound)
		return nil, erroFound
	}

	// Close the connection
	alertsToClient.CloseSend()
	if erroFound != nil {
		log.Printf("Error value: %v", erroFound)
		return nil, erroFound
	}

	// Receive the Portworx alerts obtained from the gRPC connection
	alertList, erroFound := alertsToClient.Recv()
	if erroFound != nil {
		log.Printf("Error found at this moment of Recv(): %v with error %v", alertList, erroFound)
		return nil, erroFound
	}

	alarms = alertList.GetAlerts()
	// Get all the alerts
	return alarms, nil

}
