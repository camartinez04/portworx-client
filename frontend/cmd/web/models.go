package main

import (
	api "github.com/libopenstorage/openstorage-sdk-clients/sdk/golang"
)

// TemplateData holds data sent from handlers to template
type TemplateData struct {
	StringMap             map[string]string
	IntMap                map[string]int
	FloatMap              map[string]float32
	Data                  map[string]any
	Form                  *Form
	JsonVolumeInspect     JsonVolumeInspect
	CSRFToken             string
	KeycloakToken         string
	Flash                 string
	Warning               string
	Error                 string
	IsAuthenticated       int
	JsonUsageVolume       JsonUsageVolume
	IoProfileString       string
	VolumeStatusString    string
	JsonListOfNodes       any
	JsonGetAllVolumesInfo map[string][]any
	JsonNodeInfo          NodeInfoResponse
	JsonVolumeInfo        VolumeInfoResponse
	JsonAllVolumesInfo    AllVolumesInfoResponse
	JsonAllNodesInfo      AllNodesInfoResponse
	JsonReplicaPerNode    ReplicasPerNodeResponse
	JsonClusterInfo       ClusterInfo
	JsonClusterCapacity   ClusterCapacity
	JsonSnapInfo          SnapInfoResponse
	JsonAllSnapsInfo      JsonAllCloudSnapResponse
	JsonSnapSpecific      JsonSpecificCloudSnapResponse
	JsonCloudCredsList    map[string]any
}

// SnapInfoResponse holds the response from the snap info API
type SnapInfoResponse struct {
	CloudSnapList []struct{} `json:"cloud_snap_list,omitempty"`
}

// JsonSpecificCloudSnapResponse is the response format for JSON for SpecificCloudSnapList
type JsonSpecificCloudSnapResponse struct {
	Error       bool                    `json:"error,omitempty"`
	CloudSnap   *api.SdkCloudBackupInfo `json:"cloud_snap,omitempty"`
	CloudSnapId string                  `json:"cloud_snap_id,omitempty"`
}

// JsonAllCloudSnapResponse is the response format for JSON for AllCloudSnapList
type JsonAllCloudSnapResponse struct {
	CloudSnapsList map[string]map[string][]*api.SdkCloudBackupInfo `json:"cloud_snaps_list,omitempty"`
}

type VolIdCloud struct {
	VolID         string          `json:"vol_id,omitempty"`
	CloudSnapsIDs []CloudSnapsIDs `json:"cloud_snaps_ids,omitempty"`
}

// CloudSnapsIDs struct to store the cloudsnaps information per CloudID.
type CloudSnapsIDs struct {
	CredID     string                    `json:"cred_id,omitempty"`
	CloudSnaps []*api.SdkCloudBackupInfo `json:"cloud_snaps,omitempty"`
}

// CloudSnapList is a struct for cloud snap list
type CloudSnapList []struct {
	ID            string `json:"id,omitempty"`
	SrcVolumeID   string `json:"src_volume_id,omitempty"`
	SrcVolumeName string `json:"src_volume_name,omitempty"`
	Timestamp     struct {
		Seconds int `json:"seconds,omitempty"`
	} `json:"timestamp,omitempty"`
	Metadata struct {
		CloudsnapType       string `json:"cloudsnapType,omitempty"`
		CompressedSizeBytes string `json:"compressedSizeBytes,omitempty"`
		Compression         string `json:"compression,omitempty"`
		SizeBytes           string `json:"sizeBytes,omitempty"`
		Starttime           string `json:"starttime,omitempty"`
		Status              string `json:"status,omitempty"`
		Updatetime          string `json:"updatetime,omitempty"`
		Version             string `json:"version,omitempty"`
		Volume              string `json:"volume,omitempty"`
		Volumename          string `json:"volumename,omitempty"`
	} `json:"metadata,omitempty"`
	Status int `json:"status,omitempty"`
}

// JsonUsageVolume holds the json data for the usage volume
type JsonUsageVolume struct {
	VolumeUsage            int     `json:"volume_usage,omitempty"`
	AvailableSpace         int     `json:"available_space,omitempty"`
	TotalSize              int     `json:"total_size,omitempty"`
	VolumeUsagePercent     float64 `json:"volume_usage_percent,omitempty"`
	VolumeAvailablePercent float64 `json:"volume_available_percent,omitempty"`
}

