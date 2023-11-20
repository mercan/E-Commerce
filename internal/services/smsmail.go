package services

import (
	"fmt"
	"github.com/mercan/ecommerce/internal/config"
	"github.com/mercan/ecommerce/internal/repositories/redis"
	"github.com/mercan/ecommerce/internal/utils"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type ISMSMailService interface {
	SendVerificationPhone(phoneNumber string) error
	SendVerificationEmail(firstName string, email string) error
	SendForgotPasswordEmail(email string) error
}

type SMSMailService struct {
	twilioClient                     *twilio.RestClient
	twilioMessageServiceSID          string
	twilioFromNumber                 string
	redisRepository                  *redis.Repository
	sendgridAPIKey                   string
	sendgridFromName                 string
	sendgridFromEmail                string
	sendgridVerificationTemplateID   string
	sendgridForgotPasswordTemplateID string
}

func NewSMSMailService() ISMSMailService {
	return &SMSMailService{
		twilioClient: twilio.NewRestClientWithParams(twilio.ClientParams{
			Username: config.GetTwilioConfig().AccountSID,
			Password: config.GetTwilioConfig().AuthToken,
		}),
		twilioMessageServiceSID:          config.GetTwilioConfig().MessageServiceSID,
		twilioFromNumber:                 config.GetTwilioConfig().FromNumber,
		redisRepository:                  redis.NewRepository(),
		sendgridAPIKey:                   config.GetSendgridConfig().APIKey,
		sendgridFromName:                 config.GetServerConfig().AppName,
		sendgridFromEmail:                config.GetSendgridConfig().FromEmail,
		sendgridVerificationTemplateID:   config.GetSendgridConfig().VerificationTemplateID,
		sendgridForgotPasswordTemplateID: config.GetSendgridConfig().ForgotPasswordTemplateID,
	}
}

func (service *SMSMailService) SendVerificationPhone(phoneNumber string) error {
	params := &openapi.CreateMessageParams{}
	verificationCode := utils.GenerateVerificationCode()
	message := "Your verification code is: " + verificationCode + "\nExpires in 10 minutes.\n\nEcommerce Demo API"

	params.SetBody(message)
	params.SetFrom(service.twilioFromNumber)
	params.SetMessagingServiceSid(service.twilioMessageServiceSID)
	params.SetTo(phoneNumber)

	if _, err := service.twilioClient.Api.CreateMessage(params); err != nil {
		return err
	}

	if err := service.redisRepository.SetVerificationPhone(phoneNumber, verificationCode); err != nil {
		return err
	}

	return nil
}

func (service *SMSMailService) SendVerificationEmail(firstName string, email string) error {
	m := mail.NewV3Mail()
	e := mail.NewEmail(service.sendgridFromName, service.sendgridFromEmail)
	verificationCode := utils.GenerateVerificationCode()

	m.SetFrom(e)
	m.SetTemplateID(service.sendgridVerificationTemplateID)

	p := mail.NewPersonalization()
	to := mail.NewEmail(firstName, email)

	p.AddTos(to)
	p.SetDynamicTemplateData("firstName", firstName)
	p.SetDynamicTemplateData("verificationCode", verificationCode)
	m.AddPersonalizations(p)

	request := sendgrid.GetRequest(service.sendgridAPIKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)

	if _, err := sendgrid.API(request); err != nil {
		return err
	}

	if err := service.redisRepository.SetVerificationEmail(email, verificationCode); err != nil {
		return err
	}

	return nil
}

func (service *SMSMailService) SendForgotPasswordEmail(email string) error {
	m := mail.NewV3Mail()
	e := mail.NewEmail(service.sendgridFromName, service.sendgridFromEmail)
	forgotPasswordToken := utils.GenerateForgotPasswordToken()
	forgotPasswordLink := "http://localhost:" + config.GetServerConfig().Port + "/auth/forgot-password/" + forgotPasswordToken

	fmt.Println(forgotPasswordLink)
	m.SetFrom(e)
	m.SetTemplateID(service.sendgridVerificationTemplateID)

	p := mail.NewPersonalization()
	to := mail.NewEmail("", email)

	p.AddTos(to)
	p.SetDynamicTemplateData("ForgotPasswordLink", forgotPasswordLink)
	m.AddPersonalizations(p)

	request := sendgrid.GetRequest(service.sendgridAPIKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)

	if _, err := sendgrid.API(request); err != nil {
		return err
	}

	if err := service.redisRepository.SetForgotPasswordToken(email, forgotPasswordToken); err != nil {
		return err
	}

	return nil
}
