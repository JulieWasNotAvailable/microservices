package entities

import (
	_ "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Status represents the publication status of a beat
type Status string

const (
	StatusInModeration Status = "processing"
	StatusDraft        Status = "draft"
)

// @Description entitites.UnpublishedBeatErrorResponse
type UnpublishedBeat struct {
	ID                 uuid.UUID      `json:"id" example:"019628ef-cd76-7d2d-bf80-48b8011fad40"`
	Name               string         `json:"name" validate:"required,min=2,max=100" example:"Summer Vibes"`
	Picture            string         `json:"picture" example:"https://storage.yandexcloud.net/imagesall/019623bd-3d0b-7dc2-8a1f-f782adeb42b4"`
	BeatmakerID        uuid.UUID      `json:"beatmaker_id" validate:"required" example:"019628ef-cd76-7d2d-bf80-48b8011fad40"`
	BeatmakerName      string         `json:"beatmaker_name" validate:"required"`
	AvailableFiles     AvailableFiles `json:"available_files" gorm:"foreignKey:UnpublishedBeatID;constraint:OnDelete:CASCADE;" validate:"required"`
	Price              int            `json:"price" validate:"required,gte=1" example:"2999"`
	Tags               []Tag          `json:"tags" validate:"required" gorm:"many2many:beat_tags;constraint:OnDelete:CASCADE"` //many to many
	BPM                int            `json:"bpm" validate:"required,gte=20,lte=400" example:"120"`
	Description        string         `json:"description" validate:"min=0,max=5000" example:"Chill summer beat with tropical influences"`
	Genres             []Genre        `json:"genres" validate:"required" gorm:"many2many:beat_genres;joinForeignKey:UnpublishedBeatID;joinReferences:GenreID;constraint:OnDelete:CASCADE"`
	Moods              []Mood         `json:"moods" validate:"required" gorm:"many2many:beat_moods;constraint:OnDelete:CASCADE"`   //many to many
	KeynoteID          uint           `json:"keynote_id" validate:"required" gorm:"default:NULL" example:"2"`                      //keynote has many beats, but each beat has only one keynote`
	Keynote            Keynote        `json:"keynote" gorm:"foreignKey:KeynoteID"`                                                 //`gorm:"foreignKey:UnpublishedBeatID;constraint:OnDelete:CASCADE;" validate:"required"`
	Timestamps         []Timestamp    `json:"timestamps" gorm:"foreignKey:BeatID;constraint:OnDelete:CASCADE"` //a beat has many timestamps, but each timestamp has only one beat
	Instruments        []Instrument   `json:"instruments" gorm:"many2many:beat_instruments;constraint:OnDelete:CASCADE"`           //many to many
	Status             Status         `json:"status" swaggerignore:"true"`
	Err                string         `json:"error"`
	SentToModerationAt int64          `json:"sentToModerationAt"`
	CreatedAt          int64          `json:"createdAt"`
	UpdatedAt          int64          `json:"updatesAt"`
}

// Validate - кастомная валидация
// func (beat *UnpublishedBeat) Validate() error {
// 	if beat.Status != "draft" {
// 		return errors.New("you tried to edit beat with a status, different from draft, prohibited")
// 	}
// 	return nil
// }

// BeforeSave - хуки GORM для автоматической валидации
// func (beat *UnpublishedBeat) BeforeSave(tx *gorm.DB) error {
// 	return beat.Validate()
// }

func MigrateAll(db *gorm.DB) error {
	err := db.AutoMigrate(
		&Instrument{},
		&BeatInstrument{},
		&Keynote{},
		&Genre{},
		&BeatGenre{},
		&Mood{},
		&BeatMood{},
		&Tag{},
		&BeatTag{},
		&UnpublishedBeat{},
		&AvailableFiles{},
		&Timestamp{},
	)
	return err
}
