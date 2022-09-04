package main

import (
	"fmt"
	"net/http"
	"strings"
)

var Repo *Repository

type Repository struct {
	App *AppConfig
}

func NewRepo(a *AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewTestRepo(a *AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Cluster(w http.ResponseWriter, r *http.Request) {
	Template(w, r, "index.html", &TemplateData{})
}

func (m *Repository) Documentation(w http.ResponseWriter, r *http.Request) {
	Template(w, r, "documentation.html", &TemplateData{})
}

// Volumes serves the volumes page
func (m *Repository) Volumes(w http.ResponseWriter, r *http.Request) {

	volumesInfo, err := GetAllVolumesInfo()
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Printf("volumesInfo: %v", volumesInfo.AllVolumesInfo[0].VolumeName)
	//You have to range over the array to get the values

	Template(w, r, "volumes.html", &TemplateData{
		JsonAllVolumesInfo: volumesInfo,
	})
}

// VolumeInformation serves the volume information page
func (m *Repository) VolumeInformation(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[3]

	volumeInfoResponse, err := VolumeInfofromID(volumeID)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Printf("volumeInfo: %v", volumeInfoResponse.VolumeInfo.VolumeSizeMB)

	Template(w, r, "volume-specific.html", &TemplateData{
		JsonVolumeInfo: volumeInfoResponse,
	})

}

// Nodes serves the nodes page
func (m *Repository) Nodes(w http.ResponseWriter, r *http.Request) {

	nodesInfo, err := GetAllNodesInfo()
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Printf("nodesInfo: %v", nodesInfo.AllNodesInfo[0].NodeName)
	//You have to range over the array to get the values

	Template(w, r, "nodes.html", &TemplateData{
		JsonAllNodesInfo: nodesInfo,
	})

}

// NodeInformation serves the node information page
func (m *Repository) NodeInformation(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	nodeID := exploded[3]

	nodeInfoResponse, replicaPerNodeResponse, err := NodeInfoFromID(nodeID)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Printf("nodeInfoResponse: %v", nodeInfoResponse.NodeInfo.NodeName)
	//fmt.Printf("replicaPerNodeResponse: %v", replicaPerNodeResponse.VolumeList["1014695385474634270"].VolumeName)

	Template(w, r, "node-specific.html", &TemplateData{
		JsonNodeInfo:       nodeInfoResponse,
		JsonReplicaPerNode: replicaPerNodeResponse,
	})
}

func (m *Repository) GetCreateVolume(w http.ResponseWriter, r *http.Request) {

	res, _ := m.App.Session.Get(r.Context(), "create-volume").(CreateVolume)

	m.App.Session.Put(r.Context(), "create-volume", res)

	data := make(map[string]any)

	data["create-volume"] = res

	Template(w, r, "create-volume.html", &TemplateData{
		Form: New(nil),
		Data: data,
	})

}

func (m *Repository) Snaps(w http.ResponseWriter, r *http.Request) {
	Template(w, r, "snapshots.html", &TemplateData{})
}

func (m *Repository) SnapsInformation(w http.ResponseWriter, r *http.Request) {
	Template(w, r, "snap-specific.html", &TemplateData{})
}

func (m *Repository) StoragePools(w http.ResponseWriter, r *http.Request) {
	Template(w, r, "storage-pools.html", &TemplateData{})
}

func (m *Repository) StoragePoolsInformation(w http.ResponseWriter, r *http.Request) {
	Template(w, r, "storage-pool-specific.html", &TemplateData{})
}
