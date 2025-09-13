package app

import (
	"log"

	"github.com/dangLuan01/user-manager/internal/config"
	"github.com/dangLuan01/user-manager/internal/routes"
	"github.com/dangLuan01/user-manager/internal/validation"
	"github.com/doug-martin/goqu/v9"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Module interface {
	Routes() routes.Route
}

type Application struct {
	config *config.Config
	router *gin.Engine
}

func NewApplication(cfg *config.Config, DB *goqu.Database) *Application {

	if err := validation.InitValidator(); err != nil {
		log.Fatalf("Validation init failed %v:", err)
	}
	
	r := gin.Default()

	modules := []Module{
		NewUserModule(DB),
	}

	routes.RegisterRoute(r, getModuleRoutes(modules)...)
	return &Application{
		config: cfg,
		router: r,
	}
}

func (a *Application) Run() error {
	return a.router.Run(a.config.ServerAddress)
}

func getModuleRoutes(modules []Module) []routes.Route {
	routeList := make([]routes.Route, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}

	return routeList
}
func LoadEnv()  {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}