package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

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

// PostCreateVolume handles the POST request to /create-volume
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

	http.Redirect(w, r, result, http.StatusTemporaryRedirect)

}

// createNewVolume sends a POST request to the broker to create a new volume
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

	method = "GET"

	return volumeID, nil

}
