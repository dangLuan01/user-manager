package routes

import (
	"github.com/dangLuan01/user-manager/internal/middleware"
	v1routes "github.com/dangLuan01/user-manager/internal/routes/v1"
	"github.com/dangLuan01/user-manager/pkg/auth"
	"github.com/gin-gonic/gin"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoute(r *gin.Engine, authService auth.TokenService, routes ...Route) {
	v1api := r.Group("/api/v1")

	v1api.Use(	
		middleware.ApiKeyMiddleware(),
		middleware.RateLimiterMiddleware(), 
	)
	middleware.InitAuthMiddlware(authService)
	protected := v1api.Group("")
	protected.Use(
		middleware.AuthMiddleware(),
	)

	for _, route := range routes {

		switch route.(type) {
		case *v1routes.AuthRoutes:
			route.Register(v1api)
		default:
			route.Register(protected)
		}
		
	}

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{
			"error":"NOT FOUND",
			"path": ctx.Request.URL.Path,
		})
	})
}