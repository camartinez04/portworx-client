package main

import (
	"context"
	"encoding/json"
	"log"
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

	r.Header.Add("Authorization", "Bearer "+m.App.Session.GetString(r.Context(), "token"))

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

	volumeID := r.Header.Get("Volume-ID")

	if volumeID == "" {

		volumeID = exploded[3]
	}

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

	data := make(map[string]any)

	data["create-volume"] = res

	m.App.Session.Put(r.Context(), "create-volume", res)

	Template(w, r, "create-volume.page.html", &TemplateData{
		Form: New(nil),
		Data: data,
	})

}

// PostCreateVolume handles the POST request to /create-volume
func (m *Repository) PostCreateVolume(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		log.Println("can't parse form!")
		http.Redirect(w, r, "/frontend/cluster", http.StatusSeeOther)
		return
	}

	volumeSize, err := strconv.ParseUint(r.Form.Get("volume_size"), 10, 64)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		log.Println("invalid volume size!")
		http.Redirect(w, r, "/frontend/cluster", http.StatusSeeOther)
		return
	}

	volumeHALevel, err := strconv.ParseInt(r.Form.Get("volume_ha_level"), 10, 64)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		log.Println("invalid volume ha level!")
		http.Redirect(w, r, "/frontend/cluster", http.StatusSeeOther)
		return
	}

	volumeEncrypted, err := strconv.ParseBool(r.Form.Get("volume_encrypted"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		log.Println("invalid volume encrypted value!")
		http.Redirect(w, r, "/frontend/cluster", http.StatusSeeOther)
		return
	}

	volumeSharedv4, err := strconv.ParseBool(r.Form.Get("volume_sharedv4"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		log.Println("invalid volume sharedv4 value!")
		http.Redirect(w, r, "/frontend/cluster", http.StatusSeeOther)
		return
	}

	volumeNodiscard, err := strconv.ParseBool(r.Form.Get("volume_no_discard"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		log.Println("invalid volume nodiscard value!")
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

	log.Println("successfully created the struct createVolume!")

	log.Printf("Post to send: %v", createVolume)

	m.App.Session.Put(r.Context(), "create-volume", createVolume)

	volumeIDResp, err := createNewVolume(createVolume)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't create volume!")
		log.Println("can't create volume!")
		http.Redirect(w, r, "/frontend/cluster", http.StatusSeeOther)
		return
	}

	log.Println("successfully created the volume!")

	result := "/frontend/volume/" + volumeIDResp

	http.Redirect(w, r, result, http.StatusSeeOther)

}

// GetAllSnaps serves the snapshots page
func (m *Repository) GetAllSnaps(w http.ResponseWriter, r *http.Request) {

	volumesInfo, err := GetAllVolumesInfo()
	if err != nil {
		log.Println(err)
	}

	jsonAllSnapsInfo, err := GetAllSnapshotsInfo()
	if err != nil {
		log.Println(err)
	}

	Template(w, r, "snapshots.page.html", &TemplateData{
		JsonAllVolumesInfo: volumesInfo,
		JsonAllSnapsInfo:   jsonAllSnapsInfo,
	})
}

// SpecificSpapInformation serves the specific snapshot information page
func (m *Repository) SpecificSpapInformation(w http.ResponseWriter, r *http.Request) {

	//snapshotID := r.Header.Get("Cloud-Snap-ID")

	exploded := strings.Split(r.RequestURI, "/")

	snapshotID := exploded[3] + "/" + exploded[4]

	jsonSnapInfoResponse, err := SnapInfofromID(snapshotID)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error trying to get the snapshot information!")
	}

	Template(w, r, "snap-specific.page.html", &TemplateData{
		JsonSnapSpecific: jsonSnapInfoResponse,
	})
}

// CloudCredentials serves the cloud credentials page
func (m *Repository) CloudCredentials(w http.ResponseWriter, r *http.Request) {

	cloudCredsList, err := GetCloudCredentials()
	if err != nil {
		log.Println(err)
	}

	Template(w, r, "cloud-credentials.page.html", &TemplateData{
		JsonCloudCredsList: cloudCredsList,
	})
}

// CloudCredentialsInformation serves the specific cloud credentials information page
func (m *Repository) CloudCredentialsInformation(w http.ResponseWriter, r *http.Request) {
	Template(w, r, "cloud-credential-specific.page.html", &TemplateData{})
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

// UpdateVolumeHALevelHTTP serves the update volume ha level page
func (m *Repository) UpdateVolumeHALevelHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[3]

	replica := exploded[4]

	message, err := UpdateVolumeHALevel(volumeID, replica)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error trying to update the volume")
	}

	log.Println(message)

	http.Redirect(w, r, "/frontend/volume/"+volumeID, http.StatusSeeOther)

}

// UpdateVolumeSizeHTTP serves the update volume size page
func (m *Repository) UpdateVolumeSizeHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[3]

	newSize := exploded[4]

	message, err := ResizeVolume(volumeID, newSize)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error trying to update the volume")
	}

	log.Println(message)

	http.Redirect(w, r, "/frontend/volume/"+volumeID, http.StatusSeeOther)

}

// UpdateVolumeIOProfileHTTP serves the update volume io profile page
func (m *Repository) UpdateVolumeIOProfileHTTP(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	volumeID := exploded[3]

	ioProfile := exploded[4]

	message, err := IOProfileVolume(volumeID, ioProfile)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error trying to update the volume")
	}

	log.Println(message)

	http.Redirect(w, r, "/frontend/volume/"+volumeID, http.StatusSeeOther)

}

