.PHONY: lint test check

POSTGRES_CONTAINER="postgres-${CONTROLLER_NAME}"

lint:
	golangci-lint run --config .golangci.yml
test:
	go test -cover ./...
check: | lint test
