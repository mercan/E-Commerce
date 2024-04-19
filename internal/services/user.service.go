package services

import (
	"errors"
	"github.com/mercan/ecommerce/internal/helpers"
	"github.com/mercan/ecommerce/internal/models"
	"github.com/mercan/ecommerce/internal/repositories/mongodb"
	"github.com/mercan/ecommerce/internal/repositories/redis"
	"github.com/mercan/ecommerce/internal/validators"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type UserService interface {
	Register(user *models.User) (string, error)
	Login(user models.UserLoginRequest) (string, error)
	Logout(token string, expFloat64 float64) error
	ChangePassword(userId primitive.ObjectID, user models.UserChangePasswordRequest, token string,
		expFloat64 float64) (string, error)
	ChangeEmail(userId primitive.ObjectID, user models.UserChangeEmailRequest, token string, expFloat64 float64) (string, error)
	VerifyEmail(userId primitive.ObjectID, user models.UserVerificationRequest) error
	ResendEmailVerification(userId primitive.ObjectID) error
	VerifyPhone(userId primitive.ObjectID, user models.UserVerificationRequest) error
	ResendPhoneVerification(userId primitive.ObjectID) error
}

type UserServiceImpl struct {
	userRepo            mongodb.UserMongoRepository
	authRedisRepo       redis.AuthenticationRepository
	MailService         MailService
	SMSService          SMSService
	VerificationService VerificationService
}

func NewUserService() UserService {
	return &UserServiceImpl{
		userRepo:            mongodb.NewUserMongoRepository(),
		authRedisRepo:       redis.NewAuthenticationRedisRepository(),
		MailService:         NewMailService(),
		SMSService:          NewSMSService(),
		VerificationService: NewVerificationService(),
	}
}

func (service *UserServiceImpl) Register(user *models.User) (string, error) {
	if emailExists, err := service.userRepo.CheckEmailExists(user.Email); err != nil {
		return "", err
	} else if emailExists {
		return "", errors.New("Email already exists")
	}

	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		return "", errors.New("Password hashing failed")
	}
	user.Password = hashedPassword

	token, err := helpers.GenerateJWT(user.ID)
	if err != nil {
		return "", errors.New("Token generation failed")
	}

	if err := service.userRepo.CreateUser(user); err != nil {
		return "", err
	}

	return token, nil
}

func (service *UserServiceImpl) Login(user models.UserLoginRequest) (string, error) {
	if err := validators.ValidateStruct(user); err != nil {
		return "", err
	}

	project := bson.D{{"email", 1}, {"password", 1}}
	findOneOptions := options.FindOne().SetProjection(project)

	userDoc, err := service.userRepo.GetUserByEmail(user.Email, findOneOptions)
	if err != nil {
		return "", err
	}

	if userDoc == nil {
		return "", errors.New("Invalid email or password")
	}

	if result := helpers.VerifyPassword(userDoc.Password, user.Password); result != true {
		return "", errors.New("Invalid email or password")
	}

	token, err := helpers.GenerateJWT(userDoc.ID)
	if err != nil {
		return "", errors.New("Token generation failed")
	}

	return token, nil
}

func (service *UserServiceImpl) Logout(token string, expFloat64 float64) error {
	// Convert to time.Time type from float64
	expiration := time.Unix(int64(expFloat64), 0)
	// Calculate remaining time
	remainingTime := expiration.Sub(time.Now())

	// Jwt token is saved to blacklist in redis
	if err := service.authRedisRepo.SetBlacklistToken(token, remainingTime); err != nil {
		log.Println("Error while saving token to redis: ", err.Error())

		return err
	}

	return nil
}

func (service *UserServiceImpl) ChangePassword(userId primitive.ObjectID, user models.UserChangePasswordRequest, token string, expFloat64 float64) (string, error) {
	if err := validators.ValidateStruct(user); err != nil {
		return "", err
	}

	userDoc, err := service.userRepo.GetUserByID(userId)
	if err != nil {
		return "", err
	}

	if userDoc == nil {
		return "", errors.New("User not found")
	}

	if result := helpers.VerifyPassword(userDoc.Password, user.Password); result != true {
		return "", errors.New("Invalid password")
	}

	if user.Password == user.NewPassword {
		return "", errors.New("Old password and new password cannot be the same")
	}

	if err := service.userRepo.ChangePassword(userDoc.ID, user.NewPassword); err != nil {
		return "", err
	}

	// Convert to time.Time type from float64
	expiration := time.Unix(int64(expFloat64), 0)
	// Calculate remaining time
	remainingTime := expiration.Sub(time.Now())

	// Jwt token is saved to blacklist in redis
	if err := service.authRedisRepo.SetBlacklistToken(token, remainingTime); err != nil {
		log.Println("Error while saving token to redis: ", err.Error())

		return "", err
	}

	newToken, err := helpers.GenerateJWT(userDoc.ID)
	if err != nil {
		return "", errors.New("Token generation failed")
	}

	return newToken, nil
}

