package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/JulieWasNotAvailable/microservices/beat/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/beat/pkg/beat"
	"github.com/JulieWasNotAvailable/microservices/beat/pkg/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateBeat(service beat.Service) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Parse request body into Beat entity
        var newBeat entities.UnpublishedBeat
		mfcc := entities.MFCC{
			Col1 : 1,
			Col2 : 2,
			// BeatId: newBeat.ID,
		}
        if err := c.BodyParser(&newBeat); err != nil {
            return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
        }

        // Call service to create the beat
        createdBeat, err := service.CreateBeat(newBeat, mfcc)
        if err != nil {
            return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
        }

        return c.Status(http.StatusCreated).JSON(presenters.CreateBeatSuccessResponse(createdBeat))
    }
}

func GetAllBeats(service beat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		beats, err := service.ReadBeats()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateBeatListResponse(&beats))
	}
}

func GetBeatById(service beat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
		}
		//license join should be here
		beat, err := service.ReadBeatById(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(http.StatusNotFound).JSON(presenters.CreateBeatErrorResponse(err))
			}
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateBeatSuccessResponse(beat))
	}
}

func GetBeatsByBeatmakerId(service beat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse UUID from request parameters
		beatmakerId, err := uuid.Parse(c.Params("beatmakerId"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
		}

		// Call service to get beats by beatmaker ID
		beats, err := service.ReadBeatsByBeatmakerId(beatmakerId)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateBeatListResponse(beats))
	}
}

func GetFilteredBeats(service beat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		filters := presenters.Filters{}
		if err := c.BodyParser(&filters); err != nil {
            return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
        }

		log.Println(filters)

		//add pagination

		beats, err := service.ReadFilteredBeats(filters)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}
		return c.Status(http.StatusOK).JSON(presenters.CreateBeatListResponse(beats))
	}
}

func GetBeatsByMoodId(service beat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error { 
		moodId, err := c.ParamsInt("moodId")
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
		}

		beats, err := service.ReadBeatsByMoodId(uint(moodId))
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateBeatListResponse(beats))
	}
}

func GetBeatsWithAllMoods(service beat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error { 
		filters := presenters.Filters{}
		if err := c.BodyParser(&filters); err != nil {
            return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
        }

		beats, err := service.FindBeatsWithAllMoods(filters.Moods)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateBeatListResponse(beats))
	}
}


func GetBeatsByDate(service beat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error { 
		from, err := c.ParamsInt("from")
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
		}

		to, err := c.ParamsInt("to")
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
		}
		
		beats, err := service.ReadBeatsByDate(int64(from), int64(to))
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateBeatListResponse(beats))
	}
}

// Helper function to filter beats based on query parameters
// func filterBeats(beats *[]presenters.Beat, filters struct {
// 	Genre    string
// 	Mood     string
// 	BPM      int
// 	Key      string
// 	Page     int
// 	PageSize int
// }) []presenters.Beat {
// 	filtered := make([]presenters.Beat, 0, len(*beats))

// 	for _, beat := range *beats {
// 		// Filter by genre if specified
// 		if filters.Genre != "" {
// 			genreMatch := false
// 			for _, genre := range beat.Genres {
// 				if strings.EqualFold(genre.Name, filters.Genre) {
// 					genreMatch = true
// 					break
// 				}
// 			}
// 			if !genreMatch {
// 				continue
// 			}
// 		}

// 		// Filter by mood if specified
// 		if filters.Mood != "" {
// 			moodMatch := false
// 			for _, mood := range beat.Moods {
// 				if strings.EqualFold(mood.Name, filters.Mood) {
// 					moodMatch = true
// 					break
// 				}
// 			}
// 			if !moodMatch {
// 				continue
// 			}
// 		}

// 		// Filter by BPM if specified
// 		if filters.BPM > 0 && beat.BPM != filters.BPM {
// 			continue
// 		}

// 		// Filter by key if specified
// 		if filters.Key != "" && !strings.EqualFold(beat.Key, filters.Key) {
// 			continue
// 		}

// 		filtered = append(filtered, beat)
// 	}

// 	// Apply pagination if specified
// 	if filters.Page > 0 && filters.PageSize > 0 {
// 		start := (filters.Page - 1) * filters.PageSize
// 		end := start + filters.PageSize
// 		if start >= len(filtered) {
// 			return []presenters.Beat{}
// 		}
// 		if end > len(filtered) {
// 			end = len(filtered)
// 		}
// 		return filtered[start:end]
// 	}

// 	return filtered
// }
