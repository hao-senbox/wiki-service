package route

import (
	"wiki-service/internal/interface/http/handler"
	"wiki-service/internal/interface/middleware"
	"wiki-service/pkg/gateway"

	"github.com/gofiber/fiber/v2"
)

func SetUpWikiRoutes(app *fiber.App, serviceHandler *handler.WikiHandler, userGateway gateway.UserGateway) {
	api := app.Group("/api/v1")
	api.Use(middleware.Secured(userGateway))

	wikiGroups := api.Group("/wikis")
	{
		wikiGroups.Post("/template", serviceHandler.CreateWikiTemplate)
		wikiGroups.Get("/statistics", serviceHandler.GetStatistics)
		wikiGroups.Get("/", serviceHandler.GetWikis)
		wikiGroups.Get("/:id", serviceHandler.GetWikiByID)
		wikiGroups.Put("/:id", serviceHandler.UpdateWiki)
	}

}
