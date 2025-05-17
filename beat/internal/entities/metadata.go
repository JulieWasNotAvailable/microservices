package entities

import (
	_ "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Genre struct{
	ID uint `json:"id" swaggerignore:"true"`
	Beat []*Beat `gorm:"many2many:beat_genres;joinForeignKey:GenreID;joinReferences:BeatID"`
	Name string `json:"name" example:"Jerk"`
	CreatedAt int64
}

type BeatGenre struct {
    BeatID uuid.UUID `gorm:"primaryKey;constraint:OnDelete:CASCADE"` // Delete join entry if Track is deleted
    GenreID uint `gorm:"primaryKey"`
}

type Timestamp struct{
	ID uint  `json:"id" swaggerignore:"true"`
	BeatID uuid.UUID `json:"beatId" example:"01963e01-e46c-7996-996a-42ad3df115ac"`
	Name string
	TimeStart int64 `validate:"required,gte=1,lte=299"`
	TimeEnd int64 `validate:"required,gte=2,lte=300"`
}

type Tag struct{
	ID uint
	Beat []*Beat `gorm:"many2many:beat_tags;" swaggerignore:"true"`
	Name string 
	CreatedAt int64
}

type BeatTag struct {
    BeatID uuid.UUID `gorm:"primaryKey;constraint:OnDelete:CASCADE"` // Delete join entry if Track is deleted
    TagID uint `gorm:"primaryKey"`
}

type Mood struct{
	ID uint
	Beat []*Beat `gorm:"many2many:beat_moods;" swaggerignore:"true"`
	Name string
}

type BeatMood struct {
    BeatID uuid.UUID `gorm:"primaryKey;constraint:OnDelete:CASCADE"` // Delete join entry if Track is deleted
    MoodID uint `gorm:"primaryKey"`
}

type Keynote struct{
	ID uint
	Beats []Beat `json:"beats" gorm:"foreignKey:KeynoteID" swaggerignore:"true"`
	Name string
}

type Instrument struct{
	ID uint
	Beat []*Beat `gorm:"many2many:beat_instruments;" swaggerignore:"true"`
	Name string
}

type BeatInstrument struct {
    BeatID uuid.UUID `gorm:"primaryKey;constraint:OnDelete:CASCADE"`
    InstrumentID uint `gorm:"primaryKey"`
}
