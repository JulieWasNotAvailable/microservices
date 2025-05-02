package main

import (
	"log"
	// "os"

	"github.com/JulieWasNotAvailable/microservices/beat/api/routers"
	// "github.com/JulieWasNotAvailable/microservices/beat/consumer"
	"github.com/JulieWasNotAvailable/microservices/beat/pkg/beat"
	"github.com/JulieWasNotAvailable/microservices/beat/pkg/dbconnection"
	"github.com/JulieWasNotAvailable/microservices/beat/pkg/entities"
	"github.com/JulieWasNotAvailable/microservices/beat/pkg/metadata"
	"github.com/gofiber/fiber/v2"
)

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

	// go consumer.StartConsumerPublisher(os.Getenv("KAFKA_PUBLISH_TOPIC"), beatService)

	app.Listen(":7771")
}