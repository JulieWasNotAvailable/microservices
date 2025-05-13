package bmmetadata

import (
	"github.com/JulieWasNotAvailable/microservices/user/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/user/internal/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateMetadata(metadata *entities.Metadata) (*entities.Metadata, error)
	ReadMetadatas() (*[]presenters.Metadata, error)
	ReadMetadataById(id uuid.UUID) (*presenters.Metadata, error)
	ReadMetadataByUserId(id uuid.UUID) (*presenters.Metadata, error)
	UpdateMetadataByUserId(id uuid.UUID, metadata *presenters.Metadata) (*presenters.Metadata, error)
	UpdateMetadataById(id uuid.UUID, metadata *presenters.Metadata) (*presenters.Metadata, error)
	DeleteMetadataById(id uuid.UUID) error
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r *repository) CreateMetadata(metadata *entities.Metadata) (*entities.Metadata, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	metadata.ID = id

	result := r.DB.Create(&metadata)

	if result.Error != nil {
		return nil, result.Error
	}

	return metadata, nil
}

func (r *repository) ReadMetadatas() (*[]presenters.Metadata, error) {
	metadatasModel := &[]presenters.Metadata{}

	err := r.DB.Find(&metadatasModel).Error
	// err := db.Model(&User{}).Preload("CreditCard").Find(&users).Error

	if err != nil {
		return nil, err
	}

	return metadatasModel, nil
}

func (r *repository) ReadMetadataById(id uuid.UUID) (*presenters.Metadata, error) {
	metadata := &presenters.Metadata{}

	err := r.DB.First(&metadata, id).Error	
	
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func (r *repository) ReadMetadataByUserId(id uuid.UUID) (*presenters.Metadata, error) {
	metadata := &presenters.Metadata{}

	err := r.DB.Where("user_id = ?", id).First(&metadata).Error	
	
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func (r *repository) UpdateMetadataByUserId(id uuid.UUID, metadata *presenters.Metadata) (*presenters.Metadata, error) {
	err := r.DB.Where("user_id = ?", id).Updates(metadata).Error

	if err != nil {
		return nil, err
	}	

	return metadata, nil
}

func (r *repository) UpdateMetadataById(id uuid.UUID, metadata *presenters.Metadata) (*presenters.Metadata, error) {
	err := r.DB.Where("ID = ?", id).Updates(metadata).Error

	if err != nil {
		return nil, err
	}	

	return metadata, nil
}

func (r *repository) DeleteMetadataById(id uuid.UUID) error {
	metadata := &entities.Metadata{}

	err := r.DB.Delete(metadata, id).Error

	if err != nil {
		return err
	}

	return nil
}