package middleware

import (
	"net/http"
	"strings"

	"github.com/dangLuan01/user-manager/pkg/auth"
	"github.com/dangLuan01/user-manager/pkg/cache"
	"github.com/gin-gonic/gin"
)

var (
	jwtService auth.TokenService
	cacheService cache.RedisCacheService
)

func InitAuthMiddlware(service auth.TokenService, cache cache.RedisCacheService) {
	jwtService = service
	cacheService = cache
}

func AuthMiddleware() gin.HandlerFunc {
	return func (ctx *gin.Context)  {
		authHeder := ctx.GetHeader("Authorization")
		if authHeder == "" || !strings.HasPrefix(authHeder, "Bearer ") {

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header missing or invalid(1)",
			})
			return 
		}

		tokenString := strings.TrimPrefix(authHeder, "Bearer ")

		_, claims, err := jwtService.ParseToken(tokenString)
		if err != nil {

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header missing or invalid(2)",
			})
			return 
		}

		if jti, ok := claims["jti"].(string); ok {
			key := "blacklist:" + jti
			exists, err := cacheService.Exits(key)
			if err == nil && exists {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Token revoked",
				})
				return 
			}
		}

		payload, err := jwtService.DecryptAccessTokenPayload(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header missing or invalid(2)",
			})
			return 
		}
		ctx.Set("data", payload)
		
		ctx.Next()
		
	}
}