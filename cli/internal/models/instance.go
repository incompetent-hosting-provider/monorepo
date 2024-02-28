package models

import (
	"fmt"
)


type Instance struct {
	InstanceType string `json:"type,omitempty"`
	ID string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Image ContainerImage `json:"container image,omitempty"`
	Status string `json:"status,omitempty"`
}

func (i Instance) String() string {
	return fmt.Sprintf("[%s] %s - %s, %s", i.Status, i.ID, i.Name, i.Image.String())
}
