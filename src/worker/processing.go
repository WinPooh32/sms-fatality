package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"common/broker"
	"common/messages"
	"common/sms"

	"github.com/streadway/amqp"
)

const timeoutWrite = time.Second * 5

func writeToDB(queries *messages.Queries, msg sms.SMS) error {
	ctx, _ := context.WithTimeout(context.Background(), timeoutWrite)
	err := queries.CreateMessage(
		ctx,
		messages.CreateMessageParams{
			Phone: msg.Phone,
			Body:  msg.Body,
		},
	)
	if err != nil {
		return fmt.Errorf("create message at database: %w", err)
	}

	return nil
}

func processDelivery(delivery amqp.Delivery, queries *messages.Queries, dbConn *sql.DB) error {
	const (
		single  = false
		requeue = true
		drop    = false
	)

	msg, err := sms.Decode(delivery.Body)
	if err != nil {
		delivery.Nack(single, drop)
		return fmt.Errorf("sms.Decode: %w", err)
	}

	if !msg.PhoneValid() {
		delivery.Nack(single, drop)
		return fmt.Errorf("invalid phone number")
	}

	if !msg.BodyValid() {
		delivery.Nack(single, drop)
		return fmt.Errorf("invalid message body")
	}

	err = writeToDB(queries, msg)
	if err != nil {
		delivery.Nack(single, requeue)
		return fmt.Errorf("write message to DB: %w", err)
	}

	err = delivery.Ack(single)
	if err != nil {
		return fmt.Errorf("delivery.Ack: %w", err)
	}

	return nil
}

func work(ctx context.Context, dbConn *sql.DB, consumer broker.Consumer) {
	var (
		counter, average uint64
		ticker           = time.NewTicker(time.Minute)
	)

	// prepared database queries for messages
	queries, err := messages.Prepare(ctx, dbConn)
	if err != nil {
		log.Println("prepare sql:", err)
		return
	}

	consumerChan := consumer.Consume()

	for {
		select {
		case d, ok := <-consumerChan:
			if !ok {
				log.Println("consumer channel is closed")
				return
			}

			err := processDelivery(d, queries, dbConn)
			if err != nil {
				log.Println("process delivery: ", err)
			}

			counter++

		case <-ticker.C:
			average = (average + counter) / 2
			mps := average / 60

			log.Printf("processed %d messages for the last minute, average rate: %d mps", counter, mps)
			counter = 0

		case <-ctx.Done():
			return
		}
	}
}
