package main

import (
	"log"

	"github.com/JulieWasNotAvailable/microservices/beat/api/routers"
	_ "github.com/JulieWasNotAvailable/microservices/beat/docs"
	"github.com/JulieWasNotAvailable/microservices/beat/internal/activity"
	"github.com/JulieWasNotAvailable/microservices/beat/internal/beat"
	"github.com/JulieWasNotAvailable/microservices/beat/internal/entities"
	"github.com/JulieWasNotAvailable/microservices/beat/internal/metadata"
	"github.com/JulieWasNotAvailable/microservices/beat/pkg/consumer"
	"github.com/JulieWasNotAvailable/microservices/beat/pkg/dbconnection"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

//	@BasePath					/api
//	@title						Fiber Beat Service
//	@version					1.0
//	@description				Deals Beats, Beat_Listened, filtering
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@host						localhost:7771
func main() {
	app := fiber.New()

	pgconfig := dbconnection.GetConfigs()
	db, err := dbconnection.NewConnection(pgconfig)
	if err != nil {
		log.Fatal(err)
	}
	err = entities.MigrateAll(db)
	if err != nil {
		log.Fatal(err)
	}

	metaRepo := metadata.NewRepo(db)
	metaService := metadata.NewService(metaRepo)

	beatRepo := beat.NewRepo(db)
	beatService := beat.NewService(beatRepo)

	activityRepo := activity.NewRepo(db)
	activityService := activity.NewService(activityRepo)

	api := app.Group("/api")
	app.Use(cors.New())

	routers.SetupMetadataBeatRoutes(api, metaService)
	routers.SetupBeatRoutes(api, beatService)
	routers.SetupActivityRoutes(api, activityService)
	api.Get("/swagger/*", swagger.New(swagger.Config{}))

	go consumer.StartConsumerPublisher("publish_beat_main", beatService)

	app.Listen(":7771")
}
