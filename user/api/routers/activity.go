package routes

import (
	"github.com/JulieWasNotAvailable/microservices/user/api/handlers"
	"github.com/JulieWasNotAvailable/microservices/user/internal/activity"
	"github.com/JulieWasNotAvailable/microservices/user/pkg/guards"
	"github.com/gofiber/fiber/v2"
)

func ActivityRoutes(app fiber.Router, service activity.Service) {
	activity := app.Group("/activity")
	activity.Post("/subscribeTo/:beatmakerId", guards.Protected(), handlers.PostSubscribeToBeatmaker(service))
	activity.Get("/viewMySubscriptions", guards.Protected(), handlers.GetMySubscriptions(service))
	activity.Get("/followersNumberByBeatmakerId/:beatmakerId", handlers.GetFollowersCountByBeatmakerId(service))
	activity.Delete("unsubscribe/:beatmakerId", guards.Protected(), handlers.DeleteSub(service))
}
