package main

import (
	"log"

	"github.com/JulieWasNotAvailable/microservices/beatsUpload/router"
	s3files "github.com/JulieWasNotAvailable/microservices/beatsUpload/s3Files"
	"github.com/gofiber/fiber/v2"

	"github.com/JulieWasNotAvailable/microservices/beatsUpload/dbconnection"
)

type UpdateBeatURLRequest struct {
	Url string
}

func main () {
	app := fiber.New()

	s3client := dbconnection.S3Connect()

	presignClient := s3files.S3ConnectPresign(s3client)

	presignClient.SetupRoutes(app)
	router.SetupRoutes(app)
	
	log.Println("app running successully")

	app.Listen(":7774")
}
