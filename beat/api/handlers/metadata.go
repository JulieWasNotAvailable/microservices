package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/JulieWasNotAvailable/microservices/beat/pkg/metadata"
	"github.com/JulieWasNotAvailable/microservices/beat/api/presenters"
)

func GetAllGenres(service metadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		genres, err := service.ReadAllGenres()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(genres))
	}
}

func GetAllMoods(service metadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		moods, err := service.ReadAllMoods()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(moods))
	}
}

func GetAllKeys(service metadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		keys, err := service.ReadAllKeys()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(keys))
	}
}

func GetAllInstruments(service metadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		instruments, err := service.ReadAllInstruments()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(instruments))
	}
}

func GetAllTags(service metadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tags, err := service.ReadAllTags()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(tags))
	}
}

func GetAllTimestamps(service metadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		timestamps, err := service.ReadAllTimestamps()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(timestamps))
	}
}

func GetAllMFCCs(service metadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		mfccs, err := service.ReadAllMFCC()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(mfccs))
	}
}

func GetRandomTags(service metadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tags, err := service.ReadRandomTags()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(tags))
	}
}

func GetTagByName(service metadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		name := c.Params("name")
		tag, err := service.ReadTagsByName(name)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataSuccessResponse(tag))
	}
}

func GetTagsInTrend(service metadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tags, err := service.ReadPopularTags()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(tags))
	}
}

func GetGenresInTrend(service metadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tags, err := service.ReadPopularGenres()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(tags))
	}
}