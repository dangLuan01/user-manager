package app

import (
	"log"

	"github.com/dangLuan01/user-manager/internal/config"
	"github.com/dangLuan01/user-manager/internal/db"
	"github.com/dangLuan01/user-manager/internal/routes"
	"github.com/dangLuan01/user-manager/internal/validation"
	"github.com/dangLuan01/user-manager/pkg/auth"
	"github.com/dangLuan01/user-manager/pkg/cache"
	"github.com/doug-martin/goqu/v9"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type Module interface {
	Routes() routes.Route
}

type Application struct {
	config *config.Config
	router *gin.Engine
	modules []Module
}

type ModuleContext struct {
	DB *goqu.Database
	Redis *redis.Client
}

func NewApplication(cfg *config.Config) *Application {

	if err := validation.InitValidator(); err != nil {
		log.Fatalf("⛔ Validation init failed %v:", err)
	}
	
	r := gin.Default()

	if err := db.InitDB(); err != nil {
		log.Fatalf("⛔ Unable to connect to sql")
	}

	redisClient := config.NewRedisClient()
	cacheRedisService := cache.NewRedisCacheService(redisClient)
	tokenService := auth.NewJWTService(cacheRedisService)
	
	ctx := &ModuleContext{
		DB: db.DB,
		Redis: redisClient,
	}

	modules := []Module{
		NewUserModule(ctx),
		NewAuthModule(ctx, tokenService),
	}

	routes.RegisterRoute(r, tokenService ,getModuleRoutes(modules)...)

	return &Application{
		config: cfg,
		router: r,
		modules: modules,
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