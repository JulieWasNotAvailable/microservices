package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/JulieWasNotAvailable/microservices/unpublished/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/unpublished/internal/beatmetadata"
	"github.com/JulieWasNotAvailable/microservices/unpublished/internal/entities"
	"github.com/JulieWasNotAvailable/microservices/unpublished/internal/unpbeat"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Hello godoc
//
//	@Summary		Simple hello endpoint
//	@Description	Returns a hello message
//	@Tags			utils
//	@Produce		json
//	@Success		200	{string}	string	"hello!"
//	@Router			/ [get]
func Hello(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON("hello!")
}

// SaveBeatDraft godoc
//
//	@Summary		Save a beat draft
//	@Description	Save an unpublished beat with draft status
//	@Tags			beats
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	presenters.UnpublishedBeatResponseSuccessResponse
//	@Failure		422	{object}	presenters.UnpublishedBeatResponseErrorResponse
//	@Failure		500	{object}	presenters.UnpublishedBeatResponseErrorResponse
//	@Router			/unpbeats/makeEmptyBeat [post]
func MakeEmpty(service unpbeat.Service, metaservice beatmetadata.MetadataService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		beatmakeruuid, err := getIdFromJWT(c)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(presenters.CreateBeatErrorResponse(err))
		}

		createdBeat, err := service.CreateUnpublishedBeat(beatmakeruuid)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
		}

		var emptyAvailablefiles entities.AvailableFiles
		emptyAvailablefiles.UnpublishedBeatID = createdBeat.ID
		createdAvailbableFiles, err := metaservice.CreateAvailableFiles(&emptyAvailablefiles)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
		}

		createdBeat.AvailableFiles.ID = createdAvailbableFiles.ID
		return c.Status(http.StatusOK).JSON(presenters.CreateBeatSuccessResponse2(presenters.EntityToResponse(createdBeat)))
	}
}

// UpdateBeat godoc
//
//	@Summary		Update an unpublished beat
//	@Description	Update an existing unpublished beat entry
//	@Tags			beats
//	@Accept			json
//	@Produce		json
//	@Param			beat	body		entities.UnpublishedBeat	true	"Beat data to update"
//	@Success		200		{object}	object						"Successfully updated beat"
//	@Failure		422		{object}	object						"Unprocessable entity - invalid request body"
//	@Failure		500		{object}	object						"Internal server error"
//	@Router			/unpbeats/saveDraft [patch]
func UpdateBeat(service unpbeat.Service, mservice beatmetadata.MetadataService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody presenters.UnpublishedBeatRequest
		err := c.BodyParser(&requestBody)
		if err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(presenters.CreateBeatErrorResponse(err))
		}
		beatmakerid, err := getIdFromJWT(c)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(presenters.CreateBeatErrorResponse(err))
		}
		beatEntity := presenters.RequestToEntity(requestBody)
		beat, err := service.UpdateUnpublishedBeat(&beatEntity, beatmakerid)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateBeatSuccessResponse2(presenters.EntityToResponse(*beat)))
	}
}

// PostPublishBeat godoc
//
//	@Summary		Publish a beat. Deletes it from the current service, and posts to beat service.
//	@Description	Publish an existing beat (mock implementation)
//	@Tags			beats
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string	true	"Beat ID to publish"
//	@Success		200	{object}	map[string]string
//	@Failure		422	{object}	presenters.UnpublishedBeatResponseErrorResponse
//	@Failure		500	{object}	presenters.UnpublishedBeatResponseErrorResponse
//	@Router			/unpbeats/publishBeat [get] //UPDATE SWAGGER
func PostPublishBeat(service unpbeat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestBody := unpbeat.BeatPublishOptions{}
		err := c.BodyParser(&requestBody)
		if err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(presenters.CreateBeatErrorResponse(err))
		}
		beatmakerid, err := getIdFromJWT(c)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(presenters.CreateBeatErrorResponse(err))
		}
		beatmakerName, err := getBeatmakerNameFromJWT(c)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
		}

		beat, errArray := service.PublishBeat(beatmakerid, requestBody, beatmakerName)
		if len(errArray) != 0 {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorArrayResponse(errArray)) 
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateBeatSuccessResponse2(presenters.EntityToResponse(*beat)))
	}
}

