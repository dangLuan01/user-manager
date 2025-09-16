package app

import (
	v1handler "github.com/dangLuan01/user-manager/internal/handler/v1"
	"github.com/dangLuan01/user-manager/internal/repository"
	"github.com/dangLuan01/user-manager/internal/routes"
	v1routes "github.com/dangLuan01/user-manager/internal/routes/v1"
	v1service "github.com/dangLuan01/user-manager/internal/service/v1"
)

type AuthModule struct {
	routes routes.Route
}

func NewAuthModule(ctx *ModuleContext) *AuthModule {

	userRepo := repository.NewSqlUserRepository(ctx.DB)
	authService := v1service.NewAuthService(userRepo)
	authHandler := v1handler.NewAuthHandler(authService) 
	authRoutes := v1routes.NewAuthRoutes(authHandler)

	return &AuthModule{
		routes: authRoutes,
	}
}
func (m *AuthModule) Routes() routes.Route {
	return m.routes
}