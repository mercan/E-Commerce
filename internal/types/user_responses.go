package types

type BaseResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type UserRegisterResponse struct {
	BaseResponse
	Token string `json:"token,omitempty"`
}

type UserLoginResponse struct {
	BaseResponse
	Token string `json:"token,omitempty"`
}

type UserLogoutResponse struct {
	BaseResponse
}

type UserChangePasswordResponse struct {
	BaseResponse
	Token string `json:"token,omitempty"`
}

type UserChangeEmailResponse struct {
	BaseResponse
	Token string `json:"token,omitempty"`
}

type UserVerifyEmailResponse struct {
	BaseResponse
}

type UserResendEmailVerificationResponse struct {
	BaseResponse
}

type UserVerifyPhoneResponse struct {
	BaseResponse
}

type UserResendPhoneVerificationResponse struct {
	BaseResponse
}
