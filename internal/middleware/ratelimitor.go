package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/dangLuan01/user-manager/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Client struct {
	limiter *rate.Limiter
	lastSeen time.Time
}

var (
	clients = make(map[string]*Client)
	mu sync.Mutex
)

func getClientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()
	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}
	return ip
}

func getRateLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	client, exists := clients[ip]

	if !exists {
		requestSec := utils.GetIntEnv("RATE_LIMITER_REQUEST_SEC", 5)
		brust := utils.GetIntEnv("RATE_LIMITER_REQUEST_BRUST", 15)

		limiter := rate.NewLimiter(rate.Limit(requestSec), brust)
		client = &Client{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		clients[ip] = client

		return limiter
	}

	client.lastSeen = time.Now()

	return client.limiter
}

func CleanupClients() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, client := range clients {
			if time.Since(client.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

func RateLimiterMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientIP := getClientIP(ctx)
		limiter := getRateLimiter(clientIP)
		if !limiter.Allow() {
			log.Printf("Rate limit exceeded for IP: %s", clientIP)
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests, please try again later.",
			})
			return
		}
		ctx.Next()
	}
}