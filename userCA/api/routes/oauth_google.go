package routes

import (
	"github.com/JulieWasNotAvailable/microservices/user/api/handlers"
	"github.com/JulieWasNotAvailable/microservices/user/pkg/user"

	"github.com/gofiber/fiber/v2"
)

// Routes for fiber
func GoogleRoutes(app fiber.Router, service user.Service) {
	app.Post("/auth/google/getjwt", handlers.HandleGoogleAuth(service))
}
