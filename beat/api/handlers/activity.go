package handlers

import (
	"net/http"

	"github.com/JulieWasNotAvailable/microservices/beat/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/beat/internal/activity"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type requestBody struct {
	BeatId string `json:"beatId"`
}

type requestBodyList struct {
	Beatids []string `json:"beatids"`
}

// PostLike godoc
//	@Summary		Like a beat
//	@Description	Add a like to a beat by the authenticated user
//	@Tags			likes
//	@Accept			json
//	@Produce		json
//	@Param			beatId	body	requestBody	true	"Beat id data"
//	@Security		ApiKeyAuth
//	@Success		201	{object}	map[string]interface{}	"Like created successfully"
//	@Failure		401	{object}	map[string]interface{}	"Unauthorized"
//	@Failure		500	{object}	map[string]interface{}	"Internal server error"
//	@Router			/activity/postNewLike [post]
func PostLike(service activity.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := getIdFromJWT(c)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(presenters.CreateMetadataErrorResponse(err))
		}
		requestBody := requestBody{}
		err = c.BodyParser(&requestBody)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateErrorResponse(err))
		}
		beatuuid, err := uuid.Parse(requestBody.BeatId)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateErrorResponse(err))
		}

		like, err := service.InsertLike(userId, beatuuid)
		if err != nil {
			return c.Status(http.StatusResetContent).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusCreated).JSON(presenters.CreateMetadataListResponse(like))
	}
}

// DeleteLike godoc
//	@Summary		Remove a like
//	@Description	Remove a like from a beat by the authenticated user
//	@Tags			likes
//	@Accept			json
//	@Produce		json
//	@Param			beatId	path	string	true	"Beat ID"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	map[string]interface{}	"Like removed successfully"
//	@Failure		401	{object}	map[string]interface{}	"Unauthorized"
//	@Failure		500	{object}	map[string]interface{}	"Internal server error"
//	@Router			/activity/{beatId} [delete]
func DeleteLike(service activity.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := getIdFromJWT(c)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(presenters.CreateMetadataErrorResponse(err))
		}
		beatIdStr := c.Params("beatId")
		beatId, err := uuid.Parse(beatIdStr)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		like, err := service.DeleteLike(userId, beatId)
		if err != nil {
			return c.Status(http.StatusResetContent).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataSuccessResponse(like))
	}
}

// GetLikesByUserId godoc
//	@Summary		Get user's likes
//	@Description	Get all likes by the authenticated user
//	@Tags			likes
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	map[string]interface{}	"List of likes"
//	@Failure		401	{object}	map[string]interface{}	"Unauthorized"
//	@Failure		500	{object}	map[string]interface{}	"Internal server error"
//	@Router			/activity/viewMyLikes [get]
func GetLikesByUserId(service activity.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := getIdFromJWT(c)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		likes, err := service.GetLikesByUserId(userId)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(likes))
	}
}

// GetLikesCountByBeatId godoc
//	@Summary		Get like count for a beat
//	@Description	Get the number of likes for a specific beat
//	@Tags			likes
//	@Accept			json
//	@Produce		json
//	@Param			beatId	path		string					true	"Beat ID"
//	@Success		200		{object}	map[string]interface{}	"Like count"
//	@Failure		400		{object}	map[string]interface{}	"Bad request"
//	@Failure		500		{object}	map[string]interface{}	"Internal server error"
//	@Router			/activity/viewLikesCountByBeatId/{beatId} [get]
func GetLikesCountByBeatId(service activity.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		beatID, err := uuid.Parse(c.Params("beatId"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		count, err := service.GetLikesCountByBeatId(beatID)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(
			presenters.CreateMetadataListResponse(fiber.Map{
				"beat_id": beatID,
				"count":   count,
			}))
	}
}

// GetLikesCountByUserId godoc
//	@Summary		Get like count for a user
//	@Description	Get the number of likes given by a specific user
//	@Tags			likes
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		string					true	"User ID"
//	@Success		200		{object}	map[string]interface{}	"Like count"
//	@Failure		400		{object}	map[string]interface{}	"Bad request"
//	@Failure		500		{object}	map[string]interface{}	"Internal server error"
//	@Router			/activity/viewLikesCountByUserId/{userId} [get]
func GetLikesCountByUserId(service activity.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := uuid.Parse(c.Params("userId"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		count, err := service.GetUserLikesCount(userId)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse)
		}

		return c.Status(http.StatusOK).JSON(
			presenters.CreateMetadataListResponse(fiber.Map{
				"beat_id": userId,
				"count":   count,
			}))
	}
}

// GetTotalLikesOfBeats godoc
//	@Summary		Get total likes for multiple beats
//	@Description	Get the total number of likes for a list of beats
//	@Tags			likes
//	@Accept			json
//	@Produce		json
//	@Param			requestBodyList	body		requestBodyList			true	"List of beat IDs"
//	@Success		200				{object}	map[string]interface{}	"Total likes count"
//	@Failure		400				{object}	map[string]interface{}	"Bad request"
//	@Failure		500				{object}	map[string]interface{}	"Internal server error"
//	@Router			/activity/totalLikesCountForBeats [post]
func GetTotalLikesOfBeats(service activity.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req requestBodyList
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
		}

		beatuuids := []uuid.UUID{}
		for _, v := range req.Beatids {
			uuid, err := uuid.Parse(v)
			if err != nil {
				return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
			}
			beatuuids = append(beatuuids, uuid)
		}

		count, err := service.GetLikesCountOfBeats(beatuuids)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse)
		}

		return c.Status(http.StatusOK).JSON(
			presenters.CreateMetadataListResponse(fiber.Map{
				"beats": req.Beatids,
				"count": count,
			}))
	}
}

// PostListened godoc
//	@Summary		Record a listen
//	@Description	Record that a user listened to a beat
//	@Tags			listen
//	@Accept			json
//	@Produce		json
//	@Param			request	body	requestBody	true	"Listen data"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	map[string]interface{}	"Listen recorded"
//	@Failure		400	{object}	map[string]interface{}	"Bad request"
//	@Failure		401	{object}	map[string]interface{}	"Unauthorized"
//	@Failure		500	{object}	map[string]interface{}	"Internal server error"
//	@Router			/activity/listened [post]
func PostListened(service activity.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := getIdFromJWT(c)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(presenters.CreateErrorResponse(err))
		}

		requestBody := requestBody{}
		err = c.BodyParser(&requestBody)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateErrorResponse(err))
		}
		beatuuid, err := uuid.Parse(requestBody.BeatId)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateErrorResponse(err))
		}

		listened, err := service.InsertListened(userId, beatuuid)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateSuccessResponse(listened))
	}
}
