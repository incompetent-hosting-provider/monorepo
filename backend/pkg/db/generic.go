package db

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog/log"
)

func getTableCount(dynamoDbConn *dynamodb.Client) (int, error) {
	cnt, err := dynamoDbConn.ListTables(
		context.TODO(), &dynamodb.ListTablesInput{
			Limit: aws.Int32(5),
		})
	if err != nil {
		return 0, err
	}
	return len(cnt.TableNames), nil
}

func DoesTableExist(tableName string) (bool, error) {
	conn := GetDynamoConn()

	param := &dynamodb.DescribeTableInput{
		TableName: &tableName,
	}

	_, err := conn.DescribeTable(context.TODO(), param)

	if err != nil {
		var notFoundEx *types.ResourceNotFoundException
		if errors.As(err, &notFoundEx) {
			err = nil
		}
		return false, err
	}
	log.Debug().Msgf("Table %v was found.", tableName)
	return true, err

}

func CreateTable(tableName string, contents []types.AttributeDefinition, keySchema []types.KeySchemaElement) error {
	conn := GetDynamoConn()

	params := dynamodb.CreateTableInput{
		TableName:            &tableName,
		AttributeDefinitions: contents,
		KeySchema:            keySchema,
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	}

	_, err := conn.CreateTable(context.TODO(), &params)
	// Passs error handling to upper function
	return err
}
