package runner

import (
	"DigitalCurrency/internal/blockchain"
	"DigitalCurrency/internal/config"
	"DigitalCurrency/internal/logger"
	"context"
)

type BaseRunner struct {
	EthRunner
}

func NewBaseRunner(ctx context.Context, conf *config.EthConfig) *BaseRunner {
	client := blockchain.NewMatic(ctx, conf)

	return &BaseRunner{
		EthRunner: EthRunner{
			client:      &client.Blockchain,
			LastBlockId: 0,
			chBlock:     make(chan int, 20),
			workerCount: 2,
			conf:        conf,
			logger:      logger.BaseLogger,
		},
	}
}
