package config

import api "github.com/libopenstorage/openstorage-sdk-clients/sdk/golang"

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
	VolumeEncrypted           string                   `json:"volume_encrypted,omitempty"`
	VolumeEncryptionKey       string                   `json:"volume_encryption_key,omitempty"`
	VolumeK8sNamespace        string                   `json:"volume_k8s_namespace,omitempty"`
	VolumeK8sPVCName          string                   `json:"volume_k8s_pvc_name,omitempty"`
	VolumeSharedv4            bool                     `json:"volume_sharedv4,omitempty"`
	VolumeSharedv4ServiceSpec *api.Sharedv4ServiceSpec `json:"volume_sharedv4_service_spec,omitempty"`
	VolumeIOStrategy          *api.IoStrategy          `json:"volume_io_strategy,omitempty"`
}
