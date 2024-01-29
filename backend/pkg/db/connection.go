package db

import (
	"context"
	"incompetent-hosting-provider/backend/pkg/util"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/rs/zerolog/log"
)

var dynamoDbConn *dynamodb.Client

func InitDbConn() error {

	// If already initialized -> skip
	// If is test run -> skip
	if dynamoDbConn != nil || util.IsTestRun() {
		return nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-central-1"),
		// load endpoint
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: util.GetStringEnvWithDefault("AWS_ENDPOINT", "http://localhost:8000")}, nil
			})),
		// Load credentials
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     util.GetStringEnvWithDefault("AWS_ACCESS_KEY_ID", "dummy"),
				SecretAccessKey: util.GetStringEnvWithDefault("AWS_SECRET_ACCESS_KEY", "dummy"),
				SessionToken:    util.GetStringEnvWithDefault("AWS_SESSION_TOKEN", "dummy"),
				Source:          "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		}),
	)

	if err != nil {
		log.Fatal().Msgf("unable to load SDK config, %v", err)
		return err
	}

	// Using the Config value, create the DynamoDB client
	newDbConn := dynamodb.NewFromConfig(cfg)

	tablecount, err := getTableCount(newDbConn)
	if err != nil {
		return err
	}
	log.Debug().Msgf("Connected to instance with %d tables.", tablecount)

	dynamoDbConn = newDbConn
	return nil
}

func GetDynamoConn() *dynamodb.Client {
	// This does not have a test guard, because this should never be reached if all other test guards have been placed properly
	if dynamoDbConn == nil {
		err := InitDbConn()
		// Panic this error as this means the db connection cannot be utilized
		if err != nil {
			panic(err)
		}
	}

	return dynamoDbConn
}
