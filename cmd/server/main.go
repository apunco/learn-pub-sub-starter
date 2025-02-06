package main

import (
	"fmt"
	"os"
	"os/signal"

	enums "github.com/bootdotdev/learn-pub-sub-starter/internal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/joho/godotenv"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	gamelogic.PrintServerHelp()
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

	ch, err := rabbitMqConnection.Channel()
	if err != nil {
		fmt.Println("Failed to create a RabbitMQ channel", err)
	}

	_, _, err = pubsub.DeclareAndBind(rabbitMqConnection, routing.ExchangePerilTopic, "game_logs", "game_logs.*", int(enums.Durable))
	if err != nil {
		fmt.Println("Failed to declare and bind RabbitMQ queue game_logs", err)
	}

	message := routing.PlayingState{
		IsPaused: true,
	}

	for {
		input := gamelogic.GetInput()

		if input[0] == "pause" {
			fmt.Println("Pausing the game")
			message.IsPaused = true
			pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, message)
		} else if input[0] == "resume" {
			fmt.Println("Resuming the game")
			message.IsPaused = false
			pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, message)
		} else if input[0] == "quit" {
			break
		} else {
			fmt.Println("Uknonwn command")
		}
	}

	defer rabbitMqConnection.Close()

	fmt.Println("Connected to RabbitMQ")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Shutting down")
}
