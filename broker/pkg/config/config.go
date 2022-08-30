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
