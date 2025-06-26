package activity

import (
	"errors"

	"github.com/JulieWasNotAvailable/microservices/beat/internal/entities"
	"github.com/google/uuid"
)

type Service interface {
	InsertLike(userid uuid.UUID, beatid uuid.UUID) (*entities.Like, error)
	DeleteLike(userid uuid.UUID, beatid uuid.UUID) (*entities.Like, error)
	GetLikesByUserId(userid uuid.UUID) (*[]entities.Like, error)
	GetLikesCountByBeatId(beatid uuid.UUID) (int, error) //likes number of this specific beat
	GetUserLikesCount(userid uuid.UUID) (int, error) //likes number for a specific user (how many did he like)

	GetLikesCountOfBeats(beatids []uuid.UUID) (int, error)
	InsertListened(userId uuid.UUID, beatId uuid.UUID) (entities.Listen, error)
	BeatExists(beatId uuid.UUID) (bool, error) 
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{repository: r}
}

// InsertLike implements Service.
func (s *service) InsertLike(userid uuid.UUID, beatid uuid.UUID) (*entities.Like, error) {
	return s.repository.CreateLike(userid, beatid)
}

// DeleteLike implements Service.
func (s *service) DeleteLike(userid uuid.UUID, beatid uuid.UUID) (*entities.Like, error) {
	return s.repository.DeleteLike(userid, beatid)
}

// GetLikesByUserId implements Service.
func (s *service) GetLikesByUserId(userid uuid.UUID) (*[]entities.Like, error) {
	return s.repository.ReadLikesByUserId(userid)
}

// GetLikesCountByBeatId implements Service.
func (s *service) GetLikesCountByBeatId(beatid uuid.UUID) (int, error) {
	return s.repository.ReadLikesCountByBeatId(beatid)
}

// GetLikesCountByUserId implements Service.
func (s *service) GetUserLikesCount(userid uuid.UUID) (int, error) {
	return s.repository.ReadLikesCountByUserId(userid)
}

func (s *service) GetLikesCountOfBeats(beatids []uuid.UUID) (int, error) {
	return s.repository.ReadLikesCountOfBeats(beatids)
}

func (s *service) InsertListened(userId uuid.UUID, beatId uuid.UUID) (entities.Listen, error) {
	exists, err := s.BeatExists(beatId)
	if err != nil {
		return entities.Listen{}, err
	}
	if !exists {
		return entities.Listen{}, errors.New("beat does not exist")
	}
	_, err = s.repository.ReadListenedByUserAndBeatId(userId, beatId)
	if err != nil {
		return entities.Listen{}, err
	}
	return s.repository.CreateListened(userId, beatId)
}

func (s *service)BeatExists(beatId uuid.UUID) (bool, error) {
	return s.repository.BeatExists(beatId)
}