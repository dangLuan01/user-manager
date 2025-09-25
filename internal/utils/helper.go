package utils

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"strconv"
)

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != ""{
		return value
	}
	return defaultValue
}

func GetIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key);
	if value == ""{
		
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {

		return defaultValue
	}
	
	return intValue
}

func GenerateRandomString(lenght int) (string, error) {
	bytes := make([]byte, lenght)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}
