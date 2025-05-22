package entities

import "github.com/google/uuid"

type Like struct {
	UserID uuid.UUID `json:"userID" gorm:"primaryKey"`
	BeatID uuid.UUID `json:"beatID" gorm:"primaryKey;constraint:OnDelete:CASCADE"`
}

type Listen struct {
	UserID uuid.UUID `json:"userID" gorm:"primaryKey"`
	BeatID uuid.UUID `json:"beatId" gorm:"primaryKey;constraint:OnDelete:CASCADE"`
	// Beat   Beat      `json:"beat" gorm:"foreignKey:BeatID"`
}
