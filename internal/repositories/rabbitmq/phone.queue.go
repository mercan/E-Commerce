package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mercan/ecommerce/internal/config"
	"github.com/mercan/ecommerce/internal/services"
	"github.com/streadway/amqp"
)

type PhoneQueueManager interface {
	PublishPhoneVerification(phone string)
	ConsumePhoneVerificationQueue()
}

type PhoneQueueManagerImpl struct {
	Channel                *amqp.Channel
	PhoneVerificationQueue string
	SMSService             services.SMSService
}

func NewPhoneQueueManager() PhoneQueueManager {
	return &PhoneQueueManagerImpl{
		Channel:                channel,
		PhoneVerificationQueue: config.GetRabbitMQConfig().PhoneVerificationQueue,
		SMSService:             services.NewSMSService(),
	}
}

func (queue *PhoneQueueManagerImpl) PublishPhoneVerification(phone string) {
	body := `{"phone": "` + phone + `"}`
	err := queue.Channel.Publish(
		"",
		queue.PhoneVerificationQueue,
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

func (queue *PhoneQueueManagerImpl) ConsumePhoneVerificationQueue() {
	msgs, err := channel.Consume(
		config.GetRabbitMQConfig().PhoneVerificationQueue,
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

			log.Printf(" [X] Received Message: %s", user["phone"])
			if err := queue.SMSService.SendVerificationPhone(user["phone"]); err != nil {
				fmt.Println("Error while sending phone: ", err.Error())

				if err := d.Nack(false, true); err != nil {
					fmt.Println("Error while Nacking: ", err.Error())
				}
			}

			log.Printf(" [X] Message Sent: %s", user["phone"])
		}
	}()

	log.Printf(" [*] Phone Verification Queue is waiting for messages...")
	<-forever
}
