package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go" // Делаем удобное имя для импорта в нашем коде
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/") // Создаем подключение к RabbitMQ
	if err != nil {
		log.Fatalf("unable to open connect to RabbitMQ server. Error: %s", err)
	}

	defer func() {
		_ = conn.Close() // Закрываем подключение в случае удачной попытки
	}()
}
