#!make
include .env

TEST_OUTPUT = test/output/coverage.out
BUILD_OUTPUT = build

run:
	export GITHUB_TOKEN=$(GH_PAT) GITHUB_REPOSITORY=Neidn/uptime && go run ./...

test:
	go test ./... -v -coverprofile=$(TEST_OUTPUT) && go tool cover -html=$(TEST_OUTPUT)

build:
	go build -o $(BUILD_OUTPUT) ./...

build_run:
	export GITHUB_TOKEN=$(GH_PAT) GITHUB_REPOSITORY=Neidn/uptime && ./$(BUILD_OUTPUT)/monitor

.PHONY: run test build build_run
