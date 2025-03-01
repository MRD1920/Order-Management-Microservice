package main

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/mrd1920/oms-common/api"
	"github.com/mrd1920/oms-common/broker"
	ampq "github.com/rabbitmq/amqp091-go"
)

type consumer struct {
	service PaymentsService
}

func NewConsumer(service PaymentsService) *consumer {
	return &consumer{service}
}

func (c *consumer) Listen(ch *ampq.Channel) {
	q, err := ch.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	var forever chan struct{}

	go func() {
		for delivery := range msgs {
			log.Printf("received message: %v", delivery.Body)

			o := &pb.Order{}
			if err := json.Unmarshal(delivery.Body, o); err != nil {
				log.Printf("Failed to unmarshall order: %v", err)
				continue
			}

			paymentLink, err := c.service.CreatePayment(context.Background(), o)
			if err != nil {
				log.Printf("Failed to create payment: %v", err)
				if err := broker.HandleRetry(ch, &delivery); err != nil {
					log.Printf("Error handling retry: %v", err)
				}

				continue
			}

			log.Printf("Payment link: %v", paymentLink)

		}
	}()

	<-forever
}
