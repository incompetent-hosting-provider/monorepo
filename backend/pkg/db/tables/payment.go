package db_payment

import (
	"context"
	"errors"
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

func GetUserBalance(userSub string) (int, error) {

	// If in test run -> Skip and return dummy values
	if util.IsTestRun() {
		return 1000, nil
	}

	conn := db.GetDynamoConn()

	param := dynamodb.GetItemInput{
		Key:            map[string]types.AttributeValue{"UserSub": &types.AttributeValueMemberS{Value: userSub}},
		TableName:      aws.String(TABLE_NAME),
		ConsistentRead: aws.Bool(true),
	}

	var balance int
	val, err := conn.GetItem(context.TODO(), &param)

	// Is this ideal? I am not sure
	if err == nil {
		var notFoundEx *types.ResourceNotFoundException
		// If user not in table -> Assume that user has never added currency yet i.e. has a current balance of zero
		if errors.As(err, &notFoundEx) {
			err = nil
		}
		balance = 0
	} else {
		if val.Item["Balance"] == nil {
			log.Warn().Msg("There was an error parsing the balance for a user. This may indiciate corrupt data.")
			// Set balance to 0 to avoid panic and still return a value
			// This edge case can also occur if the user is not created in the balance table therefore throwing an internal server error would (for now) not be sufficient.
			// For the user to get here the JWT has to be valid -> the user likely exists on the keycloak side of things
			return 0, nil
		}
		balance, _ = strconv.Atoi(val.Item["Balance"].(*types.AttributeValueMemberN).Value)
	}
	return balance, err
}


func DeleteUserBalance(userSub string) error {

	// If in test run -> Skip and return dummy values
	if util.IsTestRun() {
		return nil
	}

	conn := db.GetDynamoConn()

	param := dynamodb.DeleteItemInput{
		Key:       map[string]types.AttributeValue{"UserSub": &types.AttributeValueMemberS{Value: userSub}},
		TableName: aws.String(TABLE_NAME),
	}

	_, err := conn.DeleteItem(context.TODO(), &param)

	if err != nil {
		log.Warn().Msgf("Could not delete balance item", err)
	}
	return err
}


func InsertUserBalance(userSub string) error {

	// If in test run -> Skip and return dummy values
	if util.IsTestRun() {
		return nil
	}

	conn := db.GetDynamoConn()

	param := dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"UserSub": &types.AttributeValueMemberS{Value: userSub},
			"Balance": &types.AttributeValueMemberN{Value: "183"},
		},
		TableName: aws.String(TABLE_NAME),
	}

	_, err := conn.PutItem(context.TODO(), &param)

	if err != nil {
		log.Warn().Msgf("Could not insert balance item", err)
	}
	return err
}