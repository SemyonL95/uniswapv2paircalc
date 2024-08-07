package main

import (
	"1inch-pair-testtask/internal/adapters/uniswapv2pair"
	"1inch-pair-testtask/internal/application/swapcalculator"
	"context"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log/slog"
	"math/big"
	"os"
)

const defaultRPCURL = "https://eth.llamarpc.com"

func main() {
	poolAddrParam := flag.String("PoolID", "", "--PoolID=<eth_addr_of_uniswapv2pairpool> string *required")
	inputTokenParam := flag.String("FromToken", "", "--FromToken=<eth_addr_of_input_token> string *required")
	outputTokenParam := flag.String("ToToken", "", "--ToToken=<eth_addr_of_output_token> string *required")
	inputAmountParam := flag.String("InputAmount", "", "--InputAmount=<amount_of_input_token> string *required")
	rpcURLParam := flag.String("RpcURL", defaultRPCURL, "--RpcURL=<rpc_url> string, not required, will use default, if not provided")
	isDebug := flag.Bool("debug", false, "--debug bool, not required, will use default, if not provided")

	flag.Parse()

	var logLevel = slog.LevelInfo
	if *isDebug {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

	slog.SetDefault(logger)

	if *poolAddrParam == "" || *inputTokenParam == "" || *outputTokenParam == "" || *inputAmountParam == "" {
		logger.Error("one of PoolID | FromToken | ToToken | InputAmount required flags are missing")

		return
	}

	ethClient, err := ethclient.Dial(*rpcURLParam)
	if err != nil {
		logger.Error("failed to create eth client", slog.String("err", err.Error()))

		return
	}
	pairClient := uniswapv2pair.NewClient(logger, ethClient)
	calc := swapcalculator.NewService(logger, pairClient)

	runGetAmountOutCommand(calc, *poolAddrParam, *inputTokenParam, *outputTokenParam, *inputAmountParam)

	return
}

func runGetAmountOutCommand(s *swapcalculator.Service, poolAddrParam string, inputTokenParam string, outputTokenParam string, inputAmountParam string) {
	if !common.IsHexAddress(poolAddrParam) {
		slog.Error("PoolID is not a valid address")

		return
	}
	poolAddr := common.HexToAddress(poolAddrParam)

	if !common.IsHexAddress(inputTokenParam) {
		slog.Error("FromToken is not a valid address")

		return
	}
	inputToken := common.HexToAddress(inputTokenParam)

	if !common.IsHexAddress(outputTokenParam) {
		slog.Error("ToToken is not a valid address")

		return
	}
	outputToken := common.HexToAddress(outputTokenParam)

	inputAmount := big.NewInt(0)
	_, ok := inputAmount.SetString(inputAmountParam, 10)
	if !ok {
		slog.Error("InputAmount is not a valid amount")

		return
	}

	outputAmount, err := s.GetAmountOut(context.Background(), poolAddr, inputToken, outputToken, inputAmount)
	if err != nil {
		slog.Error("failed to get amount out", slog.String("err", err.Error()))

		return
	}

	fmt.Printf("The result of swap will be: %s \n", outputAmount.String())

	return
}
