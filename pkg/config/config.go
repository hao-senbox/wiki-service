package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	MongoDB  MongoDBConfig
	Consul   ConsulConfig
	Registry RegistryConfig
	Database DatabaseConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host string
	Port string
}

// MongoDBConfig holds MongoDB configuration
type MongoDBConfig struct {
	Host     string
	Port     string
	User     string // Optional - leave empty for no auth
	Password string // Optional - leave empty for no auth
	DBName   string
}

// ConsulConfig holds Consul configuration
type ConsulConfig struct {
	Host string
	Port int
}

// RegistryConfig holds service registry configuration
type RegistryConfig struct {
	Host string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	RedisCache RedisCacheConfig
}

// RedisCacheConfig holds Redis cache configuration
type RedisCacheConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if exists (errors ignored)
	_ = godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		MongoDB: MongoDBConfig{
			Host:     getEnv("MONGO_HOST", "localhost"),
			Port:     getEnv("MONGO_PORT", "27017"),
			User:     getEnv("MONGO_USER", ""),     // Empty = no authentication
			Password: getEnv("MONGO_PASSWORD", ""), // Empty = no authentication
			DBName:   getEnv("MONGO_DB_NAME", "services_management"),
		},
		Consul: ConsulConfig{
			Host: getEnv("CONSUL_HOST", "localhost"),
			Port: getEnvAsInt("CONSUL_PORT", 8500),
		},
		Registry: RegistryConfig{
			Host: getEnv("REGISTRY_HOST", "localhost"),
		},
		Database: DatabaseConfig{
			RedisCache: RedisCacheConfig{
				Host:     getEnv("REDIS_HOST", "localhost"),
				Port:     getEnv("REDIS_PORT", "6379"),
				Password: getEnv("REDIS_PASSWORD", ""),
				DB:       getEnvAsInt("REDIS_DB", 0),
			},
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	var value int
	if _, err := fmt.Sscanf(valueStr, "%d", &value); err != nil {
		return defaultValue
	}
	return value
}
