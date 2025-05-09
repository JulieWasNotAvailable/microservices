package handlers

import (
	"errors"

	"github.com/JulieWasNotAvailable/microservices/cart/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/cart/pkg/entities"
	"github.com/JulieWasNotAvailable/microservices/cart/pkg/license"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func PostNewLicenseTemplate(service license.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := getIdFromJWT(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(presenters.CreateErrorResponse(err))
		}

		var template entities.LicenseTemplate
		if err := c.BodyParser(&template); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(presenters.CreateErrorResponse(err))
		}

		createdTemplate, err := service.InsertNewLicenseTemplate(userId, template)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(presenters.CreateErrorResponse(err))
		}

		return c.Status(fiber.StatusCreated).JSON(createdTemplate)
	}
}

func GetAllLicenseTemplatesByBeatmakerId(service license.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		beatmakerId, err := getIdFromJWT(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(presenters.CreateErrorResponse(err))
		}

		templates, err := service.GetAllLicenseTemplateByBeatmakerId(beatmakerId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(presenters.CreateErrorResponse(err))
		}

		return c.Status(fiber.StatusOK).JSON(presenters.CreateSuccessResponse(templates))
	}
}

func PatchLicenseTemplate(service license.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var template presenters.LicenseTemplate
		if err := c.BodyParser(&template); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(presenters.CreateErrorResponse(err))
		}

		beatmakerId, err := getIdFromJWT(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(presenters.CreateErrorResponse(err))
		}
		
		licenseTemplate, err := service.ReadLicenseTemplateById(template.ID)
		if licenseTemplate.UserID != beatmakerId {
			return c.Status(fiber.StatusUnauthorized).JSON(presenters.CreateErrorResponse(err))
		}

		err = service.UpdateLicenseTemplate(template)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(presenters.CreateErrorResponse(err))
			}
			return c.Status(fiber.StatusInternalServerError).JSON(presenters.CreateErrorResponse(err))
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "license template updated successfully",
		})
	}
}

func GetAllLicenseTemplate(service license.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		licenseTemplates, err := service.ReadAllLicenseTemplate()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(presenters.CreateErrorResponse(err))
		}

		return  c.Status(fiber.StatusCreated).JSON(presenters.CreateListResponse(licenseTemplates))
	}
}