package runner

import (
	"DigitalCurrency/internal/blockchain"
	"DigitalCurrency/internal/config"
	"DigitalCurrency/internal/consumer"
	"DigitalCurrency/internal/model/cache"
	"fmt"
	"math"
	"strings"
	"time"

	"go.uber.org/zap"
)

type EthRunner struct {
	client      *blockchain.Blockchain
	LastBlockId int
	chBlock     chan int
	workerCount int
	conf        *config.EthConfig
	logger      *zap.Logger
}

func NewEthRunner(client *blockchain.Blockchain, conf *config.EthConfig, logger *zap.Logger) *EthRunner {
	return &EthRunner{
		client:      client,
		LastBlockId: 0,
		chBlock:     make(chan int, 20),
		workerCount: 2,
		conf:        conf,
		logger:      logger,
	}
}

func (t *EthRunner) Start(StartId int) {
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
			BlockId -= 15
			fmt.Printf("[%v]BlockId: %v\n", t.conf.Name, BlockId)
			if err != nil {
				t.logger.Error("获取最新区块失败 error", zap.Error(err), zap.String("chain", t.conf.Name))
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
func (t *EthRunner) Handler() {
	for blocksnum := range t.chBlock {
		elogs, err := t.client.GetLogsFormat(blocksnum, blocksnum)
		if err != nil {
			t.logger.Error("获取区块信息失败 error", zap.Error(err), zap.String("chain", t.conf.Name))
			go func() {
				if strings.Contains(err.Error(), "Too Many Requests") {
					t.chBlock <- blocksnum
				}
			}()

			continue
		}
		activeCount := 0
		for _, log := range elogs {
			Amount := float64(log.Amount) / math.Pow10(6)
			transactionInstance, err := cache.TransactionCacheGet(t.client.Config.Name, log.ToAddress, log.ContractAddress, Amount)
			if err != nil {
				continue
			}
			activeCount++
			if consumer.TransactionConsumersInstance != nil {
				consumer.TransactionConsumersInstance.Producer(consumer.TransactionMessage{
					Hash:          log.TxHash,
					BlockHeight:   uint64(blocksnum),
					ConfirmedAt:   time.Now(),
					TransactionId: transactionInstance.ID,
					FromAddress:   log.FromAddress,
				})
			}
		}
		t.logger.Info("transaction handler", zap.String("chain", t.conf.Name), zap.Int("blocksnum", blocksnum), zap.Int("logsnum", len(elogs)), zap.Int("activeCount", activeCount))
	}
}
