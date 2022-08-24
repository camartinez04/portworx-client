package snapshots

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/camartinez04/portworx-cli/broker/pkg/volumes"
	"github.com/golang/protobuf/ptypes"
	api "github.com/libopenstorage/openstorage-sdk-clients/sdk/golang"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

//createSnapshot creates a local snapshot of a volume
func createSnapshot(conn *grpc.ClientConn, volumeName string) {

	volumenes := api.NewOpenStorageVolumeClient(conn)

	volumeID, err := volumes.GetVolumeID(conn, volumeName)
	if err != nil {
		log.Fatal(err)
	}

	// Take a snapshot
	snap, err := volumenes.SnapshotCreate(
		context.Background(),
		&api.SdkVolumeSnapshotCreateRequest{
			VolumeId: volumeID,
			Name:     fmt.Sprintf("snap-%v", time.Now().Unix()),
		},
	)
	if err != nil {
		gerr, _ := status.FromError(err)
		fmt.Printf("Error Code[%d] Message[%s]\n",
			gerr.Code(), gerr.Message())
		os.Exit(1)
	}
	fmt.Printf("Snapshot with id %s was create for volume %s\n",
		snap.GetSnapshotId(),
		volumeID)
	fmt.Println()

}

//createCloudSnap creates a cloud snapshot
func createCloudSnap(conn *grpc.ClientConn, volumeName string) {

	volumeID, err := volumes.GetVolumeID(conn, volumeName)
	if err != nil {
		log.Fatal(err)
	}

	// Create a backup to a cloud provider of our volume
	cloudbackups := api.NewOpenStorageCloudBackupClient(conn)

	backupCreateResp, err := cloudbackups.Create(context.Background(),
		&api.SdkCloudBackupCreateRequest{
			VolumeId:     volumeID,
			CredentialId: "f28b6b49-470a-4489-b30c-613ec5d5f801",
		})
	if err != nil {
		gerr, _ := status.FromError(err)
		fmt.Printf("Error Code[%d] Message[%s]\n",
			gerr.Code(), gerr.Message())
		os.Exit(1)
	}
	fmt.Printf("Backup started for volume %s with task id %s\n",
		volumeID,
		backupCreateResp.GetTaskId())

}

//statusCloudSnap gets the status of a cloud snapshot
func statusCloudSnap(conn *grpc.ClientConn, volumeName string) {

	volumeID, err := volumes.GetVolumeID(conn, volumeName)
	if err != nil {
		log.Fatal(err)
	}

	cloudbackups := api.NewOpenStorageCloudBackupClient(conn)

	// Now check the status of the backup
	backupStatus, err := cloudbackups.Status(context.Background(),
		&api.SdkCloudBackupStatusRequest{
			VolumeId: volumeID,
		})
	if err != nil {
		gerr, _ := status.FromError(err)
		fmt.Printf("Error Code[%d] Message[%s]\n",
			gerr.Code(), gerr.Message())
		os.Exit(1)
	}
	for taskId, status := range backupStatus.GetStatuses() {
		// There will be only one value in the map, but we use
		// a for-loop as an example.
		b, _ := json.MarshalIndent(status, "", "  ")
		fmt.Printf("Backup status for taskId: %s\n"+
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
	fmt.Println()

}

//cloudSnapHistory gets the history of a cloud snapshot
func cloudSnapHistory(conn *grpc.ClientConn, volumeName string) {

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
		fmt.Printf("Error Code[%d] Message[%s]\n",
			gerr.Code(), gerr.Message())
		os.Exit(1)
	}

	fmt.Printf("Backup history for volume %s:\n", volumeID)
	for _, history := range historyResp.GetHistoryList() {

		timestamp, _ := ptypes.Timestamp(history.GetTimestamp())
		fmt.Printf("Volume:%s \tttime:%v \tstatus:%v\n",
			history.GetSrcVolumeId(),
			timestamp,
			history.GetStatus())
	}
	fmt.Println()
}

//createCloudCredentials creates a new cloud credential for the given provider
func createS3CloudCredentials(conn *grpc.ClientConn, credName string, bucketName string, accessKey string, secretKey string, endPoint string, region string, sslDisabled bool, iamPolicyEnabled bool) {

	creds := api.NewOpenStorageCredentialsClient(conn)
	credResponse, err := creds.Create(
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
	if err != nil {
		gerr, _ := status.FromError(err)
		fmt.Printf("Error Code[%d] Message[%s]\n",
			gerr.Code(), gerr.Message())
		os.Exit(1)
	}
	credID := credResponse.GetCredentialId()
	fmt.Printf("Credential named %s created with id %s\n", credName, credID)
	fmt.Println()
}
