package db

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	domain "github.com/b4t3ou/sqs-reader-service/internal"
)

type DynamoDBOption func(client *Dynamo)

type Dynamo struct {
	table   string
	region  string
	session *session.Session
	svc     *dynamodb.DynamoDB
}

func NewDynamo(tableName string, options ...DynamoDBOption) (*Dynamo, error) {
	client := &Dynamo{
		table:  tableName,
		region: domain.DefaultAWSRegion,
	}

	for _, option := range options {
		option(client)
	}

	if err := client.setSession(); err != nil {
		return nil, err
	}

	client.svc = dynamodb.New(client.session)

	return client, nil
}

func WithDynamoLocalstackSession() DynamoDBOption {
	s := session.Must(session.NewSession(aws.NewConfig().
		WithRegion(domain.DefaultAWSRegion).
		WithEndpoint("http://localhost:4566").
		WithDisableEndpointHostPrefix(true).
		WithDisableSSL(true).
		WithCredentials(credentials.NewStaticCredentials("dummy", "dummy", "dummy")),
	))

	return func(client *Dynamo) {
		client.session = s
	}
}

func WithDynamoRegion(region string) DynamoDBOption {
	return func(client *Dynamo) {
		client.region = region
	}
}

func (db *Dynamo) Save(item domain.Event) error {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(db.table),
	}

	_, err = db.svc.PutItem(input)
	return err
}

func (db *Dynamo) Get(id string) (*domain.Event, error) {
	result, err := db.svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(db.table),
		Key:       map[string]*dynamodb.AttributeValue{"id": {S: aws.String(id)}},
	})

	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, fmt.Errorf("could not find item")
	}

	event := &domain.Event{}

	err = dynamodbattribute.UnmarshalMap(result.Item, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (db *Dynamo) setSession() error {
	if db.session != nil {
		return nil
	}

	awsConfig := &aws.Config{
		Region: aws.String(db.region),
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return err
	}

	db.session = sess
	return nil
}
