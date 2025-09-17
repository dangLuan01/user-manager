package v1service

import (
	"strings"
	"time"

	"github.com/dangLuan01/user-manager/internal/repository"
	"github.com/dangLuan01/user-manager/internal/utils"
	"github.com/dangLuan01/user-manager/pkg/auth"
	"github.com/dangLuan01/user-manager/pkg/cache"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo repository.UserRepository
	tokenService auth.TokenService
	cache cache.RedisCacheService
}

func NewAuthService(repo repository.UserRepository, tokenService auth.TokenService, cache cache.RedisCacheService) *authService {
	return &authService{
		userRepo: repo,
		tokenService: tokenService,
		cache: cache,
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

func (as *authService) Logout(ctx *gin.Context, refreshTokenString string) error {
	authHeder := ctx.GetHeader("Authorization")
	if authHeder == "" || !strings.HasPrefix(authHeder, "Bearer ") {
		
		return utils.NewError(string(utils.ErrCodeUnauthorized), "Missing Authorization header")
	}
	
	accessToken := strings.TrimPrefix(authHeder, "Bearer ")
	
	_, claims, err := as.tokenService.ParseToken(accessToken)
	if err != nil {
		
		return utils.NewError(string(utils.ErrCodeUnauthorized), "Invalid access token")
	}

	if jti, ok := claims["jti"].(string); ok {
		expUnix := claims["exp"].(float64)
		exp := time.Unix(int64(expUnix), 0)
		key := "blacklist:" + jti
		ttl := time.Until(exp)
		as.cache.Set(key,"revoked", ttl)
	}

	token, err := as.tokenService.ValidaRefreshToken(refreshTokenString)
	if err != nil {
		return utils.NewError(string(utils.ErrCodeUnauthorized),"Refresh token is invalid or revoked.")
	}

	if err := as.tokenService.RevokeRefreshToken(token.Token); err != nil {
		return utils.WrapError(string(utils.ErrCodeBadRequest), "Cannot to revoke refresh token", err)
	}

	return  nil
	
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