package v1service

import (
	"github.com/dangLuan01/user-manager/internal/models"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	GetAllUser() ([]models.User, error)
	GetUserByUUID(uuid string) (models.User, error)
	CreateUser(user models.User) (models.User, error)
	UpdateUser(uuid string, user models.User) (models.User, error)
	DeleteUser(uuid string) error
}

type AuthService interface {
	Login(ctx *gin.Context, email, password string) error
	Logout(ctx *gin.Context) error
}