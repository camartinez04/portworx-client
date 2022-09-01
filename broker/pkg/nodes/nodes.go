package nodes

import (
	"context"
	"fmt"

	"github.com/camartinez04/portworx-client/broker/pkg/config"
	"github.com/camartinez04/portworx-client/broker/pkg/helpers"
	"github.com/camartinez04/portworx-client/broker/pkg/volumes"
	api "github.com/libopenstorage/openstorage-sdk-clients/sdk/golang"
	"google.golang.org/grpc"
)

// FindVolumeNodes returns a list of nodes that have the given volume
func FindVolumeNodes(conn *grpc.ClientConn, volumeName string) (volumeNodes []string, errorFound error) {

	// Retrieves the volume ID.
	volId, errorFound := volumes.GetVolumeID(conn, volumeName)
	if errorFound != nil {
		fmt.Println(errorFound)
		return nil, errorFound
	}

	// Opens the volume client connection.
	volumes := api.NewOpenStorageVolumeClient(conn)

	// Retrieves the volume information.
	volumeInspect, errorFound := volumes.Inspect(
		context.Background(),
		&api.SdkVolumeInspectRequest{
			VolumeId: volId,
		},
	)
	if errorFound != nil {
		fmt.Println(errorFound)
		return nil, errorFound
	}

	var nodesList []string

	var nodeID string

	// Retrieves the replica sets of a volume
	for _, replicaSet := range volumeInspect.Volume.GetReplicaSets() {

		// Retrieves the node ID of a replica set
		for _, replica := range replicaSet.Nodes {

			nodesList = append(nodesList, replica)

			nodeclient := api.NewOpenStorageNodeClient(conn)

			// For each node ID, get its information
			for _, nodeID = range nodesList {

				// Retrieves the node information.
				nodeIdResponse, errorFound := nodeclient.Inspect(
					context.Background(),
					&api.SdkNodeInspectRequest{
						NodeId: nodeID,
					},
				)
				if errorFound != nil {
					fmt.Println(errorFound)
					return nil, errorFound
				}

				// Retrieves the node name.
				volumeNodes = append(volumeNodes, nodeIdResponse.GetNode().GetSchedulerNodeName())

			}
		}
	}

	// Remove duplicate node entries for HA 2 or HA 1
	volumeNodes = helpers.RemoveDuplicateStr(volumeNodes)

	return volumeNodes, nil

}

// getListOfNodes retrieves a list of nodes from the cluster
func GetListOfNodes(conn *grpc.ClientConn) (nodeListReturn map[string][]string, errorFound error) {

	nodeList := make(map[string][]string)

	// First, get all node node IDs in this cluster
	nodeclient := api.NewOpenStorageNodeClient(conn)
	nodeEnumResp, errorFound := nodeclient.Enumerate(
		context.Background(),
		&api.SdkNodeEnumerateRequest{})
	if errorFound != nil {
		fmt.Println(errorFound)
		return nil, errorFound
	}

	// For each node ID, get its information
	for _, nodeID := range nodeEnumResp.GetNodeIds() {
		node, errorFound := nodeclient.Inspect(
			context.Background(),
			&api.SdkNodeInspectRequest{
				NodeId: nodeID,
			},
		)
		if errorFound != nil {
			fmt.Println(errorFound)
			return nil, errorFound
		}

		nodeUsage, errorFound := GetNodeUsage(conn, nodeID)
		if errorFound != nil {
			fmt.Println(errorFound)
			return nil, errorFound
		}

		mySlice := []string{node.GetNode().GetSchedulerNodeName(), fmt.Sprintf("%d", nodeUsage)}

		nodeList[nodeID] = mySlice
	}

	nodeListReturn = nodeList

	return nodeListReturn, nil

}

func GetNodeUsage(conn *grpc.ClientConn, nodeID string) (nodeUsage int, errorFound error) {

	nodeclient := api.NewOpenStorageNodeClient(conn)
	nodeUsageResp, errorFound := nodeclient.VolumeUsageByNode(

		context.Background(),
		&api.SdkNodeVolumeUsageByNodeRequest{
			NodeId: nodeID,
		},
	)

	if errorFound != nil {
		fmt.Println(errorFound)
		return 0, errorFound
	}

	return nodeUsageResp.GetVolumeUsageInfo().XXX_Size(), nil

}

// GetReplicasPerNode returns a list of volumes that are on the given node
func GetReplicasPerNode(conn *grpc.ClientConn, nodeID string) (volumeList map[string]config.VolumeInfo, errorFound error) {

	volumeList = make(map[string]config.VolumeInfo)

	volumeclient := api.NewOpenStorageVolumeClient(conn)

	volumeEnumResp, errorFound := volumeclient.Enumerate(
		context.Background(),
		&api.SdkVolumeEnumerateRequest{})
	if errorFound != nil {
		fmt.Println("error enumerating volumes")
		fmt.Println(errorFound)
		return nil, errorFound
	}

	// For each volume ID, get its information
	for _, volumeReplica := range volumeEnumResp.VolumeIds {

		volumeInspect, errorFound := volumeclient.Inspect(
			context.Background(),
			&api.SdkVolumeInspectRequest{
				VolumeId: volumeReplica,
			},
		)
		if errorFound != nil {
			fmt.Println("error inspecting volume")
			fmt.Println(errorFound)
			return nil, errorFound
		}

		// Retrieves the replica sets of a volume
		for _, volumeReplicaNode := range volumeInspect.Volume.GetReplicaSets() {

			// Retrieves the node ID of a replica set
			for _, volumeReplicaNodeID := range volumeReplicaNode.Nodes {

				// If the node ID matches the node ID passed in, add the volume to the list
				if volumeReplicaNodeID == nodeID {

					volumeInfo, errorFound := volumes.GetVolumeInfo(conn, volumeReplica)
					if errorFound != nil {
						fmt.Println("error getting volume info")
						fmt.Println(errorFound)
						return nil, errorFound
					}

					volumeList[volumeReplica] = volumeInfo

				}

			}

		}

	}

	return volumeList, nil
}

