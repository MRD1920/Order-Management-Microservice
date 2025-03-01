package broker

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	MaxRetryCount = 3
	DLQ           = "dlq_main"
)

func Connect(user, password, host, port string) (*amqp.Channel, func() error) {
	address := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)

	conn, err := amqp.Dial(address)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	err = ch.ExchangeDeclare(OrderCreatedEvent, "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.ExchangeDeclare(OrderPaidEvent, "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = createDLQAndDLX(ch)
	if err != nil {
		log.Fatal(err)
	}

	return ch, conn.Close

}

func HandleRetry(ch *amqp.Channel, delivery *amqp.Delivery) error {
	if delivery.Headers == nil {
		delivery.Headers = amqp.Table{}
	}

	retryCount, ok := delivery.Headers["x-retry-count"].(int64)
	if !ok {
		retryCount = 0
	}

	retryCount++
	delivery.Headers["x-retry-count"] = retryCount

	log.Printf("retrying message %s for the %d time", delivery.Body, retryCount)

	if retryCount >= MaxRetryCount {
		// send to dead letter queue
		log.Printf("moving message to DLQ %s", DLQ)

		return ch.PublishWithContext(context.Background(), "", DLQ, false, false, amqp.Publishing{
			ContentType:  "applcation/json",
			Headers:      delivery.Headers,
			Body:         delivery.Body,
			DeliveryMode: amqp.Persistent,
		})
	}
	time.Sleep(time.Second * time.Duration(retryCount))

	return ch.PublishWithContext(context.Background(), delivery.Exchange, delivery.RoutingKey, false, false, amqp.Publishing{
		ContentType:  "applcation/json",
		Headers:      delivery.Headers,
		Body:         delivery.Body,
		DeliveryMode: amqp.Persistent,
	})

}

func createDLQAndDLX(ch *amqp.Channel) error {
	q, err := ch.QueueDeclare(
		"main_queue", //name
		true,         //durable
		false,        //delete when unused
		false,        //exclusive
		false,        //no-wait
		nil,          //arguments
	)
	if err != nil {
		return err
	}

	//Declare DLX
	dlx := "dlx_main"
	err = ch.ExchangeDeclare(
		dlx,      //name
		"fanout", //type
		true,     //durable
		false,    //auto-deleted
		false,    //internal
		false,    //no-wait
		nil,      //arguments
	)
	if err != nil {
		return err
	}

	//Bind main queue to DLX
	err = ch.QueueBind(
		q.Name, //name
		"",     //routing key
		dlx,    //exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Declare DLQ

	_, err = ch.QueueDeclare(
		DLQ,   //name
		true,  //durable
		false, //delete when unused
		false, //exclusive
		false, //no-wait
		nil,   //arguments
	)
	if err != nil {
		return err
	}

	return err
}
