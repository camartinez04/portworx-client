package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/camartinez04/portworx-client/broker/pkg/cluster"
	"github.com/camartinez04/portworx-client/broker/pkg/nodes"
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

	volume, replicas, err := volumes.InspectVolume(app.Conn, volumeName)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JsonResponse{
		Error:         false,
		VolumeInspect: volume,
		ReplicasInfo:  replicas,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getClusterCapacityHTTP http function to get the cluster capacity.
func (app *AppConfig) getClusterCapacityHTTP(w http.ResponseWriter, r *http.Request) {

	cluster, _ := cluster.ClusterCapacity(app.Conn)

	resp := JsonResponse{
		Error:           false,
		ClusterCapacity: cluster,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getClusterCapacityHTTP http function to get the cluster capacity.
func (app *AppConfig) getClusterUUIDHTTP(w http.ResponseWriter, r *http.Request) {

	uuid, _ := cluster.ClusterInfo(app.Conn)

	resp := JsonResponse{
		Error:       false,
		ClusterUUID: uuid,
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

	err = volumes.CreateVolume(app.Conn, volumeName, volumeGBSize, volumeHALevel, encryptionEnabled, sharedv4Enabled, noDiscard)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JsonResponse{
		Error:   false,
		Message: "Volume created successfully",
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

	resp := JsonResponse{
		Error:         false,
		NodesOfVolume: nodes,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

// getListOfNodesHTTP http function to get the list of nodes of the Portworx cluster.
func (app *AppConfig) getListOfNodesHTTP(w http.ResponseWriter, r *http.Request) {

	nodes, err := nodes.GetListOfNodes(app.Conn)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JsonResponse{
		Error:    false,
		NodeList: nodes,
	}

	writeJSON(w, http.StatusAccepted, resp)

}

func (app *AppConfig) getReplicasPerNodeHTTP(w http.ResponseWriter, r *http.Request) {

	nodeName := r.Header.Get("Node-Name")

	volumes, err := nodes.GetReplicasPerNode(app.Conn, nodeName)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JsonResponse{
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

	resp := JsonResponse{
		Error:          false,
		AllVolumesList: volumes,
	}

	writeJSON(w, http.StatusAccepted, resp)

}
