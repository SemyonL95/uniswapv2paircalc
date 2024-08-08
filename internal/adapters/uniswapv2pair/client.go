package uniswapv2pair

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log/slog"
	"math/big"
	"uniswapv2paircalc/internal/adapters/uniswapv2pair/pair"
)

func NewClient(log *slog.Logger, client *ethclient.Client) *Client {
	return &Client{
		log:       log.WithGroup("adapters.chain.client"),
		ethClient: client,
	}
}

type Client struct {
	log *slog.Logger

	ethClient *ethclient.Client
}

func (c *Client) GetTokenPair(ctx context.Context, poolAddr common.Address) (common.Address, common.Address, error) {
	pairClient, err := pair.NewPair(poolAddr, c.ethClient)
	if err != nil {
		c.log.DebugContext(ctx, "failed to create instance of pair contract adapter", slog.String("error", err.Error()))
		return common.Address{}, common.Address{}, fmt.Errorf("failed to create instance of pair contract adapter, err: %w", err)
	}

	token0, err := pairClient.Token0(&bind.CallOpts{Context: ctx})
	if err != nil {
		c.log.DebugContext(ctx, "failed to get token0", slog.String("error", err.Error()))

		return common.Address{}, common.Address{}, fmt.Errorf("failed to get token0, err: %w", err)
	}

	token1, err := pairClient.Token1(&bind.CallOpts{Context: ctx})
	if err != nil {
		c.log.DebugContext(ctx, "failed to get token1", slog.String("error", err.Error()))

		return common.Address{}, common.Address{}, fmt.Errorf("failed to get token1, err: %w", err)
	}

	return token0, token1, nil
}

func (c *Client) GetBalances(ctx context.Context, poolAddr common.Address) (*big.Int, *big.Int, error) {
	pairClient, err := pair.NewPair(poolAddr, c.ethClient)
	if err != nil {
		c.log.DebugContext(ctx, "failed to create instance of pair contract adapter", slog.String("error", err.Error()))

		return nil, nil, err
	}

	res, err := pairClient.GetReserves(&bind.CallOpts{Context: ctx})
	if err != nil {
		c.log.DebugContext(ctx, "failed to get reserves", slog.String("error", err.Error()))

		return nil, nil, fmt.Errorf("failed to get reserves, err: %w", err)
	}

	return res.Reserve0, res.Reserve1, nil
}
