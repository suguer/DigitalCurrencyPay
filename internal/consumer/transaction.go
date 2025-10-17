package consumer

import (
	"DigitalCurrency/internal/logger"
	"DigitalCurrency/internal/model/dao"
	"DigitalCurrency/internal/model/mdb"
	"DigitalCurrency/internal/service/deposit"
	"DigitalCurrency/internal/service/wallet"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type TransactionMessage struct {
	Hash          string    `json:"hash"`
	TransactionId uint      `json:"transaction_id"`
	FromAddress   string    `json:"from_address"`
	BlockHeight   uint64    `json:"block_height"`
	ConfirmedAt   time.Time `json:"confirmed_at"` // 确认时间
}

type TransactionConsumers struct {
	transactionConsumerChan     chan TransactionMessage
	transactionConsumerQueueKey string
	queueType                   string
	ctx                         context.Context
}

var TransactionConsumersInstance *TransactionConsumers

func NewTransactionConsumers(ctx context.Context, queueType string) *TransactionConsumers {
	TransactionConsumersInstance = &TransactionConsumers{
		queueType: queueType,
		ctx:       ctx,
	}
	switch queueType {
	case "chan":
		TransactionConsumersInstance.transactionConsumerChan = make(chan TransactionMessage, 10)
	}
	return TransactionConsumersInstance
}

func (t *TransactionConsumers) Consume() {
	switch t.queueType {
	case "chan":
		for message := range t.transactionConsumerChan {
			go t.processTransactionConsumerHandle(message)
		}
	case "redis":
		for {
			result, err := dao.Rdb.BRPop(t.ctx, 0, t.transactionConsumerQueueKey).Result()
			if err != nil {
				continue
			}
			message := result[1]
			var transactionMessage TransactionMessage
			err = json.Unmarshal([]byte(message), &transactionMessage)
			if err != nil {
				continue
			}
			go t.processTransactionConsumerHandle(transactionMessage)
		}
	}

}

func (t *TransactionConsumers) Producer(message TransactionMessage) {
	switch t.queueType {
	case "chan":
		t.transactionConsumerChan <- message
	case "redis":
		body, _ := json.Marshal(message)
		dao.Rdb.LPush(t.ctx, t.transactionConsumerQueueKey, body).Result()
	}
}

func (t *TransactionConsumers) processTransactionConsumerHandle(transactionMessage TransactionMessage) {
	fmt.Printf("transactionMessage: %+v\n", transactionMessage)
	logger.Logger.Info("transaction consumer handle", zap.Uint("transaction_id", transactionMessage.TransactionId), zap.Any("transactionMessage", transactionMessage))
	var transactionInstance mdb.Transaction
	if err := dao.Mdb.Where("id = ?", transactionMessage.TransactionId).First(&transactionInstance).Error; err != nil {
		return
	}
	// 生效财务
	transactionInstance.Status = 1
	transactionInstance.Hash = transactionMessage.Hash
	transactionInstance.ConfirmedAt = &transactionMessage.ConfirmedAt
	transactionInstance.FromAddress = transactionMessage.FromAddress
	transactionInstance.BlockHeight = transactionMessage.BlockHeight
	dao.Mdb.Save(&transactionInstance)
	if transactionInstance.Amount < 0 {

	} else {
		//释放钱包地址
		wallet.Release(transactionInstance.ToAddress)
		//增加余额
		deposit.Increment(transactionInstance.Chain, transactionInstance.ContractAddress, transactionInstance.Amount, transactionInstance.UserId)

	}
	if transactionInstance.CallbackUrl != "" {
		// 回调通知
		CallbackConsumersInstance.Producer(CallbackMessage{
			TransactionId: transactionInstance.ID,
		})
	}

}
