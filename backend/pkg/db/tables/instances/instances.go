package db_instances

import (
	"context"
	"incompetent-hosting-provider/backend/pkg/db"
	"incompetent-hosting-provider/backend/pkg/util"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/rs/zerolog/log"
)

// TODO: Ports, Type, dates, update function, change deletion to delete by user isntead of single

const (
	STATUS_VALUE_SCHEDULED = "Scheduled"
	STATUS_VALUE_RUNNING   = "Running"
	STATUS_VALUE_STOPPED   = "Stopped"
	TYPE_CUSTOM            = "Custom"
	TYPE_PRESET            = "Preset"
)

const TABLE_NAME = "instances"

type ImageSpecification struct {
	Tag  string `dynamodbav:"tag"`
	Name string `dynamodbav:"name"`
}

type InstancesTable struct {
	UserSub              string             `dynamodbav:"usersub"`
	InstanceId           string             `dynamodbav:"instanceid"`
	ContainerUUID        string             `dynamodbav:"containeruuid"`
	ContainerName        string             `dynamodbav:"containername"`
	ContainerDescription string             `dynamodbav:"containerdescription"`
	Image                ImageSpecification `dynamodbav:"image"`
	ContainerPorts       []int              `dynamodbav:"containerports"`
	InstanceStatus       string             `dynamodbav:"instancestatus"`
	CreatedAt            string             `dynamodbav:"createdat"`
	StartedAt            string             `dynamodbav:"startedat"`
	Type                 string             `dynamodbav:"type"`
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
			AttributeName: aws.String("instanceid"),
			AttributeType: types.ScalarAttributeTypeS,
		},
	}
	keySchema := []types.KeySchemaElement{
		{
			AttributeName: aws.String("instanceid"),
			KeyType:       types.KeyTypeHash,
		},
	}

	err = db.CreateTable(TABLE_NAME, contents, keySchema)

	if err != nil {
		log.Error().Msgf("Could not create table due to an error (%v)", err)
	} else {
		log.Debug().Msgf("Table %v was created.", TABLE_NAME)
	}
}

func InsertInstance(instanceItem InstancesTable) error {

	// if in test run -> Skip and return nil
	if util.IsTestRun() {
		return nil
	}

	conn := db.GetDynamoConn()

	instanceItem.InstanceId = instanceItem.UserSub + instanceItem.ContainerUUID

	marshalledItem, _ := attributevalue.MarshalMap(instanceItem)

	log.Info().Msgf("%v", marshalledItem["image"])

	log.Info().Msgf("%v", marshalledItem)

	param := dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME),
		Item:      marshalledItem,
	}

	_, err := conn.PutItem(context.TODO(), &param)

	if err != nil {
		log.Warn().Msgf("Could not insert instance item %v", err)
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
		FilterExpression: aws.String("usersub = :usersub"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":usersub": &types.AttributeValueMemberS{
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

		// Parse the items
		for _, item := range scanResult.Items {
			var parsedItem InstancesTable
			err = attributevalue.UnmarshalMap(item, &parsedItem)
			if err != nil {
				log.Warn().Msgf("Could not parse item due to an error: %v", err)
			}
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
			"usersub":       &types.AttributeValueMemberS{Value: userSub},
			"containeruuid": &types.AttributeValueMemberS{Value: containerUUID},
		},
	}

	instance, err := conn.GetItem(context.TODO(), &params)

	if err != nil {
		log.Warn().Msgf("Could not get item due to an error: %v", err)
		return InstancesTable{}, err
	}

	var parsedInstance InstancesTable

	err = attributevalue.UnmarshalMap(instance.Item, &parsedInstance)

	if err != nil {
		log.Warn().Msgf("Could not parse item due to an error: %v", err)
	}

	return parsedInstance, err
}

func DeleteInstanceById(userSub string, containerUUID string) error {

	if util.IsTestRun() {
		return nil
	}

	instanceId := userSub + containerUUID

	conn := db.GetDynamoConn()

	params := &dynamodb.DeleteItemInput{
		TableName: aws.String(TABLE_NAME),
		Key:       map[string]types.AttributeValue{"instanceid": &types.AttributeValueMemberS{Value: instanceId}},
	}

	_, err := conn.DeleteItem(context.TODO(), params)

	return err
}
