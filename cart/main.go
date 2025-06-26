package main

import (
	"log"

	"github.com/JulieWasNotAvailable/microservices/cart/api/routers"
	_ "github.com/JulieWasNotAvailable/microservices/cart/docs"
	"github.com/JulieWasNotAvailable/microservices/cart/internal/cart"
	"github.com/JulieWasNotAvailable/microservices/cart/internal/entities"
	"github.com/JulieWasNotAvailable/microservices/cart/internal/license"
	"github.com/JulieWasNotAvailable/microservices/cart/pkg/consumer"
	"github.com/JulieWasNotAvailable/microservices/cart/pkg/dbconnection"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

// @BasePath					/api
// @title						Fiber Cart Service
// @version					1.0
// @description				Deals with Cart, Licenses and License Template
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @host						localhost:7000
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
	app.Use(cors.New())

	routers.SetupCartRoutes(api, cartService, licenseService)
	routers.SetupLicenseRoutes(api, licenseService)
	api.Get("/swagger/*", swagger.New(swagger.Config{}))

	app.Use(cors.New())

	appQuit := make(chan bool)
	go consumer.StartConsumer("create_license", licenseService, appQuit)

	go func() {
		if err := app.Listen(":7775"); err != nil {
			log.Printf("Server crashed: %v", err) // Don't use Fatalf (avoids os.Exit)
		}
	}()

	<-appQuit
	app.Server().Shutdown()
}
