package models

type UserRegisterRequest struct {
	FirstName        string           `json:"first_name" validate:"required"`
	LastName         string           `json:"last_name" validate:"required"`
	Email            string           `json:"email" validate:"required,email"`
	Password         string           `json:"password" validate:"required,min=6,max=500"`
	Description      string           `json:"description" validate:"required,min=6,max=500"`
	SocialMediaLinks SocialMediaLinks `json:"social_media"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=500"`
}

type UserChangePasswordRequest struct {
	Password           string `json:"password" validate:"required,min=6,max=500"`
	NewPassword        string `json:"new_password" validate:"required,min=6,max=500"`
	NewPasswordConfirm string `json:"new_password_confirm" validate:"required,min=6,max=500"`
}

type UserChangeEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type UserVerificationRequest struct {
	Code string `json:"code" query:"code" validate:"required,min=6,max=6,number"`
}
