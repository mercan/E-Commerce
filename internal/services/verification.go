package services

import (
	"errors"
	"fmt"
	"github.com/mercan/ecommerce/internal/repositories/redis"
)

type IVerificationService interface {
	VerifyPhone(phone string, verificationCode string) error
	VerifyEmail(email string, verificationCode string) error
}

type VerificationService struct {
	redisRepo *redis.Repository
}

func NewVerificationService() IVerificationService {
	return &VerificationService{
		redisRepo: redis.NewRepository(),
	}
}

func (service *VerificationService) VerifyPhone(phone string, verificationCode string) error {
	result, err := service.redisRepo.GetVerificationPhone(phone)
	if err != nil {
		if errors.Is(err, service.redisRepo.NilError) {
			return errors.New("verification code not found")
		}

		return err
	}

	if result != verificationCode {
		return errors.New("invalid verification code")
	}

	if err := service.redisRepo.DelVerificationPhone(phone); err != nil {
		fmt.Println("Error deleting phone verification code from redis: ", err)
	}

	return nil
}

func (service *VerificationService) VerifyEmail(email, verificationCode string) error {
	result, err := service.redisRepo.GetVerificationEmail(email)
	if err != nil {
		if errors.Is(err, service.redisRepo.NilError) {
			return errors.New("verification code not found")
		}

		return err
	}

	if result != verificationCode {
		return errors.New("invalid verification code")
	}

	if err := service.redisRepo.DelVerificationEmail(email); err != nil {
		fmt.Println("Error deleting email verification code from redis: ", err)
	}

	return nil
}
