package runner

import (
	"DigitalCurrency/internal/blockchain"
	"DigitalCurrency/internal/config"
	"DigitalCurrency/internal/logger"
	"context"
)

type PolygonRunner struct {
	EthRunner
}

func NewPolygonRunner(ctx context.Context, conf *config.EthConfig) *PolygonRunner {
	client := blockchain.NewPolygon(ctx, conf)

	return &PolygonRunner{
		EthRunner: EthRunner{
			client:      &client.Blockchain,
			LastBlockId: 0,
			chBlock:     make(chan int, 20),
			workerCount: 2,
			conf:        conf,
			logger:      logger.PolygonLogger,
		},
	}
}
