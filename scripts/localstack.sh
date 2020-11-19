#!/usr/bin/env bash

aws sqs create-queue \
--endpoint-url=http://localhost:4566 \
--queue-name local-eu-west-1-test-queue
