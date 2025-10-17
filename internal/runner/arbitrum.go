package runner

import (
	"DigitalCurrency/internal/blockchain"
	"DigitalCurrency/internal/config"
	"DigitalCurrency/internal/logger"
	"context"
)

type ArbitrumRunner struct {
	EthRunner
}

func NewArbitrumRunner(ctx context.Context, conf *config.EthConfig) *ArbitrumRunner {
	client := blockchain.NewArbitrum(ctx, conf)
	return &ArbitrumRunner{
		EthRunner: EthRunner{
			client:      &client.Blockchain,
			LastBlockId: 0,
			chBlock:     make(chan int, 20),
			workerCount: 2,
			conf:        conf,
			logger:      logger.ArbitrumLogger,
		},
	}
}
