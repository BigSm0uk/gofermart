ifneq ($(wildcard .env),)
include .env
export
else
$(warning WARNING: .env file not found! Using .env.example)
include .env.example
export
endif

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

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
	go build -o ./cmd/accrual/accrual ./cmd/accrual && \
	./cmd/accrual/accrual
.PHONY: run-acc	

test: ### run test
	go test -v -race -covermode atomic -coverprofile=coverage.txt ./internal/...
.PHONY: test

integration-test: ### run integration-test !!!IMPLEMENT!!!
	go clean -testcache && go test -v ./integration-test/...
.PHONY: integration-test

bin-deps: ### install tools
	go install tool
.PHONY: bin-deps

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

mock: ### run mockgen !!!IMPLEMENT!!!
.PHONY: mock

migrate-up: ### migration up !!!IMPLEMENT!!!
.PHONY: migrate-up

migrate-down: ### migration down !!!IMPLEMENT!!!
.PHONY: migrate-down


pre-commit: swag-v1 mock format linter-golangci test ### run pre-commit
.PHONY: pre-commit