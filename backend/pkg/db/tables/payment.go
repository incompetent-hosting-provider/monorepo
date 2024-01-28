package db_payment

import (
	"context"
	"incompetent-hosting-provider/backend/pkg/db"
	"incompetent-hosting-provider/backend/pkg/util"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog/log"
)

const TABLE_NAME string = "payment"

type PaymentTable struct {
	UserSub string
	Balance int
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
	}

	keySchema := []types.KeySchemaElement{
		{
			AttributeName: aws.String("UserSub"),
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

func IncreaseBalance(userSub string, balanceDelta int) (int, error) {

	// If in test -> skip shenanigans with dynamoDB
	// We do not test DynamoDb functionality/s3
	if util.IsTestRun() {
		return balanceDelta, nil
	}

	log.Debug().Msg("Updating user balance")
	conn := db.GetDynamoConn()

	updateExpression := aws.String("SET Balance = if_not_exists(Balance, :setValue) + :amount")
	expressionAttributeValues := map[string]types.AttributeValue{
		":amount":   &types.AttributeValueMemberN{Value: strconv.Itoa(balanceDelta)},
		":setValue": &types.AttributeValueMemberN{Value: "0"},
	}

	updateItem := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(TABLE_NAME),
		Key:                       map[string]types.AttributeValue{"UserSub": &types.AttributeValueMemberS{Value: userSub}},
		UpdateExpression:          updateExpression,
		ExpressionAttributeValues: expressionAttributeValues,
		ReturnValues:              types.ReturnValueAllNew,
	}
	updatedItem, err := conn.UpdateItem(context.TODO(), updateItem)

	if err != nil {
		log.Warn().Msgf("Error writing to dynamodb (%v)", err)
		return 0, err
	}

	updatedBalance, err := strconv.Atoi(updatedItem.Attributes["Balance"].(*types.AttributeValueMemberN).Value)

	if err != nil {
		log.Warn().Msgf("Could not parse return statement (%v)", err)
	}

	return updatedBalance, nil
}
