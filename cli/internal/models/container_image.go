package models

import "fmt"
type ContainerImage struct {
	Name string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

func (c ContainerImage) String() string {
	return fmt.Sprintf("%s:%s", c.Name, c.Version)
}
