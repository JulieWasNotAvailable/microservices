package handlers

import (
	"errors"
	"log"
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

type BeatForPublishing struct {
	Beat presenters.UnpublishedBeat
	MFCC []float64
}

// Hello godoc
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
//	@Summary		Save a beat draft
//	@Description	Save an unpublished beat with draft status
//	@Tags			beats
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	presenters.UnpublishedBeatSuccessResponse
//	@Failure		422	{object}	presenters.UnpublishedBeatErrorResponse
//	@Failure		500	{object}	presenters.UnpublishedBeatErrorResponse
//	@Router			/unpbeats/makeEmptyBeat [post]
func SaveBeatDraft(service unpbeat.Service, metaservice beatmetadata.MetadataService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var emptyBeat entities.UnpublishedBeat
		beatmakeruuid, err := getIdFromJWT(c)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		emptyBeat.Status = entities.StatusDraft
		emptyBeat.BeatmakerID = beatmakeruuid

		createdBeat, err := service.CreateUnpublishedBeat(emptyBeat)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		var emptyAvailablefiles entities.AvailableFiles
		emptyAvailablefiles.UnpublishedBeatID = createdBeat.ID
		createdAvailbableFiles, err := metaservice.CreateAvailableFiles(&emptyAvailablefiles)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		beatReponse := presenters.UnpublishedBeat{
			ID:             createdBeat.ID,
			Status:         string(createdBeat.Status),
			BeatmakerID:    beatmakeruuid,
			AvailableFiles: createdAvailbableFiles,
		}
		return c.Status(http.StatusOK).JSON(presenters.CreateBeatSuccessResponse2(beatReponse))
	}
}

// UpdateBeat godoc
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
		var requestBody entities.UnpublishedBeat
		err := c.BodyParser(&requestBody)

		if err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(presenters.CreateBeatErrorResponse(err))
		}

		beat, err := service.UpdateUnpublishedBeat(&requestBody)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateBeatSuccessResponse2(*beat))
	}
}

// PostPublishBeat godoc
//	@Summary		Publish a beat. Deletes it from the current service, and posts to beat service.
//	@Description	Publish an existing beat (mock implementation)
//	@Tags			beats
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string	true	"Beat ID to publish"
//	@Success		200	{object}	map[string]string
//	@Failure		422	{object}	presenters.UnpublishedBeatErrorResponse
//	@Failure		500	{object}	presenters.UnpublishedBeatErrorResponse
//	@Router			/unpbeats/publishBeat/{id} [get] //UPDATE SWAGGER
func PostPublishBeat(service unpbeat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestBody := unpbeat.BeatPublishOptions{}
		err := c.BodyParser(&requestBody)
		if err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(presenters.CreateBeatErrorResponse(err))
		}
		log.Println(requestBody)
		beatmakerid, err := getIdFromJWT(c)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}
		beatmakerName, err := getBeatmakerNameFromJWT(c)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		beat, err := service.PublishBeat(beatmakerid, requestBody, beatmakerName)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateBeatSuccessResponse2(*beat))
	}
}

// SendToModeration godoc
//	@Summary		Send beat to moderation
//	@Description	Update beat status to 'in_moderation' and set moderation timestamp
//	@Tags			admin
//	@Produce		json
//	@Param			id	path		string	true	"Beat ID to send to moderation"
//	@Success		200	{object}	presenters.UnpublishedBeatSuccessResponse
//	@Failure		400	{object}	presenters.UnpublishedBeatErrorResponse
//	@Failure		500	{object}	presenters.UnpublishedBeatErrorResponse
//	@Router			/unpbeats/sendToModeration/{id} [get]
func SendToModeration(service unpbeat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		beatId := c.Params("id")
		if beatId == "" {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(errors.New("id cannot be empty")))
		}
		uuid, err := uuid.Parse(beatId)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}
		requestBody := entities.UnpublishedBeat{
			ID:                 uuid,
			Status:             entities.StatusInModeration,
			SentToModerationAt: time.Now().Unix(),
		}
		beat, err := service.UpdateUnpublishedBeat(&requestBody)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateBeatSuccessResponse2(*beat))
	}
}

