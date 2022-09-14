package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/camartinez04/portworx-client/broker/pkg/cluster"
	"github.com/camartinez04/portworx-client/broker/pkg/nodes"
	"github.com/camartinez04/portworx-client/broker/pkg/snapshots"
	"github.com/camartinez04/portworx-client/broker/pkg/volumes"
)

// GetVolumeIDHTTP http function to get the volume ID.
func (app *AppConfig) getVolumeIDsHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	context.Background()

	volumeName := exploded[2]

	// http://localhost:8080/getvolumeid

	volumeID, err := volumes.GetVolumeID(app.Conn, volumeName)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JsonResponse{
		Error:    false,
		VolumeID: volumeID,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

func (app *AppConfig) getInspectVolumeHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeName := exploded[2]

	volume, replicas, volumenodes, status, ioprofile, err := volumes.InspectVolume(app.Conn, volumeName)
	if err != nil {
		app.errorJSON(w, err)
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
func (app *AppConfig) getPXClusterCapacityHTTP(w http.ResponseWriter, r *http.Request) {

	cluster, used, available, percentused, percentavailable, err := cluster.ClusterCapacity(app.Conn)
	if err != nil {
		app.errorJSON(w, err)
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
func (app *AppConfig) getPXClusterHTTP(w http.ResponseWriter, r *http.Request) {

	uuid, status, name, err := cluster.ClusterInfo(app.Conn)
	if err != nil {
		app.errorJSON(w, err)
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
func (app *AppConfig) postCreateNewVolumeHTTP(w http.ResponseWriter, r *http.Request) {

	volumeName := r.Header.Get("Volume-Name")

	volumeGBSize, err := strconv.ParseUint((r.Header.Get("Volume-Size")), 10, 64)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	volumeIOProfile := r.Header.Get("Volume-IO-Profile")

	volumeHALevel, err := strconv.ParseInt((r.Header.Get("Volume-Ha-Level")), 10, 64)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	encryptionEnabled, err := strconv.ParseBool(r.Header.Get("Volume-Encryption-Enabled"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	sharedv4Enabled, err := strconv.ParseBool(r.Header.Get("Volume-Sharedv4-Enabled"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	noDiscard, err := strconv.ParseBool(r.Header.Get("Volume-No-Discard"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	newVolumeID, err := volumes.CreateVolume(app.Conn, volumeName, volumeGBSize, volumeIOProfile, volumeHALevel, encryptionEnabled, sharedv4Enabled, noDiscard)
	if err != nil {
		app.errorJSON(w, err)
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
func (app *AppConfig) getNodesOfVolumeHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeName := exploded[2]

	nodes, err := nodes.FindVolumeNodes(app.Conn, volumeName)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JsonNodesOfVolume{
		Error:         false,
		NodesOfVolume: nodes,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getListOfNodesHTTP http function to get the list of nodes of the Portworx cluster.
func (app *AppConfig) getListOfNodesHTTP(w http.ResponseWriter, r *http.Request) {

	nodeList, err := nodes.GetListOfNodes(app.Conn)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JsonNodeList{
		Error:    false,
		NodeList: nodeList,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

func (app *AppConfig) getReplicasPerNodeHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	nodeID := exploded[2]

	volumes, err := nodes.GetReplicasPerNode(app.Conn, nodeID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JsonVolumeList{
		Error:      false,
		VolumeList: volumes,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getVolumeUsageHTTP http function to get the volume usage.
func (app *AppConfig) getVolumeUsageHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeName := exploded[2]

	var volUsageFloat, availSpaceFloat, totalSizeFloat float64

	var volUsagePercentFloat, volAvailablePercentFloat float32

	volumeUsage, availableSpace, totalSize, err := volumes.RetrieveVolumeUsage(app.Conn, volumeName)
	if err != nil {
		app.errorJSON(w, err)
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
func (app *AppConfig) getAllVolumesHTTP(w http.ResponseWriter, r *http.Request) {

	volumes, err := volumes.GetAllVolumes(app.Conn)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JsonAllVolumesList{
		Error:          false,
		AllVolumesList: volumes,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getAllVolumesCompleteHTTP http function to get the list of volumes with inspect information included.
func (app *AppConfig) getAllVolumesCompleteHTTP(w http.ResponseWriter, r *http.Request) {

	volumes, err := volumes.GetAllVolumesComplete(app.Conn)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JsonApiVolumesList{
		Error:          false,
		ApiVolumesList: volumes,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getNodeInfoHTTP http function to get the node information.
func (app *AppConfig) getNodeInfoHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	nodeID := exploded[2]

	nodeInformation, err := nodes.GetNodeInfo(app.Conn, nodeID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JsonGetNodeInfo{
		Error:    false,
		NodeInfo: nodeInformation,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getAllNodesInfoHTTP http function to get the list of nodes with inspect information included.
func (app *AppConfig) getAllNodesInfoHTTP(w http.ResponseWriter, r *http.Request) {

	nodesInformation, err := nodes.GetAllNodesInfo(app.Conn)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JsonGetAllNodesInfo{
		Error:        false,
		AllNodesInfo: nodesInformation,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getAllVolumesInfoHTTP http function to get the list of nodes with inspect information included.
func (app *AppConfig) getAllVolumesInfoHTTP(w http.ResponseWriter, r *http.Request) {

	volumesInformation, err := volumes.GetAllVolumesInfo(app.Conn)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JsonGetAllVolumesInfo{
		Error:          false,
		AllVolumesInfo: volumesInformation,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getNodeVolumesHTTP http function to get the node information.
func (app *AppConfig) getVolumeInfoHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[2]

	volumeInformation, err := volumes.GetVolumeInfo(app.Conn, volumeID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JsonGetVolumeInfo{
		Error:      false,
		VolumeInfo: volumeInformation,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// patchUpdateVolumeSizeHTTP http function to update a Portworx Volume Size.
func (app *AppConfig) patchUpdateVolumeSizeHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[2]

	volumeGBSize, err := strconv.ParseUint((r.Header.Get("Volume-Size")), 10, 64)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_, err = volumes.UpdateVolumeSize(app.Conn, volumeID, volumeGBSize)
	if err != nil {
		app.errorJSON(w, err)
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
func (app *AppConfig) patchUpdateVolumeIOProfileHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[2]

	volumeIOProfile := r.Header.Get("Volume-IO-Profile")

	_, err := volumes.UpdateVolumeIOProfile(app.Conn, volumeID, volumeIOProfile)
	if err != nil {
		app.errorJSON(w, err)
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
func (app *AppConfig) patchUpdateVolumeHALevelHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[2]

	volumeHALevel, err := strconv.ParseInt((r.Header.Get("Volume-Ha-Level")), 10, 64)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_, err = volumes.UpdateVolumeHALevel(app.Conn, volumeID, volumeHALevel)
	if err != nil {
		app.errorJSON(w, err)
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
func (app *AppConfig) patchUpdateVolumeSharedv4HTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[2]

	sharedv4Enabled, err := strconv.ParseBool(r.Header.Get("Volume-Sharedv4-Enabled"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_, err = volumes.UpdateVolumeSharedv4(app.Conn, volumeID, sharedv4Enabled)
	if err != nil {
		app.errorJSON(w, err)
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
func (app *AppConfig) patchUpdateVolumeSharedvService4HTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[2]

	sharedv4Service, err := strconv.ParseBool(r.Header.Get("Volume-Sharedv4-Service"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_, err = volumes.UpdateVolumeSharedv4Service(app.Conn, volumeID, sharedv4Service)
	if err != nil {
		app.errorJSON(w, err)
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
func (app *AppConfig) patchUpdateVolumeNoDiscardHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[2]

	noDiscard, err := strconv.ParseBool(r.Header.Get("Volume-No-Discard"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_, err = volumes.UpdateVolumeNoDiscard(app.Conn, volumeID, noDiscard)
	if err != nil {
		app.errorJSON(w, err)
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
func (app *AppConfig) deleteVolumeHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[2]

	volume, err := volumes.DeleteVolume(app.Conn, volumeID)
	if err != nil {
		app.errorJSON(w, err)
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
func (app *AppConfig) getCloudSnapsHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[2]

	cloudSnaps, err := snapshots.GetCloudSnaps(app.Conn, volumeID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JsonCloudSnapList{
		Error:         false,
		CloudSnapList: cloudSnaps,
	}

	writeJSON(w, http.StatusOK, resp)

}

func (app *AppConfig) getInspectAWSCloudCredentialHTTP(w http.ResponseWriter, r *http.Request) {

	cloudCredentialID := r.Header.Get("Cloud-Credential-ID")

	cloudCred, err := snapshots.AWSInspectS3CloudCredential(app.Conn, cloudCredentialID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JsonCredentialInspect{
		Error:             false,
		CredentialInspect: *cloudCred,
	}

	writeJSON(w, http.StatusOK, resp)

}
