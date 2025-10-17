package blockchain

import (
	"DigitalCurrency/internal/blockchain/model"
	"DigitalCurrency/internal/config"
	"DigitalCurrency/internal/util"
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

type Arbitrum struct {
	Blockchain
}

func NewArbitrum(ctx context.Context, conf *config.EthConfig) *Arbitrum {
	return &Arbitrum{
		Blockchain: *NewBlockchain(ctx, conf),
	}
}

func (t *Arbitrum) EthGasFeeEstimate(fromAddress, toAddress, contractAddress string, amount uint64) (model.EthGasFee, error) {
	GasLimit := uint64(21000)
	contractAddr := common.HexToAddress(contractAddress)
	toAddr := common.HexToAddress(toAddress)
	intAmount := new(big.Int).SetUint64(amount)

	data := util.EthContractTransferDataEncode(toAddr.Hex(), intAmount)

	msg := ethereum.CallMsg{
		From: common.HexToAddress(fromAddress),
		To:   &contractAddr,
		Data: common.FromHex(data),
	}
	GasLimit, err := t.client.EstimateGas(t.ctx, msg)
	return model.EthGasFee{
		GasLimit:             GasLimit,
		MaxPriorityFeePerGas: 0,
		MaxFeePerGas:         util.ToWei(0.1, 18),
	}, err
}
