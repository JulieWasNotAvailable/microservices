package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/JulieWasNotAvailable/microservices/beatsUpload/updateurl"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/checkFileUpdateUrl", updateurl.CheckFileAvailability)
}