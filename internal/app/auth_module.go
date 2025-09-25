package app

import (
	v1handler "github.com/dangLuan01/user-manager/internal/handler/v1"
	"github.com/dangLuan01/user-manager/internal/repository"
	"github.com/dangLuan01/user-manager/internal/routes"
	v1routes "github.com/dangLuan01/user-manager/internal/routes/v1"
	v1service "github.com/dangLuan01/user-manager/internal/service/v1"
	"github.com/dangLuan01/user-manager/pkg/auth"
	"github.com/dangLuan01/user-manager/pkg/cache"
	"github.com/dangLuan01/user-manager/pkg/mail"
	"github.com/dangLuan01/user-manager/pkg/rabbitmq"
)

type AuthModule struct {
	routes routes.Route
}

func NewAuthModule(ctx *ModuleContext, tokenService auth.TokenService, cacheService cache.RedisCacheService, mailService mail.EmailProviderService, rabbitmqService rabbitmq.RabbitMQService) *AuthModule {

	userRepo := repository.NewSqlUserRepository(ctx.DB)
	authService := v1service.NewAuthService(userRepo, tokenService, cacheService, mailService, rabbitmqService)
	authHandler := v1handler.NewAuthHandler(authService) 
	authRoutes := v1routes.NewAuthRoutes(authHandler)

	return &AuthModule{
		routes: authRoutes,
	}
}
func (m *AuthModule) Routes() routes.Route {
	return m.routes
}