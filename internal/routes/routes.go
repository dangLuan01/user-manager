package routes

import (
	"github.com/dangLuan01/user-manager/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoute(r *gin.Engine, routes ...Route) {
	api := r.Group("/api/v1")

	api.Use(	
		middleware.ApiKeyMiddleware(),
		middleware.RateLimiterMiddleware(), 
		middleware.AuthMiddleware(),
	)

	for _, route := range routes {
		route.Register(api)
	}

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{
			"error":"NOT FOUND",
			"path": ctx.Request.URL.Path,
		})
	})
}