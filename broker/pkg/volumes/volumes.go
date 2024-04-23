package volumes

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/camartinez04/portworx-client/broker/pkg/config"
	"github.com/camartinez04/portworx-client/broker/pkg/helpers"
	api "github.com/libopenstorage/openstorage-sdk-clients/sdk/golang"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// GetVolumeID Retrieve the Portworx volume ID having its Name and a gRPC connection to the Portworx API.
func GetVolumeID(conn *grpc.ClientConn, volumeName string) (volumeID string, errorFound error) {

	// Opens the volume client connection.
	volumes := api.NewOpenStorageVolumeClient(conn)

	// Returns a list of volume IDs that matches the filter.
	volume, err := volumes.EnumerateWithFilters(
		context.Background(),
		&api.SdkVolumeEnumerateWithFiltersRequest{
			Name: volumeName,
		})
	if err != nil {
		errorFound = fmt.Errorf("volume id not found of %s", volumeName)
		log.Println(errorFound)
		return "", errorFound
	}

	// Handles possible cases with the slice of volumes.
	switch {
	case len(volume.VolumeIds) == 1:
		volumeID = volume.VolumeIds[0]
	case len(volume.VolumeIds) == 0:
		errorFound = fmt.Errorf("no volume found given \"%s\" as volume name", volumeName)
		log.Println(errorFound)
		return "", errorFound
	case len(volume.VolumeIds) > 1:
		errorFound = fmt.Errorf("more than one volume found given \"%s\" as volume name", volumeName)
		log.Println(errorFound)
		return "", errorFound
	}

	return volumeID, nil

}

// CreateVolume Creates a new Portworx volume, if Sharedv4 enabled, sets to service ClusterIP by default.
func CreateVolume(conn *grpc.ClientConn, volumeName string, volumeGBSize uint64, volumeIOProfile string, volumeHALevel int64, encryptionEnabled bool, sharedv4Enabled bool, noDiscard bool) (string, error) {

	// Opens the volume client connection.
	volumes := api.NewOpenStorageVolumeClient(conn)

	// default volume IO profile to auto
	intIOProfile := api.IoProfile_IO_PROFILE_AUTO

	// Verifies if the volume already exists. If it does, returns an error.
	_, err := GetVolumeID(conn, volumeName)
	if err == nil {
		newError := fmt.Sprintf("a volume called \"%s\" already exists! volume will not be created", volumeName)
		log.Println(newError)
		return "", errors.New(newError)

	}

	if volumeIOProfile == "db_remote" {
		intIOProfile = api.IoProfile_IO_PROFILE_DB_REMOTE
	}

	if volumeIOProfile == "db" {
		intIOProfile = api.IoProfile_IO_PROFILE_DB
	}

	if volumeIOProfile == "sequential" {
		intIOProfile = api.IoProfile_IO_PROFILE_SEQUENTIAL
	}

	if volumeIOProfile == "sync_shared" {
		intIOProfile = api.IoProfile_IO_PROFILE_SYNC_SHARED
	}

	if sharedv4Enabled {
		// Creates the volume  witho sharedv4 service enabled.
		volume, err := volumes.Create(
			context.Background(),
			&api.SdkVolumeCreateRequest{
				Name: volumeName,
				Spec: &api.VolumeSpec{
					Size:      volumeGBSize * 1024 * 1024 * 1024,
					HaLevel:   volumeHALevel,
					IoProfile: intIOProfile,
					Cos:       api.CosType_HIGH,
					Format:    api.FSType_FS_TYPE_EXT4,
					Encrypted: encryptionEnabled,
					Sharedv4:  sharedv4Enabled,
					Nodiscard: noDiscard,
					Sharedv4ServiceSpec: &api.Sharedv4ServiceSpec{
						Type: api.Sharedv4ServiceType_SHAREDV4_SERVICE_TYPE_CLUSTERIP,
					},
					IoStrategy: &api.IoStrategy{
						AsyncIo:  true,
						EarlyAck: true,
					},
				},
			})
		if err != nil {
			gerr, _ := status.FromError(err)
			newError := fmt.Sprintf("error code[%d] message[%s]", gerr.Code(), gerr.Message())
			log.Println(newError)
			return "", errors.New(newError)
		}

		log.Printf("Volume %s of %dGi created with id %s", volumeName, volumeGBSize, volume.GetVolumeId())

		newVolumeID := volume.GetVolumeId()

		return newVolumeID, nil

	} else {

		// Creates the volume without sharedv4 service.
		volume, err := volumes.Create(
			context.Background(),
			&api.SdkVolumeCreateRequest{
				Name: volumeName,
				Spec: &api.VolumeSpec{
					Size:      volumeGBSize * 1024 * 1024 * 1024,
					HaLevel:   volumeHALevel,
					IoProfile: intIOProfile,
					Cos:       api.CosType_HIGH,
					Format:    api.FSType_FS_TYPE_EXT4,
					Encrypted: encryptionEnabled,
					Sharedv4:  sharedv4Enabled,
					Nodiscard: noDiscard,
					IoStrategy: &api.IoStrategy{
						AsyncIo:  true,
						EarlyAck: true,
					},
				},
			})
		if err != nil {
			gerr, _ := status.FromError(err)
			newError := fmt.Sprintf("error code[%d] message[%s]", gerr.Code(), gerr.Message())
			log.Println(newError)
			return "", errors.New(newError)
		}

		log.Printf("Volume %s of %dGi created with id %s", volumeName, volumeGBSize, volume.GetVolumeId())

		newVolumeID := volume.GetVolumeId()

		return newVolumeID, nil
	}

}

