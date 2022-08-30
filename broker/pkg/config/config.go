package config

// NodeInfo struct to store the node information.
type NodeInfo struct {
	NodeName          string  `json:"node_name,omitempty"`
	NodeStatus        string  `json:"node_status,omitempty"`
	NodeAvgLoad       int64   `json:"node_avg_load,omitempty"`
	NumberOfPools     int     `json:"number_of_pools,omitempty"`
	SizeNodePool      uint64  `json:"size_node_pool,omitempty"`
	UsedNodePool      uint64  `json:"used_node_pool,omitempty"`
	FreeNodePool      uint64  `json:"free_node_pool,omitempty"`
	NodeMemTotal      uint64  `json:"node_mem_total,omitempty"`
	NodeMemUsed       uint64  `json:"node_mem_used,omitempty"`
	NodeMemFree       uint64  `json:"node_mem_free,omitempty"`
	PercentUsedPool   float64 `json:"percent_used_pool,omitempty"`
	PercentUsedMemory float64 `json:"percent_used_memory,omitempty"`
	StoragelessNode   bool    `json:"storageless,omitempty"`
}

// VolumeInfo struct to store the volume information.
type VolumeInfo struct {
	VolumeName             string  `json:"volume_name,omitempty"`
	VolumeID               string  `json:"volume_id,omitempty"`
	VolumeReplicas         int     `json:"volume_replicas,omitempty"`
	VolumeStatus           string  `json:"volume_status,omitempty"`
	VolumeSize             uint64  `json:"volume_size,omitempty"`
	VolumeUsed             uint64  `json:"volume_used,omitempty"`
	VolumeAvailable        uint64  `json:"volume_available,omitempty"`
	VolumeUsedPercent      float64 `json:"volume_used_percent,omitempty"`
	VolumeType             string  `json:"volume_type,omitempty"`
	VolumeAttachedOn       string  `json:"volume_attached_on,omitempty"`
	VolumeAttachStatus     string  `json:"volume_attach_status,omitempty"`
	VolumeDevicePath       string  `json:"volume_device_path,omitempty"`
	VolumeAggregationLevel uint32  `json:"volume_aggregation_level,omitempty"`
}
