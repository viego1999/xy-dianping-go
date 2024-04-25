package receiver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"sync"
	"trpc.group/trpc-go/trpc-go/log"
	"xy-dianping-go/internal/constants"
	"xy-dianping-go/internal/models"
	"xy-dianping-go/internal/service"
)

type RabbitMqReceiver struct {
	conn                *amqp.Connection
	voucherOrderService service.VoucherOrderService
}

func NewMqReceiver(conn *amqp.Connection, orderService service.VoucherOrderService) MqReceiver {
	return &RabbitMqReceiver{conn, orderService}
}

func (receiver *RabbitMqReceiver) ReceiveSeckillOrder(ctx context.Context) {
	log.Info("Start listening for seckill voucher order queue messages.")
	ch, err := receiver.conn.Channel()
	if err != nil {
		panic(fmt.Sprintf("Failed to open a channel: %+v", err))
	}
	defer func(ch *amqp.Channel) {
		if err := ch.Close(); err != nil {
			panic(fmt.Sprintf("Failed to close a channel: %+v", err))
		}
	}(ch)

	// 是指预取计数
	err = ch.Qos(1, 0, false)
	if err != nil {
		panic(fmt.Sprintf("Qos error: %+v", err))
	}

	msgs, err := ch.Consume(
		constants.SECKILL_QUEUE,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(fmt.Sprintf("chaneel consume error: %+v", err))
	}
	var wg sync.WaitGroup
	for msg := range msgs { // 会一直运行监听 channel 上的消息，直到 channel 被关闭或者所在的 goroutine 停止
		wg.Add(1)
		go func(wg *sync.WaitGroup, ch *amqp.Channel, delivery amqp.Delivery) {
			defer wg.Done()
			log.Debugf("Received order message: %s", delivery.Body)

			// 解析消息
			var order models.VoucherOrder
			if err = json.Unmarshal(delivery.Body, &order); err != nil {
				panic(fmt.Sprintf("Failed to unmarshal order: %s", err))
			}

			// 异步创建协程进行订单创建
			if err = receiver.voucherOrderService.CreateVoucherOrder(ctx, &order); err != nil {
				log.Warnf("Order processing failed, retrying: %s", err)
				// 双重逻辑
				if err = receiver.voucherOrderService.CreateVoucherOrder(ctx, &order); err != nil {
					log.Errorf("Order processing failed after retry: %s", err)
					// todo: 第二次处理失败，则更改 Redis 中的数据（也可以将消息放入异常订单数据库或队列中特殊处理）-如回滚库存等操作。。。

				}
			}

			// 手动确认消息处理完成
			if err = ch.Ack(delivery.DeliveryTag, false); err != nil {
				log.Errorf("Failed to ack message: %s", err)
			}
		}(&wg, ch, msg)
	}
	wg.Wait()
}
