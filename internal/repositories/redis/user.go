package redis

import (
	"context"
	"github.com/mercan/ecommerce/internal/config"
	"github.com/redis/go-redis/v9"
	"time"
)

var ctx = context.Background()

type Repository struct {
	Client   *redis.Client
	NilError error
}

func NewRepository() *Repository {
	return &Repository{
		Client:   client,
		NilError: redis.Nil,
	}
}

func (repository *Repository) SetBlacklistToken(token string, expiration time.Duration) error {
	return repository.Client.Set(ctx, "token:"+token, token, expiration).Err()
}

func (repository *Repository) GetToken(token string) (string, error) {
	return repository.Client.Get(ctx, "token:"+token).Result()
}

func (repository *Repository) IsTokenInBlacklist(token string) (bool, error) {
	result, err := repository.Client.Exists(ctx, "token:"+token).Result()

	if err != nil {
		return false, err
	}

	return result == 1, nil
}

func (repository *Repository) SetVerificationPhone(phone string, verificationCode string) error {
	expiration := config.GetTimeConfig().PhoneExpireTime * time.Second

	return repository.Client.Set(ctx, "phone:"+phone, verificationCode, expiration).Err()
}

func (repository *Repository) GetVerificationPhone(phone string) (string, error) {
	return repository.Client.Get(ctx, "phone:"+phone).Result()
}

func (repository *Repository) DelVerificationPhone(phone string) error {
	return repository.Client.Del(ctx, "phone:"+phone).Err()
}

func (repository *Repository) SetVerificationEmail(email string, verificationCode string) error {
	expiration := config.GetTimeConfig().EmailExpireTime * time.Second

	return repository.Client.Set(ctx, "email:"+email, verificationCode, expiration).Err()
}

func (repository *Repository) GetVerificationEmail(email string) (string, error) {
	return repository.Client.Get(ctx, "email:"+email).Result()
}

func (repository *Repository) DelVerificationEmail(email string) error {
	return repository.Client.Del(ctx, "email:"+email).Err()
}

func (repository *Repository) SetForgotPasswordToken(email string, token string) error {
	expiration := config.GetTimeConfig().ForgotPasswordExpireTime * time.Second

	return repository.Client.Set(ctx, "forgot-password:"+email, token, expiration).Err()
}

func (repository *Repository) GetForgotPasswordToken(email string) (string, error) {
	return repository.Client.Get(ctx, "forgot-password:"+email).Result()
}

func (repository *Repository) DelForgotPasswordToken(email string) error {
	return repository.Client.Del(ctx, "forgot-password:"+email).Err()
}
