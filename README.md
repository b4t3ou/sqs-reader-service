# sqs-reader-service

Simple service which is reading data from SQS and saving to DynamoDB

## Testing

### Unit test

```shell script
make test-unit
```
### Test all

For running all the unit and integration tests you need docker.
The integration tests connecting to [Localstack](https://github.com/localstack/localstack) which is an AWS service emulator

Start localstack with docker-compose
```shell script
docker-compose up
```

Create SQS queue and Dynamo table
```shell script
scripts/localstack.sh
```

Run the tests
```shell script
make test-all
```

## Run the service locally

Start localstack with docker-compose
```shell script
docker-compose up
```

Start the service
```shell script
go run main.go
```

## Run the service against AWS

* you need a proper AWS config in your machine
* create the env config in the config folder `my_env.yaml`
* create an SQS queue on AWS and update your config with the queue name
* create a Dynamo table on AWS and update your config with the queue name (HASH = id field)
* export your ENV `export ENV=my_env`

Start the service
```shell script
go run main.go
```
