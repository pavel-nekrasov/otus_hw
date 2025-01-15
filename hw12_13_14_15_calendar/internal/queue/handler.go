package queue

import (
	"github.com/streadway/amqp"
)

type (
	DownstreamHandler func(data []byte) error
	Handler           struct {
		handler DownstreamHandler
	}
)

func NewHandler(handler DownstreamHandler) *Handler {
	return &Handler{handler: handler}
}

func (h *Handler) Handle(message amqp.Delivery) {
	if err := h.handler(message.Body); err != nil {
		message.Ack(false)
		return
	}
	message.Ack(true)
}
