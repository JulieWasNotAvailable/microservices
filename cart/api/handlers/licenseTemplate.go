package handlers

import (
	"errors"
	"log"

	"github.com/JulieWasNotAvailable/microservices/cart/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/cart/internal/entities"
	"github.com/JulieWasNotAvailable/microservices/cart/internal/license"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// PostNewLicenseTemplate creates a new license template
//	@Summary		Create license template
//	@Description	Creates a new license template (beatmaker only)
//	@Tags			LicenseTemplate
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			template	body		entities.LicenseTemplate	true	"License template data"
//	@Success		201			{object}	entities.LicenseTemplate	"Created license template"
//	@Failure		400			{object}	presenters.ErrorResponse	"Invalid request body"
//	@Failure		401			{object}	presenters.ErrorResponse	"Unauthorized"
//	@Failure		500			{object}	presenters.ErrorResponse	"Internal server error"
//	@Router			/license/newLicenseTemplate [post]
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

// GetAllLicenseTemplatesByBeatmakerId retrieves all templates for authenticated beatmaker
//	@Summary		Get beatmaker's license templates
//	@Description	Returns all license templates for the authenticated beatmaker
//	@Tags			LicenseTemplate
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	presenters.SuccessResponse	"List of beatmaker's license templates"
//	@Failure		401	{object}	presenters.ErrorResponse	"Unauthorized"
//	@Failure		500	{object}	presenters.ErrorResponse	"Internal server error"
//	@Router			/license/allLicenseTemplatesByBeatmakerJWT [get]
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

// PatchLicenseTemplate updates a license template
//	@Summary		Update license template
//	@Description	Updates an existing license template (beatmaker only)
//	@Tags			LicenseTemplate
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			template	body		presenters.LicenseTemplate	true	"License template update data"
//	@Success		200			{object}	map[string]string			"Success message"
//	@Failure		400			{object}	presenters.ErrorResponse	"Invalid request body"
//	@Failure		401			{object}	presenters.ErrorResponse	"Unauthorized or template ownership mismatch"
//	@Failure		404			{object}	presenters.ErrorResponse	"Template not found"
//	@Failure		500			{object}	presenters.ErrorResponse	"Internal server error"
//	@Router			/license/licenseTemplate [patch]
func PatchLicenseTemplate(service license.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var template presenters.LicenseTemplate
		if err := c.BodyParser(&template); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(presenters.CreateErrorResponse(err))
		}
		log.Println("patching template: ", template)
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

// GetAllLicenseTemplate retrieves all license templates (admin only)
//	@Summary		Get all license templates
//	@Description	Returns all license templates in the system (admin only)
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Success		201	{object}	presenters.ListResponse		"List of all license templates"
//	@Failure		500	{object}	presenters.ErrorResponse	"Internal server error"
//	@Router			/license/allLicenseTemplates [get]
func GetAllLicenseTemplate(service license.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		licenseTemplates, err := service.ReadAllLicenseTemplate()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(presenters.CreateErrorResponse(err))
		}

		return  c.Status(fiber.StatusCreated).JSON(presenters.CreateListResponse(licenseTemplates))
	}
}