package entities

import (
	_ "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Beat struct {
	ID              uuid.UUID       	`json:"id" example:"019628ef-cd76-7d2d-bf80-48b8011fad40"`
	Name            string    			`json:"name" validate:"required,min=2,max=60" example:"Summer Vibes"`
	Picture         string    			`json:"picture" example:"https://storage.yandexcloud.net/imagesall/019623bd-3d0b-7dc2-8a1f-f782adeb42b4"`
	BeatmakerID     uuid.UUID       	`json:"beatmakerId" validate:"required" example:"019628ef-cd76-7d2d-bf80-48b8011fad40"`
	AvailableFiles 	AvailableFiles      `gorm:"foreignKey:BeatID;" validate:"required"`
	URL             string    			`json:"url" validate:"required" example:"https://storage.yandexcloud.net/mp3beats/019623bd-3d0b-7dc2-8a1f-f782adeb42b4"`
	Price           int       			`json:"price" validate:"required, gte=1" example:"2999"`
	Tags            []Tag      			`json:"tags" validate:"required" gorm:"many2many:beat_tags;"` //many to many
	BPM             int       			`json:"bpm" validate:"required,gte=20,lte=400" example:"120"`
	Description     string    			`json:"description" validate:"min=2,max=500" example:"Chill summer beat with tropical influences"`
	Genres         	[]Genre      		`json:"genres" validate:"required" gorm:"many2many:beat_genres;joinForeignKey:BeatID;joinReferences:GenreID;constraint:OnDelete:CASCADE"`   
	Moods          	[]Mood      		`json:"moods" validate:"required" gorm:"many2many:beat_moods;constraint:OnDelete:CASCADE"`       //many to many
	KeynoteID       uint       		`json:"keynoteId" validate:"required" example:"11"`    //keynote has many beats, but each beat has only one keynote`
	Timestamps    	[]Timestamp         `json:"timestamps" validate:"required" gorm:"foreignKey:BeatID"` //a beat has many timestamps, but each timestamp has only one beat
	Instruments   	[]Instrument        `json:"instruments" gorm:"many2many:beat_instruments;constraint:OnDelete:CASCADE"` //many to many
	MFCC 			MFCC 				`gorm:"foreignKey:BeatId;constraint:OnDelete:CASCADE;" validate:"required" `
	CreatedAt      int64 				`json:"created_at"`
}

type AvailableFiles struct{
	ID uuid.UUID
	MP3Url string
	WAVUrl string
	ZIPUrl string
	BeatID uuid.UUID 
}

func MigrateAll(db *gorm.DB) error {
	err := db.AutoMigrate(
		&Instrument{},
		&BeatInstrument{},
		&MFCC{},
		&Keynote{},
		&Genre{},
		&BeatGenre{},
		&Mood{},
		&BeatMood{},
		&Tag{},
		&BeatTag{},
		&Beat{},
		&Timestamp{},
		&AvailableFiles{},
	)
	return err
}