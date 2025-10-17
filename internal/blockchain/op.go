package blockchain

import (
	"DigitalCurrency/internal/config"
	"context"
)

type Op struct {
	Blockchain
}

func NewOp(ctx context.Context, conf *config.EthConfig) *Op {
	return &Op{
		Blockchain: *NewBlockchain(ctx, conf),
	}
}