func (service *UserServiceImpl) ChangeEmail(userId primitive.ObjectID, user models.UserChangeEmailRequest, token string, expFloat64 float64) (string, error) {
	if err := validators.ValidateStruct(user); err != nil {
		return "", err
	}

	userDoc, err := service.userRepo.GetUserByID(userId)
	if err != nil {
		return "", err
	}

	if userDoc == nil {
		return "", errors.New("User not found")
	}

	if userDoc.Email == user.Email {
		return "", errors.New("Old email and new email cannot be the same")
	}

	existingEmail, err := service.userRepo.CheckEmailExists(user.Email)
	if err != nil {
		return "", err
	}

	if existingEmail {
		return "", errors.New("Email already exists")
	}

	if err := service.userRepo.ChangeEmail(userDoc.ID, user.Email); err != nil {
		return "", err
	}

	// Convert to time.Time type from float64
	expiration := time.Unix(int64(expFloat64), 0)
	// Calculate remaining time
	remainingTime := expiration.Sub(time.Now())

	// Jwt token is saved to blacklist in redis
	if err := service.authRedisRepo.SetBlacklistToken(token, remainingTime); err != nil {
		log.Println("Error while saving token to redis: ", err.Error())

		return "", err
	}

	newToken, err := helpers.GenerateJWT(userDoc.ID)
	if err != nil {
		return "", errors.New("Token generation failed")
	}

	return newToken, nil
}

func (service *UserServiceImpl) VerifyEmail(userId primitive.ObjectID, user models.UserVerificationRequest) error {
	if err := validators.ValidateStruct(user); err != nil {
		return err
	}

	userDoc, err := service.userRepo.GetUserByID(userId)
	if err != nil {
		return err
	}

	if userDoc == nil {
		return errors.New("User not found")
	}

	if userDoc.EmailVerified == true {
		return errors.New("Email already verified")
	}

	if err := service.VerificationService.VerifyEmail(userDoc.Email, user.Code); err != nil {
		return err
	}

	if err := service.userRepo.UpdateEmailVerificationStatus(userDoc.ID); err != nil {
		return err
	}

	return nil
}

func (service *UserServiceImpl) ResendEmailVerification(userId primitive.ObjectID) error {
	userDoc, err := service.userRepo.GetUserByID(userId)
	if err != nil {
		return err
	}

	if userDoc == nil {
		return errors.New("User not found")
	}

	if userDoc.EmailVerified == true {
		return errors.New("Email already verified")
	}

	if err := service.MailService.SendVerificationEmail(userDoc.FirstName, userDoc.Email); err != nil {
		return err
	}

	return nil
}

func (service *UserServiceImpl) VerifyPhone(userId primitive.ObjectID, user models.UserVerificationRequest) error {
	if err := validators.ValidateStruct(user); err != nil {
		return err
	}

	userDoc, err := service.userRepo.GetUserByID(userId)
	if err != nil {
		return err
	}

	if userDoc == nil {
		return errors.New("User not found")
	}

	if userDoc.PhoneNumberVerified == true {
		return errors.New("Phone number already verified")
	}

	if err := service.VerificationService.VerifyPhone(userDoc.PhoneNumber, user.Code); err != nil {
		return err
	}

	if err := service.userRepo.UpdatePhoneVerificationStatus(userDoc.ID); err != nil {
		return err
	}

	return nil
}

func (service *UserServiceImpl) ResendPhoneVerification(userId primitive.ObjectID) error {
	userDoc, err := service.userRepo.GetUserByID(userId)
	if err != nil {
		return err
	}

	if userDoc == nil {
		return errors.New("User not found")
	}

	if userDoc.PhoneNumberVerified == true {
		return errors.New("Phone number already verified")
	}

	if err := service.SMSService.SendVerificationPhone(userDoc.PhoneNumber); err != nil {
		return err
	}

	return nil
}
