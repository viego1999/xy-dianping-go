package receiver

import "context"

type MqReceiver interface {
	ReceiveSeckillOrder(ctx context.Context)
}
