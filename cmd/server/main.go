package main

import (
	"fmt"
	"os"
	"os/signal"

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

	defer rabbitMqConnection.Close()

	fmt.Println("Connected to RabbitMQ")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Shutting down")
}
