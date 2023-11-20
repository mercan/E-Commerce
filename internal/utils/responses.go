package utils

type FailureAuthResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Code    int    `json:"code"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Code    int    `json:"code"`
}

type SuccessAuthResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Token   string `json:"token"`
}
