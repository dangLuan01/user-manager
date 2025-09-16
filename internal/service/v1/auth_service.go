package v1service

import (
	"github.com/dangLuan01/user-manager/internal/repository"
	"github.com/dangLuan01/user-manager/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *authService {
	return &authService{
		userRepo: repo,
	}
}

func (as *authService) Login(ctx *gin.Context, email, password string) error {

	email = utils.NormailizeString(email)
	found, err := as.userRepo.FindByEmail(email)

	if err != nil {

		return utils.NewError(string(utils.ErrCodeUnauthorized), "Invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(password)); err != nil {

		return utils.NewError(string(utils.ErrCodeUnauthorized), "Invalid email or password")
	}

	return  nil
}

func (as *authService) Logout(ctx *gin.Context) error {

	panic("")
}