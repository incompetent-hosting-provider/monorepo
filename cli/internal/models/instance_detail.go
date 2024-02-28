package models

import (
	"fmt"
	"strings"
	"time"
)

type InstanceDetail struct {
	Instance
	StartedAtString string `json:"started_at,omitempty"`
	CreatedAtString string `json:"created_at,omitempty"`
	OpenPorts []int `json:"open_ports,omitempty"`
	Description string `json:"description,omitempty"`
}

func (i InstanceDetail) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("ID: %s", i.ID))
	builder.WriteString(fmt.Sprintf("\nName: %s", i.Name))
	builder.WriteString(fmt.Sprintf("\nImage: %s", i.Image))
	builder.WriteString(fmt.Sprintf("\nDescription: %s", i.Description))
	createdAt := i.CreatedAt()
	if createdAt != nil {
		builder.WriteString(fmt.Sprintf("\nCreated at: %s", createdAt.Local().Format(time.RFC1123)))
	}

	startedAt := i.StartedAt()
	if startedAt != nil {
		builder.WriteString(fmt.Sprintf("\nStarted at: %s", startedAt.Local().Format(time.RFC1123)))
	}

	if len(i.OpenPorts) > 0 {
		builder.WriteString("\nOpen Ports: ")
		for index, port := range i.OpenPorts {
			isLast := index == len(i.OpenPorts) - 1
			if isLast {
				builder.WriteString(fmt.Sprintf("%d ", port))
				} else {
				builder.WriteString(fmt.Sprintf("%d, ", port))
			}
		}
	}

	return builder.String()
}

func (i InstanceDetail) StartedAt() *time.Time {
	return convertTimeString(i.StartedAtString)
}

func (i InstanceDetail) CreatedAt() *time.Time {
	return convertTimeString(i.CreatedAtString)
}

func convertTimeString(timeString string) *time.Time {
	if timeString == "" || timeString == "N/A" {
		return nil
	} else {
		time, err := time.Parse(time.RFC3339, timeString)
		if err != nil {
			return nil
		}
		return &time
	}
}