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
		fmt.Println("Error loading .env file", err)
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
		fmt.Println("error creating and binding a pause queue", err)
	}

	moveQ := fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, username)
	armyMoveCh, _, err := pubsub.DeclareAndBind(rabbitMqConnection, routing.ExchangePerilTopic, moveQ, "army_moves.*", int(enums.Transient))
	if err != nil {
		fmt.Println("error creating and binding army moves que", err)
	}

	gs := gamelogic.NewGameState(username)
	pubsub.SubscribeJSON(rabbitMqConnection, routing.ExchangePerilDirect, fmt.Sprintf("pause.%s", username), routing.PauseKey, enums.Transient, handlerPause(gs))
	pubsub.SubscribeJSON(rabbitMqConnection, routing.ExchangePerilTopic, moveQ, "army_moves.*", enums.Transient, handlerMove(gs))
	for {
		input := gamelogic.GetInput()

		if input[0] == "spawn" {
			gs.CommandSpawn(input)
		} else if input[0] == "move" {
			move, err := gs.CommandMove(input)
			if err != nil {
				fmt.Println("error making a move from input", input)
			}

			err = pubsub.PublishJSON(armyMoveCh, routing.ExchangePerilTopic, moveQ, move)
			if err != nil {
				fmt.Println("error publishing move", input)
			}
			fmt.Println("move published", input)
		} else if input[0] == "status" {
			gs.CommandStatus()
		} else if input[0] == "help" {
			gamelogic.PrintClientHelp()
		} else if input[0] == "spam" {
			fmt.Println("spam not implemented yet")
		} else if input[0] == "quit" {
			gamelogic.PrintQuit()
			break
		} else {
			fmt.Println("Invalid input")
		}
	}
}

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) enums.AckType {
	return func(ps routing.PlayingState) enums.AckType {
		defer fmt.Println("> ")
		gs.HandlePause(ps)
		return enums.Ack
	}
}

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) enums.AckType {
	return func(move gamelogic.ArmyMove) enums.AckType {
		defer fmt.Println("> ")
		outcome := gs.HandleMove(move)

		switch outcome {
		case gamelogic.MoveOutComeSafe, gamelogic.MoveOutcomeMakeWar:
			fmt.Println("Ack message")
			return enums.Ack
		case gamelogic.MoveOutcomeSamePlayer:
			fmt.Println("Nack discard message")
			return enums.NackDiscard
		default:
			fmt.Println("Nack discard message")
			return enums.NackDiscard
		}
	}
}
