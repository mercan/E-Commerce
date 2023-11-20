package types

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserBaseResponse struct {
	ErrorResponse
	Code int `json:"code"`
}

type UserRegisterResponse struct {
	UserBaseResponse
	Token string `json:"token"`
}

type UserLoginResponse struct {
	UserBaseResponse
	Token string `json:"token"`
}

type UserLogoutResponse struct {
	UserBaseResponse
}

type UserChangePasswordResponse struct {
	UserBaseResponse
	Token string `json:"token"`
}

type UserChangeEmailResponse struct {
	UserBaseResponse
	Token string `json:"token"`
}

type UserVerifyEmailResponse struct {
	UserBaseResponse
}

type UserResendEmailVerificationResponse struct {
	UserBaseResponse
}

type UserVerifyPhoneResponse struct {
	UserBaseResponse
}

type UserResendPhoneVerificationResponse struct {
	UserBaseResponse
}
