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

	result := r.DB.Scopes(WithBasicPreloads()).Find(&unpublishedBeats)
	if result.Error != nil {
		return nil, result.Error
	}

	return &unpublishedBeats, nil
}

// ReadUnpublishedById reads a single unpublished beat by ID
func (r *repository) ReadUnpublishedById(id uuid.UUID) (*presenters.UnpublishedBeat, error) {
	var unpublished presenters.UnpublishedBeat
	
	result := r.DB.Where("id = ?", id).Scopes(WithBasicPreloads()).First(&unpublished)
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
		Where("beatmaker_id = ?", userId).Scopes(WithBasicPreloads()).
		Find(&unpublishedBeats)

	if result.Error != nil {
		return nil, result.Error
	}

	return &unpublishedBeats, nil
}

func (r *repository) ReadUnpublishedByBeatmakerandStatus(userId uuid.UUID, status string) (*[]presenters.UnpublishedBeat, error) {
	var unpublishedBeats []presenters.UnpublishedBeat

	result := r.DB.Model(unpublishedBeats).
		Where("beatmaker_id = ?", userId).Where("status = ?", status).Scopes(WithBasicPreloads()).
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
	Scopes(WithBasicPreloads()).
	Find(&unpublishedBeats)

	if result.Error != nil {
		return nil, result.Error
	}

	return &unpublishedBeats, nil
}

func (r *repository) UpdateUnpublishedById(unpublished *presenters.UnpublishedBeat) (*presenters.UnpublishedBeat, error) {
	emptyModel := entities.AvailableFiles{}

	beatInitial := &entities.UnpublishedBeat{}
	err := r.DB.Where("id = ?", unpublished.ID).First(&beatInitial).Error
		if err != nil{
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("unpublished beat not found or not owned by user")
		} else {
			return nil, err
		}
	}
	
	//it needs to be a transaction
	//i need to check that person tried to edit tags, and then delete the tags from beat_tags table
	err = r.DB.Transaction(func(tx *gorm.DB) error {
		
		if unpublished.Tags != nil{
			err := r.DB.Model(&unpublished).Association("Tags").Replace(unpublished.Tags)
			if err != nil{
				return err
			}
		}

		if unpublished.Genres != nil{
			err := r.DB.Model(&unpublished).Association("Genres").Replace(unpublished.Genres)
			if err != nil{
				return err
			}
		}

		if unpublished.Instruments != nil{
			err := r.DB.Model(&unpublished).Association("Instruments").Replace(unpublished.Instruments)
			if err != nil{
				return err
			}
		}

		if unpublished.Timestamps != nil{
			err := r.DB.Model(&unpublished).Association("Timestamps").Replace(unpublished.Timestamps)
			if err != nil{
				return err
			}
		}

		if unpublished.AvailableFiles != emptyModel{
			err := r.DB.Model(&beatInitial).Association("AvailableFiles").Replace(unpublished.AvailableFiles)
			if err != nil{
				return err
			}
		}
		
		err := r.DB.Updates(unpublished).Error	
		if err != nil {
			return err
		}

		return nil	
	})
	if err != nil{
		return nil, err
	}

	var updated presenters.UnpublishedBeat
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
