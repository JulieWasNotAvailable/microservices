package main

import (
	"log"

	"github.com/JulieWasNotAvailable/microservices/beatsUpload/api/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/JulieWasNotAvailable/microservices/beatsUpload/docs"
)

type UpdateBeatURLRequest struct {
	Url string
}

// @BasePath /api
// @title Fiber Presigner Service
// @version 1.0
// @description Deals with presigned requests. Pushes updates to Beats and User microservice, when files are uploaded.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @host localhost:7773
func main () {
	app := fiber.New()

	router.SetupRoutes(app)
	
	log.Println("app running successully")
	app.Get("/swagger/*", swagger.New(swagger.Config{}))

	app.Listen(":7774")
}
