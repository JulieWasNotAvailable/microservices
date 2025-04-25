package main

import (
	"log"

	"github.com/JulieWasNotAvailable/microservices/unpublished/api/routes"
	"github.com/JulieWasNotAvailable/microservices/unpublished/internal/consumer"
	"github.com/JulieWasNotAvailable/microservices/unpublished/pkg/beatmetadata"
	"github.com/JulieWasNotAvailable/microservices/unpublished/pkg/dbconnection"
	"github.com/JulieWasNotAvailable/microservices/unpublished/pkg/entities"
	"github.com/JulieWasNotAvailable/microservices/unpublished/pkg/unpbeat"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	_ "github.com/JulieWasNotAvailable/microservices/unpublished/docs"
)

// @BasePath /api
// @title Fiber Unpublished Beats Service
// @version 1.0
// @description Deals with unpublished beats and moderation. Has its own duplicates of Tags, Genre, and other metadata, except MFCC characteristics. When user publishes a beat, sends a kafka message to get MFCC characteristics from the track. If the user publishes a new tag, tag is updated in unpublished service and then in beats service.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @host localhost:7772
func main () {
	pgconfig := dbconnection.GetConfigs()
	db, err := dbconnection.NewConnection(pgconfig) 
	if err != nil {
		log.Fatal("Database Connection Error $s", err)
	}
	//migrate
	err = entities.MigrateAll(db)
	if err != nil {
		log.Fatal("Couldn't migrate", err)
	}

	unpBeatRepo := unpbeat.NewRepo(db)
	metadataBeatRepo := beatmetadata.NewRepo(db)
	unpBeatService := unpbeat.NewService(unpBeatRepo)
	metadataBeatService := beatmetadata.NewMetadataService(metadataBeatRepo)

	app := fiber.New()
	api := app.Group("/api")
	mfcc_channel := make(chan consumer.KafkaMessage)
	delete_approve_channel := make(chan consumer.KafkaMessage)

	routes.SetupUnpublishedBeatRoutes(api, unpBeatService, metadataBeatService, mfcc_channel, delete_approve_channel)
	routes.SetupMetadataBeatRoutes(api, metadataBeatService)
	api.Get("/swagger/*", swagger.New(swagger.Config{}))
	api.Get("", func(c *fiber.Ctx) error {
		return c.JSON("Welcome to unpublished beats service")
	} )
	
	go consumer.StartConsumer("beat_mfcc2", mfcc_channel)
	go consumer.StartConsumer("delete_approve", delete_approve_channel)
	go consumer.StartConsumerFileUpdate("beat_files_updates", unpBeatService, metadataBeatService)

	app.Listen(":7772")
}