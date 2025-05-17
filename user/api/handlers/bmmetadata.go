package handlers

import (
	"errors"
	"net/http"

	"github.com/JulieWasNotAvailable/microservices/user/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/user/internal/bmmetadata"
	"github.com/JulieWasNotAvailable/microservices/user/internal/entities"
	"github.com/JulieWasNotAvailable/microservices/user/internal/user"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

// AddMetadata godoc
// @Summary Add new metadata
// @Description Add new metadata entry to the system. In request "id" should be eliminated.
// @Tags metadata
// @Accept json
// @Produce json
// @Param metadata body entities.Metadata true "Metadata to add"
// @Success 200 {object} presenters.MetadataSuccessResponse
// @Failure 400 {object} presenters.MetadataErrorResponse
// @Failure 422 {object} presenters.MetadataErrorResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata [post]
func AddMetadata(bmservice bmmetadata.Service, uservice user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestBody := entities.Metadata{}

		err := c.BodyParser(&requestBody)
		if err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(presenters.CreateMetadataErrorResponse(err))
		}
		
		result, err := uservice.FetchUserById(requestBody.UserID)
		if err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		if result != nil && result.RoleID!= 2 {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "this user is not a beatmaker. metadata is only available for beatmakers"})
		}

		res, err := bmservice.InsertMetadata(&requestBody)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.JSON(presenters.CreateMetadataSuccessResponse(res))
	}
}

// @Tags metadata
// @Produce json
// @Success 200 {object} presenters.MetadataSuccessResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadatas [get]
func GetMetadatas(service bmmetadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		result, err := service.FetchMetadatas()
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.JSON(presenters.CreateMetadataListResponse(result))
	}
}

// GetMetadataById godoc
// @Summary Get metadata by ID
// @Description Get metadata entry by its ID
// @Tags metadata
// @Accept json
// @Produce json
// @Param id path string true "Metadata ID"
// @Success 200 {object} presenters.MetadataSuccessResponse
// @Failure 400 {object} presenters.MetadataErrorResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadataById/{id} [get]
func GetMetadataById(service bmmetadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.CreateMetadataErrorResponse(errors.New("need id")))
		}

		uuid, err := uuid.Parse(id)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.CreateMetadataErrorResponse(err))
		}

		result, err := service.FetchMetadataById(uuid)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.JSON(presenters.CreateMetadataSuccessResponse2(result))
	}
}

func UpdateMetadataById(service bmmetadata.Service) fiber.Handler{
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuid, err := uuid.Parse(id)
		if err != nil{
			return err
		}

		requestBody := presenters.Metadata{}
		err = c.BodyParser(&requestBody)
		if err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(presenters.CreateMetadataErrorResponse(err))
		}
		
		result, err := service.UpdateMetadataById(uuid, &requestBody)
		if err != nil{
			return err
		}

		return c.JSON(result)
	}
}
// RemoveMetadata godoc
// @Summary Delete metadata
// @Description Delete metadata entry by its ID. You need to be loged in. You can delete anybody's metadata.
// @Tags metadata
// @Accept json
// @Produce json
// @Param id path string true "Metadata ID"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "Successful deletion response"
// @Failure 400 {object} presenters.MetadataErrorResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadataById/{id} [delete]
func RemoveMetadata(service bmmetadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.CreateMetadataErrorResponse(errors.New("need id")))
		}

		uuid, err := uuid.Parse(id)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.CreateMetadataErrorResponse(err))
		}

		err = service.RemoveMetadataById(uuid)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.JSON(&fiber.Map{
			"status": true,
			"data":   "deleted successfully",
			"err":    nil,
		})
	}
}