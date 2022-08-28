package nodes

import (
	"context"
	"fmt"

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
	volumeNodes = removeDuplicateStr(volumeNodes)

	return volumeNodes, nil

}

// removeDuplicateStr removes duplicate strings from a slice of strings
func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// getListOfNodes retrieves a list of nodes from the cluster
func GetListOfNodes(conn *grpc.ClientConn) (nodeListReturn map[string]string, errorFound error) {

	nodeList := make(map[string]string)

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

		nodeList[nodeID] = node.Node.GetSchedulerNodeName()
	}

	nodeListReturn = nodeList

	return nodeListReturn, nil

}

// GetReplicasPerNode returns a list of volume replicas located on a node
func GetReplicasPerNode(conn *grpc.ClientConn, nodeName string) (volumeList []any, errorFound error) {

	nodeclient := api.NewOpenStorageNodeClient(conn)

	nodeEnumResp, errorFound := nodeclient.Enumerate(
		context.Background(),
		&api.SdkNodeEnumerateRequest{})
	if errorFound != nil {
		fmt.Println("error ennumerating nodes")
		fmt.Println(errorFound)
		return nil, errorFound
	}

	for _, nodeID := range nodeEnumResp.GetNodeIds() {
		node, errorFound := nodeclient.Inspect(
			context.Background(),
			&api.SdkNodeInspectRequest{
				NodeId: nodeID,
			},
		)
		if errorFound != nil {
			fmt.Println("error inspecting node")
			fmt.Println(errorFound)
			return nil, errorFound
		}

		if node.Node.GetSchedulerNodeName() == nodeName {

			nodeIDReturn := node.Node.Id

			volumeclient := api.NewOpenStorageVolumeClient(conn)

			volumeEnumResp, errorFound := volumeclient.Enumerate(
				context.Background(),
				&api.SdkVolumeEnumerateRequest{})
			if errorFound != nil {
				fmt.Println("error enumerating volumes")
				fmt.Println(errorFound)
				return nil, errorFound
			}

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

				for _, volumeReplicaNode := range volumeInspect.Volume.GetReplicaSets() {

					for _, volumeReplicaNodeID := range volumeReplicaNode.Nodes {

						if volumeReplicaNodeID == nodeIDReturn {

							volumeList = append(volumeList, volumeReplica)

						}

					}

				}

			}
		}
	}

	return volumeList, nil
}

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
