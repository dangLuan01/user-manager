package app

import (
	v1handler "github.com/dangLuan01/user-manager/internal/handler/v1"
	"github.com/dangLuan01/user-manager/internal/repository"
	"github.com/dangLuan01/user-manager/internal/routes"
	v1routes "github.com/dangLuan01/user-manager/internal/routes/v1"
	v1service "github.com/dangLuan01/user-manager/internal/service/v1"
	"github.com/dangLuan01/user-manager/pkg/auth"
	"github.com/dangLuan01/user-manager/pkg/cache"
)

type AuthModule struct {
	routes routes.Route
}

func NewAuthModule(ctx *ModuleContext, tokenService auth.TokenService, cache cache.RedisCacheService) *AuthModule {

	userRepo := repository.NewSqlUserRepository(ctx.DB)
	authService := v1service.NewAuthService(userRepo, tokenService, cache)
	authHandler := v1handler.NewAuthHandler(authService) 
	authRoutes := v1routes.NewAuthRoutes(authHandler)

	return &AuthModule{
		routes: authRoutes,
	}
}
func (m *AuthModule) Routes() routes.Route {
	return m.routes
}