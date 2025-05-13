package main

import (
	"log"
	"os"

	"github.com/JulieWasNotAvailable/microservices/beat/api/routers"
	"github.com/JulieWasNotAvailable/microservices/beat/internal/beat"
	"github.com/JulieWasNotAvailable/microservices/beat/internal/entities"
	"github.com/JulieWasNotAvailable/microservices/beat/internal/metadata"
	"github.com/JulieWasNotAvailable/microservices/beat/pkg/consumer"
	"github.com/JulieWasNotAvailable/microservices/beat/pkg/dbconnection"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/JulieWasNotAvailable/microservices/beat/docs"
)

// @BasePath /api
// @title Fiber Beat Service
// @version 1.0
// @description Deals Beats, Beat_Listened, filtering
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @host localhost:7771
func main() {
	app := fiber.New()
	
	pgconfig := dbconnection.GetConfigs()
	db, err := dbconnection.NewConnection(pgconfig)
	if err != nil{
		log.Fatal(err)
	}
	err = entities.MigrateAll(db)
	if err != nil{
		log.Fatal(err)
	}

	metaRepo := metadata.NewRepo(db)
	metaService := metadata.NewService(metaRepo)

	beatRepo := beat.NewRepo(db)
	beatService := beat.NewService(beatRepo)

	api := app.Group("/api")
	routers.SetupMetadataBeatRoutes(api, metaService)
	routers.SetupBeatRoutes(api, beatService)
	api.Get("/swagger/*", swagger.New(swagger.Config{}))

	go consumer.StartConsumerPublisher(os.Getenv("KAFKA_PUBLISH_TOPIC"), beatService)

	app.Listen(":7771")
}