package routes

import (
	"github.com/JulieWasNotAvailable/microservices/unpublished/api/handlers"
	"github.com/JulieWasNotAvailable/microservices/unpublished/internal"
	"github.com/JulieWasNotAvailable/microservices/unpublished/internal/consumer"
	"github.com/JulieWasNotAvailable/microservices/unpublished/pkg/beatmetadata"
	"github.com/JulieWasNotAvailable/microservices/unpublished/pkg/unpbeat"
	"github.com/gofiber/fiber/v2"
)

func SetupUnpublishedBeatRoutes(app fiber.Router, service unpbeat.Service, mservice beatmetadata.MetadataService, mfcc_channel <-chan consumer.KafkaMessage, delete_approve_channel <-chan consumer.KafkaMessage){
	unp := app.Group("/unpbeats")
	unp.Post("/saveDraft", internal.ProtectedRequiresBeatmaker(), handlers.SaveBeatDraft(service))
	unp.Get("/sortByStatus/:status", internal.ProtectedRequiresBeatmaker(), handlers.GetBeatsSortByStatusAndJWT(service))
	unp.Get("/unpublishedBeatsByBeatmakerJWT", internal.ProtectedRequiresBeatmaker(), handlers.GetUnpublishedBeatsByBeatmakerId(service))
	unp.Get("/allUnpublishedBeats", handlers.GetAllUnpublishedBeats(service))
	unp.Patch("/updateUnpublishedBeat", handlers.UpdateBeat(service, mservice))
	unp.Get("/beatsForModerationByDate/:from/:to", handlers.GetBeatsInModeration(service))
	unp.Get("/sendToModeration/:id", handlers.SendToModeration(service))
	unp.Get("/publishBeat/:id", handlers.PostPublishBeat(service, mfcc_channel, delete_approve_channel))
}