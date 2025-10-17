package blockchain

import (
	"DigitalCurrency/internal/config"
	"context"
)

type Matic struct {
	Blockchain
}

func NewMatic(ctx context.Context, conf *config.EthConfig) *Matic {
	return &Matic{
		Blockchain: *NewBlockchain(ctx, conf),
	}
}
