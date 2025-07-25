package unpbeat

import (
	"errors"
	"fmt"
	"time"

	"github.com/JulieWasNotAvailable/microservices/unpublished/internal/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUnpublished(unpublished entities.UnpublishedBeat) (entities.UnpublishedBeat, error)
	ReadUnpublished() (*[]entities.UnpublishedBeat, error)
	ReadUnpublishedById(id uuid.UUID) (*entities.UnpublishedBeat, error)
	ReadUnpublishedByUser(userId uuid.UUID) ([]entities.UnpublishedBeat, error)
	ReadUnpublishedInModeration(from int64, to int64) (*[]entities.UnpublishedBeat, error)
	ReadUnpublishedByBeatmakerandStatus(userId uuid.UUID, status string) (*[]entities.UnpublishedBeat, error)
	UpdateUnpublishedById(unpublished *entities.UnpublishedBeat) (*entities.UnpublishedBeat, error)
	DeleteUnpublishedById(id uuid.UUID) error

	CheckTagExists(name string) (uint, error)
	CreateTag(tag *entities.Tag) (*entities.Tag, error)
	ReadMoodByName(name string) (*entities.Mood, error)
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r repository) CreateUnpublished(emptyBeat entities.UnpublishedBeat) (entities.UnpublishedBeat, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return entities.UnpublishedBeat{}, err
	}

	emptyBeat.ID = id
	time := time.Now().Unix()
	emptyBeat.CreatedAt = time
	emptyBeat.Price = 1000
	result := r.DB.Create(&emptyBeat)

	if result.Error != nil {
		return entities.UnpublishedBeat{}, result.Error
	}

	created := entities.UnpublishedBeat{}
	err = r.DB.Where("id = ?", emptyBeat.ID).Find(&created).Error
	if err != nil {
		return entities.UnpublishedBeat{}, err
	}

	return created, nil
}

func (r *repository) ReadUnpublished() (*[]entities.UnpublishedBeat, error) {
	var unpublishedBeats []entities.UnpublishedBeat

	result := r.DB.Scopes(WithBasicPreloads()).Find(&unpublishedBeats)
	if result.Error != nil {
		return nil, result.Error
	}

	return &unpublishedBeats, nil
}

// ReadUnpublishedById reads a single unpublished beat by ID
func (r *repository) ReadUnpublishedById(id uuid.UUID) (*entities.UnpublishedBeat, error) {
	var unpublished entities.UnpublishedBeat

	result := r.DB.Where("id = ?", id).Scopes(WithBasicPreloads()).First(&unpublished)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("unpublished beat not found")
		}
		return nil, result.Error
	}

	return &unpublished, nil
}

func (r *repository) ReadUnpublishedByUser(userId uuid.UUID) ([]entities.UnpublishedBeat, error) {
	var unpublishedBeats []entities.UnpublishedBeat

	result := r.DB.Model(unpublishedBeats).
		Where("beatmaker_id = ?", userId).Scopes(WithBasicPreloads()).
		Find(&unpublishedBeats)

	if result.Error != nil {
		return nil, result.Error
	}

	return unpublishedBeats, nil
}

func (r *repository) ReadUnpublishedByBeatmakerandStatus(userId uuid.UUID, status string) (*[]entities.UnpublishedBeat, error) {
	var unpublishedBeats []entities.UnpublishedBeat

	result := r.DB.Model(unpublishedBeats).
		Where("beatmaker_id = ?", userId).Where("status = ?", status).Scopes(WithBasicPreloads()).
		Find(&unpublishedBeats)

	if result.Error != nil {
		return nil, result.Error
	}

	return &unpublishedBeats, nil
}

func (r *repository) ReadUnpublishedInModeration(from int64, to int64) (*[]entities.UnpublishedBeat, error) {
	var unpublishedBeats []entities.UnpublishedBeat

	result := r.DB.Model(unpublishedBeats).
		Where("status = ?", "in_moderation").Where("sent_to_moderation_at >= ? AND sent_to_moderation_at <= ?", from, to).
		Scopes(WithBasicPreloads()).
		Find(&unpublishedBeats)

	if result.Error != nil {
		return nil, result.Error
	}

	return &unpublishedBeats, nil
}

