package beat

import (
	"log"
	"time"

	"github.com/JulieWasNotAvailable/microservices/beat/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/beat/pkg/entities"
	"github.com/google/uuid"
)

type Service interface {
	CreateBeat(beat entities.UnpublishedBeat, mfcc entities.MFCC) (entities.Beat, error)

	// helpers
	ReadBeats() (*[]entities.Beat, error)
	ReadBeatById(id uuid.UUID) (*presenters.Beat, error)
	ReadBeatsByBeatmakerId(id uuid.UUID) (*[]presenters.Beat, error)

	ReadFilteredBeats(filters presenters.Filters) (*[]presenters.Beat, error)
	ReadBeatsByMoodId(moodId uint) (*[]presenters.Beat, error)
	ReadBeatsByDate(from int64, to int64) (*[]presenters.Beat, error)
	FindBeatsWithAllMoods(moodIDs []uint) ([]presenters.Beat, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{repository: r}
}

func (s *service) CreateBeat(unpublishedBeat entities.UnpublishedBeat, mfcc entities.MFCC) (entities.Beat, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		log.Println(err)
	}

	beat := entities.Beat{
		ID:             uuid,
		Name:           unpublishedBeat.Name,
		Picture:        unpublishedBeat.Picture,
		BeatmakerID:    unpublishedBeat.BeatmakerID,
		AvailableFiles: unpublishedBeat.AvailableFiles,
		URL:            unpublishedBeat.AvailableFiles.MP3Url,
		Price:          unpublishedBeat.Price,
		Tags:           unpublishedBeat.Tags,
		BPM:            unpublishedBeat.BPM,
		Description:    unpublishedBeat.Description,
		Genres:         unpublishedBeat.Genres,
		Moods:          unpublishedBeat.Moods,
		KeynoteID:      unpublishedBeat.KeynoteID,
		Timestamps:     unpublishedBeat.Timestamps,
		Instruments:    unpublishedBeat.Instruments,
		MFCC:           mfcc,
		CreatedAt:      time.Now().Unix(), // или unpublishedBeat.CreatedAt, если нужно сохранить оригинальное значение
	}

	log.Println(beat)

	return s.repository.CreateBeat(beat)
}

// ReadBeatById implements Service.
func (s *service) ReadBeatById(id uuid.UUID) (*presenters.Beat, error) {
	return s.repository.ReadBeatById(id)
}

// ReadBeats implements Service.
func (s *service) ReadBeats() (*[]entities.Beat, error) {
	return s.repository.ReadBeats()
}

// ReadBeatsByBeatmakerId implements Service.
func (s *service) ReadBeatsByBeatmakerId(id uuid.UUID) (*[]presenters.Beat, error) {
	return s.repository.ReadBeatsByBeatmakerId(id)
}

func (s *service) ReadFilteredBeats(filters presenters.Filters) (*[]presenters.Beat, error) {
	return s.repository.ReadFilteredBeats(filters)
}

func (s *service) ReadBeatsByMoodId(moodId uint) (*[]presenters.Beat, error) {
	return s.repository.ReadBeatsByMoodId(moodId)
}

func (s *service) ReadBeatsByDate(from int64, to int64) (*[]presenters.Beat, error){
	return s.repository.ReadBeatsByDate(from, to)
}

func (s *service)FindBeatsWithAllMoods(moodIDs []uint) ([]presenters.Beat, error){
	return s.repository.FindBeatsWithAllMoods(moodIDs)
}
