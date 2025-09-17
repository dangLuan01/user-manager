package auth

import (
	"github.com/dangLuan01/user-manager/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	GenerateAccessToken(user models.User) (string, error)
	GenerateRefreshToken(user models.User) (RefreshToken, error)
	ParseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error)
	DecryptAccessTokenPayload(tokenString string) (*EncryptedPayload, error)
	StoreRefreshToken(token RefreshToken) error
	ValidaRefreshToken(token string) (RefreshToken, error)
	RevokeRefreshToken(token string) error
}