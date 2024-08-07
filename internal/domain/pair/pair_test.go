package pair

import (
	"1inch-pair-testtask/internal/domain"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

type testcase struct {
	testName          string
	inputToken        common.Address
	outputToken       common.Address
	inputAmount       *big.Int
	expectedAmountOut *big.Int
	expectErr         error
}

func TestTokenPair_CalculateAmountOut(t *testing.T) {
	correctRes := big.NewInt(0)
	correctRes.SetString("2361854408", 10)

	testCases := []testcase{
		{
			testName:          "success case",
			inputToken:        common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
			outputToken:       common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
			inputAmount:       big.NewInt(1e18),
			expectedAmountOut: correctRes,
			expectErr:         nil,
		},
		{
			testName:          "not valid address",
			inputToken:        common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc4"),
			outputToken:       common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
			inputAmount:       big.NewInt(1e18),
			expectedAmountOut: correctRes,
			expectErr:         domain.ErrWrongToken,
		},
		{
			testName:          "not valid amount",
			inputToken:        common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
			outputToken:       common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
			inputAmount:       big.NewInt(0),
			expectedAmountOut: correctRes,
			expectErr:         domain.ErrWrongAmount,
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			// Prepare TokenPair object (you might use mocks for GetToken method)
			tp := getTp()

			amountOut, err := tp.CalculateAmountOut(tc.inputToken, tc.outputToken, tc.inputAmount)
			if tc.expectErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedAmountOut, amountOut)
			}
		})
	}
}

func getTp() *TokenPair {
	m := make(map[common.Address]*big.Int)
	b1 := big.NewInt(0)
	b1.SetString("43045503496330", 10)

	b2 := big.NewInt(0)
	b2.SetString("18169626400042346564892", 10)

	m[common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")] = b1
	m[common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")] = b2

	tp := TokenPair{
		m: m,
	}

	return &tp
}
