package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mercan/ecommerce/internal/config"
	"github.com/mercan/ecommerce/internal/services"
)

var smsMailService = services.NewSMSMailService()

func ConsumeEmailVerificationQueue() {
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
			err := smsMailService.SendVerificationEmail(user["firstName"], user["email"])
			if err != nil {
				fmt.Println("Error while sending email: ", err.Error())
				if err := d.Nack(false, true); err != nil {
					fmt.Println("Error while Nacking: ", err.Error())
				}
			}
		}
	}()

	log.Printf(" [*] Email Verification Queue is waiting for messages...")
	<-forever
}

func ConsumePhoneVerificationQueue() {
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
			err := smsMailService.SendVerificationPhone(user["phone"])
			if err != nil {
				fmt.Println("Error while sending phone: ", err.Error())
				if err := d.Nack(false, true); err != nil {
					fmt.Println("Error while Nacking: ", err.Error())
				}
			}
		}
	}()

	log.Printf(" [*] Phone Verification Queue is waiting for messages...")
	<-forever
}
