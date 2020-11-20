#!/usr/bin/env bash

aws sqs create-queue \
--endpoint-url=http://localhost:4566 \
--queue-name local-eu-west-1-test-queue

# Create dynamo table
aws dynamodb create-table \
--endpoint-url=http://localhost:4566 \
--region eu-west-1 \
--table-name local-eu-west-1-events \
--attribute-definitions AttributeName=id,AttributeType=S \
--key-schema AttributeName=id,KeyType=HASH \
--provisioned-throughput ReadCapacityUnits=10,WriteCapacityUnits=10 \

