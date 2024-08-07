.PHONY: test build test-with-e2e

build:
	go build -o build/amountoutcalc ./cmd/cli/main.go

test:
	go test ./...

test-with-e2e:
	$INTEGRATION_TEST="on" go test ./...