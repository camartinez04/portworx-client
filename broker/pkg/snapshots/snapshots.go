package snapshots

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/camartinez04/portworx-client/broker/pkg/volumes"
	"github.com/golang/protobuf/ptypes"
	api "github.com/libopenstorage/openstorage-sdk-clients/sdk/golang"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// CreateLocalSnap creates a local snapshot of a volume
func CreateLocalSnap(conn *grpc.ClientConn, volumeID string) (snapID string, errorFound error) {

	volumenes := api.NewOpenStorageVolumeClient(conn)

	// Take a snapshot
	snap, errorFound := volumenes.SnapshotCreate(
		context.Background(),
		&api.SdkVolumeSnapshotCreateRequest{
			VolumeId: volumeID,
			Name:     fmt.Sprintf("snap-%v", time.Now().Unix()),
		},
	)
	if errorFound != nil {
		log.Printf("Error creating snapshot: %v", errorFound)
		return "", errorFound
	}

	log.Printf("Snapshot with id %s was create for volume %s", snap.GetSnapshotId(), volumeID)

	snapID = snap.GetSnapshotId()

	return snapID, nil

}

// CreateCloudSnap creates a cloud snapshot
func CreateCloudSnap(conn *grpc.ClientConn, volumeID string, credentialID string) (taskID string, errorFound error) {

	// Create a backup to a cloud provider of our volume
	cloudbackups := api.NewOpenStorageCloudBackupClient(conn)

	backupCreateResp, errorFound := cloudbackups.Create(context.Background(),
		&api.SdkCloudBackupCreateRequest{
			VolumeId:     volumeID,
			CredentialId: credentialID,
		})
	if errorFound != nil {
		log.Printf("Error creating backup: %v", errorFound)
		return "", errorFound
	}

	log.Printf("Backup created with task id: %s\n", backupCreateResp.GetTaskId())

	taskID = backupCreateResp.GetTaskId()

	return taskID, nil

}

// StatusCloudSnap gets the status of a cloud snapshot
func StatusCloudSnap(conn *grpc.ClientConn, volumeID string) (jsonStatus string, errorFound error) {

	cloudbackups := api.NewOpenStorageCloudBackupClient(conn)

	// Now check the status of the backup
	backupStatus, errorFound := cloudbackups.Status(context.Background(),
		&api.SdkCloudBackupStatusRequest{
			VolumeId: volumeID,
		})
	if errorFound != nil {
		log.Printf("Error getting backup status: %v", errorFound)
		return "", errorFound
	}

	for taskId, status := range backupStatus.GetStatuses() {
		// There will be only one value in the map, but we use
		// a for-loop as an example.
		b, errorFound := json.MarshalIndent(status, "", "  ")
		if errorFound != nil {
			log.Printf("Error marshalling status: %v", errorFound)
			return "", errorFound
		}

		jsonStatus = fmt.Sprintf("Backup status for taskId: %s\n"+
			"Volume: %s\n"+
			"Type: %s\n"+
			"Status: %s\n"+
			"Full JSON Response: %s\n",
			taskId,
			status.GetSrcVolumeId(),
			status.GetOptype().String(),
			status.GetStatus().String(),
			string(b))
	}

	return jsonStatus, nil

}

// GetCloudSnaps gets the cloud snapshots of a volume
func GetCloudSnaps(conn *grpc.ClientConn, volumeID string) (cloudSnapsMap map[string][]*api.SdkCloudBackupInfo, errorFound error) {

	cloudbackups := api.NewOpenStorageCloudBackupClient(conn)

	// Get the list of cloud credentials
	credIDsList, errorFound := ListCloudCredentialIDs(conn)
	if errorFound != nil {
		log.Printf("Error getting cloud credentials: %v", errorFound)
		return nil, errorFound
	}

	// If there are no cloud credentials, return nil but without error too.
	if credIDsList == nil {
		log.Printf("No cloud credentials found")
		return nil, nil
	}

	cloudSnapsMap = make(map[string][]*api.SdkCloudBackupInfo)

	// Iterate over the list of cloud credentials
	for _, credID := range credIDsList {

		// Now check the status of the backup
		backupStatus, errorFound := cloudbackups.EnumerateWithFilters(
			context.Background(),
			&api.SdkCloudBackupEnumerateWithFiltersRequest{
				SrcVolumeId:  volumeID,
				CredentialId: credID,
			})
		if errorFound != nil {
			log.Printf("Error getting backup status: %v", errorFound)
			return nil, errorFound
		}

		backupList := backupStatus.GetBackups()

		cloudSnapsMap[credID] = backupList

	}

	return cloudSnapsMap, nil

}