// SendToModeration godoc
//
//	@Summary		Send beat to moderation
//	@Description	Update beat status to 'in_moderation' and set moderation timestamp
//	@Tags			admin
//	@Produce		json
//	@Param			id	path		string	true	"Beat ID to send to moderation"
//	@Success		200	{object}	presenters.UnpublishedBeatResponseSuccessResponse
//	@Failure		400	{object}	presenters.UnpublishedBeatResponseErrorResponse
//	@Failure		500	{object}	presenters.UnpublishedBeatResponseErrorResponse
//	@Router			/unpbeats/sendToModeration/{id} [get]
func SendToModeration(service unpbeat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		beatId := c.Params("id")
		if beatId == "" {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(errors.New("id cannot be empty")))
		}
		uuid, err := uuid.Parse(beatId)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
		}
		requestBody := presenters.UnpublishedBeatRequest{
			ID:                 uuid,
			Status:             string(entities.StatusInModeration),
			SentToModerationAt: time.Now().Unix(),
		}
		beatmakerId, err := getIdFromJWT(c)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(presenters.CreateBeatErrorResponse(err))
		}
		beatEntity := presenters.RequestToEntity(requestBody)
		beat, err := service.UpdateUnpublishedBeat(&beatEntity, beatmakerId)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateBeatSuccessResponse2(presenters.EntityToResponse(*beat)))
	}
}

// GetUnpublishedBeatsByBeatmakerId godoc
//
//	@Summary		Get user's unpublished beats
//	@Description	Get all unpublished beats for the authenticated beatmaker
//	@Tags			beats
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	presenters.UnpublishedBeatResponseListResponse
//	@Failure		500	{object}	presenters.UnpublishedBeatResponseErrorResponse
//	@Router			/unpbeats/unpublishedBeatsByBeatmakerJWT [get]
func GetUnpublishedBeatsByBeatmakerId(service unpbeat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uuid, err := getIdFromJWT(c)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(presenters.CreateBeatErrorResponse(err))
		}

		beats, err := service.GetUnpublishedBeatsByUser(uuid)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.EntityListToResponse(beats))
	}
}

// GetBeatsSortByStatusAndJWT godoc
//
//	@Summary		Get beats by status for current user
//	@Description	Get unpublished beats filtered by status for the authenticated beatmaker
//	@Tags			beats
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			status	path		string	true	"Status to filter by"
//	@Success		200		{object}	presenters.UnpublishedBeatResponseListResponse
//	@Failure		500		{object}	presenters.UnpublishedBeatResponseErrorResponse
//	@Router			/unpbeats/sortByStatus/{status} [get]
func GetBeatsSortByStatusAndJWT(service unpbeat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		status := c.Params("status")
		uuid, err := getIdFromJWT(c)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(presenters.CreateBeatErrorResponse(err))
		}

		beats, err := service.GetUnpublishedByBeatmakerandStatus(uuid, status)
		if err != nil {
			return c.Status(http.StatusBadGateway).JSON(presenters.CreateBeatErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.EntityListToResponse(*beats))
	}
}

// GetBeatsInModeration godoc
//
//	@Summary		Get beats in moderation by date range. Warning: uses unix data.
//	@Description	Get beats in moderation status within specified date range
//	@Tags			admin
//	@Produce		json
//	@Param			from	path		int	true	"Start timestamp"
//	@Param			to		path		int	true	"End timestamp"
//	@Success		200		{object}	presenters.UnpublishedBeatResponseListResponse
//	@Failure		400		{object}	presenters.UnpublishedBeatResponseErrorResponse
//	@Failure		500		{object}	presenters.UnpublishedBeatResponseErrorResponse
//	@Router			/unpbeats/beatsForModerationByDate/{from}/{to} [get]
func GetBeatsInModeration(service unpbeat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fromInt, err := c.ParamsInt("from")
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
		}
		toInt, err := c.ParamsInt("to")
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
		}

		beats, err := service.GetUnpublishedInModeration(int64(fromInt), int64(toInt))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
		}
		return c.Status(http.StatusOK).JSON(presenters.EntityListToResponse(*beats))
	}
}

