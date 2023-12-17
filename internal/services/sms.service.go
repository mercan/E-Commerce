package services

import (
	"github.com/mercan/ecommerce/internal/config"
	"github.com/mercan/ecommerce/internal/repositories/redis"
	"github.com/mercan/ecommerce/internal/utils"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type SMSService interface {
	SendVerificationPhone(phoneNumber string) error
}

type SMSServiceImpl struct {
	authRedisRepo           redis.AuthenticationRepository
	twilioClient            *twilio.RestClient
	twilioMessageServiceSID string
	twilioFromNumber        string
}

func NewSMSService() SMSService {
	return &SMSServiceImpl{
		authRedisRepo: redis.NewAuthenticationRedisRepository(),
		twilioClient: twilio.NewRestClientWithParams(twilio.ClientParams{
			Username: config.GetTwilioConfig().AccountSID,
			Password: config.GetTwilioConfig().AuthToken,
		}),
		twilioMessageServiceSID: config.GetTwilioConfig().MessageServiceSID,
		twilioFromNumber:        config.GetTwilioConfig().FromNumber,
	}
}

func (service *SMSServiceImpl) SendVerificationPhone(phoneNumber string) error {
	params := &openapi.CreateMessageParams{}
	verificationCode := utils.GenerateVerificationCode()
	message := "Your verification code is: " + verificationCode + "\nExpires in 10 minutes.\n\nEcommerce Demo API"

	params.SetBody(message)
	//params.SetFrom(service.twilioFromNumber)
	params.SetMessagingServiceSid(service.twilioMessageServiceSID)
	params.SetTo(phoneNumber)

	if _, err := service.twilioClient.Api.CreateMessage(params); err != nil {
		return err
	}

	if err := service.authRedisRepo.SetVerificationPhone(phoneNumber, verificationCode); err != nil {
		return err
	}

	return nil
}
