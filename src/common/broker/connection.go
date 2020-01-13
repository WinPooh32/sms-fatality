package broker

import (
	"fmt"

	"github.com/streadway/amqp"
)

func connect(addr string, prefetchCount int) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, nil, fmt.Errorf("connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("open a RabbitMQ channel: %w", err)
	}

	err = ch.Qos(prefetchCount, 0, false)
	if err != nil {
		return nil, nil, fmt.Errorf("set mq channel qos: %w", err)
	}

	return conn, ch, nil
}

func declareQueue(channel *amqp.Channel, queueName string, args amqp.Table) (*amqp.Queue, error) {
	const durable = true

	q, err := channel.QueueDeclare(
		queueName, // name
		durable,   // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		args,      // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("declare a queue: %w", err)
	}

	return &q, nil
}

func startConsume(channel *amqp.Channel, queueName string, args amqp.Table) (<-chan amqp.Delivery, error) {
	msgs, err := channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		args,      // args
	)
	if err != nil {
		return nil, fmt.Errorf("register a consumer: %w", err)
	}

	return msgs, nil
}

func asPublisher(publ *Connection) (Publisher, error) {
	var err error

	publ.conn, publ.mqch, err = connect(publ.addr, DefaultPrefetchCount)
	if err != nil {
		return nil, fmt.Errorf("ConnectAsPublisher: %w", err)
	}

	publ.queue, err = declareQueue(publ.mqch, DefaultQueueName, DefaultQueueArgs)
	if err != nil {
		return nil, fmt.Errorf("ConnectAsPublisher: %w", err)
	}

	// put this channel into confirm mode
	err = publ.mqch.Confirm(false)
	if err != nil {
		return nil, fmt.Errorf("put the channel into confirm mode: %w", err)
	}

	publ.confirm = publ.mqch.NotifyPublish(make(chan amqp.Confirmation))

	return publ, nil
}

func asConsumer(consume *Connection) (Consumer, error) {
	var err error

	consume.conn, consume.mqch, err = connect(consume.addr, DefaultPrefetchCount)
	if err != nil {
		return nil, fmt.Errorf("ConnectAsConsumer: %w", err)
	}

	consume.queue, err = declareQueue(consume.mqch, DefaultQueueName, DefaultQueueArgs)
	if err != nil {
		return nil, fmt.Errorf("ConnectAsConsumer: %w", err)
	}

	consume.delivery, err = startConsume(consume.mqch, consume.queue.Name, DefaultQueueArgs)
	if err != nil {
		return nil, fmt.Errorf("ConnectAsConsumer: %w", err)
	}

	return consume, nil
}
