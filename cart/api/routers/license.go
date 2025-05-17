package routers

import (
	"github.com/JulieWasNotAvailable/microservices/cart/api/handlers"
	"github.com/JulieWasNotAvailable/microservices/cart/pkg/guards"
	"github.com/JulieWasNotAvailable/microservices/cart/internal/license"
	"github.com/gofiber/fiber/v2"
)

func SetupLicenseRoutes(app fiber.Router, service license.Service) {
	license := app.Group("/license")
	license.Post("/newLicense", guards.ProtectedRequiresBeatmaker(), handlers.PostNewLicense(service))
	license.Get("/licensesForBeat/:beatId", handlers.GetLicensesForBeat(service))

	license.Post("/newLicenseTemplate", guards.ProtectedRequiresBeatmaker(), handlers.PostNewLicenseTemplate(service))
	license.Get("/allLicenseTemplatesByBeatmakerJWT", guards.ProtectedRequiresBeatmaker(), handlers.GetAllLicenseTemplatesByBeatmakerId(service))
	license.Patch("/licenseTemplate", guards.ProtectedRequiresBeatmaker(), handlers.PatchLicenseTemplate(service))

	//admin
	license.Get("/allLicenses", handlers.GetAllLicense(service))	
	license.Get("/allLicenseTemplates", handlers.GetAllLicenseTemplate(service))	
}