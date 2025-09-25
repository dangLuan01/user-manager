package v1service

import (
	"github.com/dangLuan01/user-manager/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserService interface {
	GetAllUser() ([]models.User, error)
	GetUserByUUID(uuid uuid.UUID) (models.User, error)
	CreateUser(user models.User) (models.User, error)
	UpdateUser(uuid uuid.UUID, user models.User) (models.User, error)
	DeleteUser(uuid uuid.UUID) error
}

type AuthService interface {
	Login(ctx *gin.Context, email, password string) (string, string, int, error)
	Logout(ctx *gin.Context, refreshTokenString string) error
	RefreshToken(ctx *gin.Context, token string) (string, string, int, error)
	RequestForgotPassword(ctx *gin.Context, email string) (string, error)
	RequestResetPassword(ctx *gin.Context, token, password string) error
}