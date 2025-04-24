package routes

import (
	"github.com/JulieWasNotAvailable/microservices/unpublished/api/handlers"
	"github.com/JulieWasNotAvailable/microservices/unpublished/pkg/beatmetadata"
	"github.com/gofiber/fiber/v2"
)

func SetupMetadataBeatRoutes(app fiber.Router, service beatmetadata.MetadataService) {
    meta := app.Group("/metadata")

    // meta.Post("/files", handlers.PostFile(service))
    meta.Get("/files", handlers.GetAllFiles(service))
    meta.Get("/filesByBeatId/:beatId", handlers.GetAvailableFilesByBeatId(service))
    meta.Patch("/filesByBeatId/:beatId", handlers.UpdateAvailableFilesByBeatId(service))
    // meta.Patch("/files", handlers.UpdateFiles(service))
    // meta.Delete("/singleFile/:fileId/:fileType", handlers.DeleteFileById(service))

    meta.Post("/instruments", handlers.PostInstrument(service))
    meta.Get("/instruments", handlers.GetInstruments(service))
    
    meta.Post("/genres", handlers.PostGenre(service))
    meta.Get("/genres", handlers.GetGenres(service))
    
    meta.Post("/timestamps", handlers.PostTimestamp(service))
    meta.Get("/timestamps", handlers.GetTimestamps(service))
    
    meta.Post("/tags", handlers.PostTag(service))
    meta.Get("/tags", handlers.GetTags(service))
    
    meta.Post("/moods", handlers.PostMood(service))
    meta.Get("/moods", handlers.GetMoods(service))
    
    meta.Post("/keynotes", handlers.PostKeynote(service))
    meta.Get("/keynotes", handlers.GetKeynotes(service))
}