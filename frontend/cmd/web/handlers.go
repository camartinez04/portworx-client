package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func (m *Repository) VolumeInformation(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeName := exploded[3]

	volumeInspect := InspectVolume(volumeName)

	io_profile := volumeInspect.IoProfileString

	status := volumeInspect.VolumeStatusString

	volumeUsage := UsageVolume(volumeName)

	Template(w, r, "volume-specific.html", &TemplateData{
		JsonVolumeInspect:  volumeInspect,
		JsonUsageVolume:    volumeUsage,
		IoProfileString:    io_profile,
		VolumeStatusString: status,
	})

}

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

func (m *Repository) NodeInformation(w http.ResponseWriter, r *http.Request) {
	Template(w, r, "node-specific.html", &TemplateData{})
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
