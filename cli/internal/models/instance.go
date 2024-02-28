package models

import "time"


type InstanceStatus string
const (
	Running InstanceStatus = "RUNNING"
	Terminated InstanceStatus = "TERMINATED"
	Pending InstanceStatus = "PENDING"
)

type ContainerImage struct {
	Name string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

type Instance struct {
	InstanceType string `json:"type,omitempty"`
	ID string `json:"instance_id,omitempty"`
	Name string `json:"name,omitempty"`
	Image ContainerImage `json:"container_image,omitempty"`
	Status InstanceStatus `json:"status,omitempty"`
}

type InstanceDetail struct {
	Instance
	StartedAt time.Time `json:"started_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	OpenPorts []int `json:"open_ports,omitempty"`
	Description string `json:"description,omitempty"`
}

type InstancePreset struct {
	ID string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`	
	Description string `json:"description,omitempty"`
}