// UpdateVolume updates a Portworx volume size.
func UpdateVolumeSize(conn *grpc.ClientConn, volumeID string, volumeGBSize uint64) (volumeUpdate string, errorFound error) {

	// Opens the volume client connection.
	volumes := api.NewOpenStorageVolumeClient(conn)

	// Updates the volume.
	volume, errorFound := volumes.Update(
		context.Background(),
		&api.SdkVolumeUpdateRequest{
			VolumeId: volumeID,
			Spec: &api.VolumeSpecUpdate{
				SizeOpt: &api.VolumeSpecUpdate_Size{
					Size: volumeGBSize * 1024 * 1024 * 1024,
				},
			},
		})
	if errorFound != nil {
		log.Printf("error updating volume %s", volumeID)
		return "", errorFound
	}

	volumeUpdate = volume.String()

	log.Printf("Volume %s updated to size %d", volumeUpdate, volumeGBSize)

	return volumeUpdate, nil

}

// UpdateVolume updates a Portworx IO profile.
func UpdateVolumeIOProfile(conn *grpc.ClientConn, volumeID string, volumeIOProfile string) (volumeUpdate string, errorFound error) {

	// default volume IO profile to auto
	intIOProfile := api.IoProfile_IO_PROFILE_AUTO

	if volumeIOProfile == "db_remote" {
		intIOProfile = api.IoProfile_IO_PROFILE_DB_REMOTE
	}

	if volumeIOProfile == "db" {
		intIOProfile = api.IoProfile_IO_PROFILE_DB
	}

	if volumeIOProfile == "sequential" {
		intIOProfile = api.IoProfile_IO_PROFILE_SEQUENTIAL
	}

	if volumeIOProfile == "sync_shared" {
		intIOProfile = api.IoProfile_IO_PROFILE_SYNC_SHARED
	}

	// Opens the volume client connection.
	volumes := api.NewOpenStorageVolumeClient(conn)

	// Updates the volume.
	volume, errorFound := volumes.Update(
		context.Background(),
		&api.SdkVolumeUpdateRequest{
			VolumeId: volumeID,
			Spec: &api.VolumeSpecUpdate{
				IoProfileOpt: &api.VolumeSpecUpdate_IoProfile{
					IoProfile: intIOProfile,
				},
			},
		},
	)
	if errorFound != nil {
		log.Printf("error updating volume %s", volumeID)
		return "", errorFound
	}

	volumeUpdate = volume.String()

	log.Printf("Volume %s updated to io-profile %s", volumeUpdate, volumeIOProfile)

	return volumeUpdate, nil

}

