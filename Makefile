ifneq ($(wildcard .env),)
include .env
export
else
$(warning WARNING: .env file not found! Using .env.example)
include .env.example
export
endif

deps: ### deps tidy + verify
	go mod tidy && go mod verify
.PHONY: deps

format: ### Run code formatter
	go fmt ./...
.PHONY: format

swag-v1: ### swag init
	swag init -g internal/app/router/router.go
.PHONY: swag-v1

run-gof: deps ### run gophermart
	go mod download && \
	go run ./cmd/gophermart
.PHONY: run-gof

run-acc: deps ### run accrual
	go mod download && \
	go run ./cmd/accrual
.PHONY: run-acc	

test: ### run test
	go test -v -race -covermode atomic -coverprofile=coverage.txt ./internal/...
.PHONY: test

integration-test: ### run integration-test
	go clean -testcache && go test -v ./integration-test/...
.PHONY: integration-test

bin-deps: ### install tools
	go install tool
.PHONY: bin-deps

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

pre-commit: swag-v1 mock format linter-golangci test ### run pre-commit
.PHONY: pre-commit