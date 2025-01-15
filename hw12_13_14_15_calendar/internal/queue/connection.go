package queue

import (
	"fmt"
	"sync"

	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/streadway/amqp"
)

type Connection struct {
	uri      string
	amqpConn *amqp.Connection
	mut      sync.Mutex
}

func NewConnection(conf config.QueueServerConf) *Connection {
	uri := fmt.Sprintf("amqp://%v:%v@%v:%v/", conf.User, conf.Password, conf.Host, conf.Port)

	return &Connection{uri: uri}
}

func (c *Connection) Connect() error {
	var err error
	c.mut.Lock()
	defer c.mut.Unlock()

	c.amqpConn, err = amqp.Dial(c.uri)
	if err != nil {
		return fmt.Errorf("connection failure: %w", err)
	}
	return nil
}

func (c *Connection) NewChannel() (*amqp.Channel, error) {
	return c.amqpConn.Channel()
}

func (c *Connection) Close() error {
	c.mut.Lock()
	defer c.mut.Unlock()

	if c.amqpConn != nil {
		if err := c.amqpConn.Close(); err != nil {
			return fmt.Errorf("amqp connection close failed: %w", err)
		}
	}

	return nil
}
