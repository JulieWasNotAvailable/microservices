package routes

import (
	"github.com/JulieWasNotAvailable/microservices/user/api/handlers"
	"github.com/gofiber/fiber/v2"
)

func WelcomeRouter(app fiber.Router){
	app.Get("/", handlers.Welcome)
}