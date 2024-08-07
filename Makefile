.PHONY: test build test-with-e2e

build:
	go build -o build/uniswapv2paircalc ./cmd/cli/main.go

test:
	go test ./...

test-with-e2e:
	$INTEGRATION_TEST="on" go test ./...