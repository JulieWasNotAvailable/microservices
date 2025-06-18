package main

import (
	"log"

	router "github.com/JulieWasNotAvailable/microservices/beatsUpload/api/routers"
	_ "github.com/JulieWasNotAvailable/microservices/beatsUpload/docs"

	// "github.com/JulieWasNotAvailable/microservices/beatsUpload/internal/beat"
	// "github.com/JulieWasNotAvailable/microservices/beatsUpload/pkg/consumer"
	// "github.com/JulieWasNotAvailable/microservices/beatsUpload/pkg/dbconnection"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

type UpdateBeatURLRequest struct {
	Url string
}

//	@BasePath					/api
//	@title						Beats Upload Service
//	@version					1.0
//	@description				Deals with presigned requests. Pushes updates to Beats and User microservice, when files are uploaded.
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@host						localhost:7773
func main() {
	app := fiber.New()

	api := app.Group("/api")
	app.Use(cors.New())

	router.SetupRoutes(app)

	log.Println("app running successfully")
	api.Get("/swagger/*", swagger.New(swagger.Config{}))

	// s3 := dbconnection.S3Connect()
	// beatService := beat.NewService(s3)

	// consumer.StartConsumer("beat_deleted", beatService)

	app.Listen(":7774")
}