// UpdateVolume updates a Portworx volume HA level.
func UpdateVolumeHALevel(conn *grpc.ClientConn, volumeID string, volumeHALevel int64) (volumeUpdate string, errorFound error) {

	// Opens the volume client connection.
	volumes := api.NewOpenStorageVolumeClient(conn)

	// Updates the volume.
	volume, errorFound := volumes.Update(
		context.Background(),
		&api.SdkVolumeUpdateRequest{
			VolumeId: volumeID,
			Spec: &api.VolumeSpecUpdate{
				HaLevelOpt: &api.VolumeSpecUpdate_HaLevel{
					HaLevel: volumeHALevel,
				},
			},
		})
	if errorFound != nil {
		log.Printf("error updating volume %s", volumeID)
		return "", errorFound
	}

	volumeUpdate = volume.String()

	log.Printf("Volume %s updated to %d replicas", volumeUpdate, volumeHALevel)

	return volumeUpdate, nil

}

// testing UpdateVolumeReplicaSet function
func UpdateVolumeReplicaSet(conn *grpc.ClientConn, volumeID string, poolUuids []string) (volumeUpdate string, errorFound error) {

	// Opens the volume client connection.
	volumes := api.NewOpenStorageVolumeClient(conn)

	// Updates the volume.
	volume, errorFound := volumes.Update(
		context.Background(),
		&api.SdkVolumeUpdateRequest{
			VolumeId: volumeID,
			Spec: &api.VolumeSpecUpdate{
				ReplicaSet: &api.ReplicaSet{
					PoolUuids: poolUuids,
				},
			},
		})
	if errorFound != nil {
		log.Printf("error updating volume %s", volumeID)
		return "", errorFound
	}

	volumeUpdate = volume.String()

	log.Printf("Volume %s updated to Replica Set on Pool UUIDs: %v", volumeID, poolUuids)

	return volumeUpdate, nil

}

// UpdateVolume updates if a Portworx volume would be Sharedv4 or not.
func UpdateVolumeSharedv4(conn *grpc.ClientConn, volumeID string, sharedv4Enabled bool) (volumeUpdate string, errorFound error) {

	// Opens the volume client connection.
	volumes := api.NewOpenStorageVolumeClient(conn)

	// Updates the volume.
	volume, errorFound := volumes.Update(
		context.Background(),
		&api.SdkVolumeUpdateRequest{
			VolumeId: volumeID,
			Spec: &api.VolumeSpecUpdate{
				Sharedv4Opt: &api.VolumeSpecUpdate_Sharedv4{
					Sharedv4: sharedv4Enabled,
				},
			},
		})
	if errorFound != nil {
		log.Printf("error updating volume %s", volumeID)
		return "", errorFound
	}

	volumeUpdate = volume.String()

	log.Printf("Volume %s updated to sharedv4 %t", volumeUpdate, sharedv4Enabled)

	return volumeUpdate, nil

}