// JsonVolumeInspect holds the volume inspect json
type JsonVolumeInspect struct {
	VolumeInspect struct {
		ID     string `json:"id,omitempty"`
		Source struct {
		} `json:"source,omitempty"`
		Locator struct {
			Name         string `json:"name,omitempty"`
			VolumeLabels struct {
				DisableIoProfileProtection string `json:"disable_io_profile_protection,omitempty"`
			} `json:"volume_labels,omitempty"`
		} `json:"locator,omitempty"`
		Ctime struct {
			Seconds int `json:"seconds,omitempty"`
			Nanos   int `json:"nanos,omitempty"`
		} `json:"ctime,omitempty"`
		Spec struct {
			Size         int64 `json:"size,omitempty"`
			Format       int   `json:"format,omitempty"`
			BlockSize    int   `json:"block_size,omitempty"`
			HaLevel      int   `json:"ha_level,omitempty"`
			Cos          int   `json:"cos,omitempty"`
			IoProfile    int   `json:"io_profile,omitempty"`
			VolumeLabels struct {
				DisableIoProfileProtection string `json:"disable_io_profile_protection,omitempty"`
			} `json:"volume_labels,omitempty"`
			ReplicaSet struct {
			} `json:"replica_set,omitempty"`
			AggregationLevel       int  `json:"aggregation_level,omitempty"`
			Scale                  int  `json:"scale,omitempty"`
			QueueDepth             int  `json:"queue_depth,omitempty"`
			ForceUnsupportedFsType bool `json:"force_unsupported_fs_type,omitempty"`
			IoStrategy             struct {
			} `json:"io_strategy,omitempty"`
			Xattr      int `json:"xattr,omitempty"`
			ScanPolicy struct {
			} `json:"scan_policy,omitempty"`
		} `json:"spec,omitempty"`
		Usage       int64  `json:"usage,omitempty"`
		Format      int    `json:"format,omitempty"`
		Status      int    `json:"status,omitempty"`
		State       int    `json:"state,omitempty"`
		AttachedOn  string `json:"attached_on,omitempty"`
		DevicePath  string `json:"device_path,omitempty"`
		ReplicaSets []struct {
			Nodes     []string `json:"nodes,omitempty"`
			PoolUuids []string `json:"pool_uuids,omitempty"`
		} `json:"replica_sets,omitempty"`
		RuntimeState []struct {
			RuntimeState struct {
				ID                  string `json:"ID,omitempty"`
				PXReplReAddNodeMid  string `json:"PXReplReAddNodeMid,omitempty"`
				PXReplReAddPools    string `json:"PXReplReAddPools,omitempty"`
				ReplNodePools       string `json:"ReplNodePools,omitempty"`
				ReplRemoveMids      string `json:"ReplRemoveMids,omitempty"`
				ReplicaSetCreateMid string `json:"ReplicaSetCreateMid,omitempty"`
				ReplicaSetCurr      string `json:"ReplicaSetCurr,omitempty"`
				ReplicaSetCurrMid   string `json:"ReplicaSetCurrMid,omitempty"`
				RuntimeState        string `json:"RuntimeState,omitempty"`
			} `json:"runtime_state,omitempty"`
		} `json:"runtime_state,omitempty"`
		AttachTime struct {
			Seconds int `json:"seconds,omitempty"`
			Nanos   int `json:"nanos,omitempty"`
		} `json:"attach_time,omitempty"`
		DetachTime struct {
			Seconds int `json:"seconds,omitempty"`
			Nanos   int `json:"nanos,omitempty"`
		} `json:"detach_time,omitempty"`
		FpConfig struct {
			SetupOn int  `json:"setup_on,omitempty"`
			Status  int  `json:"status,omitempty"`
			Dirty   bool `json:"dirty,omitempty"`
		} `json:"fpConfig,omitempty"`
		MountOptions struct {
			Options struct {
				Discard string `json:"discard,omitempty"`
			} `json:"options,omitempty"`
		} `json:"mount_options,omitempty"`
		Sharedv4MountOptions struct {
			Options struct {
				Actimeo string `json:"actimeo,omitempty"`
				Proto   string `json:"proto,omitempty"`
				Retrans string `json:"retrans,omitempty"`
				Soft    string `json:"soft,omitempty"`
				Timeo   string `json:"timeo,omitempty"`
				Vers    string `json:"vers,omitempty"`
			} `json:"options,omitempty"`
		} `json:"sharedv4_mount_options,omitempty"`
		PrevState int `json:"prev_state,omitempty"`
	} `json:"volume_inspect,omitempty"`
	ReplicasInfo       []string `json:"replicas_info,omitempty"`
	VolumeNodes        []string `json:"volume_nodes,omitempty"`
	VolumeStatusString string   `json:"volume_status_string,omitempty"`
	IoProfileString    string   `json:"io_profile_string,omitempty"`
}

