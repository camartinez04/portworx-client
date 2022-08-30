package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	api "github.com/libopenstorage/openstorage-sdk-clients/sdk/golang"
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

	volumesInfo := GetAllVolumesInfo()

	fmt.Printf("%+v", volumesInfo)

	Template(w, r, "volumes.html", &TemplateData{
		JsonGetAllVolumesInfo: volumesInfo,
	})
}

func GetAllVolumesInfo() (volumesInfo map[string][]any) {

	var allVolumesList map[string]map[string]*api.SdkVolumeInspectResponse

	url := brokerURL + "/getallvolumescomplete"

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

	json.Unmarshal(body, &allVolumesList)

	volumesInfo = make(map[string][]any)

	for _, volumes := range allVolumesList {

		for volumeID, volumeContent := range volumes {

			volumeName := volumeContent.GetName()
			volumeReplicas := len(volumeContent.Volume.ReplicaSets[0].GetNodes())
			volumeStatus := volumeContent.Volume.GetStatus().String()
			volumeAttachedOn := volumeContent.Volume.GetAttachedOn()
			volumeDevicePath := volumeContent.Volume.GetDevicePath()
			volumeTotalSize := volumeContent.Volume.GetSpec().GetSize()
			volumeUsage := volumeContent.Volume.GetUsage()
			volumeAvailableSpace := volumeTotalSize - volumeUsage
			volumePercentageUsed := float64(volumeUsage) / float64(volumeTotalSize) * 100
			volumePercentageUsedInt := int(volumePercentageUsed)

			volumesInfo[volumeID] = []any{volumeName, volumeReplicas, volumeStatus, volumeAttachedOn, volumeDevicePath, volumeTotalSize, volumeUsage, volumeAvailableSpace, volumePercentageUsed, volumePercentageUsedInt}

		}
	}

	return volumesInfo

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
