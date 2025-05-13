package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Metadata struct {
	ID					uuid.UUID
	VkUrl           	string `json:"vkUrl,omitempty" example:"vk.com/i_love_bunnies"`
	TelegramUrl      	string `json:"telegramUrl,omitempty" example:"tg.com/i_love_bunnies"`
	InstagramUrl     	string `json:"instagramUrl,omitempty" example:"insta.com/i_love_bunnies"`
	Description      	string `json:"description,omitempty" example:"the best beatmaker ever"`
	SubscriptionTypeID 	int    `json:"subscriptionTypeId,omitempty" example:"3"`
	UserID 				uuid.UUID `json:"userId" example:"550e8400-e29b-41d4-a716-446655440000"`
}

func MigrateMetadata(db *gorm.DB) error {
	err := db.AutoMigrate(&Metadata{})
	return err
}