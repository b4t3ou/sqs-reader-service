// +build integration

package db

import (
	"os"
	"path/filepath"
	"testing"

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

func TestNewDynamo_Save(t *testing.T) {
	client, err := NewDynamo(testConfig.DynamoTable, WithDynamoLocalstackSession())
	assert.Nil(t, err)

	err = client.Save(domain.Event{})
	assert.NotNil(t, err)

	id := uuid.NewV4().String()

	err = client.Save(domain.Event{
		ID:      id,
		Message: "test message",
	})
	assert.Nil(t, err)

	data, err := client.Get(id)
	assert.Nil(t, err)

	assert.Equal(t, id, data.ID)
	assert.Equal(t, "test message", data.Message)
}
