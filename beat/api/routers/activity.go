package routers

import (
	"github.com/JulieWasNotAvailable/microservices/beat/api/handlers"
	"github.com/JulieWasNotAvailable/microservices/beat/internal/activity"
	"github.com/JulieWasNotAvailable/microservices/beat/pkg/guards"
	"github.com/gofiber/fiber/v2"
)

func SetupActivityRoutes(app fiber.Router, service activity.Service) {
	activity := app.Group("/activity")
	activity.Post("/postNewLike", guards.Protected(), handlers.PostLike(service))
	activity.Delete("/:beatId", guards.Protected(), handlers.DeleteLike(service))
	activity.Get("/viewMyLikes", guards.Protected(), handlers.GetLikesByUserId(service))
	activity.Get("/viewLikesCountByBeatId/:beatId", handlers.GetLikesCountByBeatId(service))
	activity.Get("/viewLikesCountByUserId/:userId", handlers.GetLikesCountByUserId(service))
	activity.Post("/totalLikesCountForBeats", handlers.GetTotalLikesOfBeats(service))

	activity.Post("/listened", guards.Protected(), handlers.PostListened(service))
}
