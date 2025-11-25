package http

import (
	"wiki-service/internal/interface/http/handler"
	"wiki-service/internal/interface/http/route"
	"wiki-service/internal/interface/middleware"
	"wiki-service/pkg/gateway"

	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRouter sets up the Fiber router
func SetupRouter(
	wikiHandler *handler.WikiHandler,
	auditMiddleware *middleware.AuditMiddleware,
	userGateway gateway.UserGateway,
) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Services Management v1.0",
	})

	// Apply global middlewares
	app.Use(fiberLogger.New())
	app.Use(middleware.LoggingMiddleware())
	app.Use(middleware.CORSMiddleware())
	app.Use(auditMiddleware.Log())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	route.SetUpWikiRoutes(app, wikiHandler, userGateway)

	return app
}
