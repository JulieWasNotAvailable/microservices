package entities

import (
	_ "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Genre struct {
	ID         uint    `json:"id" swaggerignore:"true"`
	Beat       []*Beat `gorm:"many2many:beat_genres;joinForeignKey:GenreID;joinReferences:BeatID"`
	Name       string  `json:"name" example:"Jerk"`
	PictureUrl string  `json:"picture_url"`
	CreatedAt  int64   `json:"createdAt"`
}

type BeatGenre struct {
	BeatID  uuid.UUID `json:"beatId" gorm:"primaryKey;constraint:OnDelete:CASCADE"` // Delete join entry if Track is deleted
	GenreID uint      `json:"genreId" gorm:"primaryKey"`
}

type Timestamp struct {
	ID        uint      `json:"id" swaggerignore:"true"`
	BeatID    uuid.UUID `json:"beatId" example:"01963e01-e46c-7996-996a-42ad3df115ac"`
	Name      string    `json:"title"`
	TimeStart int64     `json:"start_time" validate:"required,gte=1,lte=299"`
	TimeEnd   int64     `json:"end_time" validate:"required,gte=2,lte=300"`
}

type Tag struct {
	ID        uint    `json:"id"`
	Beat      []*Beat `gorm:"many2many:beat_tags;" swaggerignore:"true"`
	Name      string  `json:"name"`
	CreatedAt int64   `json:"createdAt"`
}

type BeatTag struct {
	BeatID uuid.UUID `json:"beatId" gorm:"primaryKey;constraint:OnDelete:CASCADE"` // Delete join entry if Track is deleted
	TagID  uint      `json:"tagId" gorm:"primaryKey"`
}

type Mood struct {
	ID   uint    `json:"id"`
	Beat []*Beat `json:"beat" gorm:"many2many:beat_moods;" swaggerignore:"true"`
	Name string  `json:"name"`
}

type BeatMood struct {
	BeatID uuid.UUID `json:"beatId" gorm:"primaryKey;constraint:OnDelete:CASCADE"` // Delete join entry if Track is deleted
	MoodID uint      `json:"moodId" gorm:"primaryKey"`
}

type Keynote struct {
	ID    uint   `json:"id"`
	Beats []Beat `json:"beats" gorm:"foreignKey:KeynoteID" swaggerignore:"true"`
	Name  string `json:"name"`
}

type Instrument struct {
	ID   uint    `json:"id"`
	Beat []*Beat `json:"beat" gorm:"many2many:beat_instruments;" swaggerignore:"true"`
	Name string  `json:"name"`
}

type BeatInstrument struct {
	BeatID       uuid.UUID `json:"beatId" gorm:"primaryKey;constraint:OnDelete:CASCADE"`
	InstrumentID uint      `json:"instrumentId" gorm:"primaryKey"`
}
