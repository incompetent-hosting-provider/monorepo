package db_instances

import (
	"context"
	"incompetent-hosting-provider/backend/pkg/db"
	"incompetent-hosting-provider/backend/pkg/util"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/rs/zerolog/log"
)

// TODO: Ports, Type, dates, update function, change deletion to delete by user isntead of single

const (
	STATUS_VALUE_SCHEDULED = "Scheduled"
	STATUS_VALUE_RUNNING   = "Running"
	STATUS_VALUE_STOPPED   = "Stopped"
)

const TABLE_NAME = "instances"

type InstancesTable struct {
	UserSub              string
	InstanceId           string
	ContainerUUID        string
	ContainerName        string
	ContainerDescription string
	ImageName            string
	ImageTag             string
	ContainerPorts       []int
	InstanceStatus       string
	CreatedAt            string
	StartedAt            string
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
			AttributeName: aws.String("InstanceId"),
			AttributeType: types.ScalarAttributeTypeS,
		},
	}
	keySchema := []types.KeySchemaElement{
		{
			AttributeName: aws.String("InstanceId"),
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
		ContainerName:        "test",
		ContainerDescription: "test",
		ImageName:            "test",
		ImageTag:             "test",
		ContainerUUID:        "kjasdjkas",
		ContainerPorts:       []int{1, 2, 4, 5},
		InstanceStatus:       "Starting",
	})
}

func InsertInstance(instanceItem InstancesTable) error {

	// if in test run -> Skip and return nil
	if util.IsTestRun() {
		return nil
	}

	conn := db.GetDynamoConn()

	instanceItem.InstanceId = instanceItem.UserSub + instanceItem.ContainerUUID

	//portItem, _ := dynamodbattribute.MarshalList(instanceItem.ContainerPorts)

	param := dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"UserSub":              &types.AttributeValueMemberS{Value: instanceItem.UserSub},
			"InstanceId":           &types.AttributeValueMemberS{Value: instanceItem.InstanceId},
			"ContainerName":        &types.AttributeValueMemberS{Value: instanceItem.ContainerName},
			"ContainerDescription": &types.AttributeValueMemberS{Value: instanceItem.ContainerDescription},
			"ImageName":            &types.AttributeValueMemberS{Value: instanceItem.ImageName},
			"ImageTag":             &types.AttributeValueMemberS{Value: instanceItem.ImageTag},
			"ContainerUUID":        &types.AttributeValueMemberS{Value: instanceItem.ContainerUUID},
			//"ContainerPorts":       portItem,
			"InstanceStatus": &types.AttributeValueMemberS{Value: instanceItem.InstanceStatus},
			"CreatedAt":      &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
			"StartedAt":      &types.AttributeValueMemberS{Value: "N/A"},
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

func GetAllUserInstances(usersub string) ([]InstancesTable, error) {

	// If in test run -> Skip and return dummy values
	if util.IsTestRun() {
		return nil, nil
	}

	conn := db.GetDynamoConn()

	params := &dynamodb.ScanInput{
		TableName:        aws.String(TABLE_NAME),
		FilterExpression: aws.String("UserSub = :UserSub"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":UserSub": &types.AttributeValueMemberS{
				Value: usersub,
			},
		},
	}

	var result []InstancesTable

	for {

		log.Info().Msg("iteration")

		scanResult, err := conn.Scan(context.TODO(), params)
		if err != nil {
			log.Warn().Msgf("Error in continued scanning: %v", err)
			return nil, err
		}

		// Print the items
		for _, item := range scanResult.Items {
			parsedItem := parseScanResultToStruct(item)
			log.Warn().Msgf("%v", parsedItem)
			result = append(result, parsedItem)
		}

		if params.ExclusiveStartKey == nil {
			break
		}

		params.ExclusiveStartKey = scanResult.LastEvaluatedKey
	}

	return result, nil
}

func GetInstanceById(userSub string, containerUUID string) (InstancesTable, error) {
	if util.IsTestRun() {
		return InstancesTable{}, nil
	}

	conn := db.GetDynamoConn()

	params := dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]types.AttributeValue{
			"UserSub":       &types.AttributeValueMemberS{Value: userSub},
			"ContainerUUID": &types.AttributeValueMemberS{Value: containerUUID},
		},
	}

	instance, err := conn.GetItem(context.TODO(), &params)

	if err != nil {
		log.Warn().Msgf("Could not get item due to an error: %v", err)
		return InstancesTable{}, err
	}

	return parseScanResultToStruct(instance.Item), err
}

func DeleteInstanceById(userSub string, containerUUID string) error {

	if util.IsTestRun() {
		return nil
	}

	instanceId := userSub + containerUUID

	conn := db.GetDynamoConn()

	params := &dynamodb.DeleteItemInput{
		TableName: aws.String(TABLE_NAME),
		Key:       map[string]types.AttributeValue{"InstanceId": &types.AttributeValueMemberS{Value: instanceId}},
	}

	_, err := conn.DeleteItem(context.TODO(), params)

	return err
}
