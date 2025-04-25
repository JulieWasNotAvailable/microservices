package main

import (
	"log"

	"github.com/JulieWasNotAvailable/microservices/user/api/routes"
	_ "github.com/JulieWasNotAvailable/microservices/user/docs"
	"github.com/JulieWasNotAvailable/microservices/user/internal/consumer"
	"github.com/JulieWasNotAvailable/microservices/user/pkg/bmmetadata"
	"github.com/JulieWasNotAvailable/microservices/user/pkg/dbconnection"
	"github.com/JulieWasNotAvailable/microservices/user/pkg/entities"
	"github.com/JulieWasNotAvailable/microservices/user/pkg/user"
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
	err = entities.MigrateUser(db)
	if err != nil {
		log.Fatal("Cannot Migrate User Error $s", err)
	}
	err = entities.MigrateMetadata(db)
	if err != nil {
		log.Fatal("Cannot Migrate Metadata Error $s", err)
	}
	err = entities.MigrateRole(db)
	if err != nil {
		log.Fatal("Cannot Migrate Role Error $s", err)
	}

	userRepo := user.NewRepo(db)
	metadataRepo := bmmetadata.NewRepo(db)
	userService := user.NewService(userRepo)
	metadataService := bmmetadata.NewService(metadataRepo)

	app := fiber.New()
	api := app.Group("/api")
	app.Use(cors.New())

	// app.Get("/swagger/*", swagger.HandlerDefault)
	api.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		// Prefill OAuth ClientId on Authorize popup
		// OAuth: &swagger.OAuthConfig{
		// 	AppName:  "OAuth Provider",
		// 	ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		// },
		// Ability to change OAuth2 redirect uri location
		// OAuth2RedirectUrl: "http://localhost:8080/swagger/oauth2-redirect.html",
	}))

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

	go consumer.StartConsumer("profilepic_url_updates", userService)

	app.Listen(":7773")
}