// GetSpecificCloudSnapshot gets a specific cloud snapshot having the CloudSnap ID
func GetSpecificCloudSnapshot(conn *grpc.ClientConn, cloudSnapID string) (cloudSnap *api.SdkCloudBackupInfo, errorFound error) {

	cloudbackups := api.NewOpenStorageCloudBackupClient(conn)

	// Get the list of cloud credentials
	credIDsList, errorFound := ListCloudCredentialIDs(conn)
	if errorFound != nil {
		log.Printf("Error getting cloud credentials: %v", errorFound)
		return nil, errorFound
	}

	// Iterate over the list of cloud credentials
	for _, credID := range credIDsList {

		cloudSnapInfo, errorFound := cloudbackups.EnumerateWithFilters(
			context.Background(),
			&api.SdkCloudBackupEnumerateWithFiltersRequest{
				CloudBackupId: cloudSnapID,
				CredentialId:  credID,
			})
		if errorFound != nil {
			log.Printf("Error getting backup status: %v", errorFound)
			return nil, errorFound
		}

		cloudSnapList := cloudSnapInfo.GetBackups()

		if len(cloudSnapList) == 0 {
			log.Printf("No cloud snapshot found with ID %s", cloudSnapID)
			return nil, nil
		}

		cloudSnap = cloudSnapList[0]

	}

	return cloudSnap, nil

}

// DeleteCloudSnap deletes a cloud snapshot
func DeleteCloudSnap(conn *grpc.ClientConn, credentialID string, cloudSnapID string) (errorFound error) {

	cloudbackups := api.NewOpenStorageCloudBackupClient(conn)

	// Delete the cloud snapshot
	_, errorFound = cloudbackups.Delete(context.Background(),
		&api.SdkCloudBackupDeleteRequest{
			CredentialId: credentialID,
			BackupId:     cloudSnapID,
		})
	if errorFound != nil {
		log.Printf("Error deleting backup: %v", errorFound)
		return errorFound
	}

	return nil

}

// CloudSnapHistory gets the history of a cloud snapshot
func CloudSnapHistory(conn *grpc.ClientConn, volumeName string) {

	volumeID, err := volumes.GetVolumeID(conn, volumeName)
	if err != nil {
		log.Fatal(err)
	}

	cloudbackups := api.NewOpenStorageCloudBackupClient(conn)

	// Backup History
	historyResp, err := cloudbackups.History(context.Background(),
		&api.SdkCloudBackupHistoryRequest{
			SrcVolumeId: volumeID,
		})
	if err != nil {
		gerr, _ := status.FromError(err)
		log.Printf("Error Code[%d] Message[%s]\n",
			gerr.Code(), gerr.Message())
		os.Exit(1)
	}

	log.Printf("Backup history for volume %s:\n", volumeID)
	for _, history := range historyResp.GetHistoryList() {

		timestamp, _ := ptypes.Timestamp(history.GetTimestamp())
		log.Printf("Volume:%s \tttime:%v \tstatus:%v\n",
			history.GetSrcVolumeId(),
			timestamp,
			history.GetStatus())
	}
}

// AllCloudSnapsCluster retrieves all the cloud snapshots of all the volumes in the Portworx cluster
func AllCloudSnapsCluster(conn *grpc.ClientConn) (cloudSnaps map[string]map[string][]*api.SdkCloudBackupInfo, errorFound error) {

	cloudSnaps = make(map[string]map[string][]*api.SdkCloudBackupInfo)

	// Get all the volumes in the cluster into a slice of strings
	volumes, errorFound := volumes.GetAllVolumes(conn)
	if errorFound != nil {
		log.Fatal(errorFound)
		return nil, errorFound
	}

	// Iterate over the volumes list and get the cloud snapshots of each volume, populating the map as well.
	for _, volume := range volumes {

		// Get the cloud snapshots of the volume
		snapsOfVolume, errorFound := GetCloudSnaps(conn, volume)
		if errorFound != nil {
			log.Fatal(errorFound)
			return nil, errorFound
		}

		cloudSnaps[volume] = make(map[string][]*api.SdkCloudBackupInfo)

		for credID, snaps := range snapsOfVolume {

			cloudSnaps[volume][credID] = snaps

		}
	}

	return cloudSnaps, nil
}

