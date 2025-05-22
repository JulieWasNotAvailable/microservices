package handlers

import (
	"errors"
	"net/http"

	"github.com/JulieWasNotAvailable/microservices/user/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/user/internal/activity"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

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
