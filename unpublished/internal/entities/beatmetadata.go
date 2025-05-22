package entities

import "github.com/google/uuid"

type AvailableFiles struct {
	ID                uuid.UUID `json:"id"`
	MP3Url            string    `json:"mp3url"`
	WAVUrl            string    `json:"wavurl"`
	ZIPUrl            string    `json:"zipurl"`
	UnpublishedBeatID uuid.UUID `json:"unpublishedBeatId"`
}

// @Description entities.Genre
type Genre struct {
	ID              uint               `json:"id"`
	UnpublishedBeat []*UnpublishedBeat `json:"-" gorm:"many2many:beat_genres;joinForeignKey:GenreID;joinReferences:UnpublishedBeatID" swaggerignore:"true"`
	Name            string             `json:"name" example:"Jerk"`
	CreatedAt       int64              `json:"createdAt"`
}

type BeatGenre struct {
	UnpublishedBeatID uuid.UUID `json:"unpublishedBeatId" gorm:"primaryKey;constraint:OnDelete:CASCADE"` // Delete join entry if Track is deleted
	GenreID           uint      `json:"genreId" gorm:"primaryKey"`
}

type Timestamp struct {
	ID        uint      `json:"id" swaggerignore:"true"`
	BeatID    uuid.UUID `json:"unpublishedbeatId" example:"01963e01-e46c-7996-996a-42ad3df115ac"`
	Name      string    `json:"name"`
	TimeStart int64     `json:"timeStart" validate:"required,gte=1,lte=299"`
	TimeEnd   int64     `json:"timeEnd" validate:"required,gte=2,lte=300"`
}

type Tag struct {
	ID              uint               `json:"id"`
	UnpublishedBeat []*UnpublishedBeat `json:"unpublishedBeat" gorm:"many2many:beat_tags;" swaggerignore:"true"`
	Name            string             `gorm:"unique;not null"`
	CreatedAt       int64              `json:"-"`
}

type BeatTag struct {
	UnpublishedBeatID uuid.UUID `json:"unpublishedBeatId" gorm:"primaryKey;constraint:OnDelete:CASCADE"` // Delete join entry if Track is deleted
	TagID             uint      `json:"tagId" gorm:"primaryKey"`
}

type Mood struct {
	ID              uint               `json:"id"`
	UnpublishedBeat []*UnpublishedBeat `json:"unpublishedBeat" gorm:"many2many:beat_moods;" swaggerignore:"true"`
	Name            string             `json:"name"`
}

type BeatMood struct {
	UnpublishedBeatID uuid.UUID `json:"unpublishedBeatId" gorm:"primaryKey;constraint:OnDelete:CASCADE"` // Delete join entry if Track is deleted
	MoodID            uint      `json:"moodId" gorm:"primaryKey"`
}

type Keynote struct {
	ID    uint              `json:"id"`
	Beats []UnpublishedBeat `json:"unpublishedbeats" gorm:"foreignKey:KeynoteID" swaggerignore:"true"`
	Name  string            `json:"name"`
}

type Instrument struct {
	ID              uint               `json:"id"`
	UnpublishedBeat []*UnpublishedBeat `json:"unpublishedBeat" gorm:"many2many:beat_instruments;" swaggerignore:"true"`
	Name            string             `json:"name"`
}

type BeatInstrument struct {
	UnpublishedBeatID uuid.UUID `json:"unpublishedBeatId" gorm:"primaryKey;constraint:OnDelete:CASCADE"`
	InstrumentID      uint      `json:"instrumentId" gorm:"primaryKey"`
}
