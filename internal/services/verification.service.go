package services

import (
	"errors"
	"fmt"
	"github.com/mercan/ecommerce/internal/repositories/redis"
)

type VerificationService interface {
	VerifyPhone(phone string, verificationCode string) error
	VerifyEmail(email string, verificationCode string) error
}

type VerificationServiceImpl struct {
	authRedisRepo redis.AuthenticationRepository
}

func NewVerificationService() VerificationService {
	return &VerificationServiceImpl{
		authRedisRepo: redis.NewAuthenticationRedisRepository(),
	}
}

func (service *VerificationServiceImpl) VerifyEmail(email, verificationCode string) error {
	result, err := service.authRedisRepo.GetVerificationEmail(email)
	if err != nil {
		if errors.Is(err, service.authRedisRepo.NilError()) {
			return errors.New("verification code not found")
		}

		return err
	}

	if result != verificationCode {
		return errors.New("invalid verification code")
	}

	if err := service.authRedisRepo.DelVerificationEmail(email); err != nil {
		fmt.Println("Error deleting email verification code from redis: ", err)
	}

	return nil
}

func (service *VerificationServiceImpl) VerifyPhone(phone string, verificationCode string) error {
	result, err := service.authRedisRepo.GetVerificationPhone(phone)
	if err != nil {
		if errors.Is(err, service.authRedisRepo.NilError()) {
			return errors.New("verification code not found")
		}

		return err
	}

	if result != verificationCode {
		return errors.New("invalid verification code")
	}

	if err := service.authRedisRepo.DelVerificationPhone(phone); err != nil {
		fmt.Println("Error deleting phone verification code from redis: ", err)
	}

	return nil
}
