package user

import (
	"github.com/JulieWasNotAvailable/microservices/user/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/user/pkg/entities"
	"github.com/google/uuid"
)

type Service interface {
	InsertUser(user *entities.User) (*entities.User, error)
	FetchUsers() (*[]presenters.User, error)
	FetchUserById(id uuid.UUID) (*presenters.User, error)
	FetchUserByEmail(email string) (*presenters.User, error)
	UpdateUser(user *presenters.User) (*presenters.User, error)
	UpdateBeatmaker(userID uuid.UUID, userData *presenters.User, metadata *presenters.Metadata) (*presenters.User, error)
	RemoveUser(id uuid.UUID) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

// InsertUser implements Service.
func (s *service) InsertUser(user *entities.User) (*entities.User, error) {
	return s.repository.CreateUser(user)
}

// FetchUsers implements Service.
func (s *service) FetchUsers() (*[]presenters.User, error) {
	return s.repository.ReadUsers()
}

// FetchUserById implements Service.
func (s *service) FetchUserById(id uuid.UUID) (*presenters.User, error) {
	return s.repository.ReadUserById(id)
}

// FetchUserByEmail implements Service.
func (s *service) FetchUserByEmail(email string) (*presenters.User, error) {
	return s.repository.ReadUserByEmail(email)
}

func (s *service) UpdateUser(user *presenters.User) (*presenters.User, error){
	return s.repository.UpdateUser(user)
}

//UpdateUser implements Service.
func (s *service) UpdateBeatmaker(userID uuid.UUID, userData *presenters.User, metadata *presenters.Metadata) (*presenters.User, error) {
	return s.repository.UpdateBeatmaker(userID, userData, metadata)
}

// RemoveUser implements Service.
func (s *service) RemoveUser(id uuid.UUID) error {
	return s.repository.DeleteUser(id)
}

// NewService is used to create a single instance of the service
