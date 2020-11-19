deps:
	go mod verify

lint:
	golangci-lint run

test: lint
	go test -cover -failfast ./...

