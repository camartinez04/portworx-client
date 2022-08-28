package main

import (
	"time"
)

// User is the user model
type User struct {
	ID         int
	FirstName  string
	LastName   string
	Email      string
	Password   string
	AcessLevel int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Rooms is the rooms model
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restriction is the restrictions model
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Reservation is the reservation model
type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room
	Processed int
}

// Room Restriction is the room restriction model
type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	ReservationID int
	RestrictionID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Room
	Reservation   Reservation
	Restriction   Restriction
}

// MailData holds an email message
type MailData struct {
	To       string
	From     string
	Subject  string
	Content  string
	Template string
}

// TemplateData holds data sent from handlers to template
type TemplateData struct {
	StringMap          map[string]string
	IntMap             map[string]int
	FloatMap           map[string]float32
	Data               map[string]interface{}
	JsonVolumeInspect  JsonVolumeInspect
	CSRFToken          string
	Flash              string
	Warning            string
	Error              string
	IsAuthenticated    int
	JsonUsageVolume    JsonUsageVolume
	IoProfileString    string
	VolumeStatusString string
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
