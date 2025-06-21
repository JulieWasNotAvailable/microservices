package handlers

import (
	"errors"

	"github.com/JulieWasNotAvailable/microservices/cart/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/cart/internal/entities"
	"github.com/JulieWasNotAvailable/microservices/cart/internal/license"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var ErrUserNotOwner = errors.New("user does not own the template")

type LicenseList struct {
	BeatId   uuid.UUID          `json:"beatId"`
	Licenses []entities.License `json:"licenses"`
}

// GetLicensesForBeat retrieves all licenses for a specific beat
//
//	@Summary		Get licenses by beat ID
//	@Description	Returns all available licenses for a specific beat
//	@Tags			License
//	@Accept			json
//	@Produce		json
//	@Param			beatId	path		string						true	"Beat ID in UUID format"
//	@Success		200		{object}	presenters.SuccessResponse	"List of licenses for the beat"
//	@Failure		400		{object}	presenters.ErrorResponse	"Invalid beat ID format"
//	@Failure		500		{object}	presenters.ErrorResponse	"Internal server error"
//	@Router			/license/licensesForBeat/{beatId} [get]
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

// PostNewLicense creates a new license for a beat
//
//	@Summary		Create new license
//	@Description	Creates a new license for a beat (beatmaker only)
//	@Tags			License
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			license	body		entities.License			true	"License creation data"
//	@Success		200		{object}	entities.License			"Created license details"
//	@Failure		400		{object}	presenters.ErrorResponse	"Invalid request body"
//	@Failure		401		{object}	presenters.ErrorResponse	"Unauthorized or template ownership mismatch"
//	@Failure		500		{object}	presenters.ErrorResponse	"Internal server error"
//	@Router			/license/newLicense [post]
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

		return c.Status(fiber.StatusOK).JSON(presenters.CreateSuccessResponse(newLicense))
	}
}

// PostNewLicense creates several new licenses for a beat
//
//	@Summary		Create licenses
//	@Description	Creates a new license for a beat (beatmaker only)
//	@Tags			License
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			license	body		entities.License			true	"License creation data"
//	@Success		200		{object}	entities.License			"Created license details"
//	@Failure		400		{object}	presenters.ErrorResponse	"Invalid request body"
//	@Failure		401		{object}	presenters.ErrorResponse	"Unauthorized or template ownership mismatch"
//	@Failure		500		{object}	presenters.ErrorResponse	"Internal server error"
//	@Router			/license/newLisenseList [post]
func PostNewLicenseList(service license.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody LicenseList
		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(presenters.CreateErrorResponse(err))
		}

		beatmakerId, err := getIdFromJWT(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(presenters.CreateErrorResponse(err))
		}

		newLicenses, err := service.InsertNewLicenseList(requestBody.BeatId, beatmakerId, requestBody.Licenses)
		if err != nil {
			if errors.Is(err, ErrUserNotOwner) {
				return c.Status(fiber.StatusUnauthorized).JSON(presenters.CreateErrorResponse(err))
			}
			return err
		}

		return c.Status(fiber.StatusOK).JSON(presenters.CreateSuccessResponse(newLicenses))
	}
}

// GetAllLicense retrieves all licenses (admin only)
//
//	@Summary		Get all licenses
//	@Description	Returns all licenses in the system (admin only)
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	presenters.ListResponse		"List of all licenses"
//	@Failure		500	{object}	presenters.ErrorResponse	"Internal server error"
//	@Router			/license/allLicenses [get]
func GetAllLicense(service license.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		licenses, err := service.ReadAllLicense()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(presenters.CreateErrorResponse(err))
		}

		return c.Status(fiber.StatusOK).JSON(presenters.CreateListResponse(licenses))
	}
}
