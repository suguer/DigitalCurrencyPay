package runner

import (
	"DigitalCurrency/internal/blockchain"
	"DigitalCurrency/internal/config"
	"DigitalCurrency/internal/logger"
	"context"
)

type OpRunner struct {
	EthRunner
}

func NewOpRunner(ctx context.Context, conf *config.EthConfig) *OpRunner {
	client := blockchain.NewOp(ctx, conf)

	return &OpRunner{
		EthRunner: EthRunner{
			client:      &client.Blockchain,
			LastBlockId: 0,
			chBlock:     make(chan int, 20),
			workerCount: 2,
			conf:        conf,
			logger:      logger.OpLogger,
		},
	}
}
