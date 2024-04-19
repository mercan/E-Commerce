package rabbitmq

import (
	"log"

	"github.com/mercan/ecommerce/internal/config"
	"github.com/streadway/amqp"
)

var connection, channel = Connect()

func Connect() (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(config.GetRabbitMQConfig().URI)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	queueDeclare(ch, config.GetRabbitMQConfig().EmailVerificationQueue)
	queueDeclare(ch, config.GetRabbitMQConfig().PhoneVerificationQueue)

	log.Println("Connected to RabbitMQ")
	return conn, ch
}

func queueDeclare(channel *amqp.Channel, queueName string) {
	_, err := channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}
}

func Close() {
	if err := connection.Close(); err != nil {
		panic(err)
	}

	if err := channel.Close(); err != nil {
		panic(err)
	}
}
