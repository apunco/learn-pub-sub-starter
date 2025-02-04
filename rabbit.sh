#!/bin/bash

case "$1" in
start)
    echo "Starting RabbitMQ container..."
    if [ "$(docker ps -aq -f name=rabbitmq)" ]; then
        # If the container exists but is not running, start it
        if [ "$(docker ps -aq -f status=exited -f name=rabbitmq)" ]; then
            docker start rabbitmq
            echo "Started existing RabbitMQ container."
        else
            echo "RabbitMQ container is already running."
        fi
    else
        # If the container does not exist, create and start a new one
        docker run -d --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.13-management
        echo "Created and started new RabbitMQ container."
    fi
    ;;
stop)
    echo "Stopping RabbitMQ container..."
    docker stop rabbitmq
    ;;
logs)
    echo "Fetching logs for RabbitMQ container..."
    docker logs -f rabbitmq
    ;;
*)
    echo "Usage: $0 {start|stop|logs}"
    exit 1
    ;;
esac
