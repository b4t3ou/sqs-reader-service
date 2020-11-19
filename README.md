# sqs-reader-service

Simple service which is reading data from SQS and saving to DynamoDB

## Testing

### Unit test

```shell script
make test-unit
```
### Test all

For running all the unit and integration tests you need docker.
The integration tests connection to localstacl which is an AWS service emulator

Start localstack with docker-compose
```shell script
docker-compose up
```

Create SQS queue
```shell script
scripts/localstack.sh
```

Run the tests
```shell script
test-all
```
