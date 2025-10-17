package runner

import (
	"DigitalCurrency/internal/blockchain"
	"DigitalCurrency/internal/config"
	"DigitalCurrency/internal/logger"
	"context"
)

type AvaxRunner struct {
	EthRunner
}

func NewAvaxRunner(ctx context.Context, conf *config.EthConfig) *AvaxRunner {
	client := blockchain.NewAvax(ctx, conf)
	return &AvaxRunner{
		EthRunner: EthRunner{
			client:      &client.Blockchain,
			LastBlockId: 0,
			chBlock:     make(chan int, 20),
			workerCount: 2,
			conf:        conf,
			logger:      logger.AvaxLogger,
		},
	}
}
