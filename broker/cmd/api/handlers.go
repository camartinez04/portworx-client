package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/camartinez04/portworx-client/broker/pkg/cluster"
	"github.com/camartinez04/portworx-client/broker/pkg/nodes"
	"github.com/camartinez04/portworx-client/broker/pkg/snapshots"
	"github.com/camartinez04/portworx-client/broker/pkg/volumes"
)

func NewHandlers(App *AppConfig) *AppConfig {
	Application := App
	return Application
}

// GetVolumeIDHTTP http function to get the volume ID.
func (App *AppConfig) getVolumeIDsHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	context.Background()

	volumeName := exploded[3]

	// http://localhost:8080/getvolumeid

	volumeID, err := volumes.GetVolumeID(App.Conn, volumeName)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonResponse{
		Error:    false,
		VolumeID: volumeID,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

func (App *AppConfig) getInspectVolumeHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeName := exploded[3]

	volume, replicas, volumenodes, status, ioprofile, err := volumes.InspectVolume(App.Conn, volumeName)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonVolumeInspect{
		Error:              false,
		VolumeInspect:      volume,
		ReplicasInfo:       replicas,
		VolumeNodes:        volumenodes,
		VolumeStatusString: status,
		IoProfileString:    ioprofile,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getClusterCapacityHTTP http function to get the cluster capacity.
func (App *AppConfig) getPXClusterCapacityHTTP(w http.ResponseWriter, r *http.Request) {

	cluster, used, available, percentused, percentavailable, err := cluster.PXClusterCapacity(App.Conn)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonClusterCapacity{
		Error:                   false,
		ClusterCapacity:         cluster,
		ClusterUsed:             used,
		ClusterAvailable:        available,
		ClusterPercentUsed:      percentused,
		ClusterPercentAvailable: percentavailable,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getClusterCapacityHTTP http function to get the cluster capacity.
func (App *AppConfig) getPXClusterHTTP(w http.ResponseWriter, r *http.Request) {

	uuid, status, name, err := cluster.PXClusterInfo(App.Conn)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonClusterInfo{
		Error:         false,
		ClusterUUID:   uuid,
		ClusterStatus: status,
		ClusterName:   name,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// postCreateNewVolumeHTTP http function to create a new volume.
func (App *AppConfig) postCreateNewVolumeHTTP(w http.ResponseWriter, r *http.Request) {

	volumeName := r.Header.Get("Volume-Name")

	volumeGBSize, err := strconv.ParseUint((r.Header.Get("Volume-Size")), 10, 64)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	volumeIOProfile := r.Header.Get("Volume-IO-Profile")

	volumeHALevel, err := strconv.ParseInt((r.Header.Get("Volume-Ha-Level")), 10, 64)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	encryptionEnabled, err := strconv.ParseBool(r.Header.Get("Volume-Encryption-Enabled"))
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	sharedv4Enabled, err := strconv.ParseBool(r.Header.Get("Volume-Sharedv4-Enabled"))
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	noDiscard, err := strconv.ParseBool(r.Header.Get("Volume-No-Discard"))
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	newVolumeID, err := volumes.CreateVolume(App.Conn, volumeName, volumeGBSize, volumeIOProfile, volumeHALevel, encryptionEnabled, sharedv4Enabled, noDiscard)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonResponse{
		Error:    false,
		Message:  "Volume created successfully",
		VolumeID: newVolumeID,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getNodesOfVolumeHTTP http function to get the nodes of a volume.
func (App *AppConfig) getNodesOfVolumeHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeName := exploded[3]

	nodes, err := nodes.FindVolumeNodes(App.Conn, volumeName)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonNodesOfVolume{
		Error:         false,
		NodesOfVolume: nodes,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getListOfNodesHTTP http function to get the list of nodes of the Portworx cluster.
func (App *AppConfig) getListOfNodesHTTP(w http.ResponseWriter, r *http.Request) {

	nodeList, err := nodes.GetListOfNodes(App.Conn)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonNodeList{
		Error:    false,
		NodeList: nodeList,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getReplicasPerNodeHTTP http function to get the replicas per node.
func (App *AppConfig) getReplicasPerNodeHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	nodeID := exploded[3]

	volumes, err := nodes.GetReplicasPerNode(App.Conn, nodeID)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonVolumeList{
		Error:      false,
		VolumeList: volumes,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getVolumeUsageHTTP http function to get the volume usage.
func (App *AppConfig) getVolumeUsageHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeName := exploded[3]

	var volUsageFloat, availSpaceFloat, totalSizeFloat float64

	var volUsagePercentFloat, volAvailablePercentFloat float32

	volumeUsage, availableSpace, totalSize, err := volumes.RetrieveVolumeUsage(App.Conn, volumeName)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	volUsageFloat = float64(volumeUsage / 1024 / 1024)

	availSpaceFloat = float64(availableSpace / 1024 / 1024)

	totalSizeFloat = float64(totalSize / 1024 / 1024)

	volUsagePercentFloat = float32(volUsageFloat / totalSizeFloat * 100)

	volAvailablePercentFloat = float32(availSpaceFloat / totalSizeFloat * 100)

	resp := JsonVolumeUsage{
		Error:                  false,
		VolumeUsage:            volUsageFloat,
		AvailableSpace:         availSpaceFloat,
		TotalSize:              totalSizeFloat,
		VolumeUsagePercent:     volUsagePercentFloat,
		VolumeAvailablePercent: volAvailablePercentFloat,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getAllVolumesHTTP http function to get the list of volumes.
func (App *AppConfig) getAllVolumesHTTP(w http.ResponseWriter, r *http.Request) {

	volumes, err := volumes.GetAllVolumes(App.Conn)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonAllVolumesList{
		Error:          false,
		AllVolumesList: volumes,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getAllVolumesCompleteHTTP http function to get the list of volumes with inspect information included.
func (App *AppConfig) getAllVolumesCompleteHTTP(w http.ResponseWriter, r *http.Request) {

	volumes, err := volumes.GetAllVolumesComplete(App.Conn)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonApiVolumesList{
		Error:          false,
		ApiVolumesList: volumes,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getNodeInfoHTTP http function to get the node information.
func (App *AppConfig) getNodeInfoHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	nodeID := exploded[3]

	nodeInformation, err := nodes.GetNodeInfo(App.Conn, nodeID)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonGetNodeInfo{
		Error:    false,
		NodeInfo: nodeInformation,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getAllNodesInfoHTTP http function to get the list of nodes with inspect information included.
func (App *AppConfig) getAllNodesInfoHTTP(w http.ResponseWriter, r *http.Request) {

	nodesInformation, err := nodes.GetAllNodesInfo(App.Conn)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonGetAllNodesInfo{
		Error:        false,
		AllNodesInfo: nodesInformation,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getAllVolumesInfoHTTP http function to get the list of nodes with inspect information included.
func (App *AppConfig) getAllVolumesInfoHTTP(w http.ResponseWriter, r *http.Request) {

	volumesInformation, err := volumes.GetAllVolumesInfo(App.Conn)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonGetAllVolumesInfo{
		Error:          false,
		AllVolumesInfo: volumesInformation,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getNodeVolumesHTTP http function to get the node information.
func (App *AppConfig) getVolumeInfoHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[3]

	volumeInformation, err := volumes.GetVolumeInfo(App.Conn, volumeID)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonGetVolumeInfo{
		Error:      false,
		VolumeInfo: volumeInformation,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// patchUpdateVolumeSizeHTTP http function to update a Portworx Volume Size.
func (App *AppConfig) patchUpdateVolumeSizeHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[3]

	volumeGBSize, err := strconv.ParseUint((r.Header.Get("Volume-Size")), 10, 64)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	_, err = volumes.UpdateVolumeSize(App.Conn, volumeID, volumeGBSize)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonResponse{
		Error:    false,
		Message:  "Volume size updated to " + strconv.FormatUint(volumeGBSize, 10) + "GB",
		VolumeID: volumeID,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// patchUpdateVolumeSizeHTTP http function to update a Portworx Volume Size.
func (App *AppConfig) patchUpdateVolumeIOProfileHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[3]

	volumeIOProfile := r.Header.Get("Volume-IO-Profile")

	_, err := volumes.UpdateVolumeIOProfile(App.Conn, volumeID, volumeIOProfile)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonResponse{
		Error:    false,
		Message:  "Volume IO Profile updated to " + volumeIOProfile,
		VolumeID: volumeID,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// patchUpdateVolumeHALevelHTTP http function to update a Portworx Volume HA Level.
func (App *AppConfig) patchUpdateVolumeHALevelHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[3]

	volumeHALevel, err := strconv.ParseInt((r.Header.Get("Volume-Ha-Level")), 10, 64)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	_, err = volumes.UpdateVolumeHALevel(App.Conn, volumeID, volumeHALevel)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonResponse{
		Error:    false,
		Message:  "Volume HA Level updated to " + strconv.FormatInt(volumeHALevel, 10),
		VolumeID: volumeID,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// patchUpdateVolumeSharedv4HTTP http function to define a Portworx Volume as Sharedv4 or not.
func (App *AppConfig) patchUpdateVolumeSharedv4HTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[3]

	sharedv4Enabled, err := strconv.ParseBool(r.Header.Get("Volume-Sharedv4-Enabled"))
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	_, err = volumes.UpdateVolumeSharedv4(App.Conn, volumeID, sharedv4Enabled)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonResponse{
		Error:    false,
		Message:  "Volume Sharedv4 updated to " + strconv.FormatBool(sharedv4Enabled),
		VolumeID: volumeID,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// patchUpdateVolumeSharedvService4HTTP http function to enable or disable a Portworx Volume Sharedv4 Service.
func (App *AppConfig) patchUpdateVolumeSharedvService4HTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[3]

	sharedv4Service, err := strconv.ParseBool(r.Header.Get("Volume-Sharedv4-Service"))
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	_, err = volumes.UpdateVolumeSharedv4Service(App.Conn, volumeID, sharedv4Service)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonResponse{
		Error:    false,
		Message:  "Volume Sharedv4 Service updated to " + strconv.FormatBool(sharedv4Service),
		VolumeID: volumeID,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// patchUpdateVolumeNoDiscardHTTP http function to enable or disable a Portworx Volume No Discard.
func (App *AppConfig) patchUpdateVolumeNoDiscardHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[3]

	noDiscard, err := strconv.ParseBool(r.Header.Get("Volume-No-Discard"))
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	_, err = volumes.UpdateVolumeNoDiscard(App.Conn, volumeID, noDiscard)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonResponse{
		Error:    false,
		Message:  "Volume No Discard updated to " + strconv.FormatBool(noDiscard),
		VolumeID: volumeID,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// deleteVolumeHTTP http function to delete a Portworx Volume.
func (App *AppConfig) deleteVolumeHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[3]

	volume, err := volumes.DeleteVolume(App.Conn, volumeID)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	log.Printf("Volume %s deleted successfully, %s", volumeID, volume)

	resp := JsonResponse{
		Error:    false,
		Message:  "Volume deleted successfully",
		VolumeID: volumeID,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getCloudSnapsHTTP http function to get a list of Portworx CloudSnaps.
func (App *AppConfig) getCloudSnapsHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[3]

	var cloudCredList []string

	cloudSnaps, err := snapshots.GetCloudSnaps(App.Conn, volumeID, cloudCredList)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonCloudSnapList{
		Error:         false,
		CloudSnapList: cloudSnaps,
	}

	writeJSON(w, http.StatusOK, resp)

}

// getInspectAWSCloudCredentialHTTP http function to get a Portworx AWS Cloud Credential.
func (App *AppConfig) getInspectAWSCloudCredentialHTTP(w http.ResponseWriter, r *http.Request) {

	cloudCredentialID := r.Header.Get("Cloud-Credential-ID")

	cloudCred, err := snapshots.AWSInspectS3CloudCredential(App.Conn, cloudCredentialID)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	AccessKey := cloudCred.GetAwsCredential().AccessKey
	Endpoint := cloudCred.GetAwsCredential().Endpoint
	Region := cloudCred.GetAwsCredential().Region

	resp := JsonCredentialInspect{
		Error:             false,
		CredentialInspect: *cloudCred,
		AccessKey:         AccessKey,
		Endpoint:          Endpoint,
		Region:            Region,
	}

	writeJSON(w, http.StatusOK, resp)

}

// postCreateAWSCloudCredentialHTTP http function to create a Portworx AWS Cloud Credential.
func (App *AppConfig) postCreateAWSCloudCredentialHTTP(w http.ResponseWriter, r *http.Request) {

	credName := r.Header.Get("Cloud-Credential-Name")
	credAccessKey := r.Header.Get("Cloud-Credential-Access-Key")
	credBucketName := r.Header.Get("Cloud-Credential-Bucket-Name")
	credSecretKey := r.Header.Get("Cloud-Credential-Secret-Key")
	credRegion := r.Header.Get("Cloud-Credential-Region")
	credEndpoint := r.Header.Get("Cloud-Credential-Endpoint")
	credSSL, err := strconv.ParseBool(r.Header.Get("Cloud-Credential-Disable-SSL"))
	if err != nil {
		App.errorJSON(w, err)
		return
	}
	credIAM, err := strconv.ParseBool(r.Header.Get("Cloud-Credential-IAM-Policy-Enabled"))
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	cloudCredentialID, err := snapshots.AWSCreateS3CloudCredential(App.Conn, credName, credBucketName, credAccessKey, credSecretKey, credEndpoint, credRegion, credSSL, credIAM)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	//Validate the Cloud Credential
	err = snapshots.AWSValidateS3CloudCredential(App.Conn, cloudCredentialID)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	//Inspect the Cloud Credential
	credDetails, err := snapshots.AWSInspectS3CloudCredential(App.Conn, cloudCredentialID)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	AccessKey := credDetails.GetAwsCredential().AccessKey
	Endpoint := credDetails.GetAwsCredential().Endpoint
	Region := credDetails.GetAwsCredential().Region

	resp := JsonCredentialInspect{
		Error:             false,
		CredentialInspect: *credDetails,
		AccessKey:         AccessKey,
		Endpoint:          Endpoint,
		Region:            Region,
	}

	writeJSON(w, http.StatusCreated, resp)

}

// deleteAWSCloudCredentialHTTP http function to delete a Portworx AWS Cloud Credential.
func (App *AppConfig) deleteAWSCloudCredentialHTTP(w http.ResponseWriter, r *http.Request) {

	cloudCredentialID := r.Header.Get("Cloud-Credential-ID")

	err := snapshots.AWSDeleteS3CloudCredential(App.Conn, cloudCredentialID)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonResponse{
		Error:   false,
		Message: "Cloud Credential deleted successfully",
		CredID:  cloudCredentialID,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getPXClusterAlarmsHTTP http function to get a list of Portworx Cluster Alarms.
func (App *AppConfig) getPXClusterAlarmsHTTP(w http.ResponseWriter, r *http.Request) {

	alarms, err := cluster.PXClusterAlarms(App.Conn)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonAlarmList{
		Error:     false,
		AlarmList: alarms,
	}

	writeJSON(w, http.StatusOK, resp)

}

// postCreateLocalSnapHTTP http function to create a Portworx Local Snapshot.
func (App *AppConfig) postCreateLocalSnapHTTP(w http.ResponseWriter, r *http.Request) {

	volumeID := r.Header.Get("Volume-ID")

	createSnap, err := snapshots.CreateLocalSnap(App.Conn, volumeID)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonResponse{
		Error:   false,
		Message: "Local Snapshot created successfully",
		SnapID:  createSnap,
	}

	writeJSON(w, http.StatusAccepted, resp)
}

// postCreateCloudSnapHTTP http function to create a Portworx CloudSnap.
func (App *AppConfig) postCreateCloudSnapHTTP(w http.ResponseWriter, r *http.Request) {

	cloudCredentialID := r.Header.Get("Cloud-Credential-ID")

	volumeID := r.Header.Get("Volume-ID")

	createCloudSnap, err := snapshots.CreateCloudSnap(App.Conn, volumeID, cloudCredentialID)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonCloudSnap{
		Error:   false,
		Message: "Cloud Snap of volume %s created successfully" + volumeID,
		TaskID:  createCloudSnap,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getAllCloudSnapsHTTP http function to get a list of all Portworx CloudSnaps.
func (App *AppConfig) getAllCloudSnapsHTTP(w http.ResponseWriter, r *http.Request) {

	cloudSnaps, err := snapshots.AllCloudSnapsCluster(App.Conn)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonAllCloudSnapList{
		Error:          false,
		CloudSnapsList: cloudSnaps,
	}

	writeJSON(w, http.StatusInternalServerError, resp)
}

// getAllCloudCredentialIDsHTTP http function to get a list of all Portworx Cloud Credential IDs.
func (App *AppConfig) getAllCloudCredentialIDsHTTP(w http.ResponseWriter, r *http.Request) {

	cloudCreds, err := snapshots.ListCloudCredentialIDs(App.Conn)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonAllCloudCredsList{
		Error:          false,
		CloudCredsList: cloudCreds,
	}

	writeJSON(w, http.StatusInternalServerError, resp)
}

// getSpecificCloudSnapshotHTTP http function to get a specific Portworx CloudSnap.
func (App *AppConfig) getSpecificCloudSnapshotHTTP(w http.ResponseWriter, r *http.Request) {

	cloudSnapID := r.Header.Get("Cloud-Snap-ID")

	credentialID := r.Header.Get("Cloud-Credential-ID")

	cloudSnap, err := snapshots.GetSpecificCloudSnapshot(App.Conn, cloudSnapID, credentialID)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonSpecificCloudSnap{
		Error:       false,
		CloudSnap:   cloudSnap,
		CloudSnapId: cloudSnapID,
	}

	writeJSON(w, http.StatusInternalServerError, resp)
}

// deleteCloudSnapHTTP http function to delete a Portworx CloudSnap.
func (App *AppConfig) deleteCloudSnapHTTP(w http.ResponseWriter, r *http.Request) {

	cloudSnapID := r.Header.Get("Cloud-Snap-ID")

	err := snapshots.DeleteCloudSnap(App.Conn, cloudSnapID)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonResponse{
		Error:   false,
		Message: "Cloud Snap deleted successfully",
		SnapID:  cloudSnapID,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

func (App *AppConfig) postLoginHTTP(w http.ResponseWriter, r *http.Request) {

	_ = App.Session.RenewToken(r.Context())

	username := r.Header.Get("Username")

	password := r.Header.Get("Password")

	//log.Printf("username: %s", username)

	// Authenticate the user
	rq := &LoginRequest{username, password}

	jwt, err := App.NewKeycloak.gocloak.Login(r.Context(),
		App.NewKeycloak.clientId,
		App.NewKeycloak.clientSecret,
		App.NewKeycloak.realm,
		rq.Username,
		rq.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	rs := &LoginResponse{
		AccessToken:  jwt.AccessToken,
		RefreshToken: jwt.RefreshToken,
		ExpiresIn:    jwt.ExpiresIn,
	}

	//log.Printf("jwt: %v", jwt.AccessToken)

	rsJs, _ := json.Marshal(rs)

	_, _ = w.Write(rsJs)

	KeycloakToken = jwt.AccessToken

	KeycloakRefreshToken = jwt.RefreshToken

	resp := JsonResponse{
		Error:   false,
		Message: "Login successful",
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getLogoutHTTP logs a user out
func (App *AppConfig) getLogoutHTTP(w http.ResponseWriter, r *http.Request) {

	_ = App.Session.RenewToken(r.Context())

	App.NewKeycloak.gocloak.Logout(r.Context(),
		App.NewKeycloak.clientId,
		App.NewKeycloak.clientSecret,
		App.NewKeycloak.realm,
		KeycloakRefreshToken)

	KeycloakToken = ""
	KeycloakRefreshToken = ""

	_ = App.Session.Destroy(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

// patchUpdateVolumeReplicaSetHTTP http function to update a Portworx Volume Replica Set.
func (App *AppConfig) patchUpdateVolumeReplicaSetHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[3]

	poolReplicaSet1 := r.Header.Get("PoolUUID-ReplicaSet-1")

	poolReplicaSet2 := r.Header.Get("PoolUUID-ReplicaSet-2")

	poolReplicaSet3 := r.Header.Get("PoolUUID-ReplicaSet-3")

	var PoolUuids []string

	if poolReplicaSet1 != "" {
		PoolUuids = append(PoolUuids, poolReplicaSet1)
	}

	if poolReplicaSet2 != "" {
		PoolUuids = append(PoolUuids, poolReplicaSet2)
	}

	if poolReplicaSet3 != "" {
		PoolUuids = append(PoolUuids, poolReplicaSet3)
	}

	_, err := volumes.UpdateVolumeReplicaSet(App.Conn, volumeID, PoolUuids)
	if err != nil {
		App.errorJSON(w, err)
		return
	}

	resp := JsonResponse{
		Error:   false,
		Message: "Volume Replica Set updated successfully",
	}

	writeJSON(w, http.StatusAccepted, resp)

}
