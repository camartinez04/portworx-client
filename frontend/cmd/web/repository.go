package main

import (
	"log"
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

// Cluster serves the cluster page
func (m *Repository) Cluster(w http.ResponseWriter, r *http.Request) {

	clusterInfo, clusterCapacity, err := GetClusterInfo()
	if err != nil {
		log.Println(err)
	}

	Template(w, r, "index.page.html", &TemplateData{
		JsonClusterInfo:     clusterInfo,
		JsonClusterCapacity: clusterCapacity,
	})
}

// Documentation serves the documentation page
func (m *Repository) Documentation(w http.ResponseWriter, r *http.Request) {
	Template(w, r, "documentation.page.html", &TemplateData{})
}

// Volumes serves the volumes page
func (m *Repository) Volumes(w http.ResponseWriter, r *http.Request) {

	volumesInfo, err := GetAllVolumesInfo()
	if err != nil {
		log.Println(err)
	}

	//fmt.Printf("volumesInfo: %v", volumesInfo.AllVolumesInfo[0].VolumeName)
	//You have to range over the array to get the values

	Template(w, r, "volumes.page.html", &TemplateData{
		JsonAllVolumesInfo: volumesInfo,
	})
}

// VolumeInformation serves the volume information page
func (m *Repository) VolumeInformation(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[3]

	volumeInfoResponse, err := VolumeInfofromID(volumeID)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error trying to create the volume")
	}

	// Retrieve from the context if a new volume was created
	res, _ := m.App.Session.Get(r.Context(), "create-volume").(CreateVolume)

	// Get the volume name from the retrieved query above
	volName := res.VolumeName

	// If the volume name is not empty, then we have a new volume
	if volName != "" {

		//This means that the volume was created, therefore we will show the success message
		m.App.Session.Put(r.Context(), "flash", "volume created successfully")

		//remove the context after the volume was created
		m.App.Session.Remove(r.Context(), "create-volume")

	}

	//Server the page with the volume information and the messages that could come from the creation of a new volume
	Template(w, r, "volume-specific.page.html", &TemplateData{
		JsonVolumeInfo: volumeInfoResponse,
		Flash:          m.App.Session.PopString(r.Context(), "flash"),
		Error:          m.App.Session.PopString(r.Context(), "error"),
		Warning:        m.App.Session.PopString(r.Context(), "warning"),
	})

}

// Nodes serves the nodes page
func (m *Repository) Nodes(w http.ResponseWriter, r *http.Request) {

	nodesInfo, err := GetAllNodesInfo()
	if err != nil {
		log.Println(err)
	}

	//fmt.Printf("nodesInfo: %v", nodesInfo.AllNodesInfo[0].NodeName)
	//You have to range over the array to get the values

	Template(w, r, "nodes.page.html", &TemplateData{
		JsonAllNodesInfo: nodesInfo,
	})

}

// NodeInformation serves the node information page
func (m *Repository) NodeInformation(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	nodeID := exploded[3]

	nodeInfoResponse, replicaPerNodeResponse, err := NodeInfoFromID(nodeID)
	if err != nil {
		log.Println(err)
	}

	//fmt.Printf("nodeInfoResponse: %v", nodeInfoResponse.NodeInfo.NodeName)
	//fmt.Printf("replicaPerNodeResponse: %v", replicaPerNodeResponse.VolumeList["1014695385474634270"].VolumeName)

	Template(w, r, "node-specific.page.html", &TemplateData{
		JsonNodeInfo:       nodeInfoResponse,
		JsonReplicaPerNode: replicaPerNodeResponse,
	})
}

// CreateVolume serves the create volume page
func (m *Repository) CreateVolume(w http.ResponseWriter, r *http.Request) {

	res, _ := m.App.Session.Get(r.Context(), "create-volume").(CreateVolume)

	m.App.Session.Put(r.Context(), "create-volume", res)

	data := make(map[string]any)

	data["create-volume"] = res

	Template(w, r, "create-volume.page.html", &TemplateData{
		Form: New(nil),
		Data: data,
	})

}

func (m *Repository) Snaps(w http.ResponseWriter, r *http.Request) {
	Template(w, r, "snapshots.page.html", &TemplateData{})
}

func (m *Repository) SnapsInformation(w http.ResponseWriter, r *http.Request) {
	Template(w, r, "snap-specific.page.html", &TemplateData{})
}

func (m *Repository) StoragePools(w http.ResponseWriter, r *http.Request) {
	Template(w, r, "storage-pools.page.html", &TemplateData{})
}

func (m *Repository) StoragePoolsInformation(w http.ResponseWriter, r *http.Request) {
	Template(w, r, "storage-pool-specific.page.html", &TemplateData{})
}

// DeleteVolume serves the delete volume page
func (m *Repository) DeleteVolume(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[3]

	message, err := DeleteVolume(volumeID)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error trying to delete the volume")
	}

	volumesInfo, err := GetAllVolumesInfo()
	if err != nil {
		log.Println(err)
	}

	Template(w, r, "volumes-after-delete.page.html", &TemplateData{
		Flash:              message,
		Error:              m.App.Session.PopString(r.Context(), "error"),
		JsonAllVolumesInfo: volumesInfo,
	})

	//http.Redirect(w, r, "/frontend/volumes", http.StatusSeeOther)

}
