package main

import (
	"log"

	"github.com/JulieWasNotAvailable/microservices/beatsUpload/api/router"
	_ "github.com/JulieWasNotAvailable/microservices/beatsUpload/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

type UpdateBeatURLRequest struct {
	Url string
}

// @BasePath /api
// @title Beats Upload Service
// @version 1.0
// @description Deals with presigned requests. Pushes updates to Beats and User microservice, when files are uploaded.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @host localhost:7773
func main() {
	app := fiber.New()

	api := app.Group("/api")
	app.Use(cors.New())

	router.SetupRoutes(app)

	log.Println("app running successfully")
	api.Get("/swagger/*", swagger.New(swagger.Config{}))

	app.Listen(":7774")
}