// UpdateVolume updates Sharedv4 volume Service.
func UpdateVolumeSharedv4Service(conn *grpc.ClientConn, volumeID string, sharedv4Service bool) (volumeUpdate string, errorFound error) {

	// Opens the volume client connection.
	volumes := api.NewOpenStorageVolumeClient(conn)

	if sharedv4Service {

		isVolumeSharedv4, errorFound := volumeSharedv4(conn, volumeID)
		if errorFound != nil {
			log.Printf("error checking if volume is sharedv4 %s", volumeID)
			return "", errorFound
		}

		if !isVolumeSharedv4 {
			log.Printf("Volume %s is not Sharedv4", volumeID)
			return "", errors.New("volume is not sharedv4")
		}

		// Updates the volume.
		volume, errorFound := volumes.Update(
			context.Background(),
			&api.SdkVolumeUpdateRequest{
				VolumeId: volumeID,
				Spec: &api.VolumeSpecUpdate{
					Sharedv4ServiceSpecOpt: &api.VolumeSpecUpdate_Sharedv4ServiceSpec{
						Sharedv4ServiceSpec: &api.Sharedv4ServiceSpec{
							Type: api.Sharedv4ServiceType_SHAREDV4_SERVICE_TYPE_CLUSTERIP,
						},
					},
				},
			})
		if errorFound != nil {
			log.Printf("error updating volume %s", volumeID)
			return "", errorFound
		}

		volumeUpdate = volume.String()

		log.Printf("Volume %s updated its sharedv4 service to %t", volumeUpdate, sharedv4Service)

		return volumeUpdate, nil

	}

	// Updates the volume.
	isVolumeSharedv4, errorFound := volumeSharedv4(conn, volumeID)
	if errorFound != nil {
		log.Printf("error checking if volume is sharedv4 %s", volumeID)
		return "", errorFound
	}

	if !isVolumeSharedv4 {
		log.Printf("Volume %s is not Sharedv4", volumeID)
		return "", errors.New("volume is not sharedv4")
	}
	// Updates the volume.
	volume, errorFound := volumes.Update(
		context.Background(),
		&api.SdkVolumeUpdateRequest{
			VolumeId: volumeID,
			Spec: &api.VolumeSpecUpdate{
				Sharedv4ServiceSpecOpt: &api.VolumeSpecUpdate_Sharedv4ServiceSpec{
					Sharedv4ServiceSpec: &api.Sharedv4ServiceSpec{
						Type: api.Sharedv4ServiceType_SHAREDV4_SERVICE_TYPE_INVALID,
					},
				},
			},
		})
	if errorFound != nil {
		log.Printf("error updating volume %s", volumeID)
		return "", errorFound
	}

	volumeUpdate = volume.String()

	log.Printf("Volume %s updated its sharedv4 service to %t", volumeUpdate, sharedv4Service)

	return volumeUpdate, nil

}

// UpdateVolume updates if a Portworx volume No discard attribute.
func UpdateVolumeNoDiscard(conn *grpc.ClientConn, volumeID string, noDiscard bool) (volumeUpdate string, errorFound error) {

	// Opens the volume client connection.
	volumes := api.NewOpenStorageVolumeClient(conn)

	// Updates the volume.
	volume, errorFound := volumes.Update(
		context.Background(),
		&api.SdkVolumeUpdateRequest{
			VolumeId: volumeID,
			Spec: &api.VolumeSpecUpdate{
				NodiscardOpt: &api.VolumeSpecUpdate_Nodiscard{
					Nodiscard: noDiscard,
				},
			},
		})
	if errorFound != nil {
		log.Printf("error updating volume %s", volumeID)
		return "", errorFound
	}

	volumeUpdate = volume.String()

	log.Printf("Volume %s updated no discard to %t", volumeUpdate, noDiscard)

	return volumeUpdate, nil

}

// DeleteVolume deletes a Portworx volume.
func DeleteVolume(conn *grpc.ClientConn, volumeID string) (string, error) {

	// Opens the volume client connection.
	volumes := api.NewOpenStorageVolumeClient(conn)

	// Deletes the volume.
	_, err := volumes.Delete(
		context.Background(),
		&api.SdkVolumeDeleteRequest{
			VolumeId: volumeID,
		})
	if err != nil {
		gerr, _ := status.FromError(err)
		newError := fmt.Sprintf("error code[%d] message[%s]", gerr.Code(), gerr.Message())
		log.Println(newError)
		return "", errors.New(newError)
	}

	log.Printf("Volume %s deleted", volumeID)

	return volumeID, nil
}

