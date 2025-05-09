package handlers

import (
	"errors"

	"github.com/JulieWasNotAvailable/microservices/cart/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/cart/pkg/entities"
	"github.com/JulieWasNotAvailable/microservices/cart/pkg/license"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetLicensesForBeat(service license.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		beatId, err := uuid.Parse(c.Params("beatId"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(presenters.CreateErrorResponse(err))
		}

		licenses, err := service.GetLicenseByBeatId(beatId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(presenters.CreateErrorResponse(err))
		}

		return c.Status(fiber.StatusOK).JSON(presenters.CreateSuccessResponse(licenses))
	}
}

func PostNewLicense(service license.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestBody := entities.License{}
		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(presenters.CreateErrorResponse(err))
		}

		beatmakerId, err := getIdFromJWT(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(presenters.CreateErrorResponse(err))
		} 

		licenseTemplate, err := service.ReadLicenseTemplateById(requestBody.LicenseTemplateID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(presenters.CreateErrorResponse(err))
		}
	
		if licenseTemplate.UserID != beatmakerId {
			return c.Status(fiber.StatusUnauthorized).JSON(presenters.CreateErrorResponse(errors.New("this license template does not belong to the beatmaker")))
		}

		requestBody.UserID = beatmakerId
		newLicense, err := service.InsertNewLicense(requestBody)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(presenters.CreateErrorResponse(err))
		}

		return c.Status(fiber.StatusOK).JSON(newLicense)
	}
}

func GetAllLicense(service license.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		licenses, err := service.ReadAllLicense()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(presenters.CreateErrorResponse(err))
		}

		return  c.Status(fiber.StatusOK).JSON(presenters.CreateListResponse(licenses))
	}
}