// Slice of AllNodesInfoResponse
type AllNodesInfoResponse struct {
	AllNodesInfo []NodeInfo `json:"all_nodes_info,omitempty"`
}

// Slice of VolumeInfoResponse
type NodeInfoResponse struct {
	NodeInfo NodeInfo `json:"node_info,omitempty"`
}

// Map of ReplicasPerNodeResponse
type ReplicasPerNodeResponse struct {
	VolumeList map[string]VolumeInfo `json:"volume_list,omitempty"`
}

// NodeInfo struct to store the node information.
type NodeInfo struct {
	NodeName             string             `json:"node_name,omitempty"`
	NodeID               string             `json:"node_id,omitempty"`
	NodeStatus           string             `json:"node_status,omitempty"`
	NodeAvgLoad          int64              `json:"node_avg_load,omitempty"`
	NumberOfPools        int                `json:"number_of_pools,omitempty"`
	SizeNodePool         uint64             `json:"size_node_pool,omitempty"`
	UsedNodePool         uint64             `json:"used_node_pool,omitempty"`
	FreeNodePool         uint64             `json:"free_node_pool,omitempty"`
	NodeMemTotal         uint64             `json:"node_mem_total,omitempty"`
	NodeMemUsed          uint64             `json:"node_mem_used,omitempty"`
	NodeMemFree          uint64             `json:"node_mem_free,omitempty"`
	PercentUsedPool      float64            `json:"percent_used_pool,omitempty"`
	PercentAvailablePool float64            `json:"percent_available_pool,omitempty"`
	PercentUsedMemory    float64            `json:"percent_used_memory,omitempty"`
	StoragelessNode      bool               `json:"storageless,omitempty"`
	StoragePools         []*api.StoragePool `json:"storage_pools,omitempty"`
}

// ClusterInfo struct to store the cluster information.
type ClusterInfo struct {
	ClusterUUID   string `json:"cluster_uuid,omitempty"`
	ClusterStatus string `json:"cluster_status,omitempty"`
	ClusterName   string `json:"cluster_name,omitempty"`
}

// ClusterCapacity struct to store the cluster capacity.
type ClusterCapacity struct {
	ClusterCapacity         uint64  `json:"cluster_capacity,omitempty"`
	ClusterUsed             uint64  `json:"cluster_used,omitempty"`
	ClusterAvailable        uint64  `json:"cluster_available,omitempty"`
	ClusterPercentUsed      float64 `json:"cluster_percent_used,omitempty"`
	ClusterPercentAvailable float64 `json:"cluster_percent_available,omitempty"`
}

// Slice of VolumeInfoAllVolumesInfoResponse
type AllVolumesInfoResponse struct {
	AllVolumesInfo []VolumeInfo `json:"all_volumes_info,omitempty"`
}

// Slice of VolumeInfoResponse
type VolumeInfoResponse struct {
	VolumeInfo VolumeInfo `json:"volume_info,omitempty"`
}

