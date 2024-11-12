package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	conn      net.Conn
	ctx       context.Context
	cancelCtx context.CancelFunc
	address   string
	timeout   time.Duration
	inSc      *bufio.Scanner
	out       io.Writer
	socketSc  *bufio.Scanner
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		inSc:    bufio.NewScanner(in),
		out:     out,
	}
}

func (c *telnetClient) Connect() error {
	var err error
	c.ctx, c.cancelCtx = context.WithTimeout(context.Background(), c.timeout)
	dialer := &net.Dialer{}
	c.conn, err = dialer.DialContext(c.ctx, "tcp", c.address)
	if err != nil {
		return err
	}
	c.socketSc = bufio.NewScanner(c.conn)
	return nil
}

func (c *telnetClient) Close() error {
	defer c.cancelCtx()
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *telnetClient) Send() error {
	if c.inSc.Scan() {
		_, err := c.conn.Write([]byte(fmt.Sprintln(c.inSc.Text())))
		if err != nil {
			return err
		}
	} else if c.inSc.Err() != nil {
		return c.inSc.Err()
	}
	return nil
}

func (c *telnetClient) Receive() error {
	if c.socketSc.Scan() {
		_, err := c.out.Write([]byte(fmt.Sprintln(c.socketSc.Text())))
		if err != nil {
			return err
		}
	} else if c.socketSc.Err() != nil {
		return c.socketSc.Err()
	}
	return nil
}
