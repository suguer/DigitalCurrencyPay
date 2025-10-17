package runner

import (
	"DigitalCurrency/internal/config"
	"context"
)

func InitRunner(ctx context.Context) {
	if config.Conf.BlockChain.Tron.Enable == 1 {
		service := NewTronRunner(ctx, &config.Conf.BlockChain.Tron)
		go service.Start(0)
	}
	if config.Conf.BlockChain.TronShasta.Enable == 1 {
		service := NewTronRunner(ctx, &config.Conf.BlockChain.TronShasta)
		go service.Start(0)
	}
	if config.Conf.BlockChain.Avax.Enable == 1 {
		service := NewAvaxRunner(ctx, &config.Conf.BlockChain.Avax)
		go service.Start(0)
	}
	if config.Conf.BlockChain.Arbitrum.Enable == 1 {
		service := NewArbitrumRunner(ctx, &config.Conf.BlockChain.Arbitrum)
		go service.Start(0)
	}
	if config.Conf.BlockChain.Matic.Enable == 1 {
		service := NewMaticRunner(ctx, &config.Conf.BlockChain.Matic)
		go service.Start(0)
	}
	if config.Conf.BlockChain.Base.Enable == 1 {
		service := NewBaseRunner(ctx, &config.Conf.BlockChain.Base)
		go service.Start(0)
	}
	if config.Conf.BlockChain.Op.Enable == 1 {
		service := NewOpRunner(ctx, &config.Conf.BlockChain.Op)
		go service.Start(0)
	}
	if config.Conf.BlockChain.Ethereum.Enable == 1 {
		service := NewEthereumRunner(ctx, &config.Conf.BlockChain.Ethereum)
		go service.Start(0)
	}

}
