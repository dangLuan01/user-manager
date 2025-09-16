package middleware

import (
	"net/http"
	"strings"

	"github.com/dangLuan01/user-manager/pkg/auth"
	"github.com/gin-gonic/gin"
)

var jwtService auth.TokenService

func InitAuthMiddlware(service auth.TokenService) {
	jwtService = service
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
		_, _, err := jwtService.ParseToken(tokenString)
		if err != nil {

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header missing or invalid(2)",
			})
			return 
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