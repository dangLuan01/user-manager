package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/dangLuan01/user-manager/internal/models"
	"github.com/dangLuan01/user-manager/internal/utils"
	"github.com/dangLuan01/user-manager/pkg/cache"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	cache cache.RedisCacheService
}

type Claim struct {
	jwt.RegisteredClaims
}

type EncryptedPayload struct {
	UserUUID 	uuid.UUID `json:"user_uuid"`
	Email 		string `json:"email"`
	Role 		int8 `json:"role"`
}

type RefreshToken struct {
	Token 		string `json:"token"`
	UserUUID 	uuid.UUID `json:"user_uuid"`
	ExpiresAt 	time.Time `json:"expires_at"`
	Revoked 	bool `json:"revoked"`
}

var (
	jwtSecret = []byte(utils.GetEnv("JWT_SECRET","12345678901234567890123456789012"))
	jwtEncryptKey = []byte(utils.GetEnv("JWT_ENCRYPT_KEY","12345678901234567890123456789012"))
)
const (
	AccessTokenTTL = 30 * time.Minute
	RefreshTokenTTL = 2 * 24 * time.Hour
)

func NewJWTService(cache cache.RedisCacheService) TokenService {
	return &JWTService{
		cache: cache,
	}
}

func (js *JWTService) GenerateAccessToken(user models.User) (string, error) {
	payload := &EncryptedPayload{
		UserUUID: 	user.UUID,
		Email: 		user.Email,
		Role: 		user.Level,
	}

	rawData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	encrypted, err := utils.EncrytAES(rawData, jwtEncryptKey)
	if err != nil {
		return "", err
	}
	
	claims := jwt.MapClaims{
		"data": encrypted,
		"jti": uuid.NewString(),
		"exp": jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
		"iat": jwt.NewNumericDate(time.Now()),
		"iss": user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func (js *JWTService) ParseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, nil, utils.NewError(string(utils.ErrCodeUnauthorized), "Invalid token")
	}

	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, utils.NewError(string(utils.ErrCodeUnauthorized), "Invalid token")
	}

	return token, claim, nil
}

func (js *JWTService) DecryptAccessTokenPayload(tokenString string) (*EncryptedPayload, error) {
	_,claims, err := js.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	encrypteddata, ok := claims["data"].(string)
	if !ok {
		return nil, utils.NewError(string(utils.ErrCodeUnauthorized), "Encoded data not found")
	}

	decryptedBytes, err := utils.DecrytAES(encrypteddata, jwtEncryptKey)
	if err != nil {
		return nil, utils.WrapError(string(utils.ErrCodeBadRequest), "Cannot decode data",err)
	}

	var payload EncryptedPayload

	if err := json.Unmarshal(decryptedBytes, &payload); err != nil {
		return nil, utils.WrapError(string(utils.ErrCodeBadRequest), "Invalid data format",err)
	}

	return &payload, nil
}

func (js *JWTService) GenerateRefreshToken(user models.User) (RefreshToken, error) {
	tokenBytes := make([]byte, 32)

	if _, err := rand.Read(tokenBytes); err != nil {
		return RefreshToken{}, err
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)

	return RefreshToken{
		Token: token,
		UserUUID: user.UUID,
		ExpiresAt: time.Now().Add(RefreshTokenTTL),
		Revoked: false,
	}, nil
}

func (js *JWTService) StoreRefreshToken(token RefreshToken) error {
	cacheKey := "refresh_token:" + token.Token
	return js.cache.Set(cacheKey, token, RefreshTokenTTL)
}

func (js *JWTService) ValidaRefreshToken(token string) (RefreshToken, error) {
	cacheKey := "refresh_token:" + token

	var refreshToken RefreshToken
	if err := js.cache.Get(cacheKey, &refreshToken); err != nil || refreshToken.Revoked || refreshToken.ExpiresAt.Before(time.Now()) {
		return RefreshToken{}, utils.WrapError(string(utils.ErrCodeInternal), "Cannot get refresh token", err)
	}

	return refreshToken, nil
}

func (js *JWTService) RevokeRefreshToken(token string) error {
	cacheKey := "refresh_token:" + token

	var refreshToken RefreshToken
	if err := js.cache.Get(cacheKey, &refreshToken); err != nil {
		return utils.WrapError(string(utils.ErrCodeInternal), "Cannot get refresh token", err)
	}

	refreshToken.Revoked = true

	return js.cache.Set(cacheKey, refreshToken, time.Until(refreshToken.ExpiresAt))
}