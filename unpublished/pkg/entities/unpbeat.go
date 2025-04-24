package entities

import (
	"errors"
	_ "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Status represents the publication status of a beat
type Status string

const (
	StatusInModeration Status = "processing"
	StatusDraft       Status = "draft"
)

//@Description entitites.UnpublishedBeatErrorResponse
type UnpublishedBeat struct {
	ID              uuid.UUID       	`json:"id" example:"019628ef-cd76-7d2d-bf80-48b8011fad40"`
	Name            string    			`json:"name" validate:"required,min=2,max=60" example:"Summer Vibes"`
	Picture         string    			`json:"picture" example:"https://storage.yandexcloud.net/imagesall/019623bd-3d0b-7dc2-8a1f-f782adeb42b4"`
	BeatmakerID     uuid.UUID       	`json:"beatmakerId" validate:"required" example:"019628ef-cd76-7d2d-bf80-48b8011fad40"`
	AvailableFiles 	AvailableFiles      `gorm:"foreignKey:UnpublishedBeatID;constraint:OnDelete:CASCADE;" validate:"required" `
	URL             string    			`json:"url" validate:"required" example:"https://storage.yandexcloud.net/mp3beats/019623bd-3d0b-7dc2-8a1f-f782adeb42b4"`
	Price           int       			`json:"price" validate:"required" example:"2999"`
	Tags            []Tag      			`json:"tags" validate:"required" gorm:"many2many:tag_beats;"` //many to many
	BPM             int       			`json:"bpm" validate:"required,gte=20,lte=400" example:"120"`
	Description     string    			`json:"description" validate:"min=2,max=500" example:"Chill summer beat with tropical influences"`
	Genres         	[]Genre      		`json:"genres" validate:"required" gorm:"many2many:genre_beats;"`       //many to many
	Moods          	[]Mood      		`json:"moods" validate:"required" gorm:"many2many:mood_beats;"`       //many to many
	KeynoteID       *uint       			`json:"keynoteId" validate:"required" example:"11"`    //keynote has many beats, but each beat has only one keynote`
	Timestamps    	[]Timestamp         `json:"timestamps" validate:"required" gorm:"foreignKey:BeatID"` //a beat has many timestamps, but each timestamp has only one beat
	Instruments   	[]Instrument        `json:"instruments" gorm:"many2many:instrument_beats"` //many to many
	Status          Status    			`json:"status" example:"draft"`
	Err 			string	
	SentToModerationAt int64 			`json:"sent_to_moderation_at"`
	CreatedAt      int64 				`json:"created_at"`
	UpdatedAt      int64 				`json:"updated_at"`
}

// Validate - кастомная валидация
func (beat *UnpublishedBeat) Validate() error {
	if beat.Status != "draft" {
		return errors.New("beat status can only be draft or in moderation")
	}
	return nil
}

// BeforeSave - хуки GORM для автоматической валидации
func (beat *UnpublishedBeat) BeforeSave(tx *gorm.DB) error {
	return beat.Validate()
}

type AvailableFiles struct{
	ID uuid.UUID
	MP3Url string
	WAVUrl string
	ZIPUrl string
	UnpublishedBeatID uuid.UUID 
}

//@Description entities.Genre
type Genre struct{
	ID uint `json:"id" swaggerignore:"true"`
	UnpublishedBeat []*UnpublishedBeat `gorm:"many2many:genre_beats;" swaggerignore:"true"`
	Name string `json:"name" example:"Jerk"`
}

type Timestamp struct{
	ID uint  `json:"id" swaggerignore:"true"`
	BeatID uuid.UUID `json:"unpublishedbeatId" example:"01963e01-e46c-7996-996a-42ad3df115ac"`
	Name string
	TimeStart int64 `validate:"required,gte=1,lte=299"`
	TimeEnd int64 `validate:"required,gte=2,lte=300"`
}

type Tag struct{
	ID uint
	UnpublishedBeat []*UnpublishedBeat `gorm:"many2many:tag_beats;" swaggerignore:"true"`
	Name string 
}

type Mood struct{
	ID uint
	UnpublishedBeat []*UnpublishedBeat `gorm:"many2many:mood_beats;" swaggerignore:"true"`
	Name string
}

type Keynote struct{
	ID uint
	Beats []UnpublishedBeat `json:"unpublishedbeats" gorm:"foreignKey:KeynoteID" swaggerignore:"true"`
	Name string
}

type Instrument struct{
	ID uint
	UnpublishedBeat []*UnpublishedBeat `gorm:"many2many:instrument_beats;" swaggerignore:"true"`
	Name string
}

func MigrateAll(db *gorm.DB) error {
	err := db.AutoMigrate(
		&Instrument{},
		&Keynote{},
		&AvailableFiles{},
		&Genre{},
		&Timestamp{},
		&Mood{},
		&Tag{},
		&UnpublishedBeat{},
	)
	return err
}