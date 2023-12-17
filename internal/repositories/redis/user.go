package redis

import (
	"context"
	"github.com/mercan/ecommerce/internal/config"
	"github.com/redis/go-redis/v9"
	"time"
)

type AuthenticationRepository interface {
	SetBlacklistToken(token string, expiration time.Duration) error
	IsTokenInBlacklist(token string) (bool, error)
	SetVerificationPhone(phone string, verificationCode string) error
	GetVerificationPhone(phone string) (string, error)
	DelVerificationPhone(phone string) error
	SetVerificationEmail(email string, verificationCode string) error
	GetVerificationEmail(email string) (string, error)
	DelVerificationEmail(email string) error
	SetForgotPasswordToken(email string, token string) error
	GetForgotPasswordToken(email string) (string, error)
	DelForgotPasswordToken(email string) error
	NilError() error
}

type AuthenticationRedisRepository struct {
	Ctx    context.Context
	Client *redis.Client
}

func NewAuthenticationRedisRepository() AuthenticationRepository {
	return &AuthenticationRedisRepository{
		Ctx:    context.Background(),
		Client: client,
	}
}

func (ar *AuthenticationRedisRepository) NilError() error {
	return redis.Nil
}

func (ar *AuthenticationRedisRepository) SetBlacklistToken(token string, expiration time.Duration) error {
	return ar.Client.Set(ar.Ctx, "token:"+token, token, expiration).Err()
}

func (ar *AuthenticationRedisRepository) IsTokenInBlacklist(token string) (bool, error) {
	result, err := ar.Client.Exists(ar.Ctx, "token:"+token).Result()

	if err != nil {
		return false, err
	}

	return result == 1, nil
}

func (ar *AuthenticationRedisRepository) SetVerificationPhone(phone string, verificationCode string) error {
	expiration := config.GetTimeConfig().PhoneExpireTime * time.Second

	return ar.Client.Set(ar.Ctx, "phone:"+phone, verificationCode, expiration).Err()
}

func (ar *AuthenticationRedisRepository) GetVerificationPhone(phone string) (string, error) {
	return ar.Client.Get(ar.Ctx, "phone:"+phone).Result()
}

func (ar *AuthenticationRedisRepository) DelVerificationPhone(phone string) error {
	return ar.Client.Del(ar.Ctx, "phone:"+phone).Err()
}

func (ar *AuthenticationRedisRepository) SetVerificationEmail(email string, verificationCode string) error {
	expiration := config.GetTimeConfig().EmailExpireTime * time.Second

	return ar.Client.Set(ar.Ctx, "email:"+email, verificationCode, expiration).Err()
}

func (ar *AuthenticationRedisRepository) GetVerificationEmail(email string) (string, error) {
	return ar.Client.Get(ar.Ctx, "email:"+email).Result()
}

func (ar *AuthenticationRedisRepository) DelVerificationEmail(email string) error {
	return ar.Client.Del(ar.Ctx, "email:"+email).Err()
}

func (ar *AuthenticationRedisRepository) SetForgotPasswordToken(email string, token string) error {
	expiration := config.GetTimeConfig().ForgotPasswordExpireTime * time.Second

	return ar.Client.Set(ar.Ctx, "forgot-password:"+email, token, expiration).Err()
}

func (ar *AuthenticationRedisRepository) GetForgotPasswordToken(email string) (string, error) {
	return ar.Client.Get(ar.Ctx, "forgot-password:"+email).Result()
}

func (ar *AuthenticationRedisRepository) DelForgotPasswordToken(email string) error {
	return ar.Client.Del(ar.Ctx, "forgot-password:"+email).Err()
}
