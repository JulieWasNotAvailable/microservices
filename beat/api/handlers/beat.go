package handlers

import (
	"errors"
	"net/http"

	"github.com/JulieWasNotAvailable/microservices/beat/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/beat/internal/beat"
	"github.com/JulieWasNotAvailable/microservices/beat/internal/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateBeat creates a new beat
//	@Summary		Create a beat
//	@Description	Creates a new empty beat
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			beat	body		entities.UnpublishedBeat		true	"Beat creation data"
//	@Success		201		{object}	presenters.BeatSuccessResponse	"Created beat details"
//	@Failure		400		{object}	presenters.BeatErrorResponse	"Invalid request body"
//	@Failure		500		{object}	presenters.BeatErrorResponse	"Internal server error"
//	@Router			/beat/exampleBeat [post]
func CreateBeat(service beat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse request body into Beat entity
		newBeat := entities.UnpublishedBeat{
			Name:      "1146",
			KeynoteID: 1,
		}
		mfccfloat := []float64{1, 1.234567325, 2.567865430, 3.456754321, 3.456754321, 3.456754321,
			3.456754322, 3.456754323, 3.456754324, 3.456754325, 3.456754326,
			3.456754327, 3.456754328, 3.456754329, 3.456754330, 3.456754331,
			3.456754332, 3.456754333, 3.456754334, 3.456754335, 3.456754336,
			3.456754337, 3.456754338, 3.456754339, 3.456754340, 3.456754341,
			3.456754342, 3.456754343, 3.456754344, 3.456754345, 3.456754346,
			3.456754347, 3.456754348, 3.456754349, 3.456754350, 3.456754351,
			3.456754352, 3.456754353, 3.456754354, 3.456754355, 3.456754356,
			3.456754357, 3.456754358, 3.456754359, 3.456754360, 3.456754361,
			3.456754362, 3.456754363, 3.456754364, 3.456754365, 3.456754366,
			3.456754367, 3.456754368, 3.456754369, 3.456754370, 3.456754371,
			3.456754372, 3.456754373, 3.456754374, 3.456754375, 3.456754376,
			3.456754377, 3.456754378, 3.456754379, 3.456754380}
		if err := c.BodyParser(&newBeat); err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
		}

		// Call service to create the beat
		createdBeat, err := service.CreateBeat(newBeat, mfccfloat)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		return c.Status(http.StatusCreated).JSON(presenters.CreateBeatSuccessResponse(createdBeat))
	}
}

// GetAllBeats retrieves all beats
//	@Summary		Get all beats
//	@Description	Returns all beats in the system
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	presenters.BeatListResponse		"List of all beats"
//	@Failure		500	{object}	presenters.BeatErrorResponse	"Internal server error"
//	@Router			/beat/all [get]
func GetAllBeats(service beat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		beats, err := service.ReadBeats()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(presenters.CreateBeatListResponse(&beats))
	}
}

// GetBeatById retrieves a beat by ID
//	@Summary		Get beat by ID
//	@Description	Returns a single beat by its ID
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			beatId	path		string							true	"Beat ID in UUID format"
//	@Success		200		{object}	presenters.BeatSuccessResponse	"Beat details"
//	@Failure		400		{object}	presenters.BeatErrorResponse	"Invalid ID format"
//	@Failure		404		{object}	presenters.BeatErrorResponse	"Beat not found"
//	@Failure		500		{object}	presenters.BeatErrorResponse	"Internal server error"
//	@Router			/beat/byBeatId/{beatId} [get]
func GetBeatById(service beat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("beatId"))
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

// GetBeatsByBeatmakerId retrieves beats by beatmaker ID
//	@Summary		Get beats by beatmaker
//	@Description	Returns all beats for a specific beatmaker
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			beatmakerId	path		string							true	"Beatmaker ID in UUID format"
//	@Success		200			{object}	presenters.BeatListResponse		"List of beats"
//	@Failure		400			{object}	presenters.BeatErrorResponse	"Invalid ID format"
//	@Failure		500			{object}	presenters.BeatErrorResponse	"Internal server error"
//	@Router			/beat/byBeatmakerId/{beatmakerId} [get]
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