// inspectVolume generates a json string with Volume information equivalent of pxctl volume inspect <volume> --json
func InspectVolume(conn *grpc.ClientConn, volumeName string) (apiVolumeInspect api.Volume, apiVolumeReplicas, volumeNodes []string, apiVolumeStatus, apiIoProfile string, errorFound error) {

	// Retrieves the volume ID.
	volId, errorFound := GetVolumeID(conn, volumeName)
	if errorFound != nil {
		log.Println(errorFound)
		return apiVolumeInspect, apiVolumeReplicas, volumeNodes, "", "", errorFound
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
		log.Println(errorFound)
		return apiVolumeInspect, apiVolumeReplicas, volumeNodes, "", "", errorFound
	}

	apiVolumeInspect = *volumeInspect.Volume

	apiVolumeReplicas = apiVolumeInspect.ReplicaSets[0].Nodes

	apiVolumeStatus = apiVolumeInspect.GetStatus().String()

	apiIoProfile = apiVolumeInspect.Spec.GetIoProfile().String()

	volumeNodes = make([]string, len(apiVolumeReplicas))

	// Opens the node client connection.
	nodeclient := api.NewOpenStorageNodeClient(conn)

	// For each node ID, get its information
	for _, nodeID := range apiVolumeInspect.ReplicaSets[0].Nodes {

		// Retrieves the node information.
		nodeIdResponse, errorFound := nodeclient.Inspect(
			context.Background(),
			&api.SdkNodeInspectRequest{
				NodeId: nodeID,
			},
		)
		if errorFound != nil {
			log.Println(errorFound)
			return apiVolumeInspect, apiVolumeReplicas, volumeNodes, "", "", errorFound
		}

		// Retrieves the node name.
		volumeNodes = append(volumeNodes, nodeIdResponse.GetNode().GetSchedulerNodeName())

	}

	volumeNodes = helpers.DeleteEmpty(volumeNodes)

	return apiVolumeInspect, apiVolumeReplicas, volumeNodes, apiVolumeStatus, apiIoProfile, nil

}

// RetrieveVolumeUsage retrieves the usage of a Portworx volume.
func RetrieveVolumeUsage(conn *grpc.ClientConn, volumeName string) (volumeUsage, availableSpace, totalSize uint64, errorFound error) {

	// Retrieves the volume ID.
	volId, errorFound := GetVolumeID(conn, volumeName)
	if errorFound != nil {
		log.Println(errorFound)
		return volumeUsage, availableSpace, totalSize, errorFound
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
		log.Println(errorFound)
		return volumeUsage, availableSpace, totalSize, errorFound
	}

	volumeUsage = volumeInspect.Volume.GetUsage()

	totalSize = volumeInspect.Volume.GetSpec().GetSize()

	availableSpace = totalSize - volumeUsage

	return volumeUsage, availableSpace, totalSize, nil

}

// GetAllVolumes returns a list of all volumes.
func GetAllVolumes(conn *grpc.ClientConn) (volumeList []string, errorFound error) {

	// Opens the volume client connection.
	volumes := api.NewOpenStorageVolumeClient(conn)

	// Retrieves the volume information.
	volumeEnumerate, errorFound := volumes.Enumerate(
		context.Background(),
		&api.SdkVolumeEnumerateRequest{},
	)
	if errorFound != nil {
		log.Println(errorFound)
		return volumeList, errorFound
	}

	volumeList = volumeEnumerate.VolumeIds

	return volumeList, nil

}

// GetAllVolumesComplete returns a list of all volumes with its corresponding SdkVolumeInspectResponse Struct.
func GetAllVolumesComplete(conn *grpc.ClientConn) (volumesMap map[string]*api.SdkVolumeInspectResponse, errorFound error) {

	var volumeList []string

	// Initializes the volumes map.
	volumesMap = make(map[string]*api.SdkVolumeInspectResponse)

	// Opens the volume client connection.
	volumes := api.NewOpenStorageVolumeClient(conn)

	// Retrieves the volume information.
	volumeEnumerate, errorFound := volumes.Enumerate(
		context.Background(),
		&api.SdkVolumeEnumerateRequest{},
	)
	if errorFound != nil {
		log.Println(errorFound)
		return volumesMap, errorFound
	}

	volumeList = volumeEnumerate.VolumeIds

	// For each volume ID, get its information and fill it into the Map of volumes.
	for _, volume := range volumeList {
		volInspect, errorFound := volumeInspectFromID(conn, volume)
		if errorFound != nil {
			log.Println(errorFound)
			return volumesMap, errorFound
		}

		volumesMap[volume] = volInspect
	}

	return volumesMap, nil

}

