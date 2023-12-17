package services

import (
	"github.com/mercan/ecommerce/internal/config"
	"github.com/mercan/ecommerce/internal/repositories/redis"
	"github.com/mercan/ecommerce/internal/utils"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailService interface {
	SendVerificationEmail(firstName string, email string) error
	SendForgotPasswordEmail(email string) error
}

type MailServiceImpl struct {
	authRedisRepo                    redis.AuthenticationRepository
	sendgridAPIKey                   string
	sendgridFromName                 string
	sendgridFromEmail                string
	sendgridVerificationTemplateID   string
	sendgridForgotPasswordTemplateID string
}

func NewMailService() MailService {
	return &MailServiceImpl{
		authRedisRepo:                    redis.NewAuthenticationRedisRepository(),
		sendgridAPIKey:                   config.GetSendgridConfig().APIKey,
		sendgridFromName:                 config.GetServerConfig().AppName,
		sendgridFromEmail:                config.GetSendgridConfig().FromEmail,
		sendgridVerificationTemplateID:   config.GetSendgridConfig().VerificationTemplateID,
		sendgridForgotPasswordTemplateID: config.GetSendgridConfig().ForgotPasswordTemplateID,
	}
}

func (service *MailServiceImpl) SendVerificationEmail(firstName string, email string) error {
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

	if err := service.authRedisRepo.SetVerificationEmail(email, verificationCode); err != nil {
		return err
	}

	return nil
}

func (service *MailServiceImpl) SendForgotPasswordEmail(email string) error {
	m := mail.NewV3Mail()
	e := mail.NewEmail(service.sendgridFromName, service.sendgridFromEmail)
	forgotPasswordToken := utils.GenerateForgotPasswordToken()
	forgotPasswordLink := "http://localhost:" + config.GetServerConfig().Port + "/auth/forgot-password/" + forgotPasswordToken

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

	if err := service.authRedisRepo.SetForgotPasswordToken(email, forgotPasswordToken); err != nil {
		return err
	}

	return nil
}
