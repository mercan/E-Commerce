package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mercan/ecommerce/internal/config"
	"github.com/mercan/ecommerce/internal/services"
	"github.com/streadway/amqp"
)

type EmailQueueManager interface {
	PublishEmailVerification(firstName, email string)
	ConsumeEmailVerificationQueue()
}

type EmailQueueManagerImpl struct {
	Channel                *amqp.Channel
	EmailVerificationQueue string
	MailService            services.MailService
}

func NewEmailQueueManager() EmailQueueManager {
	return &EmailQueueManagerImpl{
		Channel:                channel,
		EmailVerificationQueue: config.GetRabbitMQConfig().EmailVerificationQueue,
		MailService:            services.NewMailService(),
	}
}

func (queue *EmailQueueManagerImpl) PublishEmailVerification(firstName, email string) {
	body := `{"firstName": "` + firstName + `", "email": "` + email + `"}`
	err := queue.Channel.Publish(
		"",
		queue.EmailVerificationQueue,
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

func (queue *EmailQueueManagerImpl) ConsumeEmailVerificationQueue() {
	msgs, err := channel.Consume(
		config.GetRabbitMQConfig().EmailVerificationQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var user map[string]string
			if err := json.Unmarshal(d.Body, &user); err != nil {
				fmt.Println("Error while unmarshalling: ", err.Error())
			}

			log.Printf(" [X] Received Message Name: %s Email: %s", user["firstName"], user["email"])
			if err := queue.MailService.SendVerificationEmail(user["firstName"], user["email"]); err != nil {
				fmt.Println("Error while sending email: ", err.Error())

				if err := d.Nack(false, true); err != nil {
					fmt.Println("Error while Nacking: ", err.Error())
				}
			}

			log.Printf(" [X] Message Sent Name: %s Email: %s", user["firstName"], user["email"])
		}
	}()

	log.Printf(" [*] Email Verification Queue is waiting for messages...")
	<-forever
}