// GetUnpublishedBeatsByBeatmakerId godoc
//	@Summary		Get user's unpublished beats
//	@Description	Get all unpublished beats for the authenticated beatmaker
//	@Tags			beats
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	presenters.UnpublishedBeatListResponse
//	@Failure		500	{object}	presenters.UnpublishedBeatErrorResponse
//	@Router			/unpbeats/unpublishedBeatsByBeatmakerJWT [get]
func GetUnpublishedBeatsByBeatmakerId(service unpbeat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uuid, err := getIdFromJWT(c)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		beats, err := service.GetUnpublishedBeatsByUser(uuid)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateBeatListSuccessResponse(beats))
	}
}

// GetBeatsSortByStatusAndJWT godoc
//	@Summary		Get beats by status for current user
//	@Description	Get unpublished beats filtered by status for the authenticated beatmaker
//	@Tags			beats
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			status	path		string	true	"Status to filter by"
//	@Success		200		{object}	presenters.UnpublishedBeatListResponse
//	@Failure		500		{object}	presenters.UnpublishedBeatErrorResponse
//	@Router			/unpbeats/sortByStatus/{status} [get]
func GetBeatsSortByStatusAndJWT(service unpbeat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		status := c.Params("status")
		log.Println(status)
		uuid, err := getIdFromJWT(c)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		beats, err := service.GetUnpublishedByBeatmakerandStatus(uuid, status)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateBeatListSuccessResponse(*beats))
	}
}

// GetBeatsInModeration godoc
//	@Summary		Get beats in moderation by date range. Warning: uses unix data.
//	@Description	Get beats in moderation status within specified date range
//	@Tags			admin
//	@Produce		json
//	@Param			from	path		int	true	"Start timestamp"
//	@Param			to		path		int	true	"End timestamp"
//	@Success		200		{object}	presenters.UnpublishedBeatListResponse
//	@Failure		400		{object}	presenters.UnpublishedBeatErrorResponse
//	@Failure		500		{object}	presenters.UnpublishedBeatErrorResponse
//	@Router			/unpbeats/beatsForModerationByDate/{from}/{to} [get]
func GetBeatsInModeration(service unpbeat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fromInt, err := c.ParamsInt("from")
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}
		toInt, err := c.ParamsInt("to")
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		beats, err := service.GetUnpublishedInModeration(int64(fromInt), int64(toInt))
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateBeatListSuccessResponse(*beats))
	}
}

// GetAllUnpublishedBeats godoc
//	@Summary		Get all unpublished beats
//	@Description	Get all unpublished beats in the system
//	@Tags			admin
//	@Produce		json
//	@Success		200	{object}	presenters.UnpublishedBeatListResponse
//	@Failure		500	{object}	presenters.UnpublishedBeatErrorResponse
//	@Router			/unpbeats/allUnpublishedBeats [get]
func GetAllUnpublishedBeats(service unpbeat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		beats, err := service.GetAllUnpublishedBeats()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateBeatListSuccessResponse(beats))
	}
}

// GetUnpublishedBeatById godoc
//	@Summary		Get an unpublished beat by ID
//	@Description	Retrieves an unpublished beat with the specified ID
//	@Tags			Beats
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Beat ID (UUID format)"
//	@Success		200	{object}	presenters.UnpublishedBeatSuccessResponse
//	@Failure		422	{object}	presenters.UnpublishedBeatErrorResponse
//	@Failure		500	{object}	presenters.UnpublishedBeatErrorResponse
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
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateBeatSuccessResponse2(*beat))
	}
}

// DeleteUnpublishedBeatById godoc
//	@Summary		Delete an unpublished beat by ID
//	@Description	Deletes an unpublished beat with the specified ID
//	@Tags			Beats
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Beat ID (UUID format)"
//	@Success		200	{object}	presenters.UnpublishedBeatSuccessResponse
//	@Failure		422	{object}	presenters.UnpublishedBeatErrorResponse
//	@Failure		500	{object}	presenters.UnpublishedBeatErrorResponse
//	@Router			/unpbeats/deleteUnpublishedBeatById/{id} [delete]
func DeleteUnpublishedBeatById(service unpbeat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		beatId := c.Params("id")
		beatuuid, err := uuid.Parse(beatId)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		err = service.DeleteUnpublishedBeat(beatuuid)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
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
