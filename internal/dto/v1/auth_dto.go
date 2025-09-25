package v1dto

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginResponse struct {
	AccessToken 	string 	`json:"access_token"`
	RefreshToken 	string	`json:"refresh_token"`
	ExpiresIn 		int 	`json:"expires_in"`
}

type RefreshTokenInput struct {
	RefreshToken 	string `json:"refresh_token"`
}

type RequestPasswordInput struct {
	Email    string `json:"email" binding:"required,email"`
}

type RequestResetInput struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}