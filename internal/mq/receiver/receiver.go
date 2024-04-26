package receiver

import "context"

type MqReceiver interface {
	// ReceiveSeckillOrder 循环消费秒杀队列的订单消息并进行异步下单（插入Order到 DB 中）
	ReceiveSeckillOrder(ctx context.Context)
}
