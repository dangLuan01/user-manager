package v1service

import (
	"github.com/dangLuan01/user-manager/internal/repository"
	"github.com/dangLuan01/user-manager/internal/utils"
	"github.com/dangLuan01/user-manager/pkg/auth"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo repository.UserRepository
	tokenService auth.TokenService
}

func NewAuthService(repo repository.UserRepository, tokenService auth.TokenService) *authService {
	return &authService{
		userRepo: repo,
		tokenService: tokenService,
	}
}

func (as *authService) Login(ctx *gin.Context, email, password string) (string, string, int, error) {

	email = utils.NormailizeString(email)
	user, err := as.userRepo.FindByEmail(email)

	if err != nil {

		return "", "", 0, utils.NewError(string(utils.ErrCodeUnauthorized), "Invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {

		return "", "", 0, utils.NewError(string(utils.ErrCodeUnauthorized), "Invalid email or password")
	}

	accessToken, err := as.tokenService.GenerateAccessToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError(string(utils.ErrCodeBadRequest), "Unable to create access token", err)
	}


	refreshToken, err := as.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError(string(utils.ErrCodeBadRequest), "Unable to create refresh token", err)
	}

	if err := as.tokenService.StoreRefreshToken(refreshToken); err != nil {
		return "", "", 0, utils.WrapError(string(utils.ErrCodeBadRequest), "Cannot save refresh token", err)
	}
	
	return  accessToken, refreshToken.Token, int(auth.AccessTokenTTL.Seconds()), nil
}

func (as *authService) Logout(ctx *gin.Context) error {

	panic("")
}

func (as *authService) RefreshToken(ctx *gin.Context, refreshTokenString string) (string, string, int, error) {

	token, err := as.tokenService.ValidaRefreshToken(refreshTokenString)
	if err != nil {
		return "","", 0, utils.NewError(string(utils.ErrCodeUnauthorized),"Refresh token is invalid or revoked.")
	}

	user, err := as.userRepo.FindBYUUID(token.UserUUID)
	if err != nil {
		return "","", 0, utils.NewError(string(utils.ErrCodeUnauthorized),"User not found.")
	}

	accessToken, err := as.tokenService.GenerateAccessToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError(string(utils.ErrCodeBadRequest), "Unable to create access token", err)
	}

	refreshToken, err := as.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError(string(utils.ErrCodeBadRequest), "Unable to create refresh token", err)
	}

	if err := as.tokenService.RevokeRefreshToken(refreshTokenString); err != nil {
		return "", "", 0, utils.WrapError(string(utils.ErrCodeBadRequest), "Cannot to revoke refresh token", err)
	}

	if err := as.tokenService.StoreRefreshToken(refreshToken); err != nil {
		return "", "", 0, utils.WrapError(string(utils.ErrCodeBadRequest), "Cannot save refresh token", err)
	}

	return  accessToken, refreshToken.Token, int(auth.AccessTokenTTL.Seconds()), nil
}