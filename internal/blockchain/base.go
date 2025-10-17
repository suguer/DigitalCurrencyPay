package blockchain

import (
	"DigitalCurrency/internal/config"
	"context"
)

type Base struct {
	Blockchain
}

func NewBase(ctx context.Context, conf *config.EthConfig) *Base {
	return &Base{
		Blockchain: *NewBlockchain(ctx, conf),
	}
}
