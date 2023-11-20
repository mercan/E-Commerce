package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID                  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName           string             `json:"first_name,omitempty" bson:"first_name,omitempty" validate:"required"`
	LastName            string             `json:"last_name,omitempty" bson:"last_name,omitempty" validate:"required"`
	Gender              string             `json:"gender" validate:"required,eq=Male|eq=Female"`
	Email               string             `json:"email,omitempty" bson:"email,omitempty" validate:"required,email"`
	EmailVerified       interface{}        `json:"email_verified,omitempty" bson:"email_verified,omitempty"`
	Password            string             `json:"password,omitempty" bson:"password,omitempty" validate:"required,min=6,max=500"`
	PhoneNumber         string             `json:"phone_number,omitempty" bson:"phone_number,omitempty" validate:"required,e164"`
	PhoneNumberVerified interface{}        `json:"phone_number_verified,omitempty" bson:"phone_number_verified,omitempty"`
	Address             []Address          `json:"address,omitempty" bson:"address,omitempty" validate:"required"`
	LoginHistory        []LoginHistory     `json:"login_history,omitempty" bson:"login_history,omitempty"`
	CreatedAt           time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt           time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type Address struct {
	Country string `json:"country,omitempty" bson:"country,omitempty" validate:"required"`
	City    string `json:"city,omitempty" bson:"city,omitempty" validate:"required"`
	Street  string `json:"street,omitempty" bson:"street,omitempty" validate:"required"`
	ZipCode string `json:"zip_code,omitempty" bson:"zip_code,omitempty" validate:"required"`
}

type LoginHistory struct {
	IP         string      `json:"ip,omitempty" bson:"ip,omitempty"`
	Device     string      `json:"device,omitempty" bson:"device,omitempty"`
	Platform   string      `json:"platform,omitempty" bson:"platform,omitempty"`
	Browser    string      `json:"browser,omitempty" bson:"browser,omitempty"`
	City       string      `json:"city,omitempty" bson:"city,omitempty"`
	Region     string      `json:"region,omitempty" bson:"region,omitempty"`
	Country    string      `json:"country,omitempty" bson:"country,omitempty"`
	Successful interface{} `json:"successful,omitempty" bson:"successful,omitempty"`
	CreatedAt  time.Time   `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

func NewUser() *User {
	return &User{
		ID:                  primitive.NewObjectID(),
		FirstName:           "",
		LastName:            "",
		Gender:              "",
		Email:               "",
		EmailVerified:       false,
		Password:            "",
		PhoneNumber:         "",
		PhoneNumberVerified: false,
		Address:             []Address{},
		LoginHistory:        []LoginHistory{},
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
}

type UserRegisterInput struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Gender    string `json:"gender" validate:"required,eq=Male|eq=Female"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"required,e164"`
	Password  string `json:"password" validate:"required,min=6,max=500"`
}

type UserLoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=500"`
}

type UserChangePasswordInput struct {
	Password           string `json:"password" validate:"required,min=6,max=500"`
	NewPassword        string `json:"new_password" validate:"required,min=6,max=500"`
	NewPasswordConfirm string `json:"new_password_confirm" validate:"required,min=6,max=500"`
}

type UserChangeEmailInput struct {
	Email string `json:"email" validate:"required,email"`
}

type UserVerificationInput struct {
	Code string `json:"code" query:"code" validate:"required,min=6,max=6,number"`
}
