package main

import (
	"log"

	"github.com/JulieWasNotAvailable/microservices/user/api/routers"
	_ "github.com/JulieWasNotAvailable/microservices/user/docs"

	"github.com/JulieWasNotAvailable/microservices/user/pkg/consumer"
	"github.com/JulieWasNotAvailable/microservices/user/internal/activity"
	"github.com/JulieWasNotAvailable/microservices/user/internal/bmmetadata"
	"github.com/JulieWasNotAvailable/microservices/user/internal/entities"
	"github.com/JulieWasNotAvailable/microservices/user/internal/user"
	"github.com/JulieWasNotAvailable/microservices/user/pkg/dbconnection"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

// @BasePath /api
// @title Fiber User Service
// @version 1.0
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @host localhost:7773
func main() {

	pgconfig := dbconnection.GetConfigs()
	db, err := dbconnection.NewConnection(pgconfig)
	if err != nil {
		log.Fatal("Database Connection Error $s", err)
	}
	err = entities.MigrateAll(db)
	if err != nil {
		log.Fatal("Cannot Migrate Error $s", err)
	}

	userRepo := user.NewRepo(db)
	metadataRepo := bmmetadata.NewRepo(db)
	activityRepo := activity.NewRepo(db)

	userService := user.NewService(userRepo)
	metadataService := bmmetadata.NewService(metadataRepo)
	activityService := activity.NewService(activityRepo)

	app := fiber.New()
	api := app.Group("/api")
	app.Use(cors.New())

	// app.Get("/swagger/*", swagger.HandlerDefault)
	api.Get("/swagger/*", swagger.New(swagger.Config{}))

	api.Get("/allGenres", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"genres": []string{"Action", "Comedy", "Drama", "Fantasy", "Horror"},
		})
	})

	api.Post("/saveAnketa", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "anketa saved",
		})
	})

	routes.UserRouter(api, userService, metadataService)
	routes.MetadataRoutes(api, metadataService, userService)
	routes.GoogleRoutes(api, userService)
	routes.WelcomeRouter(api)
	routes.ActivityRoutes(api, activityService)

	// go consumer.StartConsumer("profilepic_url_updates", userService)

	app.Listen(":7773")
}
