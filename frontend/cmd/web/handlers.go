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
	Template(w, r, "volumes.html", &TemplateData{})
}

func (m *Repository) VolumeInfo(w http.ResponseWriter, r *http.Request) {

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

	nodeList := ListOfNodes()

	Template(w, r, "nodes.html", &TemplateData{
		JsonListOfNodes: nodeList,
	})

	fmt.Printf("%+v", nodeList)
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

	fmt.Printf("%+v", JsonListOfNodes)

	return JsonListOfNodes

}

func (m *Repository) NodeInfo(w http.ResponseWriter, r *http.Request) {
	Template(w, r, "node-specific.html", &TemplateData{})
}

func (m *Repository) Snaps(w http.ResponseWriter, r *http.Request) {
	Template(w, r, "snapshots.html", &TemplateData{})
}

func (m *Repository) SnapsInfo(w http.ResponseWriter, r *http.Request) {
	Template(w, r, "snap-specific.html", &TemplateData{})
}

func (m *Repository) StoragePools(w http.ResponseWriter, r *http.Request) {
	Template(w, r, "storage-pools.html", &TemplateData{})
}

func (m *Repository) StoragePoolsInfo(w http.ResponseWriter, r *http.Request) {
	Template(w, r, "storage-pool-specific.html", &TemplateData{})
}
