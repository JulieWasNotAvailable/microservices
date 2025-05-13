package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/JulieWasNotAvailable/microservices/beat/internal/metadata"
	"github.com/JulieWasNotAvailable/microservices/beat/api/presenters"
)

// GetAllGenres retrieves all music genres
// @Summary Get all genres
// @Description Returns a list of all available music genres
// @Tags Metadata
// @Produce json
// @Success 200 {object} presenters.MetadataListResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata/genres [get]
func GetAllGenres(service metadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		genres, err := service.ReadAllGenres()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(genres))
	}
}

// GetAllMoods retrieves all mood categories
// @Summary Get all moods
// @Description Returns a list of all available mood categories
// @Tags Metadata
// @Produce json
// @Success 200 {object} presenters.MetadataListResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata/moods [get]
func GetAllMoods(service metadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		moods, err := service.ReadAllMoods()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(moods))
	}
}

// GetAllKeys retrieves all musical keys
// @Summary Get all keys
// @Description Returns a list of all available musical keys
// @Tags Metadata
// @Produce json
// @Success 200 {object} presenters.MetadataListResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata/keys [get]
func GetAllKeys(service metadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		keys, err := service.ReadAllKeys()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(keys))
	}
}

// GetAllInstruments retrieves all instruments
// @Summary Get all instruments
// @Description Returns a list of all available instruments
// @Tags Metadata
// @Produce json
// @Success 200 {object} presenters.MetadataListResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata/instruments [get]
func GetAllInstruments(service metadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		instruments, err := service.ReadAllInstruments()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(instruments))
	}
}

// GetAllTags retrieves all tags
// @Summary Get all tags
// @Description Returns a list of all available tags
// @Tags Tags
// @Produce json
// @Success 200 {object} presenters.MetadataListResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata/tags [get]
func GetAllTags(service metadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tags, err := service.ReadAllTags()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(tags))
	}
}

// GetAllTimestamps retrieves all timestamps
// @Summary Get all timestamps
// @Description Returns a list of all available timestamps
// @Tags Timestamp
// @Produce json
// @Success 200 {object} presenters.MetadataListResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata/timestamps [get]
func GetAllTimestamps(service metadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		timestamps, err := service.ReadAllTimestamps()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(timestamps))
	}
}

// GetAllMFCCs retrieves all MFCC data
// @Summary Get all MFCCs
// @Description Returns a list of all available MFCC data
// @Tags Metadata
// @Produce json
// @Success 200 {object} presenters.MetadataListResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata/mfccs [get]
func GetAllMFCCs(service metadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		mfccs, err := service.ReadAllMFCC()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(mfccs))
	}
}

// GetRandomTags retrieves random tags
// @Summary Get random tags
// @Description Returns a list of randomly selected tags
// @Tags Tags
// @Produce json
// @Success 200 {object} presenters.MetadataListResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata/tags/random [get]
func GetRandomTags(service metadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tags, err := service.ReadRandomTags()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(tags))
	}
}

// GetTagByName retrieves tag by name
// @Summary Get tag by name
// @Description Returns tag details for the specified name (ONLY 1 TAG WITH SPECIFIC NAME)
// @Tags Tags
// @Produce json
// @Param name path string true "Tag name"
// @Success 200 {object} presenters.MetadataSuccessResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata/tags/byName/{name} [get]
func GetTagByName(service metadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		name := c.Params("name")
		tag, err := service.ReadTagByName(name)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataSuccessResponse(tag))
	}
}

// @Summary Get MANY tags by name LIKE
// @Description Returns all of the tags for the specified name like%
// @Tags Tags
// @Produce json
// @Param name path string true "Tag name"
// @Success 200 {object} presenters.MetadataSuccessResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata/tags/byNameLike/{name} [get]
func GetTagsByNameLike(service metadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		name := c.Params("name")
		tag, err := service.ReadTagsByNameLike(name)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataSuccessResponse(tag))
	}
}

// GetTagsInTrend retrieves trending tags
// @Summary Get trending tags
// @Description Returns a list of popular genres. Takes beats that were created this month (today minus 30 days), counts, how frequently were they used in beat_genres table.
// @Tags Tags
// @Produce json
// @Success 200 {object} presenters.MetadataListResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata/tags/in_trend [get]
func GetTagsInTrend(service metadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tags, err := service.ReadPopularTags()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(tags))
	}
}

// GetGenresInTrend retrieves trending genres
// @Summary Get trending genres
// @Description Returns a list of currently popular genres
// @Tags Metadata
// @Produce json
// @Success 200 {object} presenters.MetadataListResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata/genres/popular [get]
func GetGenresInTrend(service metadata.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tags, err := service.ReadPopularGenres()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(tags))
	}
}