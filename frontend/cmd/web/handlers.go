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

	url := BrokerURL + "/broker/getallvolumesinfo"
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

	url := BrokerURL + "/broker/getvolumeinfo/" + volumeID
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

// SnapInfofromID retrieves from the broker /getspecificcloudsnapshot and sends it back as struct JsonSpecificCloudSnapResponse
func SnapInfofromID(snapID string, credentialID string) (jsonSnapInfo JsonSpecificCloudSnapResponse, errorFound error) {

	url := BrokerURL + "/broker/getspecificcloudsnapshot"

	method := "GET"

	client := &http.Client{}
	req, errorFound := http.NewRequest(method, url, nil)

	if errorFound != nil {
		log.Println(errorFound)
		return
	}
	req.Header.Add("Cloud-Snap-ID", snapID)

	req.Header.Add("Cloud-Credential-ID", credentialID)

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
	log.Println(string(body))

	json.Unmarshal(body, &jsonSnapInfo)

	return jsonSnapInfo, nil

}

// InspectVolume retrieves from the broker /inspectvolume/{volumeName} and sends it back as struct JsonVolumeInspect
func InspectVolume(volumeName string) JsonVolumeInspect {

	var jsonVolume JsonVolumeInspect

	url := BrokerURL + "/broker/getinspectvolume/" + volumeName

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

	url := BrokerURL + "/broker/getvolumeusage/" + volumeName

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

	url := BrokerURL + "/broker/getallnodesinfo"
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

	url := BrokerURL + "/broker/getlistofnodes"
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

	url := BrokerURL + "/broker/getnodeinfo/" + nodeID
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

	url = BrokerURL + "/broker/getreplicaspernode/" + nodeID

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

	url := BrokerURL + "/broker/postcreatevolume"

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
	url := BrokerURL + "/broker/getpxcluster"
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
	url = BrokerURL + "/broker/getpxclustercapacity"
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

	url := BrokerURL + "/broker/deletevolume/" + volumeID
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

	url := BrokerURL + "/broker/patchvolumesize/" + volumeID
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

	url := BrokerURL + "/broker/patchvolumehalevel/" + volumeID
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

	url := BrokerURL + "/broker/patchvolumeioprofile/" + volumeID
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

func GetAllSnapshotsInfo() (AllSnaps JsonAllCloudSnapResponse, errorFound error) {

	url := BrokerURL + "/broker/getallcloudsnaps"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Println(err)
		return
	}
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

	json.Unmarshal(body, &AllSnaps)

	//log.Printf("AllSnaps: %v", AllSnaps)

	return AllSnaps, nil

}

// createNewCredential sends a POST request to the broker to create a new volume
func createNewCredential(createCloudCredential CreateCloudCredentials) (credentialID string, errorFound error) {

	url := BrokerURL + "/broker/postcreateawscloudcreds"

	method := "POST"

	cloudCredResponse := CreateCloudCredentialsResponse{}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Add("Cloud-Credential-Name", createCloudCredential.CloudCredentialName)
	req.Header.Add("Cloud-Credential-Access-Key", createCloudCredential.CloudCredentialAccessKey)
	req.Header.Add("Cloud-Credential-Secret-Key", createCloudCredential.CloudCredentialSecretKey)
	req.Header.Add("Cloud-Credential-Region", createCloudCredential.CloudCredentialRegion)
	req.Header.Add("Cloud-Credential-Endpoint", createCloudCredential.CloudCredentialEndpoint)
	req.Header.Add("Cloud-Credential-Bucket-Name", createCloudCredential.CloudCredentialBucketName)
	req.Header.Add("Cloud-Credential-Disable-SSL", strconv.FormatBool(createCloudCredential.CloudCredentialDisableSSL))
	req.Header.Add("Cloud-Credential-IAM-Policy-Enabled", strconv.FormatBool(createCloudCredential.CloudCredentialIAMPolicyEnabled))

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

	json.Unmarshal(body, &cloudCredResponse)

	credentialID = cloudCredResponse.CredentialInspect.CredentialId

	method = "GET"

	return credentialID, nil

}

// GetCloudCredentials sends a GET request to the broker to get all cloud credentials
func GetCloudCredentials() (cloudCredsListMap map[string]any, errorFound error) {

	url := BrokerURL + "/broker/getallcloudcredsids"
	method := "GET"

	allCloudCredsIDsResponse := AllCloudCredsIDsResponse{}

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

	json.Unmarshal(body, &allCloudCredsIDsResponse)

	cloudCredsListMap = make(map[string]any)

	for _, cloudCredID := range allCloudCredsIDsResponse.CloudCredsList {

		url = BrokerURL + "/broker/getinspectawscloudcreds"

		method = "GET"

		client = &http.Client{}
		req, errorFound = http.NewRequest(method, url, nil)

		req.Header.Add("Cloud-Credential-ID", cloudCredID)

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
		if errorFound != nil {
			log.Println(errorFound)
			return
		}

		log.Println(string(body))

		var cloudCredInspect CreateCloudCredentialsResponse

		json.Unmarshal(body, &cloudCredInspect)

		cloudCredsListMap[cloudCredID] = cloudCredInspect

		log.Printf("cloudCredsListMap: %v", cloudCredsListMap)

	}

	return cloudCredsListMap, nil

}

// createCloudSnap sends a POST request to the broker to create a cloud snapshot
func createCloudSnapshot(createCloudSnap CreateCloudSnap) (taskID string, errorFound error) {

	url := BrokerURL + "/broker/postcreatecloudsnap"
	method := "POST"

	var createCloudSnapResponse CreateCloudSnapResponse
	client := &http.Client{}
	req, errorFound := http.NewRequest(method, url, nil)

	if errorFound != nil {
		log.Println(errorFound)
		return
	}
	req.Header.Add("Volume-ID", createCloudSnap.VolumeID)
	req.Header.Add("Cloud-Credential-ID", createCloudSnap.CloudCredentialID)

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

	json.Unmarshal(body, &createCloudSnapResponse)

	taskID = createCloudSnapResponse.TaskID

	return taskID, nil
}

// getClusterAlarms sends a GET request to the broker to get all cluster alarms
func getClusterAlarms() (clusterAlarms ClusterAlarms, errorFound error) {

	url := BrokerURL + "/broker/getpxclusteralarms"
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

	json.Unmarshal(body, &clusterAlarms)

	return clusterAlarms, nil

}

// deleteCloudSnapshot sends a DELETE request to the broker to delete a cloud snapshot
func deleteCloudSnapshot(cloudSnapID string) (errorFound error) {
	url := BrokerURL + "/broker/deletecloudsnap"
	method := "DELETE"

	client := &http.Client{}
	req, errorFound := http.NewRequest(method, url, nil)
	if errorFound != nil {
		log.Println(errorFound)
		return
	}
	req.Header.Add("Cloud-Snap-ID", cloudSnapID)

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
	log.Println(string(body))

	return nil
}

// postLogin sends a POST request to the broker to login the backend broker API
func postLogin(Username string, Password string) (errorFound error) {

	url := BrokerURL + "/login"
	method := "POST"

	client := &http.Client{}
	req, errorFound := http.NewRequest(method, url, nil)

	if errorFound != nil {
		log.Println(errorFound)
		return
	}
	req.Header.Add("Username", Username)
	req.Header.Add("Password", Password)

	res, errorFound := client.Do(req)
	if errorFound != nil {
		log.Println(errorFound)
		return
	}
	defer res.Body.Close()

	return nil
}
