package handlers

import (
	"errors"
	"net/http"

	"github.com/JulieWasNotAvailable/microservices/user/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/user/internal/activity"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// PostSubscribeToBeatmaker godoc
// @Summary Subscribe to a beatmaker
// @Description Subscribe the authenticated user to a beatmaker
// @Tags activity
// @Accept json
// @Produce json
// @Param beatmakerId path string true "Beatmaker ID to subscribe to"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "Subscription details"
// @Failure 400 {object} map[string]interface{} "Invalid beatmaker ID or self-subscription attempt"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /activity/subscribeTo/{beatmakerId} [post]
func PostSubscribeToBeatmaker(service activity.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := getIdFromJWT(c)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(presenters.CreateErrorResponse(err))
		}

		beatmakerId, err := uuid.Parse(c.Params("beatmakerId"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateErrorResponse(err))
		}

		if beatmakerId == userId {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateErrorResponse(errors.New("user cannot subsribe to himself")))
		}

		subscription, err := service.InsertSub(userId, beatmakerId)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateSuccessResponse(subscription))
	}
}

// GetMySubscriptions godoc
// @Summary Get user's subscriptions
// @Description Get all beatmakers the authenticated user is subscribed to
// @Tags activity
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "List of subscriptions"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /activity/viewMySubscriptions [get]
func GetMySubscriptions(service activity.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := getIdFromJWT(c)
		if err != nil {
			return c.JSON(presenters.CreateErrorResponse(err))
		}

		subscriptions, err := service.FetchSubsByUserId(userId)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateErrorResponse(err))
		}

		return c.JSON(presenters.CreateListResponse(subscriptions))
	}
}

// GetFollowersCountByBeatmakerId godoc
// @Summary Get beatmaker's followers count
// @Description Get the number of followers for a specific beatmaker
// @Tags activity
// @Accept json
// @Produce json
// @Param beatmakerId path string true "Beatmaker ID"
// @Success 200 {object} map[string]interface{} "Followers count"
// @Failure 400 {object} map[string]interface{} "Invalid beatmaker ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /activity/followersNumberByBeatmakerId/{beatmakerId} [get]
func GetFollowersCountByBeatmakerId(service activity.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		beatmakerID, err := uuid.Parse(c.Params("beatmakerId"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateErrorResponse(fiber.ErrBadRequest))
		}

		count, err := service.FetchSubsCountByBeatmakerId(beatmakerID)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateErrorResponse(err))
		}

		return c.JSON(presenters.CreateSuccessResponse(count))
	}
}

// DeleteSub godoc
// @Summary Unsubscribe from a beatmaker
// @Description Unsubscribe the authenticated user from a beatmaker
// @Tags activity
// @Accept json
// @Produce json
// @Param beatmakerId path string true "Beatmaker ID to unsubscribe from"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "Unsubscription confirmation"
// @Failure 400 {object} map[string]interface{} "Invalid beatmaker ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /activity/unsubscribe/{beatmakerId} [delete]
func DeleteSub(service activity.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := getIdFromJWT(c)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(presenters.CreateErrorResponse(err))
		}

		beatmakerId, err := uuid.Parse(c.Params("beatmakerId"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateErrorResponse(fiber.ErrBadRequest))
		}

		subscription, err := service.RemoveSub(userId, beatmakerId)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateSuccessResponse(subscription))
	}
}
