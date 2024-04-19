package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"github.com/mercan/ecommerce/internal/validators"
)

type User struct {
	ID                  primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName           string             `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName            string             `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email               string             `json:"email,omitempty" bson:"email,omitempty"`
	EmailVerified       bool               `json:"email_verified" bson:"email_verified"`
	Password            string             `json:"password,omitempty" bson:"password,omitempty"`
	PhoneNumber         string             `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
	PhoneNumberVerified bool               `json:"phone_number_verified" bson:"phone_number_verified"`
	IsActive            bool               `json:"is_active" bson:"is_active"`
	Description         string             `json:"description,omitempty" bson:"description,omitempty"`
	SocialMediaLinks    SocialMediaLinks   `json:"social_media_links,omitempty" bson:"social_media_links,omitempty"`
	Price               int                `json:"price,omitempty" bson:"price,omitempty"`
	ProfileImage        string             `json:"profile_image,omitempty" bson:"profile_image,omitempty"`
	BannerImage         string             `json:"banner_image,omitempty" bson:"banner_image,omitempty"`
	CreatedAt           time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt           time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type SocialMediaLinks struct {
	Website   string `json:"website,omitempty" bson:"website,omitempty" validate:"customURL"`
	Github    string `json:"github,omitempty" bson:"github,omitempty" validate:"customURL"`
	Instagram string `json:"instagram,omitempty" bson:"instagram,omitempty" validate:"customURL"`
	Facebook  string `json:"facebook,omitempty" bson:"facebook,omitempty" validate:"customURL"`
	Twitter   string `json:"twitter,omitempty" bson:"twitter,omitempty" validate:"customURL"`
	LinkedIn  string `json:"linkedin,omitempty" bson:"linkedin,omitempty" validate:"customURL"`
	YouTube   string `json:"youtube,omitempty" bson:"youtube,omitempty" validate:"customURL"`
}

func NewUser() *User {
	return &User{
		ID:               primitive.NewObjectID(),
		IsActive:         true,
		SocialMediaLinks: SocialMediaLinks{},
		Price:            100, // Constant value
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

func (u *User) RegisterValidation() error {
	registerStruct := UserRegisterRequest{
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Email:       u.Email,
		Password:    u.Password,
		Description: u.Description,
		SocialMediaLinks: SocialMediaLinks{
			Website:   u.SocialMediaLinks.Website,
			Github:    u.SocialMediaLinks.Github,
			Instagram: u.SocialMediaLinks.Instagram,
			Facebook:  u.SocialMediaLinks.Facebook,
			Twitter:   u.SocialMediaLinks.Twitter,
			LinkedIn:  u.SocialMediaLinks.LinkedIn,
			YouTube:   u.SocialMediaLinks.YouTube,
		},
	}

	return validators.ValidateStruct(registerStruct)
}
