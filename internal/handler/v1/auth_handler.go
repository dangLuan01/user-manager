package v1handler

import (
	v1dto "github.com/dangLuan01/user-manager/internal/dto/v1"
	v1service "github.com/dangLuan01/user-manager/internal/service/v1"
	"github.com/dangLuan01/user-manager/internal/utils"
	"github.com/dangLuan01/user-manager/internal/validation"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService v1service.AuthService
}

func NewAuthHandler(service v1service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: service,
	}
}

func (ah *AuthHandler) Login(ctx *gin.Context) {

	var input v1dto.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	if err := ah.authService.Login(ctx, input.Email, input.Password); err != nil {
		utils.ResponseError(ctx, err)
		return
	}
}

func (ah *AuthHandler) Logout(ctx *gin.Context) {
	
}