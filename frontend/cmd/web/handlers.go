package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	volume := "pvc-52574539-e72f-452f-b355-caa63e41cd9d"

	volume_inspect := InspectVolume(volume)

	Template(w, r, "index.html", &TemplateData{
		JsonVolumeInspect: volume_inspect,
	})
}

func InspectVolume(volume string) JsonVolumeInspect {

	var jsonVolume JsonVolumeInspect

	url := "http://localhost:8080/getinspectvolume/" + volume

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
