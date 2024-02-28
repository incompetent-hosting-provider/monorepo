package models

import (
	"fmt"
	"strings"
	"time"
)


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

func (c ContainerImage) String() string {
	return fmt.Sprintf("%s:%s", c.Name, c.Version)
}

type Instance struct {
	InstanceType string `json:"type,omitempty"`
	ID string `json:"instance_id,omitempty"`
	Name string `json:"name,omitempty"`
	Image ContainerImage `json:"container_image,omitempty"`
	Status InstanceStatus `json:"status,omitempty"`
}

func (i Instance) String() string {
	return fmt.Sprintf("[%s] %s - %s, %s", i.Status, i.ID, i.Name, i.Image.String())
}

type InstanceDetail struct {
	Instance
	StartedAt time.Time `json:"started_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	OpenPorts []int `json:"open_ports,omitempty"`
	Description string `json:"description,omitempty"`
}

func (i InstanceDetail) String() string {
	builder := strings.Builder{}
	// Generic header as in Instance.String()
	builder.WriteString(i.Instance.String())
	builder.WriteString(fmt.Sprintf("\nDescription: %s", i.Description))
	builder.WriteString(fmt.Sprintf("\nCreated at: %s", i.CreatedAt.Format(time.RFC3339)))
	if i.Status == Running {
		builder.WriteString(fmt.Sprintf("\nStarted at: %s", i.StartedAt.Format(time.RFC3339)))
	}

	if len(i.OpenPorts) > 0 {
		builder.WriteString("\nOpen Ports: ")
		for _, port := range i.OpenPorts {
			builder.WriteString(fmt.Sprintf("%d ", port))
		}
	}

	return builder.String()
}

type InstancePreset struct {
	ID string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`	
	Description string `json:"description,omitempty"`
}

func (i InstancePreset) String() string {
	return fmt.Sprintf("%s: %s - %s", i.ID, i.Name, i.Description)
}