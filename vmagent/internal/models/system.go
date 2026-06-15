package models

// Host information
type Host struct {
	Hostname    string `json:"hostname"`
	CurrentUser string `json:"current_user"`
	Timestamp   string `json:"timestamp"`
}

// OS information
type OS struct {
	Name         string `json:"name"`
	Edition      string `json:"edition"`
	Version      string `json:"version"`
	BuildNumber  string `json:"build_number"`
	Kernel       string `json:"kernel_version"`
	Architecture string `json:"architecture"`
	BootTime     string `json:"boot_time"`
}

// CPU information
type CPU struct {
	Model         string `json:"model"`
	PhysicalCores int    `json:"physical_cores"`
	LogicalCores  int    `json:"logical_cores"`
	Vendor        string `json:"vendor"`
	Architecture  string `json:"architecture"`
}

// Memory information
type Memory struct {
	TotalBytes     uint64  `json:"total_bytes"`
	AvailableBytes uint64  `json:"available_bytes"`
	UsedBytes      uint64  `json:"used_bytes"`
	UsedPercent    float64 `json:"used_percent"`
}

// Disk information (supports multiple disks)
type Disk struct {
	MountPoint  string  `json:"mount_point"`
	Filesystem  string  `json:"filesystem"`
	TotalBytes  uint64  `json:"total_bytes"`
	UsedBytes   uint64  `json:"used_bytes"`
	FreeBytes   uint64  `json:"free_bytes"`
	UsedPercent float64 `json:"used_percent"`
}

// NetworkInterface information
type NetworkInterface struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	MACAddress string   `json:"mac_address"`
	IPv4       []string `json:"ipv4"`
	IPv6       []string `json:"ipv6"`
}