func (r *repository) UpdateUnpublishedById(unpublished *entities.UnpublishedBeat) (*entities.UnpublishedBeat, error) {
	beatInitial := &entities.UnpublishedBeat{}
	err := r.DB.Where("id = ?", unpublished.ID).First(&beatInitial).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("unpublished beat not found or not owned by user")
		} else {
			return nil, err
		}
	}

	genres := unpublished.Genres
	for _, genre := range genres {
		err := r.DB.Where("id = ?", genre.ID).First(&entities.Genre{}).Error
		if err != nil {
			return nil, fmt.Errorf("genre %d does not exist", genre.ID)
		}
	}
	moods := unpublished.Moods
	for _, mood := range moods {
		err := r.DB.Where("id = ?", mood.ID).First(&entities.Mood{}).Error
		if err != nil {
			return nil, fmt.Errorf("mood %d does not exist", mood.ID)
		}
	}
	instruments := unpublished.Instruments
	for _, instrument := range instruments {
		err := r.DB.Where("id = ?", instrument.ID).First(&entities.Instrument{}).Error
		if err != nil {
			return nil, fmt.Errorf("instrument %d does not exist", instrument.ID)
		}
	}

	emptyModel := entities.AvailableFiles{}
	err = r.DB.Transaction(func(tx *gorm.DB) error {

		if unpublished.Tags != nil {
			err := r.DB.Model(&unpublished).Association("Tags").Replace(unpublished.Tags)
			if err != nil {
				return err
			}
		}

		if unpublished.Genres != nil {
			err := r.DB.Model(&unpublished).Association("Genres").Replace(unpublished.Genres)
			if err != nil {
				return err
			}
		}

		if unpublished.Moods != nil {
			err := r.DB.Model(&unpublished).Association("Moods").Replace(unpublished.Moods)
			if err != nil {
				return err
			}
		}

		if unpublished.Instruments != nil {
			err := r.DB.Model(&unpublished).Association("Instruments").Replace(unpublished.Instruments)
			if err != nil {
				return err
			}
		}

		if unpublished.Timestamps != nil {
			err := r.DB.Model(&unpublished).Association("Timestamps").Replace(unpublished.Timestamps)
			if err != nil {
				return err
			}
		}

		if unpublished.AvailableFiles != emptyModel {
			err := r.DB.Model(&unpublished).Association("AvailableFiles").Replace(unpublished.AvailableFiles)
			if err != nil {
				return err
			}
		}

		if unpublished.Err == "none" {
			err := r.DB.Model(&unpublished).Where("id = ?", unpublished.ID).Update("Err", "").Error
			if err != nil {
				return err
			}
		}

		err := r.DB.Updates(unpublished).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	var updated entities.UnpublishedBeat
	if err := r.DB.Where("id = ?", unpublished.ID).Scopes(WithBasicPreloads()).First(&updated).Error; err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *repository) DeleteUnpublishedById(id uuid.UUID) error {
	unpublished := entities.UnpublishedBeat{}
	result := r.DB.Where("id = ?", id).Delete(unpublished)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("unpublished beat not found")
	}

	return nil
}

func WithBasicPreloads() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Preload("Tags").Preload("Genres").
			Preload("Moods").Preload("Timestamps").
			Preload("Keynote").Preload("AvailableFiles").
			Preload("Instruments")
	}
}

func (r *repository) CheckTagExists(name string) (uint, error) {
	tag := entities.Tag{}
	result := r.DB.Where("name = ?", name).First(&tag)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, result.Error
	}
	return tag.ID, nil
}

func (r *repository) CreateTag(tag *entities.Tag) (*entities.Tag, error) {
	tag.CreatedAt = time.Now().Unix()
	result := r.DB.Create(tag)
	if result.Error != nil {
		return nil, result.Error
	}
	return tag, nil
}

func (r *repository) ReadMoodByName(name string) (*entities.Mood, error) {
	mood := entities.Mood{}
	err := r.DB.Where("LOWER(name) = LOWER(?)", name).First(&mood).Error
	if err != nil {
		return nil, err
	}
	return &mood, nil
}
