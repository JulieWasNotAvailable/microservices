package beatmetadata

import (
	"errors"
	"log"

	"github.com/JulieWasNotAvailable/microservices/unpublished/pkg/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MetadataRepository interface {
	CreateAvailableFiles(availableFiles *entities.AvailableFiles) (entities.AvailableFiles, error)
	ReadAllAvailableFiles() (*[]entities.AvailableFiles, error)
	UpdateAvailableFiles(availableFiles *entities.AvailableFiles) (entities.AvailableFiles, error)

	CreateInstrument(instrument *entities.Instrument) (*entities.Instrument, error)
	ReadAllInstruments() (*[]entities.Instrument, error)

	CreateGenre(genre *entities.Genre) (*entities.Genre, error)
	ReadAllGenres() (*[]entities.Genre, error)

	CreateTimestamp(timestamp *entities.Timestamp) (*entities.Timestamp, error)
	ReadAllTimestamps() (*[]entities.Timestamp, error)
	DeleteTimestampById(id uint) error

	CreateTag(tag *entities.Tag) (*entities.Tag, error)
	ReadAllTags() (*[]entities.Tag, error)
	ReadTagById(id uint) (*entities.Tag, error)
	DeleteTagById(id uint) error

	CreateMood(mood *entities.Mood) (*entities.Mood, error)
	ReadAllMoods() (*[]entities.Mood, error)

	CreateKeynote(keynote *entities.Keynote) (*entities.Keynote, error)
	ReadAllKeynotes() (*[]entities.Keynote, error)
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) MetadataRepository {
	return &repository{
		DB: db,
	}
}

func (r *repository) CreateAvailableFiles(availableFiles *entities.AvailableFiles) (entities.AvailableFiles, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return entities.AvailableFiles{}, err
	}
	availableFiles.ID = uuid
	err = r.DB.Create(availableFiles).Error
	if err != nil {
		return entities.AvailableFiles{}, err
	}
	return *availableFiles, nil
}

func (r *repository) ReadAllAvailableFiles() (*[]entities.AvailableFiles, error) {
	var availFiles []entities.AvailableFiles
	err := r.DB.Find(&availFiles).Error
	if err != nil {
		return nil, err
	}
	return &availFiles, nil
}

func (r *repository) UpdateAvailableFiles(availableFiles *entities.AvailableFiles) (entities.AvailableFiles, error) {
	log.Println(availableFiles)
	err := r.DB.Where("id = ?", availableFiles.ID).Updates(&availableFiles).Error
	if err != nil {
		return entities.AvailableFiles{}, err
	}

	return *availableFiles, nil
}

func (r *repository) CreateInstrument(instrument *entities.Instrument) (*entities.Instrument, error) {
	result := r.DB.Create(instrument)
	if result.Error != nil {
		return nil, result.Error
	}
	return instrument, nil
}

func (r *repository) ReadAllInstruments() (*[]entities.Instrument, error) {
	var instruments []entities.Instrument
	result := r.DB.Find(&instruments)
	if result.Error != nil {
		return nil, result.Error
	}
	return &instruments, nil
}

func (r *repository) CreateGenre(genre *entities.Genre) (*entities.Genre, error) {
	result := r.DB.Create(genre)
	if result.Error != nil {
		return nil, result.Error
	}
	return genre, nil
}

func (r *repository) ReadAllGenres() (*[]entities.Genre, error) {
	var genres []entities.Genre
	result := r.DB.Find(&genres)
	if result.Error != nil {
		return nil, result.Error
	}
	return &genres, nil
}

func (r *repository) CreateTimestamp(timestamp *entities.Timestamp) (*entities.Timestamp, error) {
	result := r.DB.Create(timestamp)
	if result.Error != nil {
		return nil, result.Error
	}
	return timestamp, nil
}

func (r *repository) ReadAllTimestamps() (*[]entities.Timestamp, error) {
	var timestamps []entities.Timestamp
	result := r.DB.Find(&timestamps)
	if result.Error != nil {
		return nil, result.Error
	}
	return &timestamps, nil
}

func (r *repository) DeleteTimestampById(id uint) error {
	timestamp := entities.Timestamp{}
	result := r.DB.Where("id = ?", id).Delete(timestamp)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("timestamp was not found")
	}

	return nil
}

func (r *repository) CreateTag(tag *entities.Tag) (*entities.Tag, error) {
	result := r.DB.Create(tag)
	if result.Error != nil {
		return nil, result.Error
	}
	return tag, nil
}

func (r *repository) ReadAllTags() (*[]entities.Tag, error) {
	var tags []entities.Tag
	result := r.DB.Find(&tags)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tags, nil
}

func (r *repository) ReadTagById(id uint) (*entities.Tag, error) {
	var tag entities.Tag
	result := r.DB.First(&tag).Where("id = ?")
	if result.Error != nil {
		return nil, result.Error
	}
	return &tag, nil
}

func (r *repository) DeleteTagById(id uint) error {
	tagModel := entities.Tag{}
	result := r.DB.Where("id = ?", id).Delete(tagModel)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("tag was not found")
	}

	return nil
}

func (r *repository) CreateMood(mood *entities.Mood) (*entities.Mood, error) {
	result := r.DB.Create(mood)
	if result.Error != nil {
		return nil, result.Error
	}
	return mood, nil
}

func (r *repository) ReadAllMoods() (*[]entities.Mood, error) {
	var moods []entities.Mood
	result := r.DB.Find(&moods)
	if result.Error != nil {
		return nil, result.Error
	}
	return &moods, nil
}

func (r *repository) CreateKeynote(keynote *entities.Keynote) (*entities.Keynote, error) {
	result := r.DB.Create(keynote)
	if result.Error != nil {
		return nil, result.Error
	}
	return keynote, nil
}

func (r *repository) ReadAllKeynotes() (*[]entities.Keynote, error) {
	var keynotes []entities.Keynote
	result := r.DB.Find(&keynotes)
	if result.Error != nil {
		return nil, result.Error
	}
	return &keynotes, nil
}
