package sender

import "xy-dianping-go/internal/models"

type MqSender interface {
	// SendSeckillMessage 发送秒杀订单消息，这里需要保证可靠传递性，失败重传，消息发送到队列失败，进行消息回退
	//
	// @param order: 秒杀订单信息
	//
	// @param reliable: 是否进行可靠传输模式
	SendSeckillMessage(order *models.VoucherOrder, reliable bool)
}
