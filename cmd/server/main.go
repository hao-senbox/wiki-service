package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wiki-service/internal/app"
)

func main() {
	container, err := app.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	addr := fmt.Sprintf("%s:%s", container.Config.Server.Host, container.Config.Server.Port)
	container.Logger.Info(fmt.Sprintf("Server starting on %s", addr))
	log.Printf("Server starting on %s", addr)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := container.App.Listen(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-done
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	container.ConsulConn.Deregister()
	if err := container.App.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited gracefully")
}
