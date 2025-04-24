package unpbeat

import (
	"errors"
	"time"

	"github.com/JulieWasNotAvailable/microservices/unpublished/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/unpublished/pkg/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUnpublished(unpublished entities.UnpublishedBeat) (entities.UnpublishedBeat, error)
	ReadUnpublished() (*[]presenters.UnpublishedBeat, error)
	ReadUnpublishedById(id uuid.UUID) (*presenters.UnpublishedBeat, error)
	ReadUnpublishedByUser(id uuid.UUID) (*[]presenters.UnpublishedBeat, error)
	ReadUnpublishedInModeration(from int64, to int64) (*[]presenters.UnpublishedBeat, error)
	ReadUnpublishedByBeatmakerandStatus(userId uuid.UUID, status string) (*[]presenters.UnpublishedBeat, error)
	UpdateUnpublishedById(unpublished *presenters.UnpublishedBeat) (*presenters.UnpublishedBeat, error)
	DeleteUnpublishedById(id uuid.UUID) error
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r repository) CreateUnpublished(unpublished entities.UnpublishedBeat) (entities.UnpublishedBeat, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return entities.UnpublishedBeat{}, err
	}

	unpublished.ID = id
	time := time.Now().Unix()
	unpublished.CreatedAt = time
	result := r.DB.Create(&unpublished)

	if result.Error != nil {
		return entities.UnpublishedBeat{}, result.Error
	}

	return unpublished, nil
}

func (r *repository) ReadUnpublished() (*[]presenters.UnpublishedBeat, error) {
	var unpublishedBeats []presenters.UnpublishedBeat

	result := r.DB.Model(unpublishedBeats).Preload("AvailableFiles").Preload("Tags").Preload("Genres").
	Preload("Moods").Preload("Timestamps").Preload("Instruments").Find(&unpublishedBeats)
	if result.Error != nil {
		return nil, result.Error
	}

	return &unpublishedBeats, nil
}

// ReadUnpublishedById reads a single unpublished beat by ID
func (r *repository) ReadUnpublishedById(id uuid.UUID) (*presenters.UnpublishedBeat, error) {
	var unpublished presenters.UnpublishedBeat
	
	result := r.DB.Where("id = ?", id).Preload("AvailableFiles").Preload("Tags").Preload("Genres").
	Preload("Moods").Preload("Timestamps").Preload("Instruments").First(&unpublished)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("unpublished beat not found")
		}
		return nil, result.Error
	}

	return &unpublished, nil
}

func (r *repository) ReadUnpublishedByUser(userId uuid.UUID) (*[]presenters.UnpublishedBeat, error) {
	var unpublishedBeats []presenters.UnpublishedBeat

	result := r.DB.Model(unpublishedBeats).
		Where("beatmaker_id = ?", userId).Preload("AvailableFiles").Preload("Tags").Preload("Genres").
		Preload("Moods").Preload("Timestamps").Preload("Instruments").
		Find(&unpublishedBeats)

	if result.Error != nil {
		return nil, result.Error
	}

	return &unpublishedBeats, nil
}

func (r *repository) ReadUnpublishedByBeatmakerandStatus(userId uuid.UUID, status string) (*[]presenters.UnpublishedBeat, error) {
	var unpublishedBeats []presenters.UnpublishedBeat

	result := r.DB.Model(unpublishedBeats).
		Where("beatmaker_id = ?", userId).Where("status = ?", status).Preload("Tags").Preload("Genres").
		Preload("Moods").Preload("Timestamps").Preload("Instruments").
		Find(&unpublishedBeats)

	if result.Error != nil {
		return nil, result.Error
	}

	return &unpublishedBeats, nil
}

func (r *repository) ReadUnpublishedInModeration(from int64, to int64) (*[]presenters.UnpublishedBeat, error){
	var unpublishedBeats []presenters.UnpublishedBeat

	result := r.DB.Model(unpublishedBeats).
	Where("status = ?", "in_moderation").Where("sent_to_moderation_at >= ? AND sent_to_moderation_at <= ?", from, to).
	Preload("Tags").Preload("AvailableFiles").Preload("Genres").
	Preload("Moods").Preload("Timestamps").Preload("Instruments").
	Find(&unpublishedBeats)

	if result.Error != nil {
		return nil, result.Error
	}

	return &unpublishedBeats, nil
}

func (r *repository) UpdateUnpublishedById(unpublished *presenters.UnpublishedBeat) (*presenters.UnpublishedBeat, error) {
	//Where("id = ?", unpublished.ID)
	// unpublishedModel := entities.UnpublishedBeat{}
	result := r.DB.Updates(unpublished)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("unpublished beat not found or not owned by user")
		}
		return nil, result.Error
	}

	var updated presenters.UnpublishedBeat
	if err := r.DB.Where("id = ?", unpublished.ID).First(&updated).Preload("AvailableFiles").Preload("Tags").Preload("Genres").
	Preload("Moods").Preload("Timestamps").Preload("Instruments").Error; err != nil {
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
