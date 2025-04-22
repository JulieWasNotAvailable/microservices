package router

import (
	"github.com/JulieWasNotAvailable/microservices/beatsUpload/api/handler"
	// "github.com/JulieWasNotAvailable/microservices/beatsUpload/internal"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/checkFileUpdateUrl", handler.CheckFileAvailability)
	api.Get("/buckets", handler.GetBuckets)
	api.Get("/objectsFromBucket/:bucket", handler.GetObjectsFromBucket)
	api.Get("/headObject/:bucket/:key", handler.GetHeadObject)
	api = app.Group("api/presigned")
	api.Post("/presignedGetRequest/:bucket", handler.GetObject)
	api.Post("/presignedPostRequest/:bucket", handler.PutObject)
	api.Post("/presignedDeleteRequest/:bucket", handler.DeleteObject)
}