// GetBeatsByJWT retrieves beats by JWT
//
//	@Summary		Get beats by JWT
//	@Description	Returns all beats
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			beatmakerId	path		string							true	"Beatmaker ID in UUID format"
//	@Success		200			{object}	presenters.BeatListResponse		"List of beats"
//	@Failure		400			{object}	presenters.BeatErrorResponse	"Invalid ID format"
//	@Failure		500			{object}	presenters.BeatErrorResponse	"Internal server error"
//	@Router			/beat/byBeatmakerByJWT [get]
func GetBeatsByJWT(service beat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse UUID from request parameters
		beatmakerId, err := getIdFromJWT(c)
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

// GetFilteredBeats retrieves beats with filters
//	@Summary		Filter beats
//	@Description	Returns beats matching the provided filters
//	@Tags			Filters
//	@Accept			json
//	@Produce		json
//	@Param			filters	body		presenters.Filters				true	"Filter criteria"
//	@Success		200		{object}	presenters.BeatListResponse		"Filtered list of beats"
//	@Failure		400		{object}	presenters.BeatErrorResponse	"Invalid filter format"
//	@Failure		500		{object}	presenters.BeatErrorResponse	"Internal server error"
//	@Router			/beat/filteredBeats [get]
func GetFilteredBeats(service beat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		filters := presenters.Filters{}
		if err := c.BodyParser(&filters); err != nil {
			return c.Status(http.StatusBadRequest).JSON(presenters.CreateBeatErrorResponse(err))
		}

		//add pagination

		beats, err := service.ReadFilteredBeats(filters)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}
		return c.Status(http.StatusOK).JSON(presenters.CreateBeatListResponse(beats))
	}
}

// GetBeatsByMoodId retrieves beats by mood ID
//	@Summary		Get beats by mood
//	@Description	Returns beats matching a specific mood ID
//	@Tags			Filters
//	@Accept			json
//	@Produce		json
//	@Param			moodId	path		int								true	"Mood ID"
//	@Success		200		{object}	presenters.BeatListResponse		"List of beats with this mood"
//	@Failure		400		{object}	presenters.BeatErrorResponse	"Invalid mood ID"
//	@Failure		500		{object}	presenters.BeatErrorResponse	"Internal server error"
//	@Router			/beat/beatsByMoodId/{moodId} [get]
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

// GetBeatsWithAllMoods retrieves beats matching all specified moods
//	@Summary		Get beats by multiple moods
//	@Description	Returns beats that match ALL the specified mood IDs
//	@Tags			Filters
//	@Accept			json
//	@Produce		json
//	@Param			filters	body		presenters.Filters				true	"Mood IDs to filter by"
//	@Success		200		{object}	presenters.BeatListResponse		"List of matching beats"
//	@Failure		400		{object}	presenters.BeatErrorResponse	"Invalid mood format"
//	@Failure		500		{object}	presenters.BeatErrorResponse	"Internal server error"
//	@Router			/beat/withAllMoods [get]
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

// GetBeatsByDate retrieves beats within a date range
//	@Summary		Get beats by date range
//	@Description	Returns beats created between the specified timestamps
//	@Tags			Filters
//	@Accept			json
//	@Produce		json
//	@Param			from	path		int								true	"Start timestamp (Unix epoch)"
//	@Param			to		path		int								true	"End timestamp (Unix epoch)"
//	@Success		200		{object}	presenters.BeatListResponse		"List of beats in date range"
//	@Failure		400		{object}	presenters.BeatErrorResponse	"Invalid timestamp format"
//	@Failure		500		{object}	presenters.BeatErrorResponse	"Internal server error"
//	@Router			/beat/beatsByDate/{from}/{to} [get]
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

// DeleteBeatById godoc
//	@Summary		Delete beat by ID
//	@Description	Deletes beat with the specified ID
//	@Tags			Beats
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Beat ID (UUID format)"
//	@Success		200	{object}	presenters.BeatSuccessResponse
//	@Failure		422	{object}	presenters.BeatErrorResponse
//	@Failure		500	{object}	presenters.BeatErrorResponse
//	@Router			/beat/deleteBeatById/{id} [delete]
func DeleteBeatById(service beat.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		beatId := c.Params("id")

		err := service.DeleteBeatById(beatId)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(presenters.CreateBeatErrorResponse(err))
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"message": "deleted successfully",
		})
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
