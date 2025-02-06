package pubsub

import (
	"encoding/json"
	"fmt"

	enums "github.com/bootdotdev/learn-pub-sub-starter/internal"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType enums.SimpleQueueType,
	handler func(T) enums.AckType,
) error {

	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, int(simpleQueueType))
	if err != nil {
		return err
	}

	deliveries, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for delivery := range deliveries {
			var genericType T
			json.Unmarshal(delivery.Body, &genericType)

			ackType := handler(genericType)

			switch ackType {
			case enums.Ack:
				fmt.Println("handling Ack")
				err = delivery.Ack(false)
				handleErr(err)
			case enums.NackRequeue:
				fmt.Println("handling Requeue")
				err = delivery.Nack(false, true)
				handleErr(err)
			case enums.NackDiscard:
				fmt.Println("handling Discard")
				err = delivery.Nack(false, false)
				handleErr(err)
			}
		}
	}()

	return nil
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
