package blockchain

import (
	"DigitalCurrency/internal/config"
	"context"
)

type Avax struct {
	Blockchain
}

func NewAvax(ctx context.Context, conf *config.EthConfig) *Avax {
	return &Avax{
		Blockchain: *NewBlockchain(ctx, conf),
	}
}
