package unpbeat

import (
	"encoding/json"
	"errors"
	"fmt"

	// "fmt"
	"regexp"
	"strings"

	"github.com/JulieWasNotAvailable/microservices/unpublished/internal/entities"
	"github.com/JulieWasNotAvailable/microservices/unpublished/pkg/producer"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Service interface {
	CreateUnpublishedBeat(beatmakeruuid uuid.UUID) (entities.UnpublishedBeat, error)
	GetAllUnpublishedBeats() ([]entities.UnpublishedBeat, error)
	GetUnpublishedBeatByID(id uuid.UUID) (*entities.UnpublishedBeat, error)
	GetUnpublishedBeatsByUser(userID uuid.UUID) ([]entities.UnpublishedBeat, error)
	GetUnpublishedInModeration(from int64, to int64) (*[]entities.UnpublishedBeat, error)
	GetUnpublishedByBeatmakerandStatus(userId uuid.UUID, status string) (*[]entities.UnpublishedBeat, error)
	UpdateUnpublishedBeat(beat *entities.UnpublishedBeat, beatmakerId uuid.UUID) (*entities.UnpublishedBeat, error)
	DeleteUnpublishedBeat(id uuid.UUID) error

	PublishBeat(beatmakerId uuid.UUID, beatOptions BeatPublishOptions, beatmakerName string) (*entities.UnpublishedBeat, []error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{repository: r}
}

type License struct {
	LicenseTemplateID uint `json:"id"`
	Price             int  `json:"price"`
}

type BeatPublishOptions struct {
	BeatId      string    `json:"beatId"`
	UserId      uuid.UUID `json:"-"`
	LicenseList []License `json:"licenseList"`
}

type CreateLicense struct {
	BeatId      uuid.UUID `json:"beatId"`
	UserId      uuid.UUID `json:"userId"`
	LicenseList []License `json:"licenseList"`
}

func (s *service) CreateUnpublishedBeat(beatmakeruuid uuid.UUID) (entities.UnpublishedBeat, error) {
	var emptyBeat entities.UnpublishedBeat
	emptyBeat.Status = entities.StatusDraft
	emptyBeat.BeatmakerID = beatmakeruuid

	createdBeat, err := s.repository.CreateUnpublished(emptyBeat)
	return createdBeat, err
}

func (s *service) GetAllUnpublishedBeats() ([]entities.UnpublishedBeat, error) {
	beats, err := s.repository.ReadUnpublished()
	return *beats, err
}

func (s *service) GetUnpublishedBeatByID(id uuid.UUID) (*entities.UnpublishedBeat, error) {
	return s.repository.ReadUnpublishedById(id)
}

func (s *service) GetUnpublishedBeatsByUser(userID uuid.UUID) ([]entities.UnpublishedBeat, error) {
	return s.repository.ReadUnpublishedByUser(userID)
}

func (s service) GetUnpublishedInModeration(from int64, to int64) (*[]entities.UnpublishedBeat, error) {
	beats, err := s.repository.ReadUnpublishedInModeration(from, to)
	return beats, err
}

func (s *service) GetUnpublishedByBeatmakerandStatus(userId uuid.UUID, status string) (*[]entities.UnpublishedBeat, error) {
	beats, err := s.repository.ReadUnpublishedByBeatmakerandStatus(userId, status)
	return beats, err
}

func (s *service) UpdateUnpublishedBeat(beat *entities.UnpublishedBeat, beatmakerId uuid.UUID) (*entities.UnpublishedBeat, error) {
	beatInitial, err := s.repository.ReadUnpublishedById(beat.ID)
	if err != nil {
		return &entities.UnpublishedBeat{}, err
	}
	if beatInitial.BeatmakerID != beatmakerId {
		return &entities.UnpublishedBeat{}, errors.New("tried to edit beat that does not belong to you")
	}
	return s.repository.UpdateUnpublishedById(beat)
}

func (s *service) UpdateUnpublishedBeatMass(beat *entities.UnpublishedBeat) (*entities.UnpublishedBeat, error) {
	toSave := entities.UnpublishedBeat{
		ID:          beat.ID,
		Name:        beat.Name,
		KeynoteID:   beat.KeynoteID,
		Picture:     beat.Picture,
		Description: beat.Description,
		Timestamps:  beat.Timestamps,
		Price:       beat.Price,
		Genres:      beat.Genres,
		Moods:       beat.Moods,
		Status:      beat.Status,
	}

	excludedWords := []string{
		"beat",
		"beats",
		"type",
		"free",
		"typebeat",
		"free for profit",
		"free for non-profit",
		"music",
		"profit",
		"insrumental",
		"for",
	}

	moods := []string{
		"chill",
		"dark",
		"sad",
		"energetic",
		"bouncy",
		"calm",
		"angry",
		"happy",
		"epic",
		"crazy",
		"intense",
		"hyper",
		"gloomy",
		"romantic",
		"excited",
		"dramatic",
		"majestic",
		"frantic",
		"scary",
		"dreamy",
		"slow",
	}

	if len(beat.Tags) != 0 {
		for _, tag := range beat.Tags {
			// Remove any year in format 20__
			yearPattern := regexp.MustCompile(`20\d{2}`)
			cleanedName := yearPattern.ReplaceAllString(tag.Name, "")

			// Remove other excluded words
			for _, word := range excludedWords {
				wordRegex := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(word) + `\b`)
				cleanedName = wordRegex.ReplaceAllString(cleanedName, "")
			}

			for _, mood := range moods {
				moodRegex := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(mood) + `\b`)
				if moodRegex.MatchString(cleanedName) {
					moodEntity, _ := s.repository.ReadMoodByName(mood)
					// log.Println("matched mood: ", moodEntity)
					// log.Println("from the tag: ", cleanedName)
					// log.Println(" ")
					toSave.Moods = append(toSave.Moods, *moodEntity)
					cleanedName = moodRegex.ReplaceAllString(cleanedName, "")
				}
			}

			if cleanedName == "" || cleanedName == " " {
				continue
			}

			// Clean up extra spaces
			spaceRegex := regexp.MustCompile(`\s+`)
			cleanedName = spaceRegex.ReplaceAllString(strings.TrimSpace(cleanedName), "_")

			id, _ := s.repository.CheckTagExists(cleanedName)
			if id != 0 {
				tag := entities.Tag{
					ID: id,
				}
				toSave.Tags = append(toSave.Tags, tag)
			} else {
				tag.Name = cleanedName
				s.repository.CreateTag(&tag)
				toSave.Tags = append(toSave.Tags, tag)
			}
		}
	}

	return s.repository.UpdateUnpublishedById(&toSave)
}

