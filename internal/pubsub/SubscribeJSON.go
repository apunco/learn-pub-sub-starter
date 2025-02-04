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
	handler func(T),
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

			handler(genericType)
			err = delivery.Ack(false)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	return nil
}