// FormVolumeNodes returns a list of nodes that have a volume replica
func FormVolumeNodes(conn *grpc.ClientConn) {

	nodeList, err := GetListOfNodes(conn)
	if err != nil {
		fmt.Println(err)
	}

	volumeIDList, err := volumes.GetAllVolumes(conn)
	if err != nil {
		fmt.Println(err)
	}

	replicasMap := make(map[string][]string)

	volumeclient := api.NewOpenStorageVolumeClient(conn)

	for _, volumeID := range volumeIDList {

		volumeEnumResp, errorFound := volumeclient.Inspect(
			context.Background(),
			&api.SdkVolumeInspectRequest{
				VolumeId: volumeID,
			},
		)
		if errorFound != nil {
			fmt.Println(errorFound)
			return
		}

		for _, volumeReplica := range volumeEnumResp.Volume.GetReplicaSets() {

			for _, volumeReplicaNode := range volumeReplica.Nodes {

				replicasMap[volumeReplicaNode] = append(replicasMap[volumeReplicaNode], volumeID)

				fmt.Println(replicasMap)

				fmt.Println(nodeList)
			}

		}

	}

}

// GetNodeInfo returns a Node's information
func GetNodeInfo(conn *grpc.ClientConn, nodeID string) (nodeInfo config.NodeInfo, errorFound error) {

	var sizeNodePool uint64
	var usedNodePool uint64
	var percentUsedPool float64
	var percentAvailablePool float64
	var percentUsedMemory float64
	var storagelessNode bool

	nodeclient := api.NewOpenStorageNodeClient(conn)

	apiNodeInfo, errorFound := nodeclient.Inspect(
		context.Background(),
		&api.SdkNodeInspectRequest{
			NodeId: nodeID,
		},
	)
	if errorFound != nil {
		fmt.Println(errorFound)
		return nodeInfo, errorFound
	}

	nodeStatus := apiNodeInfo.Node.GetStatus().String()
	nodeName := apiNodeInfo.Node.GetSchedulerNodeName()
	nodeAvgLoad := apiNodeInfo.Node.GetAvgLoad()
	nodePools := apiNodeInfo.Node.GetPools()
	nodeMemTotal := apiNodeInfo.Node.GetMemTotal() / 1024 / 1024 / 1024
	nodeMemUsed := apiNodeInfo.Node.GetMemUsed() / 1024 / 1024 / 1024
	nodeMemFree := apiNodeInfo.Node.GetMemFree() / 1024 / 1024 / 1024

	numberOfPools := len(nodePools)

	// for loop over the pools to get the total size and used size of the pools
	for _, pool := range nodePools {

		sizeNodePool = sizeNodePool + pool.GetTotalSize()
		usedNodePool = usedNodePool + pool.GetUsed()

	}

	sizeNodePool = sizeNodePool / 1024 / 1024
	usedNodePool = usedNodePool / 1024 / 1024
	freeNodePool := sizeNodePool - usedNodePool

	// prevent storageless issue when calculating percent used pool
	if numberOfPools == 0 {
		percentUsedPool = 0
		storagelessNode = true
		percentAvailablePool = 0
	} else {
		percentUsedPool = helpers.RoundFloat(((float64(usedNodePool) / float64(sizeNodePool)) * 100), 2)
		storagelessNode = false
		percentAvailablePool = 100 - percentUsedPool

	}

	percentUsedMemory = helpers.RoundFloat(((float64(nodeMemUsed) / float64(nodeMemTotal)) * 100), 2)

	nodeInfo = config.NodeInfo{
		NodeName:             nodeName,
		NodeID:               nodeID,
		NodeStatus:           nodeStatus,
		NodeAvgLoad:          nodeAvgLoad,
		NumberOfPools:        numberOfPools,
		NodeMemTotal:         nodeMemTotal,
		NodeMemUsed:          nodeMemUsed,
		NodeMemFree:          nodeMemFree,
		PercentUsedMemory:    percentUsedMemory,
		PercentUsedPool:      percentUsedPool,
		PercentAvailablePool: percentAvailablePool,
		SizeNodePool:         sizeNodePool,
		UsedNodePool:         usedNodePool,
		FreeNodePool:         freeNodePool,
		StoragelessNode:      storagelessNode,
		StoragePools:         nodePools,
	}

	return nodeInfo, nil

}

// GetAllNodesInfo returns a list with relevant Node's information
func GetAllNodesInfo(conn *grpc.ClientConn) (AllNodesInfo []config.NodeInfo, errorFound error) {

	nodeclient := api.NewOpenStorageNodeClient(conn)

	// First, get all node node IDs in this cluster
	nodeEnumResp, errorFound := nodeclient.Enumerate(
		context.Background(),
		&api.SdkNodeEnumerateRequest{})
	if errorFound != nil {
		fmt.Println(errorFound)
		return nil, errorFound
	}

	// For each node ID, get its information
	for _, nodeID := range nodeEnumResp.GetNodeIds() {

		nodeInfo, errorFound := GetNodeInfo(conn, nodeID)
		if errorFound != nil {
			fmt.Println(errorFound)
			return nil, errorFound
		}

		AllNodesInfo = append(AllNodesInfo, nodeInfo)

	}

	return AllNodesInfo, nil

}