// volumeInspectFromID retrieves the volume information from the volume ID.
func volumeInspectFromID(conn *grpc.ClientConn, volumeID string) (volumeInspect *api.SdkVolumeInspectResponse, errorFound error) {

	// Opens the volume client connection.
	volumes := api.NewOpenStorageVolumeClient(conn)

	// Retrieves the volume information.
	volumeInspect, errorFound = volumes.Inspect(
		context.Background(),
		&api.SdkVolumeInspectRequest{
			VolumeId: volumeID,
		},
	)
	if errorFound != nil {
		log.Println(errorFound)
		return volumeInspect, errorFound
	}

	return volumeInspect, nil
}

// volumeSharedv4 checks if a volume is sharedv4.
func volumeSharedv4(conn *grpc.ClientConn, volID string) (isSharedv4 bool, errorFound error) {

	// Opens the volume client connection.
	volumes := api.NewOpenStorageVolumeClient(conn)

	// Retrieves the volume information.
	volumeInspect, errorFound := volumes.Inspect(
		context.Background(),
		&api.SdkVolumeInspectRequest{
			VolumeId: volID,
		},
	)
	if errorFound != nil {
		log.Println(errorFound)
		return isSharedv4, errorFound
	}

	isSharedv4 = volumeInspect.Volume.GetSpec().GetSharedv4()

	return isSharedv4, nil

}

