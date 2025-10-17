package consumer

import (
	"DigitalCurrency/internal/model/dao"
	"DigitalCurrency/internal/service/transaction"
	"context"
	"encoding/json"
	"fmt"
)

type CallbackMessage struct {
	TransactionId uint `json:"transaction_id"`
}

type CallbackConsumers struct {
	consumerChan     chan CallbackMessage
	consumerQueueKey string
	queueType        string
	ctx              context.Context
}

var CallbackConsumersInstance *CallbackConsumers

func NewCallbackConsumers(ctx context.Context, queueType string) *CallbackConsumers {
	CallbackConsumersInstance = &CallbackConsumers{
		queueType: queueType,
		ctx:       ctx,
	}
	switch queueType {
	case "chan":
		CallbackConsumersInstance.consumerChan = make(chan CallbackMessage, 10)
	}
	return CallbackConsumersInstance
}

func (t *CallbackConsumers) Consume() {
	switch t.queueType {
	case "chan":
		for message := range t.consumerChan {
			go t.processHandle(message)
		}
	case "redis":
		for {
			result, err := dao.Rdb.BRPop(t.ctx, 0, t.consumerQueueKey).Result()
			if err != nil {
				continue
			}
			message := result[1]
			var callbackMessage CallbackMessage
			err = json.Unmarshal([]byte(message), &callbackMessage)
			if err != nil {
				continue
			}
			go t.processHandle(callbackMessage)
		}
	}

}

func (t *CallbackConsumers) Producer(message CallbackMessage) {
	switch t.queueType {
	case "chan":
		t.consumerChan <- message
	case "redis":
		body, _ := json.Marshal(message)
		dao.Rdb.LPush(t.ctx, t.consumerQueueKey, body).Result()
	}
}

func (t *CallbackConsumers) processHandle(callbackMessage CallbackMessage) {
	tx, _ := transaction.Instance(callbackMessage.TransactionId)
	fmt.Printf("通知回调callbacktx: %+v\n", tx)
}
