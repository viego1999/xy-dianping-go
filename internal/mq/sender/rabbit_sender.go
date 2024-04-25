package sender

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"trpc.group/trpc-go/trpc-go/log"
	"xy-dianping-go/internal/constants"
	"xy-dianping-go/internal/models"
)

type RabbitMqSender struct {
	conn *amqp.Connection
}

func NewMqSender(conn *amqp.Connection) MqSender {
	return &RabbitMqSender{conn}
}

func (sender *RabbitMqSender) SendSeckillMessage(order *models.VoucherOrder, reliable bool) {
	ch, err := sender.conn.Channel()
	if err != nil {
		panic(fmt.Sprintf("Failed to open a channel: %+v", err))
	}
	defer func(ch *amqp.Channel) {
		if err := ch.Close(); err != nil {
			panic(fmt.Sprintf("Failed to close a channel: %+v", err))
		}
	}(ch)
	// 将 VoucherOrder 序列化为字节
	var body []byte
	if body, err = json.Marshal(order); err != nil {
		panic(fmt.Sprintf("Failed to serialize VoucherOrder: %+v", err))
	}
	// 如果需要进行可靠传输
	if reliable {
		// 开启手动确认机制
		_ = ch.Confirm(false)
		ack, nack := ch.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))

		// 监听消息确认状态
		go func() {
			select {
			case seqNo := <-ack:
				log.Infof("Confirmed received with sequence number: %d", seqNo)
			case seqNo := <-nack:
				log.Errorf("Nack received with sequence number: %d", seqNo)
				// 发送失败，重新发送一次
				if err = ch.Publish(
					constants.SECKILL_EXCHANGE,    // exchange
					constants.SECKILL_ROUTING_KEY, // routing key
					reliable,                      // mandatory: 为 true 时，消息无法路由会返回给发送者
					false,                         // immediate
					amqp.Publishing{
						ContentType: "application/json",
						Body:        body,
					}); err != nil {
					panic(fmt.Sprintf("Failed to publish a message: %+v", err))
				}
			}
		}()

		// 监听返回消息
		returns := ch.NotifyReturn(make(chan amqp.Return, 1))
		go func() {
			for r := range returns {
				log.Infof("Message returned: reply code %d, reply text %s\n", r.ReplyCode, r.ReplyText)
			}
		}()
	}

	// 发送消息
	if err = ch.Publish(
		constants.SECKILL_EXCHANGE,    // exchange
		constants.SECKILL_ROUTING_KEY, // routing key
		reliable,                      // mandatory: 为 true 时，消息无法路由会返回给发送者
		false,                         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		}); err != nil {
		panic(fmt.Sprintf("Failed to publish a message: %+v", err))
	}
	log.Infof(" [x] Sent %s", body)
}
