package user

import (
	"github.com/JulieWasNotAvailable/microservices/user/api/presenters"
	"github.com/JulieWasNotAvailable/microservices/user/pkg/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(user *entities.User) (*entities.User, error)
	ReadUsers() (*[]presenters.User, error)
	ReadUserById(id uuid.UUID) (*presenters.User, error)
	ReadUserByEmail(email string) (*presenters.User, error)
	UpdateUser(user *presenters.User) (*presenters.User, error)
	UpdateBeatmaker(userID uuid.UUID, userData *presenters.User, metadata *presenters.Metadata) (*presenters.User, error)
	DeleteUser(id uuid.UUID) error
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

type GetByEmailRequest struct{
	Email *string
}

// CreateUser implements Repository.
func (r *repository) CreateUser(user *entities.User) (*entities.User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	user.ID = id

	err = r.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repository) ReadUsers () (*[]presenters.User, error) {
	userModels := &[]presenters.User{}

	err := r.DB.Model(userModels).Preload("Metadata").Find(&userModels).Error
	if err != nil {
		return nil, err
	}

	return userModels, nil
}

func (r *repository) ReadUserById (id uuid.UUID) (*presenters.User, error) {
	user := &presenters.User{}
	err := r.DB.Where("ID = ?", id).Preload("Metadata").First(user).Error

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repository) ReadUserByEmail (email string) (*presenters.User, error) {
	user := &presenters.User{}

	err := r.DB.Where("email = ?", email).Preload("Metadata").First(user).Error

	if err != nil{
		return nil, err
	}

	return user, nil
}

func (r *repository) UpdateUser(user *presenters.User) (*presenters.User, error) {
	err := r.DB.Updates(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repository) UpdateBeatmaker (userID uuid.UUID, userData *presenters.User, metadata *presenters.Metadata) (*presenters.User, error) {
	return userData, r.DB.Transaction(func(tx *gorm.DB) error {
        // First update the user
        err := r.DB.Where("id = ?", userID).Updates(userData).Error
		if err != nil {
            return err
        }

		err = r.DB.Where("user_id = ?", userID).Updates(metadata).Error
		if err != nil {
			return err
		}

        return nil
    })
	// return user, nil
}

func (r *repository) DeleteUser (id uuid.UUID) error {
	userModel := &entities.User{}

	err := r.DB.Delete(userModel, id).Error 

	if err != nil {
		return err
	}

	return nil
}