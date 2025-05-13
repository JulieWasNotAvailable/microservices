package routers

import (
	"github.com/JulieWasNotAvailable/microservices/cart/api/handlers"
	"github.com/JulieWasNotAvailable/microservices/cart/pkg/guards"
	"github.com/JulieWasNotAvailable/microservices/cart/internal/cart"
	"github.com/JulieWasNotAvailable/microservices/cart/internal/license"
	"github.com/gofiber/fiber/v2"
)

func SetupCartRoutes(app fiber.Router, service cart.Service, licenseService license.Service) {
	cart := app.Group("/cart")
	cart.Get("/hello", handlers.Hello(service))
	cart.Get("/addLicenseToCart/:licenseId", guards.Protected(), handlers.PostAddToCart(service, licenseService))
	cart.Get("/getByJWT", guards.Protected(), handlers.GetCartByUser(service))
	cart.Delete("/deleteLicense/:licenseId", guards.Protected(), handlers.DeleteLicenseFromCart(service))
}
