package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
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
		log.Println(errorFound)
	}
	res, errorFound := client.Do(req)
	if errorFound != nil {
		log.Println(errorFound)
	}
	defer res.Body.Close()

	body, errorFound := ioutil.ReadAll(res.Body)
	if errorFound != nil {
		log.Println(errorFound)
	}

	json.Unmarshal(body, &JsonAllVolumesInfo)
	if errorFound != nil {
		log.Println(errorFound)
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
		log.Println(errorFound)
	}
	res, errorFound := client.Do(req)
	if errorFound != nil {
		log.Println(errorFound)
	}
	defer res.Body.Close()

	body, errorFound := ioutil.ReadAll(res.Body)
	if errorFound != nil {
		log.Println(errorFound)
	}

	json.Unmarshal(body, &jsonVolumeInfo)

	return jsonVolumeInfo, nil

}

func SnapInfofromID(volumeID, snapID string) (jsonSnapInfo SnapInfoResponse, errorFound error) {

	return jsonSnapInfo, nil

}

// InspectVolume retrieves from the broker /inspectvolume/{volumeName} and sends it back as struct JsonVolumeInspect
func InspectVolume(volumeName string) JsonVolumeInspect {

	var jsonVolume JsonVolumeInspect

	url := brokerURL + "/getinspectvolume/" + volumeName

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Println(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
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
		log.Println(errorFound)
	}
	res, errorFound := client.Do(req)
	if errorFound != nil {
		log.Println(errorFound)
	}
	defer res.Body.Close()

	body, errorFound := ioutil.ReadAll(res.Body)
	if errorFound != nil {
		log.Println(errorFound)
	}

	json.Unmarshal(body, &JsonAllNodesInfo)
	if errorFound != nil {
		log.Println(errorFound)
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
		log.Println(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
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
		log.Println(errorFound)
	}
	res, errorFound := client.Do(req)
	if errorFound != nil {
		log.Println(errorFound)
	}
	defer res.Body.Close()

	body, errorFound := ioutil.ReadAll(res.Body)
	if errorFound != nil {
		log.Println(errorFound)
	}

	json.Unmarshal(body, &jsonNodeInfo)
	if errorFound != nil {
		log.Println(errorFound)
	}

	url = brokerURL + "/getreplicaspernode/" + nodeID

	client = &http.Client{}
	req, errorFound = http.NewRequest(method, url, nil)

	if errorFound != nil {
		log.Println(errorFound)
	}
	res, errorFound = client.Do(req)
	if errorFound != nil {
		log.Println(errorFound)
	}
	defer res.Body.Close()

	body, errorFound = ioutil.ReadAll(res.Body)
	if errorFound != nil {
		log.Println(errorFound)
	}

	json.Unmarshal(body, &jsonReplicaPerNode)
	if errorFound != nil {
		log.Println(errorFound)
	}

	return jsonNodeInfo, jsonReplicaPerNode, nil

}

// createNewVolume sends a POST request to the broker to create a new volume
func createNewVolume(createVolume CreateVolume) (volumeID string, errorFound error) {

	url := brokerURL + "/postcreatevolume"

	method := "POST"

	volResponse := CreateVolumeResponse{}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(body))

	json.Unmarshal(body, &volResponse)

	volumeID = volResponse.VolumeID

	method = "GET"

	return volumeID, nil

}

// GetClusterInfo handles the GET request to /getpxcluster and /getpxclustercapacity to get the cluster info
func GetClusterInfo() (jsonClusterInfo ClusterInfo, jsonClusterCapacity ClusterCapacity, errorFound error) {

	//get cluster info
	url := brokerURL + "/getpxcluster"
	method := "GET"

	client := &http.Client{}
	req, errorFound := http.NewRequest(method, url, nil)

	if errorFound != nil {
		log.Println(errorFound)
		return
	}
	res, errorFound := client.Do(req)
	if errorFound != nil {
		log.Println(errorFound)
		return
	}
	defer res.Body.Close()

	body, errorFound := ioutil.ReadAll(res.Body)
	if errorFound != nil {
		log.Println(errorFound)
		return
	}

	json.Unmarshal(body, &jsonClusterInfo)

	// get the cluster capacity
	url = brokerURL + "/getpxclustercapacity"
	method = "GET"

	req, err := http.NewRequest(method, url, nil)

	if errorFound != nil {
		log.Println(errorFound)
		return
	}
	res, errorFound = client.Do(req)
	if errorFound != nil {
		log.Println(errorFound)
		return
	}
	defer res.Body.Close()

	body, errorFound = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(errorFound)
		return
	}

	json.Unmarshal(body, &jsonClusterCapacity)

	return jsonClusterInfo, jsonClusterCapacity, nil

}

func DeleteVolume(volumeID string) (string, error) {

	url := brokerURL + "/deletevolume/" + volumeID
	method := "DELETE"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Println(err)
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}

	json.Unmarshal(body, &volumeID)

	message := "Volume " + volumeID + " deleted!"

	return message, nil
}

func ResizeVolume(volumeID string, volSize string) (string, error) {

	url := brokerURL + "/patchvolumesize/" + volumeID
	method := "PATCH"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Println(err)
		return "", err
	}

	req.Header.Add("Volume-Size", volSize)

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	log.Println(string(body))

	message := "Volume " + volumeID + " resized to " + volSize + " GB!"

	return message, nil
}

func UpdateVolumeHALevel(volumeID string, volHALevel string) (string, error) {

	url := brokerURL + "/patchvolumehalevel/" + volumeID
	method := "PATCH"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Println(err)
		return "", err
	}

	req.Header.Add("Volume-Ha-Level", volHALevel)

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	log.Println(string(body))

	message := "Volume " + volumeID + " HA level updated to " + volHALevel + "!"

	return message, nil
}

func IOProfileVolume(volumeID string, volIOProfile string) (string, error) {

	url := brokerURL + "/patchvolumeioprofile/" + volumeID
	method := "PATCH"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Println(err)
		return "", err
	}

	req.Header.Add("Volume-IO-Profile", volIOProfile)

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	log.Println(string(body))

	message := "Volume " + volumeID + " IO profile updated to " + volIOProfile + "!"

	return message, nil
}
