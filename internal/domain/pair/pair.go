package pair

import (
	"1inch-pair-testtask/internal/domain"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type BalancesFetcher interface {
	GetBalances(
		ctx context.Context,
		poolAddr common.Address,
	) (b0 *big.Int, b1 *big.Int, err error)
}

type AddressesFetcher interface {
	GetTokenPair(
		ctx context.Context,
		poolAddr common.Address,
	) (token0 common.Address, token1 common.Address, err error)
}

type TokenPair struct {
	m map[common.Address]*big.Int
}

func (t *TokenPair) GetToken(address common.Address) (*big.Int, error) {
	v, ok := t.m[address]
	if !ok {
		return nil, fmt.Errorf("token (%s) not found in pool by provided address, err: %w", address.String(), domain.ErrWrongToken)
	}

	return v, nil
}

func NewTokenPair(token0, token1 common.Address, token0amount *big.Int, token1amount *big.Int) *TokenPair {
	return &TokenPair{
		m: map[common.Address]*big.Int{
			token0: token0amount,
			token1: token1amount,
		},
	}
}

func (t *TokenPair) CalculateAmountOut(
	inputToken,
	outputToken common.Address,
	inputAmount *big.Int,
) (*big.Int, error) {
	if cmp := inputAmount.Cmp(big.NewInt(0)); cmp == -1 || cmp == 0 {
		return nil, fmt.Errorf("insufficient input amount, err: %w", domain.ErrWrongAmount)
	}

	inputTokenBalance, err := t.GetToken(inputToken)
	if err != nil {
		return nil, err
	}

	outputTokenBalance, err := t.GetToken(outputToken)
	if err != nil {
		return nil, err
	}

	if cmp := inputTokenBalance.Cmp(big.NewInt(0)); cmp == -1 || cmp == 0 {
		return nil, fmt.Errorf("input token balance is zero, err: %w", domain.ErrWrongAmount)
	}

	if cmp := outputTokenBalance.Cmp(big.NewInt(0)); cmp == -1 || cmp == 0 {
		return nil, fmt.Errorf("output token balance is zero, err: %w", domain.ErrWrongAmount)
	}

	inputAmountWithFee := big.NewInt(0)
	inputAmountWithFee.Mul(inputAmount, big.NewInt(997))

	numerator := big.NewInt(0)
	numerator = numerator.Mul(inputAmountWithFee, outputTokenBalance)

	denominator := big.NewInt(0)
	denominator = denominator.Mul(inputTokenBalance, big.NewInt(1000)).Add(denominator, inputAmountWithFee)

	amountOut := big.NewInt(0)
	amountOut = amountOut.Div(numerator, denominator)

	return amountOut, nil
}
