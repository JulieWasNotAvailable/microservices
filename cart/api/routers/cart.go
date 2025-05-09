package routers

import (
	"github.com/JulieWasNotAvailable/microservices/cart/api/handlers"
	"github.com/JulieWasNotAvailable/microservices/cart/internal"
	"github.com/JulieWasNotAvailable/microservices/cart/pkg/cart"
	"github.com/JulieWasNotAvailable/microservices/cart/pkg/license"
	"github.com/gofiber/fiber/v2"
)

func SetupCartRoutes(app fiber.Router, service cart.Service, licenseService license.Service) {
	cart := app.Group("/cart")
	cart.Get("/hello", handlers.Hello(service))
	cart.Get("/addLicenseToCart/:licenseId", internal.Protected(), handlers.PostAddToCart(service, licenseService))
	cart.Get("/getByJWT", internal.Protected(), handlers.GetCartByUser(service))
	cart.Delete("/deleteLicense/:licenseId", internal.Protected(), handlers.DeleteLicenseFromCart(service))
}
