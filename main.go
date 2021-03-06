package main

import (
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/b4t3ou/sqs-reader-service/internal/config"
	"github.com/b4t3ou/sqs-reader-service/internal/db"
	"github.com/b4t3ou/sqs-reader-service/internal/handler"
	"github.com/b4t3ou/sqs-reader-service/internal/queue"
)

func main() {
	cfgPath, _ := filepath.Abs("config")
	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	sqsClient, err := getSQSClient(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create SQS client")
	}

	dbClient, err := getDBClient(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create Dynamo client")
	}

	sqsHandler := handler.NewSQS(sqsClient, dbClient)

	log.Info().Str("queue", cfg.QueueName).Msg("starting SQS handler")
	if err := sqsHandler.Run(); err != nil {
		log.Fatal().Err(err).Msg("failed to receive message")
	}
}

func getSQSClient(cfg *config.Config) (*queue.SQSClient, error) {
	if cfg.Env == "local" {
		return queue.NewSQSClient(cfg.QueueName, queue.WithSQSLocalstackSession())
	}

	return queue.NewSQSClient(cfg.QueueName)
}

func getDBClient(cfg *config.Config) (*db.Dynamo, error) {
	if cfg.Env == "local" {
		return db.NewDynamo(cfg.DynamoTable, db.WithDynamoLocalstackSession())
	}

	return db.NewDynamo(cfg.QueueName)
}
