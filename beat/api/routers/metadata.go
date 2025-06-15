package routers

import (
	"github.com/JulieWasNotAvailable/microservices/beat/api/handlers"
	"github.com/JulieWasNotAvailable/microservices/beat/internal/metadata"
	"github.com/gofiber/fiber/v2"
)

func SetupMetadataBeatRoutes(app fiber.Router, service metadata.Service) {
	meta := app.Group("/metadata")
	meta.Get("/genres", handlers.GetAllGenres(service))
	meta.Get("/moods", handlers.GetAllMoods(service))
	meta.Get("/keys", handlers.GetAllKeys(service))
	// meta.Get("/instruments", handlers.GetAllInstruments(service))
	// meta.Get("/files", handlers.GetAllAvailableFiles(service))
	meta.Get("/mfccs", handlers.GetAllMFCCs(service))

	// Tags endpoints
	meta.Get("/tags", handlers.GetAllTags(service))
	meta.Get("/tags/random", handlers.GetRandomTags(service))
	meta.Get("/tags/byName/:name", handlers.GetTagByName(service))
	meta.Get("/tags/byNameLike/:name", handlers.GetTagsByNameLike(service))
	meta.Get("/tags/in_trend", handlers.GetTagsInTrend(service))

	// Timestamps endpoint
	meta.Get("/timestamps", handlers.GetAllTimestamps(service))

	// Genres popularity endpoint
	meta.Get("/genres/popular", handlers.GetGenresInTrend(service))
}
