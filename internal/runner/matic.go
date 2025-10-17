package runner

import (
	"DigitalCurrency/internal/blockchain"
	"DigitalCurrency/internal/config"
	"DigitalCurrency/internal/logger"
	"context"
)

type MaticRunner struct {
	EthRunner
}

func NewMaticRunner(ctx context.Context, conf *config.EthConfig) *MaticRunner {
	client := blockchain.NewMatic(ctx, conf)

	return &MaticRunner{
		EthRunner: EthRunner{
			client:      &client.Blockchain,
			LastBlockId: 0,
			chBlock:     make(chan int, 20),
			workerCount: 2,
			conf:        conf,
			logger:      logger.MaticLogger,
		},
	}
}
