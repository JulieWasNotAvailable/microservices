package routers

import (
	"github.com/JulieWasNotAvailable/microservices/beat/api/handlers"
	"github.com/JulieWasNotAvailable/microservices/beat/internal/beat"
	"github.com/gofiber/fiber/v2"
)

func SetupBeatRoutes(app fiber.Router, service beat.Service) {
	beat := app.Group("/beat")
	//helpers
	beat.Post("/exampleBeat", handlers.CreateBeat(service))
	beat.Get("/all", handlers.GetAllBeats(service))

	beat.Get("/byBeatmakerId/:beatmakerId", handlers.GetBeatsByBeatmakerId(service))
	beat.Get("/byBeatId/:beatId", handlers.GetBeatById(service))

	beat.Get("/filteredBeats", handlers.GetFilteredBeats(service))
	beat.Get("/withAllMoods", handlers.GetBeatsWithAllMoods(service))
	beat.Get("/beatsByMoodId/:moodId", handlers.GetBeatsByMoodId(service))
	beat.Get("/beatsByDate/:from/:to", handlers.GetBeatsByDate(service))
}