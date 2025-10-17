package blockchain

import (
	"DigitalCurrency/internal/config"
	"context"
)

type Ethereum struct {
	Blockchain
}

func NewEthereum(ctx context.Context, conf *config.EthConfig) *Ethereum {
	return &Ethereum{
		Blockchain: *NewBlockchain(ctx, conf),
	}
}
