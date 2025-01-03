package queue

import (
	"context"
	"fmt"

	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/config"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

type Consumer struct {
	connection   *Connection
	exchangeName string
	queueName    string
	routingKey   string
	tag          string
	channel      *amqp.Channel
	handler      Handler
}

func NewConsumer(
	conn *Connection,
	conf config.QueueConsumerConf,
	handler Handler,
) *Consumer {
	return &Consumer{
		connection:   conn,
		exchangeName: conf.Exchange,
		queueName:    conf.Queue,
		routingKey:   conf.RoutingKey,
		tag:          uuid.NewV4().String(),
		handler:      handler,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	var err error

	c.channel, err = c.connection.NewChannel()
	if err != nil {
		return fmt.Errorf("channel: %w", err)
	}

	queue, err := c.channel.QueueDeclare(
		c.queueName, // name of the queue
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // noWait
		nil,         // arguments
	)
	if err != nil {
		return fmt.Errorf("queue declare: %w", err)
	}

	if err = c.channel.QueueBind(
		queue.Name,     // name of the queue
		c.routingKey,   // bindingKey
		c.exchangeName, // sourceExchange
		false,          // noWait
		nil,            // arguments
	); err != nil {
		return fmt.Errorf("queue bind: %w", err)
	}

	message, err := c.channel.Consume(queue.Name, c.tag, false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("queue declare: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-message:
			if !ok {
				return nil
			}
			c.handler.Handle(msg)
		}
	}
}

func (c *Consumer) Close() error {
	if c.channel != nil {
		if err := c.channel.Cancel(c.tag, true); err != nil {
			return fmt.Errorf("amqp channel close failed: %w", err)
		}
	}
	return nil
}
