package handler

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestNewSQS_Run(t *testing.T) {
	h := NewSQS(sqsClientMock{
		error: fmt.Errorf("cannot receice data"),
	})
	assert.NotNil(t, h.Run())
}

func TestSQS_saveEvent(t *testing.T) {
	h := NewSQS(sqsClientMock{})
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
}
