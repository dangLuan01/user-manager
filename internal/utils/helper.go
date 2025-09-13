package utils

import (
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
