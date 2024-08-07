package swapcalculator

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"log/slog"
	"math/big"
	"uniswapv2paircalc/internal/domain"
	"uniswapv2paircalc/internal/domain/pair"
)

type PairDataProvider interface {
	pair.AddressesFetcher
	pair.BalancesFetcher
}

func NewService(log *slog.Logger, pairDataProvider PairDataProvider) *Service {
	return &Service{
		log:              log.WithGroup("application.swapcalculator.service"),
		pairDataProvider: pairDataProvider,
	}
}

type Service struct {
	log *slog.Logger

	pairDataProvider PairDataProvider
}

func (s *Service) GetAmountOut(
	ctx context.Context,
	poolAddr,
	inputToken,
	outputToken common.Address,
	inputAmount *big.Int,
) (*big.Int, error) {
	token0, token1, err := s.pairDataProvider.GetTokenPair(ctx, poolAddr)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to get token pair", slog.String("error", err.Error()))

		return nil, err
	}

	if err := s.validateAddresses(inputToken, token0, token1); err != nil {
		s.log.ErrorContext(ctx, "failed to validate input token", slog.String("error", err.Error()))

		return nil, err
	}

	if err := s.validateAddresses(outputToken, token0, token1); err != nil {
		s.log.ErrorContext(ctx, "failed to validate output token", slog.String("error", err.Error()))

		return nil, err
	}

	token0Balance, token1Balance, err := s.pairDataProvider.GetBalances(ctx, poolAddr)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to get balances", slog.String("error", err.Error()))

		return nil, err
	}

	tokenPair := pair.NewTokenPair(token0, token1, token0Balance, token1Balance)

	return tokenPair.CalculateAmountOut(inputToken, outputToken, inputAmount)
}

func (s *Service) validateAddresses(subject, expected0, expected1 common.Address) error {
	if subject != expected0 && subject != expected1 {
		return fmt.Errorf(
			"provided invalid token addr (%s) available in provided pool: (%s), (%s), err: %w",
			subject.String(),
			expected0.String(),
			expected1.String(),
			domain.ErrWrongToken,
		)
	}

	return nil
}
