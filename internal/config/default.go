package config

import (
	"io/ioutil"

	"github.com/caarlos0/env"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Env         string `env:"ENV" envDefault:"local"`
	QueueName   string `yaml:"queueName"`
	DynamoTable string `yaml:"dynamoTable"`
}

// NewConfig returns with a new Config object
func NewConfig(configPath string) (*Config, error) {
	c := &Config{}

	if err := env.Parse(c); err != nil {
		return c, err
	}

	if err := c.setServiceConfig(configPath); err != nil {
		return c, err
	}

	v := validator.New()
	if err := v.Struct(c); err != nil {
		return c, err
	}

	return c, nil
}

func (c *Config) setServiceConfig(path string) error {
	source, err := ioutil.ReadFile(path + "/" + c.Env + ".yaml")
	if err != nil {
		return err
	}

	return yaml.Unmarshal(source, &c)
}
