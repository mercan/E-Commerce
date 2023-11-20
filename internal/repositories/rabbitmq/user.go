package rabbitmq

import (
	"github.com/mercan/ecommerce/internal/config"
	"github.com/streadway/amqp"
	"log"
)

type Repository struct {
	Channel                *amqp.Channel
	EmailVerificationQueue string
	PhoneVerificationQueue string
}

func NewRepository() *Repository {
	return &Repository{
		Channel:                channel,
		EmailVerificationQueue: config.GetRabbitMQConfig().EmailVerificationQueue,
		PhoneVerificationQueue: config.GetRabbitMQConfig().PhoneVerificationQueue,
	}
}

func (repository *Repository) PublishEmailVerification(firstName, email string) {
	body := `{"firstName": "` + firstName + `", "email": "` + email + `"}`
	err := repository.Channel.Publish(
		"",
		repository.EmailVerificationQueue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		},
	)
	if err != nil {
		log.Printf(" [X] Failed to publish email verification: %s", err.Error())
	}

	log.Printf(" [X] Published Message: %s", body)
}

func (repository *Repository) PublishPhoneVerification(phone string) {
	body := `{"phone": "` + phone + `"}`
	err := repository.Channel.Publish(
		"",
		repository.PhoneVerificationQueue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		},
	)
	if err != nil {
		log.Printf(" [X] Failed to publish phone verification: %s", err.Error())
	}

	log.Printf(" [X] Sent Phone %s", body)
}