// AWSCreateS3CloudCredential creates a new AWS S3 Cloud credential
func AWSCreateS3CloudCredential(conn *grpc.ClientConn, credName string, bucketName string, accessKey string, secretKey string, endPoint string, region string, sslDisabled bool, iamPolicyEnabled bool) (credID string, errorFound error) {

	creds := api.NewOpenStorageCredentialsClient(conn)

	credResponse, errorFound := creds.Create(
		context.Background(),
		&api.SdkCredentialCreateRequest{
			Name:      credName,
			IamPolicy: iamPolicyEnabled,
			Bucket:    bucketName,
			CredentialType: &api.SdkCredentialCreateRequest_AwsCredential{
				AwsCredential: &api.SdkAwsCredentialRequest{
					AccessKey:  accessKey,
					SecretKey:  secretKey,
					Endpoint:   endPoint,
					Region:     region,
					DisableSsl: sslDisabled,
				},
			},
		})
	if errorFound != nil {
		log.Printf("Error creating credential: %v", errorFound)
		return "", errorFound
	}

	credID = credResponse.GetCredentialId()
	log.Printf("Credential named %s created with id %s\n", credName, credID)

	return credID, nil
}

// AWSValidateS3CloudCredential validates the given an AWS S3 Cloud credential
func AWSValidateS3CloudCredential(conn *grpc.ClientConn, credentialId string) error {

	creds := api.NewOpenStorageCredentialsClient(conn)
	credResponse, err := creds.Validate(
		context.Background(),
		&api.SdkCredentialValidateRequest{
			CredentialId: credentialId,
		})
	if err != nil {
		gerr, _ := status.FromError(err)
		log.Printf("Error Code[%d] Message[%s]\n",
			gerr.Code(), gerr.Message())
		os.Exit(1)
	}

	response := credResponse.String()

	log.Printf("Credential ID %s validated with response %s", credentialId, response)

	return nil

}

// AWSInspectS3CloudCredential inspects the given an AWS S3 Cloud credential
func AWSInspectS3CloudCredential(conn *grpc.ClientConn, credentialId string) (credResponse *api.SdkCredentialInspectResponse, errorFound error) {

	creds := api.NewOpenStorageCredentialsClient(conn)
	credResponse, err := creds.Inspect(
		context.Background(),
		&api.SdkCredentialInspectRequest{
			CredentialId: credentialId,
		})
	if err != nil {
		gerr, _ := status.FromError(err)
		log.Printf("Error Code[%d] Message[%s]\n",
			gerr.Code(), gerr.Message())
		os.Exit(1)
	}

	log.Printf("Credential ID %s inspected with name %s", credentialId, credResponse.GetName())

	return credResponse, nil

}

// AWSDeleteS3CloudCredentials deletes an AWS S3 Cloud credential
func AWSDeleteS3CloudCredential(conn *grpc.ClientConn, credentialId string) error {

	creds := api.NewOpenStorageCredentialsClient(conn)
	_, err := creds.Delete(
		context.Background(),
		&api.SdkCredentialDeleteRequest{
			CredentialId: credentialId,
		})
	if err != nil {
		gerr, _ := status.FromError(err)
		log.Printf("Error Code[%d] Message[%s]\n", gerr.Code(), gerr.Message())
		return err
	}

	log.Printf("Credential with ID %s has been deleted", credentialId)

	return nil
}

// ListCloudCredentialIDs lists all the cloud credentials ids in the cluster
func ListCloudCredentialIDs(conn *grpc.ClientConn) (CredIds []string, errorFound error) {

	creds := api.NewOpenStorageCredentialsClient(conn)
	credResponse, err := creds.Enumerate(
		context.Background(),
		&api.SdkCredentialEnumerateRequest{},
	)
	if err != nil {
		gerr, _ := status.FromError(err)
		log.Printf("Error Code[%d] Message[%s]\n",
			gerr.Code(), gerr.Message())
		os.Exit(1)
	}

	CredIds = credResponse.GetCredentialIds()

	return CredIds, nil
}
