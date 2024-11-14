package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sera_backend/internal/config/logger"
)

func (rbm *rbm_pool) SenderRb(ctx context.Context, exchange_name string, msg *Message) error {

	if rbm.channel == nil {
		logger.Info("RMB.CHANNEL E NULL")
	}
	if ctx == nil {
		logger.Info("ctx E NULL")
	}
	err := rbm.channel.PublishWithContext(ctx,
		"amq.direct", // exchange amq.direct
		"",           // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			Body:        msg.Data,
			ContentType: msg.ContentType,
		})

	if err != nil {
		logger.Error("DEU MERDA AQUI NA PUPLICACAO", err)
	}
	logger.Info("MENSAGEM ENVIADA")
	return nil
}
