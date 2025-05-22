package entities

import (
	"errors"

	_ "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// @Description User
type User struct {
	ID                      uuid.UUID                  `json:"id" example:"1"`
	Email                   string                     `json:"email,omitempty" validate:"required,email,min=5,max=50" example:"550e8400-e29b-41d4-a716-446655440000" unique:"true"`
	Password                string                     `json:"password,omitempty" validate:"required,min=6,max=20" example:"securepassword123"` // Never exposed in responses
	Firstname               string                     `json:"firstname,omitempty" example:"John"`
	Lastname                string                     `json:"lastname,omitempty" example:"Doe"`
	Patronymic              string                     `json:"patronymic,omitempty" example:"Smith"`
	Username                string                     `json:"username,omitempty" example:"johndoe"`
	ProfilePictureUrl       string                     `json:"profilepicture,omitempty" example:"https://storage.yandexcloud.net/imagesall/01961f2b-61b4-74ee-8e5b-26044ec630ea"`
	RoleID                  uint                       `json:"roleId" validate:"required" example:"1"`
	SubscriptionID          int                        `json:"subscriptionId,omitempty" example:"1"`
	FollowerOf              int                        `json:"followerOf,omitempty" example:"10"`
	Metadata                Metadata                   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	UserFollowsBeatmaker    []*User_Follows_Beatmakers `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	BeatmakerFollowedByUser []*User_Follows_Beatmakers `gorm:"foreignKey:BeatmakerID;constraint:OnDelete:CASCADE"`
}

// Validate - кастомная валидация
func (u *User) Validate() error {
	if u.RoleID != 1 && u.RoleID != 2 && u.RoleID != 3 {
		return errors.New("roleID must be 1, 2, or 3")
	}
	return nil
}

// BeforeSave - хуки GORM для автоматической валидации
func (u *User) BeforeSave(tx *gorm.DB) error {
	return u.Validate()
}
