package main

import (
	"fmt"
	"os"

	enums "github.com/bootdotdev/learn-pub-sub-starter/internal"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	connString := os.Getenv("RABBITMQ_CONN_STRING")
	rabbitMqConnection, err := amqp.Dial(connString)
	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ", err)
		return
	}

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Println("Failed getting username")
	}

	queueName := fmt.Sprintf("%s.%s", routing.PauseKey, username)
	_, _, err = pubsub.DeclareAndBind(rabbitMqConnection, routing.ExchangePerilDirect, queueName, routing.PauseKey, int(enums.Transient))
	if err != nil {
		fmt.Println("error creating and binding a queue", err)
	}

	forever := make(chan struct{})
	<-forever

}
