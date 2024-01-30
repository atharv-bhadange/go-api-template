package routes

import (
	"github.com/atharvbhadange/go-api-template/api/v1/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", controllers.Health)

	v1API := app.Group("/api/v1")

	SetupProductsRoutes(v1API)
}