// VolumeInfo struct to store the volume information.
type VolumeInfo struct {
	VolumeName                string                   `json:"volume_name,omitempty"`
	VolumeID                  string                   `json:"volume_id,omitempty"`
	VolumeReplicas            int                      `json:"volume_replicas,omitempty"`
	VolumeReplicaNodes        []string                 `json:"volume_replica_nodes,omitempty"`
	VolumeIOProfile           string                   `json:"volume_io_profile,omitempty"`
	VolumeIOProfileAPI        string                   `json:"volume_io_profile_api,omitempty"`
	VolumeIOPriority          string                   `json:"volume_io_priority,omitempty"`
	VolumeStatus              string                   `json:"volume_status,omitempty"`
	VolumeSizeMB              uint64                   `json:"volume_size_mb,omitempty"`
	VolumeUsedMB              uint64                   `json:"volume_used_mb,omitempty"`
	VolumeAvailable           uint64                   `json:"volume_available,omitempty"`
	VolumeUsedPercent         float64                  `json:"volume_used_percent,omitempty"`
	VolumeAvailablePercent    float64                  `json:"volume_available_percent,omitempty"`
	VolumeType                string                   `json:"volume_type,omitempty"`
	VolumeAttachedPath        []string                 `json:"volume_attached_path,omitempty"`
	VolumeAttachedOn          string                   `json:"volume_attached_on,omitempty"`
	VolumeAttachStatus        string                   `json:"volume_attach_status,omitempty"`
	VolumeDevicePath          string                   `json:"volume_device_path,omitempty"`
	VolumeAggregationLevel    uint32                   `json:"volume_aggregation_level,omitempty"`
	VolumeConsumers           []*api.VolumeConsumer    `json:"volume_consumers,omitempty"`
	VolumeEncrypted           bool                     `json:"volume_encrypted"`
	VolumeEncryptionKey       string                   `json:"volume_encryption_key,omitempty"`
	VolumeK8sNamespace        string                   `json:"volume_k8s_namespace,omitempty"`
	VolumeK8sPVCName          string                   `json:"volume_k8s_pvc_name,omitempty"`
	VolumeSharedv4            bool                     `json:"volume_sharedv4"`
	VolumeSharedv4ServiceSpec *api.Sharedv4ServiceSpec `json:"volume_sharedv4_service_spec,omitempty"`
	VolumeIOStrategy          *api.IoStrategy          `json:"volume_io_strategy,omitempty"`
}

// CreateVolume struct to store the volume creation information.
type CreateVolume struct {
	VolumeName      string `json:"volume_name,omitempty"`
	VolumeSize      uint64 `json:"volume_size,omitempty"`
	VolumeIOProfile string `json:"volume_io_profile,omitempty"`
	VolumeHALevel   int64  `json:"volume_ha_level,omitempty"`
	VolumeEncrypted bool   `json:"volume_encrypted,omitempty"`
	VolumeSharedv4  bool   `json:"volume_sharedv4,omitempty"`
	VolumeNoDiscard bool   `json:"volume_no_discard,omitempty"`
}

// CreateVolumeResponse struct to store the volume creation response.
type CreateVolumeResponse struct {
	Error    bool   `json:"error,omitempty"`
	Message  string `json:"message,omitempty"`
	VolumeID string `json:"volume_id,omitempty"`
}

// CreateCloudCredentials struct to store Cloud Credentials creation form
type CreateCloudCredentials struct {
	CloudCredentialName             string `json:"cloud_credential_name,omitempty"`
	CloudCredentialAccessKey        string `json:"cloud_credential_access_key,omitempty"`
	CloudCredentialSecretKey        string `json:"cloud_credential_secret_key,omitempty"`
	CloudCredentialBucketName       string `json:"cloud_credential_bucket_name,omitempty"`
	CloudCredentialRegion           string `json:"cloud_credential_region,omitempty"`
	CloudCredentialEndpoint         string `json:"cloud_credential_endpoint,omitempty"`
	CloudCredentialDisableSSL       bool   `json:"cloud_credential_disable_ssl,omitempty"`
	CloudCredentialIAMPolicyEnabled bool   `json:"cloud_credential_iam_policy_enabled,omitempty"`
}

// CreateCloudCredentialsResponse struct to store Cloud Credentials creation response.
type CreateCloudCredentialsResponse struct {
	Error                  bool                             `json:"error,omitempty"`
	Message                string                           `json:"message,omitempty"`
	CloudCredentialInspect api.SdkCredentialInspectResponse `json:"credential_inspect,omitempty"`
}

// CreateCloudSnap struct to store Cloud Snap creation form
type CreateCloudSnap struct {
	VolumeID          string `json:"volume_id,omitempty"`
	CloudCredentialID string `json:"cloud_credential_id,omitempty"`
}

// CreateCloudSnapResponse struct to store Cloud Snap creation response.
type CreateCloudSnapResponse struct {
	Error   bool   `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	TaskID  string `json:"task_id,omitempty"`
}

// AllCloudCredsIDsResponse struct to store all cloud credentials IDs
type AllCloudCredsIDsResponse struct {
	Error          bool     `json:"error,omitempty"`
	CloudCredsList []string `json:"cloud_creds_list,omitempty"`
}

// CloudCredentialsListResponse struct to store all cloud credentials
type CloudCredentialsListResponse struct {
	Error              bool                               `json:"error,omitempty"`
	Message            string                             `json:"message,omitempty"`
	CredentialsInspect []api.SdkCredentialInspectResponse `json:"credentials_inspect,omitempty"`
}
