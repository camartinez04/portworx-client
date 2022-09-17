package cluster

import (
	"context"
	"os"
	"time"

	"log"

	"github.com/camartinez04/portworx-client/broker/pkg/helpers"
	api "github.com/libopenstorage/openstorage-sdk-clients/sdk/golang"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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

// ClusterAlarms gRPC client to get the Portworx cluster alarms
func ClusterAlarms(conn *grpc.ClientConn) (alarms []*api.Alert, errorFound error) {

	// Create a cluster client
	cluster := api.NewOpenStorageAlertsClient(conn)

	//var clusterServer api.OpenStorageAlertsServer

	//var clusterFilters api.OpenStorageAlerts_EnumerateWithFiltersServer

	var queries api.SdkAlertsEnumerateWithFiltersRequest

	queries.Queries = []*api.SdkAlertsQuery{
		{
			Query: &api.SdkAlertsQuery_AlertTypeQuery{

				AlertTypeQuery: &api.SdkAlertsAlertTypeQuery{
					ResourceType: api.ResourceType_RESOURCE_TYPE_NODE,
					AlertType:    2,
				},
			},
			Opts: []*api.SdkAlertsOption{
				{
					Opt: &api.SdkAlertsOption_MinSeverityType{
						MinSeverityType: api.SeverityType_SEVERITY_TYPE_ALARM,
					},
				},
				{
					Opt: &api.SdkAlertsOption_TimeSpan{
						TimeSpan: &api.SdkAlertsTimeSpan{
							StartTime: timestamppb.New(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
							EndTime:   timestamppb.Now(),
						},
					},
				},
			},
		},
	}

	// Get the cluster alarms
	clusterAlarms, erroFound := cluster.EnumerateWithFilters(
		context.Background(),
		&api.SdkAlertsEnumerateWithFiltersRequest{
			Queries: []*api.SdkAlertsQuery{
				{
					Query: &api.SdkAlertsQuery_AlertTypeQuery{

						AlertTypeQuery: &api.SdkAlertsAlertTypeQuery{
							ResourceType: api.ResourceType_RESOURCE_TYPE_NODE,
							AlertType:    2,
						},
					},
					Opts: []*api.SdkAlertsOption{
						{
							Opt: &api.SdkAlertsOption_MinSeverityType{
								MinSeverityType: api.SeverityType_SEVERITY_TYPE_ALARM,
							},
						},
						{
							Opt: &api.SdkAlertsOption_TimeSpan{
								TimeSpan: &api.SdkAlertsTimeSpan{
									StartTime: timestamppb.New(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
									EndTime:   timestamppb.Now(),
								},
							},
						},
					},
				},
				{
					Query: &api.SdkAlertsQuery_AlertTypeQuery{

						AlertTypeQuery: &api.SdkAlertsAlertTypeQuery{
							ResourceType: api.ResourceType_RESOURCE_TYPE_CLUSTER,
							AlertType:    2,
						},
					},
					Opts: []*api.SdkAlertsOption{
						{
							Opt: &api.SdkAlertsOption_MinSeverityType{
								MinSeverityType: api.SeverityType_SEVERITY_TYPE_ALARM,
							},
						},
						{
							Opt: &api.SdkAlertsOption_TimeSpan{
								TimeSpan: &api.SdkAlertsTimeSpan{
									StartTime: timestamppb.New(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
									EndTime:   timestamppb.Now(),
								},
							},
						},
					},
				},
				{
					Query: &api.SdkAlertsQuery_AlertTypeQuery{

						AlertTypeQuery: &api.SdkAlertsAlertTypeQuery{
							ResourceType: api.ResourceType_RESOURCE_TYPE_DRIVE,
							AlertType:    2,
						},
					},
					Opts: []*api.SdkAlertsOption{
						{
							Opt: &api.SdkAlertsOption_MinSeverityType{
								MinSeverityType: api.SeverityType_SEVERITY_TYPE_ALARM,
							},
						},
						{
							Opt: &api.SdkAlertsOption_TimeSpan{
								TimeSpan: &api.SdkAlertsTimeSpan{
									StartTime: timestamppb.New(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
									EndTime:   timestamppb.Now(),
								},
							},
						},
					},
				},
				{
					Query: &api.SdkAlertsQuery_AlertTypeQuery{

						AlertTypeQuery: &api.SdkAlertsAlertTypeQuery{
							ResourceType: api.ResourceType_RESOURCE_TYPE_VOLUME,
							AlertType:    2,
						},
					},
					Opts: []*api.SdkAlertsOption{
						{
							Opt: &api.SdkAlertsOption_MinSeverityType{
								MinSeverityType: api.SeverityType_SEVERITY_TYPE_ALARM,
							},
						},
						{
							Opt: &api.SdkAlertsOption_TimeSpan{
								TimeSpan: &api.SdkAlertsTimeSpan{
									StartTime: timestamppb.New(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
									EndTime:   timestamppb.Now(),
								},
							},
						},
					},
				},
			},
		},
	)
	if erroFound != nil {
		log.Printf("Error %v", erroFound)
		return nil, erroFound

	}

	clusterAlarms.Context()

	alarmsGotten, erroFound := clusterAlarms.Recv()

	//clusterFilters.Send(alarmsGotten)

	//errorFound = clusterServer.EnumerateWithFilters(&queries, clusterFilters)
	//if errorFound != nil {
	//	log.Printf("Error %v", errorFound)
	//	return nil, errorFound
	//}

	log.Printf("Alerts: %v", alarmsGotten)
	if erroFound != nil {
		log.Printf("Error trying to get response %v", erroFound)
		return nil, erroFound
	}

	alerts := alarmsGotten.GetAlerts()

	//log.Printf("Cluster alarms: %v", alerts)

	clusterAlarms.CloseSend()

	return alerts, nil

}
