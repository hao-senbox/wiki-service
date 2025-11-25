package main

import (
	"fmt"
	"log"

	"wiki-service/internal/app"
)

func main() {
	// Initialize application container with all dependencies
	container, err := app.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Log server start
	addr := fmt.Sprintf("%s:%s", container.Config.Server.Host, container.Config.Server.Port)
	container.Logger.Info(fmt.Sprintf("Server starting on %s", addr))
	log.Printf("Server starting on %s", addr)

	// Start HTTP server
	if err := container.App.Listen(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
