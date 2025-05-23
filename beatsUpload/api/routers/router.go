package router

import (
	"github.com/JulieWasNotAvailable/microservices/beatsUpload/api/handlers"
	"github.com/JulieWasNotAvailable/microservices/beatsUpload/pkg/guards"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/updateURL/:entity/:filetype", handler.UpdateFile)
	api.Post("/checkFileUpdateUrl", handler.UpdateFile)
	api.Get("/buckets", handler.GetBuckets)
	api.Get("/objectsFromBucket/:bucket", handler.GetObjectsFromBucket)
	api.Get("/headObject/:bucket/:key", handler.GetHeadObject)
	
	api = app.Group("api/presigned")
	api.Post("/presignedGetRequest/:bucket", guards.Protected(), handler.GetObject)
	api.Post("/presignedPostRequest/:bucket", guards.Protected(), handler.PutObject)
	api.Post("/presignedDeleteRequest/:bucket", guards.Protected(), handler.DeleteObject)
}