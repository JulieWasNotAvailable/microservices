package main

import (
	"log"

	"github.com/JulieWasNotAvailable/microservices/cart/api/routers"
	"github.com/JulieWasNotAvailable/microservices/cart/pkg/cart"
	"github.com/JulieWasNotAvailable/microservices/cart/pkg/dbconnection"
	"github.com/JulieWasNotAvailable/microservices/cart/pkg/entities"
	"github.com/JulieWasNotAvailable/microservices/cart/pkg/license"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	pgconfig := dbconnection.GetConfigs()
	db, err := dbconnection.NewConnection(pgconfig)
	if err != nil {
		log.Println(err)
	}
	err = entities.MigrateAll(db)
	if err != nil {
		log.Println(err)
	}

	cartRepo := cart.NewRepo(db)
	cartService := cart.NewService(cartRepo)

	licenseRepo := license.NewRepo(db)
	licenseService := license.NewService(licenseRepo)

	api := app.Group("/api")
	routers.SetupCartRoutes(api, cartService, licenseService)
	routers.SetupLicenseRoutes(api, licenseService)
	
	app.Use(cors.New())

	app.Listen(":7000")
}