// CreateCloudCredentials serves the create cloud credentials page
func (m *Repository) CreateCloudCredentials(w http.ResponseWriter, r *http.Request) {

	res, _ := m.App.Session.Get(r.Context(), "create-credentials").(CreateCloudCredentials)

	data := make(map[string]any)

	data["create-credentials"] = res

	m.App.Session.Put(r.Context(), "create-credentials", res)

	Template(w, r, "create-credentials.page.html", &TemplateData{
		Form: New(nil),
		Data: data,
	})

}

// PostCreateCloudCredentials handles the POST request to /create-credentials
func (m *Repository) PostCreateCloudCredentials(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		log.Println("can't parse form!")
		http.Redirect(w, r, "/frontend/cluster", http.StatusSeeOther)
		return
	}

	cloudCredIAMPolicyEnabled, err := strconv.ParseBool(r.Form.Get("cloud_credential_iam_policy_enabled"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		log.Println("invalid disable ssl value!")
		http.Redirect(w, r, "/frontend/cluster", http.StatusSeeOther)
		return
	}

	cloudCredDisableSSL, err := strconv.ParseBool(r.Form.Get("cloud_credential_disable_ssl"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		log.Println("invalid disable ssl value!")
		http.Redirect(w, r, "/frontend/cluster", http.StatusSeeOther)
		return
	}

	createCloudCredential := CreateCloudCredentials{
		CloudCredentialName:             r.FormValue("cloud_credential_name"),
		CloudCredentialAccessKey:        r.FormValue("cloud_credential_access_key"),
		CloudCredentialSecretKey:        r.FormValue("cloud_credential_secret_key"),
		CloudCredentialEndpoint:         r.FormValue("cloud_credential_endpoint"),
		CloudCredentialBucketName:       r.FormValue("cloud_credential_bucket_name"),
		CloudCredentialRegion:           r.FormValue("cloud_credential_region"),
		CloudCredentialDisableSSL:       cloudCredDisableSSL,
		CloudCredentialIAMPolicyEnabled: cloudCredIAMPolicyEnabled,
	}

	log.Println("successfully created the struct createCloudCredential!")

	log.Printf("Post to send: %v", createCloudCredential)

	m.App.Session.Put(r.Context(), "create-credentials", createCloudCredential)

	credentialIDResp, err := createNewCredential(createCloudCredential)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't create new credential!")
		log.Println("can't create credential!")
		http.Redirect(w, r, "/frontend/cluster", http.StatusSeeOther)
		return
	}

	log.Println("successfully created the new cloud credential! ID: " + credentialIDResp)

	result := "/frontend/cloud-credentials"

	http.Redirect(w, r, result, http.StatusSeeOther)

}

// CreateCloudSnap serves the create cloud snap page
func (m *Repository) CreateCloudSnap(w http.ResponseWriter, r *http.Request) {

	res, _ := m.App.Session.Get(r.Context(), "create-cloudsnap").(CreateCloudSnap)

	data := make(map[string]any)

	data["create-cloudsnap"] = res

	m.App.Session.Put(r.Context(), "create-cloudsnap", res)

	Template(w, r, "create-cloudsnap.page.html", &TemplateData{
		Form: New(nil),
		Data: data,
	})

}

// PostCreateCloudSnap handles the POST request to /create-cloudsnap
func (m *Repository) PostCreateCloudSnap(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		log.Println("can't parse form!")
		http.Redirect(w, r, "/frontend/cluster", http.StatusSeeOther)
		return
	}

	createCloudSnap := CreateCloudSnap{
		VolumeID:          r.FormValue("volume-id"),
		CloudCredentialID: r.FormValue("cloud-credential-id"),
	}

	log.Println("successfully created the struct createCloudSnap!")

	log.Printf("Post to send: %v", createCloudSnap)

	m.App.Session.Put(r.Context(), "create-cloudsnap", createCloudSnap)

	taskID, err := createCloudSnapshot(createCloudSnap)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't create new cloudsnap!")
		log.Println("can't create cloudsnap!")
		http.Redirect(w, r, "/frontend/cluster", http.StatusSeeOther)
		return
	}

	log.Println("successfully created the new cloudsnap with Task ID: " + taskID)

	result := "/frontend/snapshots"

	http.Redirect(w, r, result, http.StatusSeeOther)

}

func (m *Repository) GetLogin(w http.ResponseWriter, r *http.Request) {

	token := m.App.Session.GetString(r.Context(), "token")

	Template(w, r, "login.page.html", &TemplateData{
		Form:          New(nil),
		KeycloakToken: token,
	})

}

// PostLogin handles the POST request to /login
func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {

	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	username := r.Form.Get("username")

	password := r.Form.Get("password")

	//log.Printf("username: %s", username)

	rq := &loginRequest{username, password}

	jwt, err := m.App.NewKeycloak.gocloak.Login(context.Background(),
		m.App.NewKeycloak.clientId,
		m.App.NewKeycloak.clientSecret,
		m.App.NewKeycloak.realm,
		rq.Username,
		rq.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	rs := &loginResponse{
		AccessToken:  jwt.AccessToken,
		RefreshToken: jwt.RefreshToken,
		ExpiresIn:    jwt.ExpiresIn,
	}

	//log.Printf("jwt: %v", jwt.AccessToken)

	rsJs, _ := json.Marshal(rs)

	_, _ = w.Write(rsJs)

	keycloakToken = jwt.AccessToken

	http.Redirect(w, r, "/frontend/cluster", http.StatusSeeOther)

}

// Logout logs a user out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())

	_ = m.App.Session.RenewToken(r.Context())

	keycloakToken = ""

	http.Redirect(w, r, "/frontend/login", http.StatusSeeOther)
}
