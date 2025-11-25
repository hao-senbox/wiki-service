package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Log after request
		duration := time.Since(start)
		log.Printf(
			"%s %s %d %s",
			c.Method(),
			c.OriginalURL(),
			c.Response().StatusCode(),
			duration,
		)

		return err
	}
}
