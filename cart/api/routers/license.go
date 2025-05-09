package routers

import (
	"github.com/JulieWasNotAvailable/microservices/cart/api/handlers"
	"github.com/JulieWasNotAvailable/microservices/cart/internal"
	"github.com/JulieWasNotAvailable/microservices/cart/pkg/license"
	"github.com/gofiber/fiber/v2"
)

func SetupLicenseRoutes(app fiber.Router, service license.Service) {
	license := app.Group("/license")
	license.Post("/newLicense", internal.ProtectedRequiresBeatmaker(), handlers.PostNewLicense(service))
	license.Get("/licensesForBeat/:beatId", handlers.GetLicensesForBeat(service))

	license.Post("/newLicenseTemplate", internal.ProtectedRequiresBeatmaker(), handlers.PostNewLicenseTemplate(service))
	license.Get("/allLicenseTemplatesByBeatmakerJWT", internal.ProtectedRequiresBeatmaker(), handlers.GetAllLicenseTemplatesByBeatmakerId(service))
	license.Patch("/licenseTemplate", internal.ProtectedRequiresBeatmaker(), handlers.PatchLicenseTemplate(service))

	//admin
	license.Get("/allLicenses", handlers.GetAllLicense(service))	
	license.Get("/allLicenseTemplates", handlers.GetAllLicenseTemplate(service))	
}