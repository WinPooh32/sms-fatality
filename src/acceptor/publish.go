package main

import (
	"fmt"
	"net/http"

	"common/broker"
	"common/sms"
)

func publish(r *http.Request) error {
	const (
		smsFieldPhone = "Phone"
		smsFieldText  = "Text"
	)

	msg := sms.SMS {
		Phone: r.FormValue(smsFieldPhone),
		Body:  r.FormValue(smsFieldText),
	}

	if len(msg.Phone) == 0 || len(msg.Body) == 0 {
		return fmt.Errorf("got empty message")
	}

	data, err := sms.Encode(msg)
	if err != nil {
		return fmt.Errorf("encode message: %w", err)
	}

	publ, ok := r.Context().Value(contextKeyPublisher).(broker.Publisher)
	if !ok {
		return fmt.Errorf("extract broker from context")
	}

	if err := publ.Publish(data); err != nil {
		return fmt.Errorf("publish message to broker: %w", err)
	}

	return nil
}

