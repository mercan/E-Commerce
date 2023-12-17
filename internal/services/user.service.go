package services

import (
	"errors"
	"github.com/mercan/ecommerce/internal/config"
	"github.com/mercan/ecommerce/internal/models"
	"github.com/mercan/ecommerce/internal/repositories/mongodb"
	"github.com/mercan/ecommerce/internal/repositories/redis"
	"github.com/mercan/ecommerce/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

type UserService interface {
	Register(user *models.User, userAgent string) (string, error)
	Login(user models.UserLoginInput, IPAddress, userAgent string) (string, error)
	Logout(token string, expFloat64 float64) error
	ChangePassword(userId primitive.ObjectID, user models.UserChangePasswordInput, token string, expFloat64 float64) (string, error)
	ChangeEmail(userId primitive.ObjectID, user models.UserChangeEmailInput, token string, expFloat64 float64) (string, error)
	VerifyEmail(userId primitive.ObjectID, user models.UserVerificationInput) error
	ResendEmailVerification(userId primitive.ObjectID) error
	VerifyPhone(userId primitive.ObjectID, user models.UserVerificationInput) error
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

func (service *UserServiceImpl) Register(user *models.User, userAgent string) (string, error) {
	if err := utils.UserRegisterValidate(user); err != nil {
		return "", err
	}

	existingPhone, err := service.userRepo.CheckPhoneExists(user.PhoneNumber)
	if err != nil {
		return "", errors.New("Internal Server Error")
	}

	if existingPhone {
		return "", errors.New("Phone number already exists")
	}

	existingEmail, err := service.userRepo.CheckEmailExists(user.Email)
	if err != nil {
		return "", errors.New("Internal Server Error")
	}

	if existingEmail {
		return "", errors.New("Email already exists")
	}

	if hashedPassword, err := utils.HashPassword(user.Password); err != nil {
		return "", errors.New("Password hashing failed")
	} else {
		user.Password = hashedPassword
	}

	user.LoginHistory[0].LoginSuccess = true
	user.LoginHistory[0].CreatedAt = time.Now()

	if location := utils.GetLocationFromIP(user.LoginHistory[0].IP); location != nil {
		user.LoginHistory[0].City = location.City
		user.LoginHistory[0].Region = location.Region
		user.LoginHistory[0].Country = location.Country
	}

	if ua := utils.ParseUserAgent(userAgent); ua != nil {
		user.LoginHistory[0].Device = ua.Device
		user.LoginHistory[0].Platform = ua.OS
		user.LoginHistory[0].Browser = ua.Browser
	}

	token, err := utils.GenerateJWT(user.ID, config.GetJWTConfig().UserRole)
	if err != nil {
		return "", errors.New("Token generation failed")
	}

	if err := service.userRepo.CreateUser(user); err != nil {
		return "", err
	}

	return token, nil
}

func (service *UserServiceImpl) Login(user models.UserLoginInput, IPAddress, userAgent string) (string, error) {
	if err := utils.ValidateStruct(user); err != nil {
		return "", err
	}

	// Bütün user datalarını almak yerine sadece email ve password'ü al.
	userDoc, err := service.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return "", err
	}

	if userDoc == nil {
		return "", errors.New("Invalid email or password")
	}

	loginHistory := models.LoginHistory{
		IP:           IPAddress,
		LoginSuccess: false,
		CreatedAt:    time.Now(),
	}

	if location := utils.GetLocationFromIP(IPAddress); location != nil {
		loginHistory.City = location.City
		loginHistory.Region = location.Region
		loginHistory.Country = location.Country
	}

	if ua := utils.ParseUserAgent(userAgent); ua != nil {
		loginHistory.Device = ua.Device
		loginHistory.Platform = ua.OS
		loginHistory.Browser = ua.Browser
	}

	if userDoc == nil {
		service.userRepo.AddLoginHistory(userDoc.ID, loginHistory)
		return "", errors.New("Invalid email or password")
	}

	if result := utils.VerifyPassword(userDoc.Password, user.Password); result != true {
		service.userRepo.AddLoginHistory(userDoc.ID, loginHistory)
		return "", errors.New("Invalid email or password")
	}

	loginHistory.LoginSuccess = true
	token, err := utils.GenerateJWT(userDoc.ID, config.GetJWTConfig().UserRole)
	if err != nil {
		service.userRepo.AddLoginHistory(userDoc.ID, loginHistory)
		return "", errors.New("Token generation failed")
	}

	service.userRepo.AddLoginHistory(userDoc.ID, loginHistory)
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

func (service *UserServiceImpl) ChangePassword(userId primitive.ObjectID, user models.UserChangePasswordInput, token string, expFloat64 float64) (string, error) {
	if err := utils.ValidateStruct(user); err != nil {
		return "", err
	}

	userDoc, err := service.userRepo.GetUserByID(userId)
	if err != nil {
		return "", err
	}

	if userDoc == nil {
		return "", errors.New("User not found")
	}

	if result := utils.VerifyPassword(userDoc.Password, user.Password); result != true {
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

	newToken, err := utils.GenerateJWT(userDoc.ID, config.GetJWTConfig().UserRole)
	if err != nil {
		return "", errors.New("Token generation failed")
	}

	return newToken, nil
}

func (service *UserServiceImpl) ChangeEmail(userId primitive.ObjectID, user models.UserChangeEmailInput, token string, expFloat64 float64) (string, error) {
	if err := utils.ValidateStruct(user); err != nil {
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

	newToken, err := utils.GenerateJWT(userDoc.ID, config.GetJWTConfig().UserRole)
	if err != nil {
		return "", errors.New("Token generation failed")
	}

	return newToken, nil
}

func (service *UserServiceImpl) VerifyEmail(userId primitive.ObjectID, user models.UserVerificationInput) error {
	if err := utils.ValidateStruct(user); err != nil {
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

func (service *UserServiceImpl) VerifyPhone(userId primitive.ObjectID, user models.UserVerificationInput) error {
	if err := utils.ValidateStruct(user); err != nil {
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
