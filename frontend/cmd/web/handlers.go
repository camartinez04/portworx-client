package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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

// GetAllVolumesInfo retrieves from the broker /getallvolumesinfo and sends it back as struct JsonAllVolumesInfo
func GetAllVolumesInfo() (JsonAllVolumesInfo AllVolumesInfoResponse, errorFound error) {

	url := brokerURL + "/getallvolumesinfo"
	method := "GET"

	client := &http.Client{}
	req, errorFound := http.NewRequest(method, url, nil)

	if errorFound != nil {
		fmt.Println(errorFound)
	}
	res, errorFound := client.Do(req)
	if errorFound != nil {
		fmt.Println(errorFound)
	}
	defer res.Body.Close()

	body, errorFound := ioutil.ReadAll(res.Body)
	if errorFound != nil {
		fmt.Println(errorFound)
	}

	json.Unmarshal(body, &JsonAllVolumesInfo)
	if errorFound != nil {
		fmt.Println(errorFound)
	}

	return JsonAllVolumesInfo, nil

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

// VolumeInfofromID retrieves from the broker /getvolumeinfo/{volumeID} and sends it back as struct VolumeInfo
func VolumeInfofromID(volumeID string) (jsonVolumeInfo VolumeInfoResponse, errorFound error) {

	url := brokerURL + "/getvolumeinfo/" + volumeID
	method := "GET"

	client := &http.Client{}
	req, errorFound := http.NewRequest(method, url, nil)

	if errorFound != nil {
		fmt.Println(errorFound)
	}
	res, errorFound := client.Do(req)
	if errorFound != nil {
		fmt.Println(errorFound)
	}
	defer res.Body.Close()

	body, errorFound := ioutil.ReadAll(res.Body)
	if errorFound != nil {
		fmt.Println(errorFound)
	}

	json.Unmarshal(body, &jsonVolumeInfo)

	return jsonVolumeInfo, nil

}

// InspectVolume retrieves from the broker /inspectvolume/{volumeName} and sends it back as struct JsonVolumeInspect
func InspectVolume(volumeName string) JsonVolumeInspect {

	var jsonVolume JsonVolumeInspect

	url := brokerURL + "/getinspectvolume/" + volumeName

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(body, &jsonVolume)

	return jsonVolume

}

// UsageVolume retrieves from the broker /usagevolume/{volumeName} and sends it back as struct JsonUsageVolume
func UsageVolume(volumeName string) JsonUsageVolume {

	var jsonUsage JsonUsageVolume

	url := brokerURL + "/getvolumeusage/" + volumeName

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(body, &jsonUsage)

	return jsonUsage

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

// GetAllNodesInfo retrieves from the broker /getallnodesinfo and sends it back as struct AllNodesInfoResponse
func GetAllNodesInfo() (JsonAllNodesInfo AllNodesInfoResponse, errorFound error) {

	url := brokerURL + "/getallnodesinfo"
	method := "GET"

	client := &http.Client{}
	req, errorFound := http.NewRequest(method, url, nil)

	if errorFound != nil {
		fmt.Println(errorFound)
	}
	res, errorFound := client.Do(req)
	if errorFound != nil {
		fmt.Println(errorFound)
	}
	defer res.Body.Close()

	body, errorFound := ioutil.ReadAll(res.Body)
	if errorFound != nil {
		fmt.Println(errorFound)
	}

	json.Unmarshal(body, &JsonAllNodesInfo)
	if errorFound != nil {
		fmt.Println(errorFound)
	}

	return JsonAllNodesInfo, nil

}

// ListOfNodes retrieves from the broker /getlistofnodes and sends it back as JsonListOfNodes of any
func ListOfNodes() (JsonListOfNodes any) {

	url := brokerURL + "/getlistofnodes"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(body, &JsonListOfNodes)

	return JsonListOfNodes

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

// NodeInfoFromID retrieves from the broker /getvolumeinfo/{volumeID} and sends it back as struct VolumeInfo
func NodeInfoFromID(nodeID string) (jsonNodeInfo NodeInfoResponse, jsonReplicaPerNode ReplicasPerNodeResponse, errorFound error) {

	url := brokerURL + "/getnodeinfo/" + nodeID
	method := "GET"

	client := &http.Client{}
	req, errorFound := http.NewRequest(method, url, nil)

	if errorFound != nil {
		fmt.Println(errorFound)
	}
	res, errorFound := client.Do(req)
	if errorFound != nil {
		fmt.Println(errorFound)
	}
	defer res.Body.Close()

	body, errorFound := ioutil.ReadAll(res.Body)
	if errorFound != nil {
		fmt.Println(errorFound)
	}

	json.Unmarshal(body, &jsonNodeInfo)

	url = brokerURL + "/getreplicaspernode/" + nodeID

	client = &http.Client{}
	req, errorFound = http.NewRequest(method, url, nil)

	if errorFound != nil {
		fmt.Println(errorFound)
	}
	res, errorFound = client.Do(req)
	if errorFound != nil {
		fmt.Println(errorFound)
	}
	defer res.Body.Close()

	body, errorFound = ioutil.ReadAll(res.Body)
	if errorFound != nil {
		fmt.Println(errorFound)
	}

	json.Unmarshal(body, &jsonReplicaPerNode)

	return jsonNodeInfo, jsonReplicaPerNode, nil

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

func (m *Repository) PostCreateVolume(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		fmt.Println("can't parse form!")
		http.Redirect(w, r, "/frontend/cluster", http.StatusSeeOther)
		return
	}

	volumeSize, err := strconv.ParseUint(r.Form.Get("volume_size"), 10, 64)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		fmt.Println("invalid volume size!")
		http.Redirect(w, r, "/frontend/cluster", http.StatusSeeOther)
		return
	}

	volumeHALevel, err := strconv.ParseInt(r.Form.Get("volume_ha_level"), 10, 64)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		fmt.Println("invalid volume ha level!")
		http.Redirect(w, r, "/frontend/cluster", http.StatusSeeOther)
		return
	}

	volumeEncrypted, err := strconv.ParseBool(r.Form.Get("volume_encrypted"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		fmt.Println("invalid volume encrypted value!")
		http.Redirect(w, r, "/frontend/cluster", http.StatusSeeOther)
		return
	}

	volumeSharedv4, err := strconv.ParseBool(r.Form.Get("volume_sharedv4"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		fmt.Println("invalid volume sharedv4 value!")
		http.Redirect(w, r, "/frontend/cluster", http.StatusSeeOther)
		return
	}

	volumeNodiscard, err := strconv.ParseBool(r.Form.Get("volume_no_discard"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		fmt.Println("invalid volume nodiscard value!")
		http.Redirect(w, r, "/frontend/cluster", http.StatusSeeOther)
		return
	}

	createVolume := CreateVolume{
		VolumeName:      r.FormValue("volume_name"),
		VolumeSize:      volumeSize,
		VolumeIOProfile: r.FormValue("volume_io_profile"),
		VolumeHALevel:   volumeHALevel,
		VolumeEncrypted: volumeEncrypted,
		VolumeSharedv4:  volumeSharedv4,
		VolumeNoDiscard: volumeNodiscard,
	}

	fmt.Println("successfully created the struct createVolume!")

	fmt.Printf("Post to send: %v", createVolume)

	m.App.Session.Put(r.Context(), "create-volume", createVolume)

	volumeIDResp, err := createNewVolume(createVolume)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't create volume!")
		fmt.Println("can't create volume!")
		http.Redirect(w, r, "/frontend/cluster", http.StatusSeeOther)
		return
	}

	fmt.Println("successfully created the volume!")

	result := "/frontend/volume/" + volumeIDResp

	http.Redirect(w, r, result, http.StatusAccepted)

}

func createNewVolume(createVolume CreateVolume) (volumeID string, errorFound error) {

	url := brokerURL + "/postcreatevolume"

	method := "POST"

	volResponse := CreateVolumeResponse{}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Volume-Name", createVolume.VolumeName)
	req.Header.Add("Volume-Size", strconv.FormatUint(createVolume.VolumeSize, 10))
	req.Header.Add("Volume-Ha-Level", strconv.FormatInt(createVolume.VolumeHALevel, 10))
	req.Header.Add("Volume-Encryption-Enabled", strconv.FormatBool(createVolume.VolumeEncrypted))
	req.Header.Add("Volume-Sharedv4-Enabled", strconv.FormatBool(createVolume.VolumeSharedv4))
	req.Header.Add("Volume-No-Discard", strconv.FormatBool(createVolume.VolumeNoDiscard))
	req.Header.Add("Volume-IO-Profile", createVolume.VolumeIOProfile)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

	json.Unmarshal(body, &volResponse)

	volumeID = volResponse.VolumeID

	return volumeID, nil

}
