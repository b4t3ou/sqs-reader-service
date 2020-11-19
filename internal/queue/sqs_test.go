// +build integration

package queue

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	domain "github.com/b4t3ou/sqs-reader-service/internal"
	"github.com/b4t3ou/sqs-reader-service/internal/config"
)

var (
	testConfig *config.Config
)

func TestMain(m *testing.M) {
	cfgPath, _ := filepath.Abs("../../config")
	os.Clearenv()

	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	testConfig = cfg

	code := m.Run()
	os.Exit(code)
}

func TestNewSQSClient_AWSSessionFails(t *testing.T) {
	_, err := NewSQSClient(testConfig.QueueName, testConfig.Env)
	assert.NotNil(t, err)
}

func TestSQSClient_Publish(t *testing.T) {
	client, err := NewSQSClient(
		testConfig.QueueName,
		testConfig.Env,
		WithSQSLocalstackSession(),
		WithSQSVisibilityTimeout(1),
		WithSQSBatchSize(10),
	)
	assert.Nil(t, err)
	id := uuid.NewV4().String()

	_, err = client.Publish(domain.Event{
		ID:        id,
		Message:   "test message",
		Timestamp: time.Now(),
	})
	assert.Nil(t, err)

	data, err := client.Receive()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(data.Messages))

	event := &domain.Event{}
	err = json.Unmarshal([]byte(*data.Messages[0].Body), event)
	assert.Nil(t, err)
	assert.Equal(t, id, event.ID)

	err = client.Delete(*data.Messages[0].ReceiptHandle)
	assert.Nil(t, err)
}
