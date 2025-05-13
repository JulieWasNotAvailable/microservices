package bmmetadata

import (
	"github.com/JulieWasNotAvailable/microservices/user/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/user/internal/entities"
	"github.com/google/uuid"
)

type Service interface {
	InsertMetadata(metadata *entities.Metadata) (*entities.Metadata, error)
	FetchMetadatas() (*[]presenters.Metadata, error)
	FetchMetadataById(id uuid.UUID) (*presenters.Metadata, error)
	FetchMetadataByUserId(id uuid.UUID) (*presenters.Metadata, error)
	RemoveMetadataById(id uuid.UUID) error
	UpdateMetadataByUserId(id uuid.UUID, metadata *presenters.Metadata) (*presenters.Metadata, error)
	UpdateMetadataById(id uuid.UUID, metadata *presenters.Metadata) (*presenters.Metadata, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) InsertMetadata(metadata *entities.Metadata) (*entities.Metadata, error) {
	return s.repository.CreateMetadata(metadata)
}

func (s *service) FetchMetadatas() (*[]presenters.Metadata, error) {
	return s.repository.ReadMetadatas()
}

func (s *service) FetchMetadataById(id uuid.UUID) (*presenters.Metadata, error) {
	return s.repository.ReadMetadataById(id)
}

func (s *service) FetchMetadataByUserId(id uuid.UUID) (*presenters.Metadata, error) {
	return s.repository.ReadMetadataByUserId(id)
}

func (s *service) UpdateMetadataByUserId(id uuid.UUID, metadata *presenters.Metadata) (*presenters.Metadata, error) {
	return s.repository.UpdateMetadataByUserId(id, metadata)
}

func (s *service) UpdateMetadataById(id uuid.UUID, metadata *presenters.Metadata) (*presenters.Metadata, error) {
	return s.repository.UpdateMetadataById(id, metadata)
}

func (s *service) RemoveMetadataById(id uuid.UUID) error {
	return s.repository.DeleteMetadataById(id)
}