func (s *service) DeleteUnpublishedBeat(id uuid.UUID) error {
	return s.repository.DeleteUnpublishedById(id)
}

func (s *service) PublishBeat(beatmakerId uuid.UUID, beatOptions BeatPublishOptions, beatmakerName string) (*entities.UnpublishedBeat, []error) {
	errArr := []error{}
	beatuuid, err := uuid.Parse(beatOptions.BeatId)
	if err != nil {
		return nil, append(errArr, err)
	}
	beat, err := s.GetUnpublishedBeatByID(beatuuid)
	if err != nil {
		return nil, append(errArr, err)
	}
	if string(beat.Status) == string(entities.StatusInModeration) {
		return nil, append(errArr, errors.New("cannot publish beat with status in_process"))
	}

	toUpdateName := entities.UnpublishedBeat{
		ID:            beatuuid,
		BeatmakerName: beatmakerName,
	}
	_, err = s.repository.UpdateUnpublishedById(&toUpdateName)
	if err != nil {
		return nil, append(errArr, err)
	}

	_, errArr = s.CheckFieldsBeforePublishing(*beat)
	if errArr != nil {
		return nil, errArr
	}

	message := CreateLicense{
		BeatId:      beat.ID,
		UserId:      beatmakerId,
		LicenseList: beatOptions.LicenseList,
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return nil, append(errArr, err)
	}
	err = producer.CreateMessage(messageBytes, "create_license")
	if err != nil {
		return nil, append(errArr, err)
	}

	beatStatusUpdate := entities.UnpublishedBeat{
		ID:     beat.ID,
		Status: entities.StatusInModeration,
	}
	updated, err := s.UpdateUnpublishedBeat(&beatStatusUpdate, beatmakerId)
	if err != nil {
		return nil, append(errArr, err)
	}

	return updated, nil
}

func (s *service) CheckFieldsBeforePublishing(beat entities.UnpublishedBeat) (status bool, err []error) {
	validate := validator.New()
	valErr := validate.Struct(beat)
	errArray := []error{}
	if valErr != nil {
		for _, err := range valErr.(validator.ValidationErrors) {
			errMsg := fmt.Errorf("%s: поле необходимо к заполнению, но было предоставлено следующее значение: '%v')",
				err.Field(),
				err.Value())
			errArray = append(errArray, errMsg)
		}
		return false, errArray
	}
	return true, nil
}
