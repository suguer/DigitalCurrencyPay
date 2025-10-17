package runner

import (
	"DigitalCurrency/internal/blockchain"
	"DigitalCurrency/internal/config"
	"DigitalCurrency/internal/logger"
	"context"
)

type EthereumRunner struct {
	EthRunner
}

func NewEthereumRunner(ctx context.Context, conf *config.EthConfig) *EthereumRunner {
	client := blockchain.NewEthereum(ctx, conf)

	return &EthereumRunner{
		EthRunner: EthRunner{
			client:      &client.Blockchain,
			LastBlockId: 0,
			chBlock:     make(chan int, 20),
			workerCount: 2,
			conf:        conf,
			logger:      logger.EthereumLogger,
		},
	}
}
