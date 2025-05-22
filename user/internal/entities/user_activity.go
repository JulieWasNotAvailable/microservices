package entities

import "github.com/google/uuid"

type User_Follows_Beatmakers struct {
	UserID      uuid.UUID `json:"userId" gorm:"primaryKey;constraint:OnDelete:CASCADE"`
	User        User      `json:"-" gorm:"foreignKey:BeatmakerID"`
	BeatmakerID uuid.UUID `json:"-" gorm:"primaryKey"`
	Beatmaker   User      `json:"beatmaker" gorm:"foreignKey:BeatmakerID"`
}
