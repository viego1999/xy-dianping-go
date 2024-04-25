package mq

import (
	"fmt"
	"github.com/streadway/amqp"
	"trpc.group/trpc-go/trpc-go/log"
	"xy-dianping-go/internal/config"
	"xy-dianping-go/internal/constants"
)

var (
	AMQPConn *amqp.Connection
)

func InitAMQPConnection() *amqp.Connection {
	conn, err := amqp.Dial(config.AMQPURI)
	if err != nil {
		panic(err)
	}
	AMQPConn = conn
	ch, err := conn.Channel()
	if err != nil {
		panic(fmt.Sprintf("Failed to open a channel: %+v", err))
	}
	defer func(ch *amqp.Channel) {
		if err := ch.Close(); err != nil {
			panic(fmt.Sprintf("Failed to close a channel: %+v", err))
		}
	}(ch)

	// 声明交换机
	if err = ch.ExchangeDeclare(constants.SECKILL_EXCHANGE, "direct", true, false, false, false, nil); err != nil {
		panic(fmt.Sprintf("Failed to decalre an excahnge: %+v", err))
	}
	// 声明队列
	queue, err := ch.QueueDeclare(
		constants.SECKILL_QUEUE, // 队列名称
		true,                    // 持久化
		false,                   // 自动删除
		false,                   // 排他性
		false,                   // 不等待
		nil,                     // 其他参数
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to declare a queue: %+v", err))
	}
	// 绑定队列到交换机
	if err = ch.QueueBind(queue.Name, constants.SECKILL_ROUTING_KEY, constants.SECKILL_EXCHANGE, false, nil); err != nil {
		panic(fmt.Sprintf("Failed to bind a queue: %+v", err))
	}
	log.Info("AMQP Conn initialization completed!")
	return AMQPConn
}
