package volumes

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/camartinez04/portworx-client/broker/pkg/config"
	"github.com/camartinez04/portworx-client/broker/pkg/helpers"
	api "github.com/libopenstorage/openstorage-sdk-clients/sdk/golang"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// getVolumeID Retrieve the Portworx volume ID having its Name and a gRPC connection to the Portworx API.
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

// UpdateVolume updates a Portworx volume.
func UpdateVolume(conn *grpc.ClientConn, volumeID string, volumeGBSize uint64, volumeIOProfile string, volumeHALevel int64, sharedv4Enabled bool, noDiscard bool) (string, error) {

	// Opens the volume client connection.
	volumes := api.NewOpenStorageVolumeClient(conn)

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

	// Updates the volume.

	if sharedv4Enabled {
		volume, err := volumes.Update(
			context.Background(),
			&api.SdkVolumeUpdateRequest{
				VolumeId: volumeID,
				Spec: &api.VolumeSpecUpdate{
					SizeOpt: &api.VolumeSpecUpdate_Size{
						Size: volumeGBSize * 1024 * 1024 * 1024,
					},
					HaLevelOpt: &api.VolumeSpecUpdate_HaLevel{
						HaLevel: volumeHALevel,
					},
					IoProfileOpt: &api.VolumeSpecUpdate_IoProfile{
						IoProfile: intIOProfile,
					},
					Sharedv4Opt: &api.VolumeSpecUpdate_Sharedv4{
						Sharedv4: sharedv4Enabled,
					},
					NodiscardOpt: &api.VolumeSpecUpdate_Nodiscard{
						Nodiscard: noDiscard,
					},
					Sharedv4ServiceSpecOpt: &api.VolumeSpecUpdate_Sharedv4ServiceSpec{
						Sharedv4ServiceSpec: &api.Sharedv4ServiceSpec{
							Type: api.Sharedv4ServiceType_SHAREDV4_SERVICE_TYPE_CLUSTERIP,
						},
					},
				},
			})
		if err != nil {
			gerr, _ := status.FromError(err)
			newError := fmt.Sprintf("error code[%d] message[%s]", gerr.Code(), gerr.Message())
			log.Println(newError)
			return "", errors.New(newError)
		}

		volumeUpdateResponse := volume.String()

		log.Printf("Volume %s updated", volumeUpdateResponse)

	} else {
		volume, err := volumes.Update(
			context.Background(),
			&api.SdkVolumeUpdateRequest{
				VolumeId: volumeID,
				Spec: &api.VolumeSpecUpdate{
					SizeOpt: &api.VolumeSpecUpdate_Size{
						Size: volumeGBSize * 1024 * 1024 * 1024,
					},
					HaLevelOpt: &api.VolumeSpecUpdate_HaLevel{
						HaLevel: volumeHALevel,
					},
					IoProfileOpt: &api.VolumeSpecUpdate_IoProfile{
						IoProfile: intIOProfile,
					},
					Sharedv4Opt: &api.VolumeSpecUpdate_Sharedv4{
						Sharedv4: sharedv4Enabled,
					},
					NodiscardOpt: &api.VolumeSpecUpdate_Nodiscard{
						Nodiscard: noDiscard,
					},
				},
			})
		if err != nil {
			gerr, _ := status.FromError(err)
			newError := fmt.Sprintf("error code[%d] message[%s]", gerr.Code(), gerr.Message())
			log.Println(newError)
			return "", errors.New(newError)
		}

		volumeUpdateResponse := volume.String()

		log.Printf("Volume %s updated", volumeUpdateResponse)

	}

	return volumeID, nil
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

// updateVolumeSize updates the size of a Portworx volume.
func UpdateVolumeSize(conn *grpc.ClientConn, volumeName string, volSize uint64) error {

	// Retrieves the volume ID.
	volId, err := GetVolumeID(conn, volumeName)
	if err != nil {
		log.Println(err)
		return err
	}

	// Opens the volume client connection.
	volumes := api.NewOpenStorageVolumeClient(conn)

	// Updates the volume size.
	volumeUpdate, err := volumes.Update(
		context.Background(),
		&api.SdkVolumeUpdateRequest{
			VolumeId: volId,
			Spec: &api.VolumeSpecUpdate{
				SizeOpt: &api.VolumeSpecUpdate_Size{
					Size: volSize * 1024 * 1024 * 1024,
				},
			},
		},
	)
	if err != nil {
		gerr, _ := status.FromError(err)
		log.Printf("Error Code[%d] Message[%s]\n",
			gerr.Code(), gerr.Message())
		return err
	}

	log.Printf("Volume %s updated size to %dGi %s\n", volumeName, volSize, volumeUpdate.String())

	return nil
}

// updateVolumeHALevel updates the HA level of a Portworx volume.
func UpdateVolumeHALevel(conn *grpc.ClientConn, volumeName string, haLevel int64) error {

	// Retrieves the volume ID.
	volId, err := GetVolumeID(conn, volumeName)
	if err != nil {
		log.Println(err)
		return err
	}

	// Opens the volume client connection.
	volumes := api.NewOpenStorageVolumeClient(conn)

	// Updates the volume HA level.
	volumeUpdate, err := volumes.Update(
		context.Background(),
		&api.SdkVolumeUpdateRequest{
			VolumeId: volId,
			Spec: &api.VolumeSpecUpdate{
				HaLevelOpt: &api.VolumeSpecUpdate_HaLevel{
					HaLevel: haLevel,
				},
			},
		},
	)
	if err != nil {
		gerr, _ := status.FromError(err)
		log.Printf("Error Code[%d] Message[%s]\n",
			gerr.Code(), gerr.Message())
		return err
	}

	log.Printf("Volume %s updated HA Level to %d replicas %s\n", volumeName, haLevel, volumeUpdate.String())

	return nil
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

	volumeName := volumeInspect.Volume.Locator.GetName()
	volumeReplicas := len(volumeInspect.Volume.ReplicaSets[0].GetNodes())
	volumeReplicaNodes := volumeInspect.Volume.ReplicaSets[0].GetNodes()
	volumeIOProfile := volumeInspect.Volume.GetLocator().GetVolumeLabels()["io_profile"]
	volumeIOProfileAPI := volumeInspect.Volume.Spec.GetIoProfile().String()
	volumeIOPriority := volumeInspect.Volume.GetLocator().GetVolumeLabels()["io_priority"]
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
	volumeEncrypted := volumeInspect.Volume.GetLocator().GetVolumeLabels()["secure"]
	volumeEncryptionKey := volumeInspect.Volume.GetLocator().GetVolumeLabels()["secret_key"]
	volumeK8sNamespace := volumeInspect.Volume.GetLocator().GetVolumeLabels()["namespace"]
	volumeK8sPVCName := volumeInspect.Volume.GetLocator().GetVolumeLabels()["pvc"]
	volumeSharedv4 := volumeInspect.Volume.Spec.GetSharedv4()
	volumeSharedv4ServiceSpec := volumeInspect.Volume.Spec.GetSharedv4ServiceSpec()
	volumeIOStrategy := volumeInspect.Volume.Spec.GetIoStrategy()

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

	// For each volume ID, get its information
	for _, volID := range volsEnumResp.GetVolumeIds() {

		volumeInfo, errorFound := GetVolumeInfo(conn, volID)
		if errorFound != nil {
			log.Println(errorFound)
			return AllVolumesInfo, errorFound
		}

		AllVolumesInfo = append(AllVolumesInfo, volumeInfo)

	}

	return AllVolumesInfo, nil

}
