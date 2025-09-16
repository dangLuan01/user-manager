package config

import (
	"fmt"

	"github.com/dangLuan01/user-manager/internal/utils"
)

type DatabaseConfig struct {
	Host string
	Port string
	User string
	Password string
	DBName string
	SSLMode string
}

type Config struct {
	ServerAddress string
	DB DatabaseConfig
}

func NewConfig() *Config {
	
	return &Config{
		ServerAddress: fmt.Sprintf(":%s", utils.GetEnv("PORT", "8080")),
		DB: DatabaseConfig {
			Host: utils.GetEnv("DB_HOST","localhost"),
			Port: utils.GetEnv("DB_PORT","3306"),
			User: utils.GetEnv("DB_USER",""),
			Password: utils.GetEnv("DB_PASSWORD",""),
			DBName: utils.GetEnv("DB_DBNAME","mysql"),
			SSLMode: utils.GetEnv("DB_SSLMODE","disable"),
		},
	}
}

func (c *Config) DNS() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
    	c.DB.User, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.DBName,
	)
}