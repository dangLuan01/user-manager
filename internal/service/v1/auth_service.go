package v1service

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	v1dto "github.com/dangLuan01/user-manager/internal/dto/v1"
	"github.com/dangLuan01/user-manager/internal/repository"
	"github.com/dangLuan01/user-manager/internal/utils"
	"github.com/dangLuan01/user-manager/pkg/auth"
	"github.com/dangLuan01/user-manager/pkg/cache"
	"github.com/dangLuan01/user-manager/pkg/mail"
	"github.com/dangLuan01/user-manager/pkg/rabbitmq"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
)

var (
	clients = make(map[string]*LoginAttempt)
	mu sync.Mutex
	LoginAttemptTTL = 5 * time.Minute
	MaxLoginAttempt = 5
)

type LoginAttempt struct {
	limiter *rate.Limiter
	lastSeen time.Time
}

type authService struct {
	userRepo repository.UserRepository
	tokenService auth.TokenService
	cache cache.RedisCacheService
	mailService mail.EmailProviderService
	rabbitmqService rabbitmq.RabbitMQService
}

func NewAuthService(repo repository.UserRepository, tokenService auth.TokenService, cacheService cache.RedisCacheService, mailService mail.EmailProviderService, rabbitmqService rabbitmq.RabbitMQService) *authService {
	return &authService{
		userRepo: repo,
		tokenService: tokenService,
		cache: cacheService,
		mailService: mailService,
		rabbitmqService: rabbitmqService,
	}
}

func (as *authService) getClientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()
	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}
	return ip
}

