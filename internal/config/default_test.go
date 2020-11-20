package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew_NoConfigFile(t *testing.T) {
	cfgPath, _ := filepath.Abs(".")
	os.Clearenv()

	os.Setenv("ENV", "test-not-exists")
	_, err := NewConfig(cfgPath)
	assert.NotNil(t, err)
}

func TestNew_Env(t *testing.T) {
	os.Clearenv()

	cfgPath, _ := filepath.Abs("../../config")
	c, err := NewConfig(cfgPath)

	assert.Nil(t, err)
	assert.Equal(t, "local", c.Env)
	assert.Equal(t, "local-eu-west-1-test-queue", c.QueueName)
	assert.Equal(t, "local-eu-west-1-events", c.DynamoTable)
}
