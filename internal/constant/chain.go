package constant

import "errors"

const (
	ChainTron       = "tron"
	ChainTronShasta = "tron-shasta"
	ChainArbitrum   = "arbitrum"
	ChainAvax       = "avax"
	ChainMatic      = "matic"
	ChainBase       = "base"
	ChainOp         = "op"
	ChainEthereum   = "ethereum"
)

const (
	CurrencyUSDT = "usdt"
	CurrencyUSDC = "usdc"
)

func GetContractAddress(chain, short string) (string, error) {
	switch chain {

	case ChainTron:
		switch short {
		case CurrencyUSDT:
			return "0x41a614f803b6fd780986a42c78ec9c7f77e6ded13c", nil
		}
	case ChainArbitrum:
		switch short {
		case CurrencyUSDC:
			return "0xaf88d065e77c8cc2239327c5edb3a432268e5831", nil
		}

	case ChainAvax:
		switch short {
		case CurrencyUSDT:
			return "0xB97EF9Ef8734C71904D8002F8b6Bc66Dd9c48a6E", nil
		}
	case ChainMatic:
		switch short {
		case CurrencyUSDT:
			return "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359", nil
		}
	case ChainBase:
		switch short {
		case CurrencyUSDT:
			return "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913", nil
		}
	case ChainOp:
		switch short {
		case CurrencyUSDT:
			return "0x0b2C639c533813f4Aa9D7837CAf62653d097Ff85", nil
		}
	}
	return short, errors.New("chain not support")
}
