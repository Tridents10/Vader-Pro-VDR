package models

// Inventory is the root object exported to JSON
type Inventory struct {
	AgentUUID    string             `json:"agent_uuid"`
	AgentVersion string             `json:"agent_version"`
	Timestamp    string             `json:"timestamp"`
	Host         Host               `json:"host"`
	OS           OS                 `json:"os"`
	CPU          CPU                `json:"cpu"`
	Memory       Memory             `json:"memory"`
	Disks        []Disk             `json:"disks"`
	Network      []NetworkInterface `json:"network"`
	Applications []Application      `json:"applications"`
}
