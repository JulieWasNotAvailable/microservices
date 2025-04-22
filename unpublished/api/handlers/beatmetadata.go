package handlers

import (
	"net/http"

	"github.com/JulieWasNotAvailable/microservices/unpublished/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/unpublished/pkg/beatmetadata"
	"github.com/JulieWasNotAvailable/microservices/unpublished/pkg/entities"
	"github.com/gofiber/fiber/v2"
)

// PostFile godoc
// @Summary Create available files
// @Description Create new available files metadata
// @Tags files
// @Accept json
// @Produce json
// @Param availableFiles body entities.AvailableFiles true "Available files data"
// @Success 201 {object} object "Successfully created files metadata"
// @Failure 400 {object} object "Invalid request body"
// @Failure 500 {object} object "Internal server error"
// @Router /metadata/files [post]
func PostFile(service beatmetadata.MetadataService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var availableFiles entities.AvailableFiles
		if err := c.BodyParser(&availableFiles); err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		createdFiles, err := service.CreateAvailableFiles(&availableFiles)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusCreated).JSON(presenters.CreateMetadataSuccessResponse(createdFiles))
	}
}

// GetAllFiles godoc
// @Summary Get all available files
// @Description Retrieve all available files metadata
// @Tags files
// @Produce json
// @Success 200 {object} object "Successfully retrieved all files metadata"
// @Failure 500 {object} object "Internal server error"
// @Router /metadata/files [get]
func GetAllFiles(service beatmetadata.MetadataService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		availableFiles, err := service.ReadAllAvailableFiles()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(availableFiles))
	}
}

// UpdateFiles godoc
// @Summary Update available files
// @Description Update existing available files metadata
// @Tags files
// @Accept json
// @Produce json
// @Param availableFiles body entities.AvailableFiles true "Updated files data"
// @Success 200 {object} object "Successfully updated files metadata"
// @Failure 400 {object} object "Invalid request body"
// @Failure 500 {object} object "Internal server error"
// @Router /metadata/files [put]
func UpdateFiles(service beatmetadata.MetadataService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var availableFiles entities.AvailableFiles
		if err := c.BodyParser(&availableFiles); err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		updatedFiles, err := service.UpdateAvailableFiles(&availableFiles)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(updatedFiles))
	}}
	
// PostInstrument godoc
// @Summary Create a new instrument
// @Description Add a new instrument to the system
// @Tags instruments
// @Accept json
// @Produce json
// @Param instrument body entities.Instrument true "Instrument to create"
// @Success 201 {object} presenters.MetadataSuccessResponse
// @Failure 422 {object} presenters.MetadataErrorResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata/instruments [post]
func PostInstrument(service beatmetadata.MetadataService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var instrument entities.Instrument
		if err := c.BodyParser(&instrument); err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		created, err := service.CreateInstrument(&instrument)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusCreated).JSON(presenters.CreateMetadataSuccessResponse(created))
	}
}

// GetInstruments godoc
// @Summary Get all instruments
// @Description Retrieve all instruments from the system
// @Tags instruments
// @Produce json
// @Success 200 {object} presenters.MetadataListResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata/instruments [get]
func GetInstruments(service beatmetadata.MetadataService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		instruments, err := service.GetAllInstruments()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(instruments))
	}
}

// PostGenre godoc
// @Summary Create a new genre
// @Description Add a new genre to the system
// @Tags genres
// @Accept json
// @Produce json
// @Param genre body entities.Genre true "Genre to create"
// @Success 201 {object} presenters.MetadataSuccessResponse
// @Failure 422 {object} presenters.MetadataErrorResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata/genres [post]
func PostGenre(service beatmetadata.MetadataService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var genre entities.Genre
		if err := c.BodyParser(&genre); err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		created, err := service.CreateGenre(&genre)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusCreated).JSON(presenters.CreateMetadataSuccessResponse(created))
	}
}

// GetGenres godoc
// @Summary Get all genres
// @Description Retrieve all genres from the system
// @Tags genres
// @Produce json
// @Success 200 {object} presenters.MetadataListResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata/genres [get]
func GetGenres(service beatmetadata.MetadataService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		genres, err := service.GetAllGenres()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(genres))
	}
}

