package broker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

var (
	DefaultPrefetchCount = 32
	DefaultQueueName = "sms"
	DefaultQueueArgs = amqp.Table {
		"x-queue-mode": "lazy",
		"x-max-length": 1000,
		"x-overflow": "reject-publish",
	}
)

type Connector interface {
	Close() error
	Reconnect()
	NotifyClose() chan *amqp.Error
	Active() bool
	HoldAlive(ctx context.Context)
}

type Publisher interface {
	Connector
	Publish(data []byte) error
}

type Consumer interface {
	Connector
	Consume() <-chan amqp.Delivery
}

type Connection struct {
	sync.Mutex
	conn      *amqp.Connection
	mqch      *amqp.Channel
	queue     *amqp.Queue
	delivery  <-chan amqp.Delivery
	confirm   <-chan amqp.Confirmation
	addr      string
}

func (c *Connection) Close() error {
	c.Lock()
	defer c.Unlock()

	if err := c.mqch.Close(); err != nil{
		return fmt.Errorf("close broker channel: %w", err)
	}
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("close broker connection: %w", err)
	}

	c.conn = nil
	c.mqch = nil
	c.queue = nil

	return nil
}

func (c *Connection) Active() bool {
	c.Lock()
	defer c.Unlock()

	if c.conn == nil {
		return false
	}

	return !c.conn.IsClosed()
}

func (c *Connection) Reconnect() {
	c.Lock()
	defer c.Unlock()

	if c.delivery != nil {
		asConsumer(c)
	}else{
		asPublisher(c)
	}
}

func (c *Connection) NotifyClose() chan *amqp.Error {
	c.Lock()
	defer c.Unlock()

	receiver := make(chan *amqp.Error)
	return c.conn.NotifyClose(receiver)
}

func (c *Connection) HoldAlive(ctx context.Context) {
	go func() {
		for {
			timer := time.NewTimer(time.Second)
			<-timer.C

			if !c.Active() {
				c.Reconnect()
				continue
			}

			notify := c.NotifyClose()

			select {
			case <-notify:
				c.Reconnect()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (c *Connection) Publish(data []byte) error {
	c.Lock()
	defer c.Unlock()

	if c.conn == nil {
		return fmt.Errorf("publish: amqp connection is nil")
	}
	if c.conn.IsClosed(){
		return fmt.Errorf("publish: amqp connection is closed")
	}

	//err := c.mqch.Tx()
	//if err != nil {
	//	return fmt.Errorf("transaction: %w", err)
	//}

	err := c.mqch.Publish(
		"",     // exchange
		 c.queue.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType:  "application/octet-stream",
			//persistent messages will be restored to durable queues and lost on non-durable queues during server restart.
			DeliveryMode: amqp.Persistent,
			Body:         data,
		},
	)

	if err != nil {
		//c.mqch.TxRollback()
		return fmt.Errorf("publish: %w", err)
	}

	//
	//err = c.mqch.TxCommit()
	//if err != nil {
	//	c.mqch.TxRollback()
	//	return fmt.Errorf("commit: %w", err)
	//}

	d, ok := <- c.confirm
	if ok {
		if !d.Ack {
			return fmt.Errorf("nack!")
		}
	}else {
		return fmt.Errorf("c.confirm closed")
	}

	return nil
}

func (c *Connection) Consume() <-chan amqp.Delivery {
	c.Lock()
	defer c.Unlock()

	return c.delivery
}

func ConnectAsPublisher(addr string) (Publisher, error) {
	publisher, err := asPublisher(&Connection{addr: addr})
	return publisher, err
}

func ConnectAsConsumer(addr string) (Consumer, error) {
	consumer, err := asConsumer(&Connection{addr: addr})
	return consumer, err
}