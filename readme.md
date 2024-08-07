# Build And Run

## Build the project

```make build```

## Run the project

```./build/amountoutcalc --PoolID=<eth_address> --FromToken=<eth_address> --ToToken=<eth_address> --InputAmount=<string number>```

### Available flags

1. PoolID: The address of the pool
2. FromToken: The address of the token to swap from
3. ToToken: The address of the token to swap to
4. InputAmount: The amount of the token to swap from
5. RpcURL (optional): The URL of the RPC node to use
6. debug (optional): Enable debug mode (logs)

## Run the tests

```make test```

## Run test with e2e

```make test-and-e2e```