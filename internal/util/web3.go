package util

import (
	"math/big"
	"strings"

	"github.com/shopspring/decimal"
)

func EthContractTransferDataEncode(address string, amount *big.Int) string {
	addr := PadLeft(address[2:], "0", 64)
	addr = strings.ToLower(addr)
	intAmount := PadLeft(amount.Text(16), "0", 64)

	return "0xa9059cbb" + addr + intAmount
}
func ToWei(iamount interface{}, decimals int) *big.Int {
	amount := decimal.NewFromFloat(0)
	switch v := iamount.(type) {
	case string:
		amount, _ = decimal.NewFromString(v)
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	result := amount.Mul(mul)

	wei := new(big.Int)
	wei.SetString(result.String(), 10)

	return wei
}