// PostTimestamp godoc
// @Summary Create a new timestamp
// @Description Add a new timestamp to the system
// @Tags timestamps
// @Accept json
// @Produce json
// @Param timestamp body entities.Timestamp true "Timestamp to create"
// @Success 201 {object} presenters.MetadataSuccessResponse
// @Failure 422 {object} presenters.MetadataErrorResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata/timestamps [post]
func PostTimestamp(service beatmetadata.MetadataService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var timestamp entities.Timestamp
		if err := c.BodyParser(&timestamp); err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		created, err := service.CreateTimestamp(&timestamp)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusCreated).JSON(presenters.CreateMetadataSuccessResponse(created))
	}
}

// GetTimestamps godoc
// @Summary Get all timestamps
// @Description Retrieve all timestamps from the system
// @Tags timestamps
// @Produce json
// @Success 200 {object} presenters.MetadataListResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata/timestamps [get]
func GetTimestamps(service beatmetadata.MetadataService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		timestamps, err := service.GetAllTimestamps()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(timestamps))
	}
}

// PostTag godoc
// @Summary Create a new tag
// @Description Add a new tag to the system
// @Tags tags
// @Accept json
// @Produce json
// @Param tag body entities.Tag true "Tag to create"
// @Success 201 {object} presenters.MetadataSuccessResponse
// @Failure 422 {object} presenters.MetadataErrorResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata/tags [post]
func PostTag(service beatmetadata.MetadataService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tag entities.Tag
		if err := c.BodyParser(&tag); err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		created, err := service.CreateTag(&tag)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusCreated).JSON(presenters.CreateMetadataSuccessResponse(created))
	}
}

// GetTags godoc
// @Summary Get all tags
// @Description Retrieve all tags from the system
// @Tags tags
// @Produce json
// @Success 200 {object} presenters.MetadataListResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata/tags [get]
func GetTags(service beatmetadata.MetadataService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tags, err := service.GetAllTags()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(tags))
	}
}

// PostMood godoc
// @Summary Create a new mood
// @Description Add a new mood to the system
// @Tags moods
// @Accept json
// @Produce json
// @Param mood body entities.Mood true "Mood to create"
// @Success 201 {object} presenters.MetadataSuccessResponse
// @Failure 422 {object} presenters.MetadataErrorResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata/moods [post]
func PostMood(service beatmetadata.MetadataService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var mood entities.Mood
		if err := c.BodyParser(&mood); err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		created, err := service.CreateMood(&mood)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusCreated).JSON(presenters.CreateMetadataSuccessResponse(created))
	}
}

// GetMoods godoc
// @Summary Get all moods
// @Description Retrieve all moods from the system
// @Tags moods
// @Produce json
// @Success 200 {object} presenters.MetadataListResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata/moods [get]
func GetMoods(service beatmetadata.MetadataService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		moods, err := service.GetAllMoods()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(moods))
	}
}

// PostKeynote godoc
// @Summary Create a new keynote
// @Description Add a new keynote to the system
// @Tags keynotes
// @Accept json
// @Produce json
// @Param keynote body entities.Keynote true "Keynote to create"
// @Success 201 {object} presenters.MetadataSuccessResponse
// @Failure 422 {object} presenters.MetadataErrorResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata/keynotes [post]
func PostKeynote(service beatmetadata.MetadataService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		 var keynote entities.Keynote
		if err := c.BodyParser(&keynote); err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		created, err := service.CreateKeynote(&keynote)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusCreated).JSON(presenters.CreateMetadataSuccessResponse(created))
	}
}

// GetKeynotes godoc
// @Summary Get all keynotes
// @Description Retrieve all keynotes from the system
// @Tags keynotes
// @Produce json
// @Success 200 {object} presenters.MetadataListResponse
// @Failure 500 {object} presenters.MetadataErrorResponse
// @Router /metadata/keynotes [get]
func GetKeynotes(service beatmetadata.MetadataService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		keynotes, err := service.GetAllKeynotes()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateMetadataErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateMetadataListResponse(keynotes))
	}
}