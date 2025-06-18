package entities

import (
	_ "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Beat struct {
	ID             uuid.UUID      `json:"id" example:"019628ef-cd76-7d2d-bf80-48b8011fad40"`
	Name           string         `json:"name" validate:"required,min=2,max=60" example:"Summer Vibes"`
	Picture        string         `json:"picture" example:"https://storage.yandexcloud.net/imagesall/019623bd-3d0b-7dc2-8a1f-f782adeb42b4"`
	BeatmakerID    uuid.UUID      `json:"beatmakerId" validate:"required" example:"019628ef-cd76-7d2d-bf80-48b8011fad40"`
	BeatmakerName  string         `json:"beatmakerName"`
	AvailableFiles AvailableFiles `json:"availableFiles" gorm:"foreignKey:BeatID;constraint:OnDelete:CASCADE" validate:"required"`
	URL            string         `json:"url" validate:"required" example:"https://storage.yandexcloud.net/mp3beats/019623bd-3d0b-7dc2-8a1f-f782adeb42b4"`
	Price          int            `json:"price" validate:"required, gte=1" example:"2999"`
	Tags           []Tag          `json:"tags" validate:"required" gorm:"many2many:beat_tags;constraint:OnDelete:CASCADE"` //many to many
	BPM            int            `json:"bpm" validate:"required,gte=20,lte=400" example:"120"`
	Description    string         `json:"description" validate:"min=2,max=5000" example:"Chill summer beat with tropical influences" gorm:"type:text;size:5000"`
	Genres         []Genre        `json:"genres" validate:"required" gorm:"many2many:beat_genres;joinForeignKey:BeatID;joinReferences:GenreID;constraint:OnDelete:CASCADE"`
	Moods          []Mood         `json:"moods" validate:"required" gorm:"many2many:beat_moods;constraint:OnDelete:CASCADE"`   //many to many
	KeynoteID      uint           `json:"keynoteId" validate:"required" example:"11"`                                          //keynote has many beats, but each beat has only one keynote`
	Keynote Keynote 
	Timestamps     []Timestamp    `json:"timestamps" validate:"required" gorm:"foreignKey:BeatID;constraint:OnDelete:CASCADE"` //a beat has many timestamps, but each timestamp has only one beat
	Instruments    []Instrument   `json:"instruments" gorm:"many2many:beat_instruments;constraint:OnDelete:CASCADE"`           //many to many
	MFCC           MFCC           `json:"-" gorm:"foreignKey:BeatID;constraint:OnDelete:CASCADE;" validate:"required" `
	Likes          []*Like        `json:"likes" gorm:"foreignKey:BeatID;constraint:OnDelete:CASCADE;"`
	Listen         []*Listen      `json:"-" gorm:"foreignKey:BeatID;constraint:OnDelete:CASCADE;"`
	Plays          int64          `json:"plays" example:"105"`
	CreatedAt      int64          `json:"created_at"`
}

type AvailableFiles struct {
	ID     uuid.UUID `json:"id"`
	MP3Url string    `json:"mp3url"`
	WAVUrl string    `json:"wavurl"`
	ZIPUrl string    `json:"zipurl"`
	BeatID uuid.UUID `json:"beatId"`
}

func MigrateAll(db *gorm.DB) error {
	err := db.AutoMigrate(
		&Instrument{},
		&BeatInstrument{},
		&Like{},
		&Listen{},
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
