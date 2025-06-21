package unpbeat

import (
	"encoding/json"
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/JulieWasNotAvailable/microservices/unpublished/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/unpublished/internal/entities"
	"github.com/JulieWasNotAvailable/microservices/unpublished/pkg/producer"
	"github.com/google/uuid"
)

type Service interface {
	CreateUnpublishedBeat(beat entities.UnpublishedBeat) (entities.UnpublishedBeat, error)
	GetAllUnpublishedBeats() ([]presenters.UnpublishedBeat, error)
	GetUnpublishedBeatByID(id uuid.UUID) (*presenters.UnpublishedBeat, error)
	GetUnpublishedBeatsByUser(userID uuid.UUID) ([]presenters.UnpublishedBeat, error)
	GetUnpublishedInModeration(from int64, to int64) (*[]presenters.UnpublishedBeat, error)
	GetUnpublishedByBeatmakerandStatus(userId uuid.UUID, status string) (*[]presenters.UnpublishedBeat, error)
	UpdateUnpublishedBeat(unpublished *entities.UnpublishedBeat) (*presenters.UnpublishedBeat, error)
	DeleteUnpublishedBeat(id uuid.UUID) error

	PublishBeat(userId uuid.UUID, beatOptions BeatPublishOptions, beatmakerName string) (*presenters.UnpublishedBeat, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{repository: r}
}

type License struct {
	LicenseTemplateID uint
	Price             int
}

type BeatPublishOptions struct {
	BeatId      string    `json:"beatId"`
	UserId      uuid.UUID `json:"-"`
	LicenseList []License `json:"licenseList"`
}

type CreateLicense struct {
	BeatId      uuid.UUID
	UserId      uuid.UUID
	LicenseList []License
}

func (s *service) CreateUnpublishedBeat(beat entities.UnpublishedBeat) (entities.UnpublishedBeat, error) {
	createdBeat, err := s.repository.CreateUnpublished(beat)
	return createdBeat, err
}

func (s *service) GetAllUnpublishedBeats() ([]presenters.UnpublishedBeat, error) {
	beats, err := s.repository.ReadUnpublished()
	return *beats, err
}

func (s *service) GetUnpublishedBeatByID(id uuid.UUID) (*presenters.UnpublishedBeat, error) {
	return s.repository.ReadUnpublishedById(id)
}

func (s *service) GetUnpublishedBeatsByUser(userID uuid.UUID) ([]presenters.UnpublishedBeat, error) {
	beats, err := s.repository.ReadUnpublishedByUser(userID)
	return *beats, err
}

func (s service) GetUnpublishedInModeration(from int64, to int64) (*[]presenters.UnpublishedBeat, error) {
	beats, err := s.repository.ReadUnpublishedInModeration(from, to)
	return beats, err
}

func (s *service) GetUnpublishedByBeatmakerandStatus(userId uuid.UUID, status string) (*[]presenters.UnpublishedBeat, error) {
	beats, err := s.repository.ReadUnpublishedByBeatmakerandStatus(userId, status)
	return beats, err
}

func (s *service) UpdateUnpublishedBeat(beat *entities.UnpublishedBeat) (*presenters.UnpublishedBeat, error) {
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
					log.Println("matched mood: ", moodEntity)
					log.Println("from the tag: ", cleanedName)
					log.Println(" ")
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

func (s *service) PublishBeat(userId uuid.UUID, beatOptions BeatPublishOptions, beatmakerName string) (*presenters.UnpublishedBeat, error) {
	beatuuid, err := uuid.Parse(beatOptions.BeatId)
	if err != nil {
		return nil, err
	}
	beat, err := s.GetUnpublishedBeatByID(beatuuid)
	if err != nil {
		return nil, err
	}
	if beat.BeatmakerID != userId {
		return nil, errors.New("unauthorized")
	}
	filename := beat.AvailableFiles.MP3Url
	if filename == "" {
		return nil, errors.New("mp3 file path is required to publish the beat")
	}
	_ = s.CheckFieldsBeforePublishing()
	toUpdateName := entities.UnpublishedBeat{
		ID:            beatuuid,
		BeatmakerName: beatmakerName,
	}
	_, err = s.repository.UpdateUnpublishedById(&toUpdateName)
	if err != nil {
		return nil, err
	}

	message := CreateLicense{
		BeatId:      beat.ID,
		UserId:      userId,
		LicenseList: beatOptions.LicenseList,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	err = producer.CreateMessage(messageBytes, "create_license")
	if err != nil {
		return nil, err
	}

	beat.Status = "processing"
	beatStatusUpdate := entities.UnpublishedBeat{
		ID:     beat.ID,
		Status: entities.StatusInModeration,
	}
	_, err = s.UpdateUnpublishedBeat(&beatStatusUpdate)
	if err != nil {
		return nil, err
	}

	return beat, nil
}

func (s *service) CheckFieldsBeforePublishing() error {
	return errors.New("unimplemented")
}
