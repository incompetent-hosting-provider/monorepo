package instances

import db_instances "incompetent-hosting-provider/backend/pkg/db/tables/instances"

func serializeInstanceResponses(instances []db_instances.InstancesTable) []InstanceInfo {
	var serializedInstances []InstanceInfo

	for _, instance := range instances {
		serializedInstances = append(serializedInstances, InstanceInfo{
			Type:           "",
			ContainerName:  instance.ContainerName,
			ContainerId:    instance.ContainerUUID,
			InstanceStatus: instance.InstanceStatus,
			ContainerImageData: ContainerImageDescription{
				Tag:       instance.ImageTag,
				ImageName: instance.ImageName,
			},
		})
	}
	return serializedInstances
}

func serializeDetailedInstanceResponse(instance db_instances.InstancesTable) InstanceInfoDetailedResponse {
	return InstanceInfoDetailedResponse{
		Type:           "",
		ContainerName:  instance.ContainerName,
		ContainerId:    instance.ContainerUUID,
		InstanceStatus: instance.InstanceStatus,
		ContainerImageData: ContainerImageDescription{
			Tag:       instance.ImageTag,
			ImageName: instance.ImageName,
		},
		CreatedAt:      instance.CreatedAt,
		StartedAt:      instance.StartedAt,
		ContainerPorts: instance.ContainerPorts,
	}
}
