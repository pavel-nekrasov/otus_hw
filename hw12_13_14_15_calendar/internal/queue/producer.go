package queue

import (
	"fmt"

	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/streadway/amqp"
)

type Producer struct {
	connection   *Connection
	exchangeName string
	routingKey   string
	channel      *amqp.Channel
}

func NewProducer(conn *Connection, conf config.QueueProducerConf) *Producer {
	return &Producer{connection: conn, exchangeName: conf.Exchange, routingKey: conf.RoutingKey}
}

func (p *Producer) Start() error {
	var err error

	p.channel, err = p.connection.NewChannel()
	if err != nil {
		return fmt.Errorf("channel: %w", err)
	}

	if err = p.channel.ExchangeDeclare(p.exchangeName, "direct", true, false, false, false, nil); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	return nil
}

func (p *Producer) Publish(data []byte) error {
	if err := p.channel.Publish(
		p.exchangeName, // publish to an exchange
		p.routingKey,   // routing to 0 or more queues
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            data,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
		},
	); err != nil {
		return fmt.Errorf("exchange Publish: %w", err)
	}
	return nil
}

func (p *Producer) Close() error {
	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			return fmt.Errorf("amqp channel close failed: %w", err)
		}
	}
	return nil
}
