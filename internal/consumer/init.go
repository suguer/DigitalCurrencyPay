package consumer

import "context"

func InitConsumers(ctx context.Context, queueType string) {
	go NewTransactionConsumers(ctx, queueType).Consume()
	go NewCallbackConsumers(ctx, queueType).Consume()
}
