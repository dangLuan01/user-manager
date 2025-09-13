package app

import (
	v1handler "github.com/dangLuan01/user-manager/internal/handler/v1"
	"github.com/dangLuan01/user-manager/internal/repository"
	"github.com/dangLuan01/user-manager/internal/routes"
	v1routes "github.com/dangLuan01/user-manager/internal/routes/v1"
	v1service "github.com/dangLuan01/user-manager/internal/service/v1"
	"github.com/doug-martin/goqu/v9"
)

type UserModule struct {
	routes routes.Route
}

func NewUserModule(DB *goqu.Database) *UserModule {

	userRepo := repository.NewSqlUserRepository(DB)
	userService := v1service.NewUserService(userRepo)
	UserHandler := v1handler.NewUserHandler(userService)
	userRoutes := v1routes.NewUserRoutes(UserHandler)

	return &UserModule{
		routes: userRoutes,
	}
}
func (m *UserModule) Routes() routes.Route {
	return m.routes
}