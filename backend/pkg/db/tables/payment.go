package db_payment

import (
	"context"
	"incompetent-hosting-provider/backend/pkg/db"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog/log"
)

const TABLE_NAME string = "payment"

type PaymentTable struct {
	UserSub        string
	CurrentBalance int
}

func IncreaseBalance(userSub string, balanceDelta int) (int, error) {
	conn := db.GetDynamoConn()

	updateExpression := "SET balance = balance + :amount"
	expressionAttributeValues := map[string]types.AttributeValue{":amount": &types.AttributeValueMemberN{Value: strconv.Itoa(balanceDelta)}}

	updateItem := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(TABLE_NAME),
		Key:                       map[string]types.AttributeValue{"UserSub": &types.AttributeValueMemberS{Value: userSub}},
		UpdateExpression:          &updateExpression,
		ExpressionAttributeValues: expressionAttributeValues,
		ReturnValues:              types.ReturnValueAllNew,
	}
	updatedItem, err := conn.UpdateItem(context.TODO(), updateItem)

	if err != nil {
		log.Warn().Msgf("Error writing to dynamodb (%v)", err)
		return 0, err
	}

	updatedBalance, err := strconv.Atoi(updatedItem.Attributes["balance"].(*types.AttributeValueMemberN).Value)

	if err != nil {
		log.Warn().Msgf("Could not parse return statement (%v)", err)
	}

	return updatedBalance, nil

}
