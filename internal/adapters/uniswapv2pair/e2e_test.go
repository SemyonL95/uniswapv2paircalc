package uniswapv2pair

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
)

const RPCURL = "https://eth.llamarpc.com"

var poolAddr = common.HexToAddress("0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc") // eth/usdc pool address
var token0 = common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")   // eth/usdc t0
var token1 = common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")   // eth/usdc t1

func Test_RunBaseE2E(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != "on" {
		t.Skip()
	}

	ethClient, err := ethclient.Dial(RPCURL)
	if err != nil {
		t.Log("failed to create eth client", err)
		t.Fail()

		return
	}

	pairClient := NewClient(slog.Default(), ethClient)

	t0, t1, err := pairClient.GetTokenPair(context.TODO(), poolAddr)
	if err != nil {
		t.Log("failed to get token pair", err)
		t.Fail()
	}

	assert.Equal(t, token0, t0)
	assert.Equal(t, token1, t1)

	b0, b1, err := pairClient.GetBalances(context.TODO(), poolAddr)
	if err != nil {
		t.Log("failed to get balances", err)
		t.Fail()
	}

	t.Log("balances", slog.String("t0", b0.String()), slog.String("t1", b1.String()))
}
