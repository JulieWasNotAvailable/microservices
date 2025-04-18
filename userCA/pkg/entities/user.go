package entities

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// @Description User
type User struct {
	ID				uuid.UUID `json:"id" example:"1"`
 	Email           string `json:"email,omitempty" example:"550e8400-e29b-41d4-a716-446655440000" unique:"true"`
	Password        string `json:"password,omitempty" example:"securepassword123"` // Never exposed in responses
	Firstname       string `json:"firstname,omitempty" example:"John"`
	Lastname        string `json:"lastname,omitempty" example:"Doe"`
	Patronymic      string `json:"patronymic,omitempty" example:"Smith"`
	Username        string `json:"username,omitempty" example:"johndoe"`
	RoleID          uint    `json:"roleId" example:"1"`
	SubscriptionID  int    `json:"subscriptionId,omitempty" example:"1"`
	UsersFavourites int    `json:"usersFavourites,omitempty" example:"5"`
	FollowerOf      int    `json:"followerOf,omitempty" example:"10"`
	Metadata  	 	Metadata   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

func MigrateUser(db *gorm.DB) error {
	err := db.AutoMigrate(&User{})
	return err
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