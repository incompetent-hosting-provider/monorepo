package models

import "fmt"

type InstancePreset struct {
	ID int `json:"id,omitempty"`
	Name string `json:"name,omitempty"`	
	Description string `json:"description,omitempty"`
}

func (i InstancePreset) String() string {
	return fmt.Sprintf("%d: %s - %s", i.ID, i.Name, i.Description)
}