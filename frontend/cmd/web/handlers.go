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

func (m *Repository) Volume(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeName := exploded[3]

	volumeInspect := InspectVolume(volumeName)

	volumeUsage := UsageVolume(volumeName)

	Template(w, r, "index.html", &TemplateData{
		JsonVolumeInspect: volumeInspect,
		JsonUsageVolume:   volumeUsage,
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
