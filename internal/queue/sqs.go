package queue

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	domain "github.com/b4t3ou/sqs-reader-service/internal"
)

type SQSClientOption func(client *SQSClient)

type SQSClient struct {
	batchSize         int64
	visibilityTimeout int64
	waitTimeSeconds   int64
	publishDelay      int64
	queueName         string
	queueURL          string
	session           *session.Session
	region            string
	*sqs.SQS
}

func NewSQSClient(queueName string, options ...SQSClientOption) (*SQSClient, error) {
	client := &SQSClient{
		batchSize:         10,
		visibilityTimeout: 10,
		region:            domain.DefaultAWSRegion,
		queueName:         queueName,
	}

	for _, option := range options {
		option(client)
	}

	if err := client.setSession(); err != nil {
		return nil, err
	}

	client.SQS = sqs.New(client.session)

	url, err := client.SQS.GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: aws.String(client.queueName)})
	if err != nil {
		return nil, err
	}

	client.queueURL = *url.QueueUrl

	return client, nil
}

func WithSQSLocalstackSession() SQSClientOption {
	s := session.Must(session.NewSession(aws.NewConfig().
		WithRegion(domain.DefaultAWSRegion).
		WithEndpoint("http://localhost:4566").
		WithDisableEndpointHostPrefix(true).
		WithDisableSSL(true).
		WithCredentials(credentials.NewStaticCredentials("dummy", "dummy", "dummy")),
	))

	return func(client *SQSClient) {
		client.session = s
	}
}

func WithSQSBatchSize(batchSize int64) SQSClientOption {
	return func(client *SQSClient) {
		client.batchSize = batchSize
	}
}

func WithSQSVisibilityTimeout(visibilityTimeout int64) SQSClientOption {
	return func(client *SQSClient) {
		client.visibilityTimeout = visibilityTimeout
	}
}

func WithSQSPublishDelay(publishDelay int64) SQSClientOption {
	return func(client *SQSClient) {
		client.publishDelay = publishDelay
	}
}

func WithSQSRegion(region string) SQSClientOption {
	return func(client *SQSClient) {
		client.region = region
	}
}

func (c *SQSClient) Publish(data interface{}) (*sqs.SendMessageOutput, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	input := &sqs.SendMessageInput{
		MessageBody:  aws.String(string(body)),
		QueueUrl:     aws.String(c.queueURL),
		DelaySeconds: aws.Int64(c.publishDelay),
	}

	return c.SendMessage(input)
}

func (c *SQSClient) Receive() (*sqs.ReceiveMessageOutput, error) {
	params := &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(c.queueURL),
		AttributeNames:        []*string{aws.String("All")},
		MaxNumberOfMessages:   aws.Int64(c.batchSize),
		MessageAttributeNames: []*string{aws.String("All")},
		VisibilityTimeout:     aws.Int64(c.visibilityTimeout),
		WaitTimeSeconds:       aws.Int64(c.waitTimeSeconds),
	}

	return c.ReceiveMessage(params)
}

func (c *SQSClient) Delete(receiptHandle string) error {
	_, err := c.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})

	return err
}

func (c *SQSClient) setSession() error {
	if c.session != nil {
		return nil
	}

	awsConfig := &aws.Config{
		Region: aws.String(c.region),
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return err
	}

	c.session = sess
	return nil
}
