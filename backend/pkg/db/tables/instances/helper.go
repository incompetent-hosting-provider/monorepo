package db_instances

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog/log"
)

func parseScanResultToStruct(result map[string]types.AttributeValue) InstancesTable {
	log.Warn().Msgf("%v", result)
	r := InstancesTable{
		UserSub:              result["UserSub"].(*types.AttributeValueMemberS).Value,
		ContainerUUID:        result["ContainerUUID"].(*types.AttributeValueMemberS).Value,
		ContainerName:        result["ContainerName"].(*types.AttributeValueMemberS).Value,
		ContainerDescription: result["ContainerDescription"].(*types.AttributeValueMemberS).Value,
		ImageName:            result["ImageName"].(*types.AttributeValueMemberS).Value,
		ImageTag:             result["ImageTag"].(*types.AttributeValueMemberS).Value,
		//ContainerPorts     : result["ContainerPorts"].(*types.AttributeValueMemberN).Value,
		InstanceStatus: result["InstanceStatus"].(*types.AttributeValueMemberS).Value,
		StartedAt:      result["StartedAt"].(*types.AttributeValueMemberS).Value,
		CreatedAt:      result["CreatedAt"].(*types.AttributeValueMemberS).Value,
	}
	return r
}
