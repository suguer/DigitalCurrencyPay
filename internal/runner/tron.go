package runner

import (
	"DigitalCurrency/internal/blockchain"
	"DigitalCurrency/internal/config"
	"DigitalCurrency/internal/constant"
	"DigitalCurrency/internal/consumer"
	"DigitalCurrency/internal/logger"
	"DigitalCurrency/internal/model/cache"
	"context"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"go.uber.org/zap"
)

type TronRunner struct {
	client      *blockchain.Tron
	LastBlockId int
	chBlock     chan int
	workerCount int
	conf        *config.EthConfig
}

func NewTronRunner(ctx context.Context, conf *config.EthConfig) *TronRunner {
	tron := TronRunner{
		client:      blockchain.NewTron(ctx, conf),
		chBlock:     make(chan int, 20),
		workerCount: 3,
		conf:        conf,
	}
	return &tron
}
func (t *TronRunner) Start(StartId int) {
	for i := 0; i < t.workerCount; i++ {
		go t.Handler()
	}
	var BlockId int
	var err error
	if StartId > 0 {
		t.LastBlockId = StartId
		t.chBlock <- StartId
	} else {
		for {
			BlockId, _, err = t.client.GetNowBlockId()
			fmt.Printf("[%v]BlockId: %v\n", t.conf.Name, BlockId)
			if err != nil {
				logger.ErrorLogger.Error("波场获取最新区块失败 error", zap.Error(err), zap.String("chain", t.conf.Name))
				time.Sleep(15 * time.Second)
				continue
			}
			if t.LastBlockId == 0 {
				t.LastBlockId = BlockId - 100
			}

			for i := t.LastBlockId; i < BlockId; i++ {
				t.chBlock <- i
			}
			t.LastBlockId = BlockId
			time.Sleep(10 * time.Second)
		}
	}
}

func (t *TronRunner) Handler() {
	for blocksnum := range t.chBlock {
		// fmt.Printf("blocksnum: %v\n", blocksnum)
		data, err := t.client.GetBlockByNum(blocksnum)
		if err != nil {
			logger.ErrorLogger.Error("波场获取区块信息失败 error", zap.Error(err), zap.String("chain", t.conf.Name), zap.Int("blocksnum", blocksnum))
			if strings.Contains(err.Error(), "unexpected end of JSON input") {
				t.chBlock <- blocksnum
			}
			continue

		}
		totalTransactionCount := 0
		activeTransactionCount := 0
		precision := 6
		if t.conf.Precision > 0 {
			precision = t.conf.Precision
		}
		for _, transaction := range data.Transactions {
			out_trade_no := transaction.TxID
			// if out_trade_no != "022d54ee2e14b1265d6fcdf9286e72b2da4863f855aa2d75cd123f7db260d887" {
			// 	continue
			// }
			// fmt.Printf("transaction: %+v\n", transaction)

			status := strings.ToLower(transaction.Ret[0].ContractRet)
			if status != "success" {
				continue
			}
			contractAddress := transaction.RawData.Contract[0].Parameter.Value.ContractAddress
			toAddress := ""
			var transactionAmount float64
			if contractAddress == "" {
				//trx转账
				if transaction.RawData.Contract[0].Type != "TransferContract" {
					continue
				}
				toAddress = transaction.RawData.Contract[0].Parameter.Value.ToAddress
				amountInt := transaction.RawData.Contract[0].Parameter.Value.Amount
				transactionAmount = float64(amountInt) / 1000000
			} else {
				contractAddr := address.HexToAddress(contractAddress)
				contractAddress = contractAddr.Hex()
				if strings.HasPrefix(contractAddress, "41") {
					contractAddress = "0x" + contractAddress
				}
				if len(transaction.RawData.Contract[0].Parameter.Value.Data) < 72 {
					continue
				}
				funcValue := transaction.RawData.Contract[0].Parameter.Value.Data[0:8]
				if funcValue != "a9059cbb" {
					continue
				}
				toAddress = strings.TrimLeft(transaction.RawData.Contract[0].Parameter.Value.Data[9:72], "0")
				amountHex := strings.TrimLeft(transaction.RawData.Contract[0].Parameter.Value.Data[72:], "0")
				if amountHex == "" {
					amountHex = "0"
				}
				bigInt := new(big.Int)
				bigInt.SetString(amountHex, 16)
				transactionAmount = float64(bigInt.Uint64()) / math.Pow10(precision)
			}
			if transactionAmount < 0.1 {
				continue
			}
			totalTransactionCount++
			// fmt.Printf("blocksnum:%v,out_trade_no: %v,contact:%v,to:%v,amount:%v\n", blocksnum, out_trade_no, contractAddress, toAddress, transactionAmount)
			transactionInstance, err := cache.TransactionCacheGet(t.client.Config.Name, toAddress, contractAddress, transactionAmount)
			if err != nil {
				continue
			}
			activeTransactionCount++
			consumer.TransactionConsumersInstance.Producer(consumer.TransactionMessage{
				Hash:          out_trade_no,
				BlockHeight:   uint64(blocksnum),
				ConfirmedAt:   time.Now(),
				TransactionId: transactionInstance.ID,
				FromAddress:   transaction.RawData.Contract[0].Parameter.Value.OwnerAddress,
			})

		}
		if t.conf.Name == constant.ChainTronShasta {
			logger.TronShastaLogger.Info("tron transaction handler", zap.String("chain", t.conf.Name), zap.Int("blocksnum", blocksnum), zap.Int("totalTransactionCount", totalTransactionCount), zap.Int("activeTransactionCount", activeTransactionCount))
		} else {
			logger.TronLogger.Info("tron transaction handler", zap.String("chain", t.conf.Name), zap.Int("blocksnum", blocksnum), zap.Int("totalTransactionCount", totalTransactionCount), zap.Int("activeTransactionCount", activeTransactionCount))
		}
	}
}
