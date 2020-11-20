package handler

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"

	domain "github.com/b4t3ou/sqs-reader-service/internal"
)

type sqsClientMock struct {
	error   error
	message *sqs.ReceiveMessageOutput
}

func (c sqsClientMock) Receive() (*sqs.ReceiveMessageOutput, error) {
	return c.message, c.error
}

func (c sqsClientMock) Delete(receiptHandle string) error {
	return nil
}

type dbMock struct {
	event domain.Event
	error error
}

func (db *dbMock) Save(item domain.Event) error {
	db.event = item
	return db.error
}

func TestNewSQS_Run(t *testing.T) {
	h := NewSQS(sqsClientMock{
		error: fmt.Errorf("cannot receice data"),
	}, &dbMock{},
	)
	assert.NotNil(t, h.Run())
}

func TestSQS_saveEvent(t *testing.T) {
	db := &dbMock{}
	h := NewSQS(sqsClientMock{}, db)
	err := h.saveEvent(&sqs.Message{
		ReceiptHandle: aws.String("foo"),
		Body:          aws.String(``),
	})
	assert.NotNil(t, err)

	err = h.saveEvent(&sqs.Message{
		Body:          aws.String(`{"id":"foo"}`),
		ReceiptHandle: aws.String("foo"),
	})
	assert.Nil(t, err)
	assert.Equal(t, "foo", db.event.ID)
}

func TestSQS_saveEvent_dbError(t *testing.T) {
	db := &dbMock{error: fmt.Errorf("failed to save data")}

	h := NewSQS(sqsClientMock{}, db)
	err := h.saveEvent(&sqs.Message{
		ReceiptHandle: aws.String("foo"),
		Body:          aws.String(``),
	})
	assert.NotNil(t, err)

	err = h.saveEvent(&sqs.Message{
		Body:          aws.String(`{"id":"foo"}`),
		ReceiptHandle: aws.String("foo"),
	})
	assert.NotNil(t, err)
}
