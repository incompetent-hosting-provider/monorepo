package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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
