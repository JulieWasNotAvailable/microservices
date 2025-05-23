package activity

import (
	"github.com/JulieWasNotAvailable/microservices/user/internal/entities"
	"github.com/google/uuid"
)

type Service interface {
	InsertSub(userId uuid.UUID, beatmakerId uuid.UUID) (entities.User_Follows_Beatmakers, error)
	FetchSubsCountByBeatmakerId(beatmakerId uuid.UUID) (int, error)
	FetchSubsByUserId(userId uuid.UUID) ([]entities.User_Follows_Beatmakers, error)
	RemoveSub(userId uuid.UUID, beatmakerId uuid.UUID) (entities.User_Follows_Beatmakers, error)

	CheckBeatmakerRole(beatmakerId uuid.UUID) (bool, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) FetchSubsByUserId(userId uuid.UUID) ([]entities.User_Follows_Beatmakers, error) {
	return s.repository.ReadSubsByUserId(userId)
}

func (s *service) FetchSubsCountByBeatmakerId(beatmakerId uuid.UUID) (int, error) {
	_, err := s.CheckBeatmakerRole(beatmakerId)
	if err != nil {
		return 0, err
	}
	return s.repository.ReadSubsCountByBeatmakerId(beatmakerId)
}

func (s *service) InsertSub(userId uuid.UUID, beatmakerId uuid.UUID) (entities.User_Follows_Beatmakers, error) {
	return s.repository.CreateSub(userId, beatmakerId)
}

func (s *service) RemoveSub(userId uuid.UUID, beatmakerId uuid.UUID) (entities.User_Follows_Beatmakers, error) {
	return s.repository.DeleteSub(userId, beatmakerId)
}

func (s *service) CheckBeatmakerRole(beatmakerId uuid.UUID) (bool, error) {
	return s.repository.CheckBeatmakerRole(beatmakerId)
}
