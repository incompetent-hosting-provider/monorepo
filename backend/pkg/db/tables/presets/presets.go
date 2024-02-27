package db_presets

import (
	"context"
	"incompetent-hosting-provider/backend/pkg/db"
	db_instances "incompetent-hosting-provider/backend/pkg/db/tables/instances"
	"incompetent-hosting-provider/backend/pkg/util"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog/log"
)

const TABLE_NAME string = "presets"

type PresetTable struct {
	Image       db_instances.ImageSpecification `dynamodbav:"image"`
	Name        string                          `dynamodbav:"name"`
	Description string                          `dynamodbav:"description"`
	PresetId    int                             `dynamodbav:"presetid"`
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
			AttributeName: aws.String("presetid"),
			AttributeType: types.ScalarAttributeTypeN,
		},
	}

	keySchema := []types.KeySchemaElement{
		{
			AttributeName: aws.String("presetid"),
			KeyType:       types.KeyTypeHash,
		},
	}

	err = db.CreateTable(TABLE_NAME, contents, keySchema)

	if err != nil {
		log.Error().Msgf("Could not create table due to an error (%v)", err)
		return
	} else {
		log.Debug().Msgf("Table %v was created.", TABLE_NAME)
	}
	insertInitData()
}

func insertInitData() {
	if util.IsTestRun() {
		return
	}

	conn := db.GetDynamoConn()

	mySQLPreset, err := attributevalue.MarshalMap(PresetTable{
		PresetId:    1,
		Name:        "Mysql",
		Description: "Your very own instance of mysql",
		Image: db_instances.ImageSpecification{
			Name: "mysql",
			Tag:  "latest",
		},
	})

	if err != nil {
		log.Warn().Msgf("Could not prepare write for preset list due to an error: %v", err)
		return
	}

	items := []types.WriteRequest{
		{
			PutRequest: &types.PutRequest{
				Item: mySQLPreset,
			},
		},
	}

	params := dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			TABLE_NAME: items,
		},
	}

	_, err = conn.BatchWriteItem(context.TODO(), &params)

	if err != nil {
		log.Warn().Msgf("Write for preset list failed due to an error: %v", err)
		return
	}
	log.Debug().Msgf("Preset table was filled with %d items.", len(params.RequestItems))

}

func GetAllPresets() ([]PresetTable, error) {

	if util.IsTestRun() {
		return nil, nil
	}

	conn := db.GetDynamoConn()

	params := &dynamodb.ScanInput{
		TableName: aws.String(TABLE_NAME),
	}

	var result []PresetTable

	for {

		scanResult, err := conn.Scan(context.TODO(), params)
		if err != nil {
			log.Warn().Msgf("Error in continued scanning: %v", err)
			return nil, err
		}

		// Parse the items
		for _, preset := range scanResult.Items {
			var parsedPreset PresetTable
			err = attributevalue.UnmarshalMap(preset, &parsedPreset)
			if err != nil {
				log.Warn().Msgf("Could not parse item due to an error: %v", err)
			}
			result = append(result, parsedPreset)
		}

		if params.ExclusiveStartKey == nil {
			break
		}

		params.ExclusiveStartKey = scanResult.LastEvaluatedKey
	}

	return result, nil
}