func (as *authService) getLoginAttempt(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	client, exists := clients[ip]
	
	if !exists {

		limiter := rate.NewLimiter(rate.Limit(float32(MaxLoginAttempt) / float32(LoginAttemptTTL.Seconds())), MaxLoginAttempt)
		client = &LoginAttempt{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		clients[ip] = client

		return limiter
	}

	client.lastSeen = time.Now()

	return client.limiter
}

func (as *authService) CheckLoginAttempt(ip string) error {
	limitter := as.getLoginAttempt(ip)

	if !limitter.Allow() {
		return utils.NewError(string(utils.ErrCodeTooManyRequest), "Too many login attempt. Please rety again later")
	}

	return  nil
}

func (as *authService) CleanupClients(ip string) {
	mu.Lock()
	defer mu.Unlock()
	delete(clients, ip)
}

func (as *authService) Login(ctx *gin.Context, email, password string) (string, string, int, error) {
	ip := as.getClientIP(ctx)

	if err := as.CheckLoginAttempt(ip); err != nil {
		return "", "", 0, err
	}

	email = utils.NormailizeString(email)
	user, err := as.userRepo.FindByEmail(email)

	if err != nil {
		as.getLoginAttempt(ip)
		return "", "", 0, utils.NewError(string(utils.ErrCodeUnauthorized), "Invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		as.getLoginAttempt(ip)
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

	as.CleanupClients(ip) 
	
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

func (as *authService) RequestForgotPassword(ctx *gin.Context, email string) (string, error) {

	rateLimitKey := fmt.Sprintf("reset:ratelimit:%s", email)

	if exists, err := as.cache.Exits(rateLimitKey); exists && err == nil {
		return "", utils.NewError(string(utils.ErrCodeTooManyRequest), "Wait before requesting anorther password reset")
	}

	email = utils.NormailizeString(email)
	user, err := as.userRepo.FindByEmail(email)

	if err != nil || user.Email == "" {
		return "", utils.NewError(string(utils.ErrCodeNotFound), "Email not found")
	}

	token, err := utils.GenerateRandomString(32)
	if err != nil {
		return "", utils.NewError(string(utils.ErrCodeInternal), "Failed to generate reset token")
	}

	err = as.cache.Set("reset:" + token, user.UUID.String(), time.Hour)
	if err != nil {
		return "", utils.NewError(string(utils.ErrCodeInternal), "Failed to store reset token")
	}

	err = as.cache.Set(rateLimitKey, "1", time.Minute)
	if err != nil {
		return "", utils.NewError(string(utils.ErrCodeInternal), "Failed to store rate limit reset password")
	}

	resetLink := fmt.Sprintf("https://yourdomain.com/reset-password?token=%s", token)
	mailContent := &mail.Email{
		To: []mail.Address{
			{Email: email},
		},
		// Subject: "Password reset request",
		// Text: fmt.Sprintf("Hi %s, \n\n You requested to reset your password. Click the link below to reset it:\n%s\n\n The link will exprie in a hours", user.Name, resetLink),
		Template_Uuid: "87473bd7-9b31-47e4-888c-4ebec21d2397",
		Template_Variables: mail.EmailParams{
			User_Email: user.Name,
			Pass_Reset_Link: resetLink,
		},
	}

	// if err := as.mailService.SendMail(ctx, mailContent); err != nil {
	// 	log.Println(err)
	// 	return "", utils.NewError(string(utils.ErrCodeInternal), "Failed to send email.")
		
	// }
	if err := as.rabbitmqService.Publish(ctx, "auth_email_queue", mailContent); err != nil {
		return "", utils.NewError(string(utils.ErrCodeInternal), "Failed to send email reset password.")
	}

	return token, nil
}

func (as *authService)RequestResetPassword(ctx *gin.Context, token, password string) error {

	var userUUIDStr string
	err := as.cache.Get("reset:" + token, &userUUIDStr)
	log.Println(userUUIDStr)
	if err == redis.Nil || userUUIDStr == "" {
		return utils.NewError(string(utils.ErrCodeInternal), "Invalid or expried token")
	}

	if err != nil {
		return utils.NewError(string(utils.ErrCodeInternal), "Failed to get reset token")
	}

	userUUID, err := uuid.Parse(userUUIDStr)
	if err != nil {
		return utils.WrapError(string(utils.ErrCodeInternal),"Uuid is invalid", err)
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return utils.WrapError(
			string(utils.ErrCodeInternal), 
			"Faile hash password", 
			err,
		)
	}

	if err := as.userRepo.UpdatePassword(userUUID, string(hashPassword)); err != nil {
		return utils.NewError(string(utils.ErrCodeInternal), "Unable update new password")
	}

	if err := as.cache.Clear("reset:" + token); err != nil {
		return utils.NewError(string(utils.ErrCodeInternal), "Failed to revoked token")
	}

	return nil
}

func (as *authService) Register(ctx *gin.Context, input v1dto.RegisterInput) error {

	rateLimitKey := fmt.Sprintf("code:ratelimit:%s", input.Email)

	if exists, err := as.cache.Exits(rateLimitKey); exists && err == nil {
		return utils.NewError(string(utils.ErrCodeTooManyRequest), "Wait before requesting anorther code")
	}

	email := utils.NormailizeString(input.Email)
	user, err := as.userRepo.FindByEmail(email)
	if err != nil || user.Email != "" {
		return utils.NewError(string(utils.ErrCodeConflict), "Email existsing!")
	}	
	
	code, err := utils.GenerateRandomInt(6)
	if err != nil {
		return utils.NewError(string(utils.ErrCodeInternal), "Unable error generate otp.")
	}

	codeKey := fmt.Sprintf("code:%s", code)
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.NewError(string(utils.ErrCodeInternal), "Unable error hash password.")
	}

	input.Password = string(hashPassword)
	if err := as.cache.Set(codeKey, input, 10 * time.Minute); err != nil{
		return utils.NewError(string(utils.ErrCodeInternal), "Unable error store otp")
	}

	err = as.cache.Set(rateLimitKey, "1", 2 * time.Minute)
	if err != nil {
		return utils.NewError(string(utils.ErrCodeInternal), "Failed to store rate limit code otp")
	}

	mailContent := &mail.Email{
		To: []mail.Address{
			{Email: input.Email},
		},
		Template_Uuid: "545b7c8a-cfa2-498b-83f3-8b284e30f318",
		Template_Variables: mail.EmailParams {
			User_Email: input.Email,
			Pass_Reset_Link: code,
		},
	}

	if err := as.rabbitmqService.Publish(ctx, "auth_email_queue", mailContent); err != nil {
		return utils.NewError(string(utils.ErrCodeInternal), "Failed to send email code confirm.")
	}

	return nil
}

func (as *authService) RegisterOTP(ctx *gin.Context, code string) error {
	var user v1dto.RegisterInput
	codeKey := fmt.Sprintf("code:%s", code)
	
	err := as.cache.Get(codeKey, &user)
	if err != nil || user.Email == "" || user.Password == ""{
		return utils.NewError(string(utils.ErrCodeInternal), "Code invalid or expried.")
	}

	uuidUser := uuid.New()
	userModel := v1dto.RegisterDTOToModel(uuidUser, user)
	if err := as.userRepo.Create(userModel); err != nil {
		return utils.WrapError(string(utils.ErrCodeInternal), "Failed to store user.", err)
	}

	if err := as.cache.Clear(codeKey); err != nil {
		return utils.NewError(string(utils.ErrCodeInternal), "Unable error clear otp.")
	}

	return nil
}



