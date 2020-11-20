package handler

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/rs/zerolog/log"

	domain "github.com/b4t3ou/sqs-reader-service/internal"
)

type SQSClient interface {
	Receive() (*sqs.ReceiveMessageOutput, error)
	Delete(receiptHandle string) error
}

type DBClient interface {
	Save(item domain.Event) error
}

type SQS struct {
	sqsClient SQSClient
	dbClient  DBClient
}

func NewSQS(sqsClient SQSClient, dbClient DBClient) SQS {
	return SQS{
		sqsClient: sqsClient,
		dbClient:  dbClient,
	}
}

func (h SQS) Run() error {
	for {
		data, err := h.sqsClient.Receive()
		if err != nil {
			return err
		}

		h.process(data.Messages)
	}
}

func (h SQS) process(messages []*sqs.Message) {
	for _, msg := range messages {
		if err := h.saveEvent(msg); err != nil {
			log.Err(err).Interface("msg", msg).Msg("failed to process event")
		}
	}
}

func (h SQS) saveEvent(message *sqs.Message) error {
	defer h.deleteMessage(message)

	var event domain.Event
	if err := json.Unmarshal([]byte(*message.Body), &event); err != nil {
		return err
	}

	log.Info().Interface("event", event).Msg("event has been received")

	return h.dbClient.Save(event)
}

func (h SQS) deleteMessage(message *sqs.Message) {
	if err := h.sqsClient.Delete(*message.ReceiptHandle); err != nil {
		log.Err(err).Interface("msg", message).Msg("failed to delete message")
	}
}
