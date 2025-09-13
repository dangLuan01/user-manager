package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func ApiKeyMiddleware() gin.HandlerFunc {
	expectedApiKey := os.Getenv("API_KEY")
	if expectedApiKey == "" {
		expectedApiKey = os.Getenv("DEFAULT_API_KEY")
	}
	return func(ctx *gin.Context)  {
		apiKey := ctx.GetHeader("X-API-Key")
		if apiKey == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "API key is required",
			})
			return 
		}	
		if apiKey != expectedApiKey {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid API key",
			})
			return 
		}
		ctx.Next()
	}
}