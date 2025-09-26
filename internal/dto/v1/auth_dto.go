package v1dto

import (
	"github.com/dangLuan01/user-manager/internal/models"
	"github.com/google/uuid"
)

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

type RegisterInput struct {
	Name 	 string `json:"name" binding:"required"`
	Email 	 string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type RequestOTPInput struct {
	OTP string `json:"otp" binding:"required,max=6"`
}

func RegisterDTOToModel(uuid uuid.UUID, user RegisterInput) models.User {
	return models.User{
		UUID: uuid,
		Name: user.Name,
		Email: user.Email,
		Password: user.Password,
		Age: 1,
		Level: 2,
		Status: 1,
	}
}