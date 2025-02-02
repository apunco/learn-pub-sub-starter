package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connString := "amqp://guest:guest@localhost:5672/"
	rabbitMqConnection, err := amqp.Dial(connString)
	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ", err)
		return
	}

	ch, err := rabbitMqConnection.Channel()
	if err != nil {
		fmt.Println("Failed to create a RabbitMQ channel", err)
	}

	message := routing.PlayingState{
		IsPaused: true,
	}
	pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, message)

	defer rabbitMqConnection.Close()

	fmt.Println("Connected to RabbitMQ")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Shutting down")
}
