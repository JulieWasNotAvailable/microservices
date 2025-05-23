package routes

import (
	"github.com/JulieWasNotAvailable/microservices/unpublished/api/handlers"
	"github.com/JulieWasNotAvailable/microservices/unpublished/pkg/guards"
	"github.com/JulieWasNotAvailable/microservices/unpublished/pkg/consumer"
	"github.com/JulieWasNotAvailable/microservices/unpublished/internal/beatmetadata"
	"github.com/JulieWasNotAvailable/microservices/unpublished/internal/unpbeat"
	"github.com/gofiber/fiber/v2"
)

func SetupUnpublishedBeatRoutes(app fiber.Router, service unpbeat.Service, mservice beatmetadata.MetadataService, mfcc_channel <-chan consumer.KafkaMessageValue, delete_approve_channel <-chan consumer.KafkaMessageValue){
	unp := app.Group("/unpbeats")
	unp.Post("/makeEmptyBeat", guards.ProtectedRequiresBeatmaker(), handlers.SaveBeatDraft(service, mservice))
	unp.Get("/sortByStatus/:status", guards.ProtectedRequiresBeatmaker(), handlers.GetBeatsSortByStatusAndJWT(service))
	unp.Get("/unpublishedBeatsByBeatmakerJWT", guards.ProtectedRequiresBeatmaker(), handlers.GetUnpublishedBeatsByBeatmakerId(service))
	unp.Get("/unpublishedBeatById/:id", handlers.GetUnpublishedBeatById(service))
	unp.Get("/allUnpublishedBeats", handlers.GetAllUnpublishedBeats(service))
	unp.Patch("/saveDraft", handlers.UpdateBeat(service, mservice))
	unp.Get("/beatsForModerationByDate/:from/:to", handlers.GetBeatsInModeration(service))
	unp.Get("/sendToModeration/:id", handlers.SendToModeration(service))
	unp.Get("/publishBeat/:id", handlers.PostPublishBeat(service, mfcc_channel, delete_approve_channel))
	unp.Delete("/deleteUnpublishedBeatById/:id", handlers.DeleteUnpublishedBeatById(service))
}
