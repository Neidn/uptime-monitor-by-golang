#!make
include .env

TEST_OUTPUT = test/output/coverage.out

run:
	export GITHUB_TOKEN=$(GH_PAT) && go run ./...

test:
	go test -v -coverprofile=$(TEST_OUTPUT) && go tool cover -html=$(TEST_OUTPUT)

.PHONY: run test
