RUN_DIR = cmd/monitor
TEST_OUTPUT = test/output/coverage.out

run:
	go run $(RUN_DIR)/main.go

test:
	go test -v -coverprofile=$(TEST_OUTPUT) && go tool cover -html=$(TEST_OUTPUT)

.PHONY: run test
