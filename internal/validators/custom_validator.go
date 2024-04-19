package validators

import (
	"github.com/go-playground/validator/v10"
	url2 "net/url"
)

var validate = validator.New()

func init() {
	// Register custom validation
	validate.RegisterValidation("customURL", customURLValidation)
}

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

func customURLValidation(fl validator.FieldLevel) bool {
	url := fl.Field().String()

	if url == "" {
		return true
	}

	u, err := url2.Parse(url)
	if err != nil {
		return false
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	if u.Host == "" {
		return false
	}

	return true
}
