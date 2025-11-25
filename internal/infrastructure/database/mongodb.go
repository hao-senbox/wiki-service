package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoConfig holds MongoDB configuration
type MongoConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// NewMongoConnection creates a new MongoDB connection
func NewMongoConnection(config MongoConfig) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var uri string

	// Build connection URI based on whether authentication is needed
	if config.User != "" && config.Password != "" {
		// With authentication
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s",
			config.User,
			config.Password,
			config.Host,
			config.Port,
		)
	} else {
		// Without authentication (local development)
		uri = fmt.Sprintf("mongodb://%s:%s", config.Host, config.Port)
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client.Database(config.DBName), nil
}
