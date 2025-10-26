package blockchain

import (
	"DigitalCurrency/internal/config"
	"context"
)

type Polygon struct {
	Blockchain
}

func NewPolygon(ctx context.Context, conf *config.EthConfig) *Polygon {
	return &Polygon{
		Blockchain: *NewBlockchain(ctx, conf),
	}
}