// GetAllUnpublishedBeats godoc
//
//	@Summary		Get all unpublished beats
//	@Description	Get all unpublished beats in the system
//	@Tags			admin
//	@Produce		json
//	@Success		200	{object}	presenters.UnpublishedBeatResponseListResponse
//	@Failure		500	{object}	presenters.UnpublishedBeatResponseErrorResponse
//	@Router			/unpbeats/allUnpublishedBeats [get]
func GetAllUnpublishedBeats(service unpbeat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		beats, err := service.GetAllUnpublishedBeats()
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.EntityListToResponse(beats))
	}
}

// GetUnpublishedBeatById godoc
//
//	@Summary		Get an unpublished beat by ID
//	@Description	Retrieves an unpublished beat with the specified ID
//	@Tags			Beats
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Beat ID (UUID format)"
//	@Success		200	{object}	presenters.UnpublishedBeatResponseSuccessResponse
//	@Failure		422	{object}	presenters.UnpublishedBeatResponseErrorResponse
//	@Failure		500	{object}	presenters.UnpublishedBeatResponseErrorResponse
//	@Router			/unpbeats/unpublishedBeatById/{id} [get]
func GetUnpublishedBeatById(service unpbeat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		beatId := c.Params("id")
		beatuuid, err := uuid.Parse(beatId)
		if err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(presenters.CreateBeatErrorResponse(errors.New("wrong id format")))
		}

		beat, err := service.GetUnpublishedBeatByID(beatuuid)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateBeatSuccessResponse2(presenters.EntityToResponse(*beat)))
	}
}

// DeleteUnpublishedBeatById godoc
//
//	@Summary		Delete an unpublished beat by ID
//	@Description	Deletes an unpublished beat with the specified ID
//	@Tags			Beats
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Beat ID (UUID format)"
//	@Success		200	{object}	presenters.UnpublishedBeatResponseSuccessResponse
//	@Failure		422	{object}	presenters.UnpublishedBeatResponseErrorResponse
//	@Failure		500	{object}	presenters.UnpublishedBeatResponseErrorResponse
//	@Router			/unpbeats/deleteUnpublishedBeatById/{id} [delete]
func DeleteUnpublishedBeatById(service unpbeat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		beatId := c.Params("id")
		beatuuid, err := uuid.Parse(beatId)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
		}

		err = service.DeleteUnpublishedBeat(beatuuid)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
		}
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"message": "deleted successfully",
		})
	}
}

func getIdFromJWT(c *fiber.Ctx) (uuid.UUID, error) {
	auth := c.GetReqHeaders()
	authHeader, ok := auth["Authorization"]
	if !ok {
		return uuid.Nil, errors.New("auth header is absent")
	}
	splitToken := strings.Split(authHeader[0], "Bearer ")
	tokenStr := splitToken[1]

	nilUuid := uuid.Nil
	token, _, err := jwt.NewParser().ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return nilUuid, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nilUuid, err
	}

	id := claims["id"].(string)
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nilUuid, err
	}

	return uuid, nil
}

func getBeatmakerNameFromJWT(c *fiber.Ctx) (string, error) {
	auth := c.GetReqHeaders()
	authHeader, ok := auth["Authorization"]
	if !ok {
		return "", errors.New("auth header is absent")
	}
	splitToken := strings.Split(authHeader[0], "Bearer ")
	tokenStr := splitToken[1]

	token, _, err := jwt.NewParser().ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", err
	}

	name := claims["username"].(string)

	return name, nil
}
