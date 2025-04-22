package unpbeat

import (
	"github.com/JulieWasNotAvailable/microservices/unpublished/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/unpublished/pkg/entities"
	"github.com/google/uuid"
)

type Service interface {
	CreateUnpublishedBeat(beat entities.UnpublishedBeat) (entities.UnpublishedBeat, error)
	GetAllUnpublishedBeats() ([]presenters.UnpublishedBeat, error)
	GetUnpublishedBeatByID(id uuid.UUID) (*presenters.UnpublishedBeat, error)
	GetUnpublishedBeatsByUser(userID uuid.UUID) ([]presenters.UnpublishedBeat, error)
	GetUnpublishedInModeration(from int64, to int64) (*[]presenters.UnpublishedBeat, error)
	GetUnpublishedByBeatmakerandStatus(userId uuid.UUID, status string) (*[]presenters.UnpublishedBeat, error)
	UpdateUnpublishedBeat(beat *presenters.UnpublishedBeat) (*presenters.UnpublishedBeat, error)
	DeleteUnpublishedBeat(id uuid.UUID) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{repository: r}
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

func (s service) GetUnpublishedInModeration(from int64, to int64) (*[]presenters.UnpublishedBeat, error){
	beats, err := s.repository.ReadUnpublishedInModeration(from, to)
	return beats, err
}

func (s *service) GetUnpublishedByBeatmakerandStatus(userId uuid.UUID, status string) (*[]presenters.UnpublishedBeat, error) {
	beats, err := s.repository.ReadUnpublishedByBeatmakerandStatus(userId, status)
	return beats, err
}

func (s *service) UpdateUnpublishedBeat(beat *presenters.UnpublishedBeat) (*presenters.UnpublishedBeat, error) {
	return s.repository.UpdateUnpublishedById(beat)
}

func (s *service) DeleteUnpublishedBeat(id uuid.UUID) error {
	return s.repository.DeleteUnpublishedById(id)
}