package db_instances

import (
	"context"
	"incompetent-hosting-provider/backend/pkg/db"
	"incompetent-hosting-provider/backend/pkg/util"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	//	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/rs/zerolog/log"
)

const TABLE_NAME = "instances"

const X = 1

type InstancesTable struct {
	UserSub              string
	InstanceId           string
	ContainerName        string
	ContainerDescription string
	ImageName            string
	ImageTag             string
	ContainerPorts       []int
}

func init() {

	// No setup needed in test run
	if util.IsTestRun() {
		return
	}

	doesTableExist, err := db.DoesTableExist(TABLE_NAME)

	if err != nil {
		panic("Could not fetch table status")
	}

	// if table already exists => we are done here
	if doesTableExist {
		log.Info().Msgf("Table %v is present", TABLE_NAME)
		return
	}

	log.Info().Msgf("Table %v does not exist. Creating...", TABLE_NAME)

	contents := []types.AttributeDefinition{
		{
			AttributeName: aws.String("UserSub"),
			AttributeType: types.ScalarAttributeTypeS,
		},
		{
			AttributeName: aws.String("ContainerId"),
			AttributeType: types.ScalarAttributeTypeS,
		},
	}

	keySchema := []types.KeySchemaElement{
		{
			AttributeName: aws.String("UserSub"),
			KeyType:       types.KeyTypeHash,
		},
		{
			AttributeName: aws.String("ContainerId"),
			KeyType:       types.KeyTypeHash,
		},
	}

	err = db.CreateTable(TABLE_NAME, contents, keySchema)

	if err != nil {
		log.Error().Msgf("Could not create table due to an error (%v)", err)
	} else {
		log.Debug().Msgf("Table %v was created.", TABLE_NAME)
	}

	InsertInstance(InstancesTable{
		UserSub:              "test",
		InstanceId:           "askjda",
		ContainerName:        "test",
		ContainerDescription: "test",
		ImageName:            "test",
		ImageTag:             "test",
		ContainerPorts:       []int{1, 2, 4, 5},
	})
}

func InsertInstance(instanceItem InstancesTable) error {

	// if in test run -> Skip and return nil
	if util.IsTestRun() {
		return nil
	}

	conn := db.GetDynamoConn()

	//portItem, _ := dynamodbattribute.MarshalList(instanceItem.ContainerPorts)

	param := dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"UserSub":              &types.AttributeValueMemberS{Value: instanceItem.UserSub},
			"InstanceId":           &types.AttributeValueMemberS{Value: instanceItem.InstanceId},
			"ContainerName":        &types.AttributeValueMemberS{Value: instanceItem.ContainerName},
			"ContainerDescription": &types.AttributeValueMemberS{Value: instanceItem.ContainerDescription},
			"ImageName":            &types.AttributeValueMemberS{Value: instanceItem.ImageName},
			"ImageTag":             &types.AttributeValueMemberS{Value: instanceItem.ImageTag},
			//"ContainerPorts":       portItem,
		},
		TableName: aws.String(TABLE_NAME),
	}

	_, err := conn.PutItem(context.TODO(), &param)

	if err != nil {
		log.Warn().Msgf("Could not insert balance item %v", err)
	} else {
		log.Warn().Msg("Created")
	}
	return err
}
