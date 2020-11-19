deps:
	go mod verify

lint:
	golangci-lint run

test-unit:
	go test -cover -failfast ./...

test-all:
	go test --tags=integration -cover -failfast ./...

