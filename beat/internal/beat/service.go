package beat

import (
	"encoding/json"
	"log"
	"reflect"
	"time"

	"github.com/JulieWasNotAvailable/microservices/beat/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/beat/internal/entities"
	"github.com/JulieWasNotAvailable/microservices/beat/pkg/producer"
	"github.com/google/uuid"
)

type Service interface {
	CreateBeat(beat entities.UnpublishedBeat, mfcc []float64) (entities.Beat, error)

	// helpers
	ReadBeats() (*[]entities.Beat, error)
	ReadBeatById(id uuid.UUID) (*presenters.Beat, error)
	ReadBeatsByBeatmakerId(id uuid.UUID) (*[]presenters.Beat, error)
	DeleteBeatById(id string) error

	ReadFilteredBeats(filters presenters.Filters) (*[]presenters.Beat, error)
	ReadBeatsByMoodId(moodId uint) (*[]presenters.Beat, error)
	ReadBeatsByDate(from int64, to int64) (*[100]presenters.Beat, error)
	FindBeatsWithAllMoods(moodIDs []uint) ([]presenters.Beat, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{repository: r}
}

func (s *service) CreateBeat(unpublishedBeat entities.UnpublishedBeat, mfccfloat []float64) (entities.Beat, error) {
	mfcc := entities.MFCC{}
	mfcc, err := fillDataFromArray(mfccfloat, &mfcc)
	if err != nil {
		log.Println(err)
	}

	beat := entities.Beat{
		ID:             unpublishedBeat.ID,
		Name:           unpublishedBeat.Name,
		Picture:        unpublishedBeat.Picture,
		BeatmakerID:    unpublishedBeat.BeatmakerID,
		BeatmakerName:  unpublishedBeat.BeatmakerName,
		AvailableFiles: unpublishedBeat.AvailableFiles,
		URL:            unpublishedBeat.AvailableFiles.MP3Url,
		Price:          unpublishedBeat.Price,
		BPM:            unpublishedBeat.BPM,
		Description:    unpublishedBeat.Description,
		Genres:         unpublishedBeat.Genres,
		Moods:          unpublishedBeat.Moods,
		Tags:           unpublishedBeat.Tags,
		KeynoteID:      unpublishedBeat.KeynoteID,
		Timestamps:     unpublishedBeat.Timestamps,
		// Instruments:    unpublishedBeat.Instruments,
		MFCC:      mfcc,
		CreatedAt: time.Now().Unix(), // или unpublishedBeat.CreatedAt, если нужно сохранить оригинальное значение
	}

	beat, err = s.repository.CreateBeat(beat)
	if err != nil {
		message := producer.KafkaMessageValue{
			ID:  unpublishedBeat.ID.String(),
			Err: "",
		}
		messageInBytes, err := json.Marshal(message)
		if err != nil {
			return entities.Beat{}, err
		}
		producer.CreateMessage(messageInBytes, "delete_approve")
		return entities.Beat{}, err
	}

	message := producer.KafkaMessageValue{
		ID:  beat.ID.String(),
		Err: "",
	}

	messageInBytes, err := json.Marshal(message)
	if err != nil {
		return entities.Beat{}, err
	}

	producer.CreateMessage(messageInBytes, "delete_approve")

	return beat, nil
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

func (s *service) ReadBeatsByDate(from int64, to int64) (*[100]presenters.Beat, error) {
	return s.repository.ReadBeatsByDate(from, to)
}

func (s *service) DeleteBeatById(id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repository.DeleteBeatById(uuid)
}

func (s *service) FindBeatsWithAllMoods(moodIDs []uint) ([]presenters.Beat, error) {
	return s.repository.FindBeatsWithAllMoods(moodIDs)
}

func fillDataFromArray(arr []float64, data *entities.MFCC) (entities.MFCC, error) {
	val := reflect.ValueOf(data).Elem() // Dereference the pointer to the struct

	// Iterate over struct fields and assign from array
	for i := 2; i < (val.NumField() - 1); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.Float64 && field.CanSet() {
			field.SetFloat(arr[i-2]) // Assign array value to struct field
		}
	}

	return *data, nil
}
