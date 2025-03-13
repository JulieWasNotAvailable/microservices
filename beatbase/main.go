package main

import (
	"log"

	"github.com/JulieWasNotAvailable/goBeatsBackend/beat"
	"github.com/JulieWasNotAvailable/goBeatsBackend/consumer"
	"github.com/JulieWasNotAvailable/goBeatsBackend/dbconnection"
	"github.com/JulieWasNotAvailable/goBeatsBackend/model"
	"github.com/JulieWasNotAvailable/goBeatsBackend/s3storage"
	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/cors"
)


func main() {

	//db connection

	pgconfig := dbconnection.NewConfig()

	pgdb, err := dbconnection.NewConnection(pgconfig) //postgres database

	if err != nil{
		log.Fatal("could not reach the DB")
	}
	
	err = model.MigrateBeats(pgdb)

	if err != nil {
		log.Fatal("could not migrate")
	}

	pgrepo := beat.Repository{
		DB: pgdb,
	}
	
	//s3 storage connection
	storage := s3storage.S3Connect()

	var CorsConfig = cors.Config{
		AllowOrigins: "*",
		AllowMethods: "*",
		AllowHeaders: "*",
	}

	app := fiber.New()
	app.Use(cors.New(CorsConfig))

	//setup routing
	pgrepo.SetupRoutes(app)
	storage.SetupRoutes(app)

	go consumer.StartConsumer("beat_url_updates")
	
	app.Listen(":7776")	
}