// GetVolumeInfo returns the volume information from the volume ID.
func GetVolumeInfo(conn *grpc.ClientConn, volumeID string) (volumeInfo config.VolumeInfo, errorFound error) {

	// Opens the volume client connection.
	volumes := api.NewOpenStorageVolumeClient(conn)

	// Retrieves the volume information.
	volumeInspect, errorFound := volumes.Inspect(
		context.Background(),
		&api.SdkVolumeInspectRequest{
			VolumeId: volumeID,
		},
	)
	if errorFound != nil {
		log.Println(errorFound)
		return volumeInfo, errorFound
	}

	var IoProfile_name = map[int32]string{
		0: "sequential",
		1: "random",
		2: "db",
		3: "db_remote",
		4: "cms",
		5: "sync_shared",
		6: "auto",
	}

	var CosType_name = map[int32]string{
		0: "none",
		1: "low",
		2: "medium",
		3: "high",
	}

	volumeIOProfilePrev := volumeInspect.Volume.Spec.GetIoProfile()
	volumeIOPriorityPrev := volumeInspect.Volume.Spec.GetCos()

	volumeName := volumeInspect.Volume.Locator.GetName()
	volumeReplicas := len(volumeInspect.Volume.ReplicaSets[0].GetNodes())
	volumeReplicaNodes := volumeInspect.Volume.ReplicaSets[0].GetNodes()
	volumeIOProfile := IoProfile_name[int32(volumeIOProfilePrev)]
	volumeIOProfileAPI := volumeInspect.Volume.Spec.GetIoProfile().String()
	volumeIOPriority := CosType_name[int32(volumeIOPriorityPrev)]
	volumeStatus := volumeInspect.Volume.GetStatus().String()
	volumeAttachedOn := volumeInspect.Volume.GetAttachedOn()
	volumeAttachedPath := volumeInspect.Volume.GetAttachPath()
	volumeDevicePath := volumeInspect.Volume.GetDevicePath()
	volumeTotalSizeMB := volumeInspect.Volume.GetSpec().GetSize() / 1024 / 1024
	volumeUsageMB := volumeInspect.Volume.GetUsage() / 1024 / 1024
	volumeAvailableSpace := volumeTotalSizeMB - volumeUsageMB
	volumePercentageUsed := helpers.RoundFloat((float64(volumeUsageMB) / float64(volumeTotalSizeMB) * 100), 2)
	volumePercentageAvailable := 100 - volumePercentageUsed
	volumeType := volumeInspect.Volume.Format.String()
	volumeAttachStatus := volumeInspect.Volume.AttachedState.String()
	volumeAggregationLevel := volumeInspect.Volume.Spec.GetAggregationLevel()
	volumeConsumers := volumeInspect.Volume.GetVolumeConsumers()
	volumeEncrypted := volumeInspect.Volume.Spec.GetEncrypted()
	volumeEncryptionKey := volumeInspect.Volume.GetLocator().GetVolumeLabels()["secret_key"]
	volumeK8sNamespace := volumeInspect.Volume.GetLocator().GetVolumeLabels()["namespace"]
	volumeK8sPVCName := volumeInspect.Volume.GetLocator().GetVolumeLabels()["pvc"]
	volumeSharedv4 := volumeInspect.Volume.Spec.GetSharedv4()
	volumeSharedv4ServiceSpec := volumeInspect.Volume.Spec.GetSharedv4ServiceSpec()
	volumeIOStrategy := volumeInspect.Volume.Spec.GetIoStrategy()
	volumeNoDiscard := volumeInspect.Volume.Spec.GetNodiscard()

	volumeInfo = config.VolumeInfo{
		VolumeName:                volumeName,
		VolumeID:                  volumeID,
		VolumeK8sNamespace:        volumeK8sNamespace,
		VolumeK8sPVCName:          volumeK8sPVCName,
		VolumeReplicas:            volumeReplicas,
		VolumeReplicaNodes:        volumeReplicaNodes,
		VolumeIOProfile:           volumeIOProfile,
		VolumeIOProfileAPI:        volumeIOProfileAPI,
		VolumeIOPriority:          volumeIOPriority,
		VolumeIOStrategy:          volumeIOStrategy,
		VolumeStatus:              volumeStatus,
		VolumeAttachedOn:          volumeAttachedOn,
		VolumeAttachedPath:        volumeAttachedPath,
		VolumeDevicePath:          volumeDevicePath,
		VolumeSizeMB:              volumeTotalSizeMB,
		VolumeUsedMB:              volumeUsageMB,
		VolumeAvailable:           volumeAvailableSpace,
		VolumeUsedPercent:         volumePercentageUsed,
		VolumeAvailablePercent:    volumePercentageAvailable,
		VolumeType:                volumeType,
		VolumeAttachStatus:        volumeAttachStatus,
		VolumeAggregationLevel:    volumeAggregationLevel,
		VolumeConsumers:           volumeConsumers,
		VolumeEncrypted:           volumeEncrypted,
		VolumeEncryptionKey:       volumeEncryptionKey,
		VolumeSharedv4:            volumeSharedv4,
		VolumeSharedv4ServiceSpec: volumeSharedv4ServiceSpec,
		VolumeNoDiscard:           volumeNoDiscard,
	}

	return volumeInfo, nil
}

// GetAllVolumesInfo returns the volume information for all volumes on the cluster.
func GetAllVolumesInfo(conn *grpc.ClientConn) (AllVolumesInfo []config.VolumeInfo, errorFound error) {

	volumeclient := api.NewOpenStorageVolumeClient(conn)

	// First, get all volumes IDs in this cluster
	volsEnumResp, errorFound := volumeclient.Enumerate(
		context.Background(),
		&api.SdkVolumeEnumerateRequest{})
	if errorFound != nil {
		log.Println(errorFound)
		return AllVolumesInfo, errorFound
	}

	// Create a slice of volume IDs
	volIDsSlice := volsEnumResp.GetVolumeIds()

	// Then, sort the slice of volume IDs
	sort.Strings(volIDsSlice)

	// For each volume ID, get its information
	for _, volID := range volIDsSlice {

		volumeInfo, errorFound := GetVolumeInfo(conn, volID)
		if errorFound != nil {
			log.Println(errorFound)
			return AllVolumesInfo, errorFound
		}

		AllVolumesInfo = append(AllVolumesInfo, volumeInfo)

	}

	return AllVolumesInfo, nil

}
