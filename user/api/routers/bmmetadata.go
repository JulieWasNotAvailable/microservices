package routes

import (
	"github.com/JulieWasNotAvailable/microservices/user/api/handlers"
	"github.com/JulieWasNotAvailable/microservices/user/pkg/guards"
	"github.com/JulieWasNotAvailable/microservices/user/internal/bmmetadata"
	"github.com/JulieWasNotAvailable/microservices/user/internal/user"
	"github.com/gofiber/fiber/v2"
)

func MetadataRoutes(app fiber.Router, bmservice bmmetadata.Service, uservice user.Service) {
    app.Post("/metadata", handlers.AddMetadata(bmservice, uservice))
    app.Get("/metadatas", handlers.GetMetadatas(bmservice)) 
    app.Get("/metadataById/:id", handlers.GetMetadataById(bmservice))
	app.Patch("/metadataById/:id", handlers.UpdateMetadataById(bmservice))
    app.Delete("/metadataById/:id", guards.Protected(), handlers.RemoveMetadata(bmservice))
}
