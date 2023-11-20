package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/mercan/ecommerce/internal/models"
)

var validate = validator.New()

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

func UserRegisterValidate(user *models.User) error {
	registerValue := models.UserRegisterInput{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Gender:    user.Gender,
		Email:     user.Email,
		Phone:     user.PhoneNumber,
		Password:  user.Password,
	}

	return validate.Struct(registerValue)
}
