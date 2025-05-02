package entities

import "github.com/google/uuid"

type UnpublishedBeat struct {
	ID                 uuid.UUID      `json:"id" example:"019628ef-cd76-7d2d-bf80-48b8011fad40"`
	Name               string         `json:"name" validate:"required,min=2,max=60" example:"Summer Vibes"`
	Picture            string         `json:"picture" example:"https://storage.yandexcloud.net/imagesall/019623bd-3d0b-7dc2-8a1f-f782adeb42b4"`
	BeatmakerID        uuid.UUID      `json:"beatmakerId" validate:"required" example:"019628ef-cd76-7d2d-bf80-48b8011fad40"`
	AvailableFiles     AvailableFiles `validate:"required" gorm:"foreignKey:UnpublishedBeatID;constraint:OnDelete:CASCADE;"`
	Price              int            `json:"price" validate:"required" example:"2999"`
	Tags               []Tag          `json:"tags" validate:"required" gorm:"many2many:tag_beats;"` //many to many
	BPM                int            `json:"bpm" validate:"required,gte=20,lte=400" example:"120"`
	Description        string         `json:"description" validate:"min=2,max=500" example:"Chill summer beat with tropical influences"`
	Genres             []Genre        `json:"genres" validate:"required" gorm:"many2many:genre_beats;"` //many to many
	Moods              []Mood         `json:"moods" validate:"required" gorm:"many2many:mood_beats;"`   //many to many
	KeynoteID          uint           `json:"keynoteId" validate:"required" example:"11"`               //keynote has many beats, but each beat has only one keynote`
	Timestamps         []Timestamp    `json:"timestamps" validate:"required" gorm:"foreignKey:BeatID"`  //a beat has many timestamps, but each timestamp has only one beat
	Instruments        []Instrument   `json:"instruments" gorm:"many2many:instrument_beats"`            //many to many
	Status             string         `json:"status" example:"draft"`
	Err                string
	SentToModerationAt int64 `json:"sent_to_moderation_at"`
	CreatedAt          int64 `json:"created_at"`
	UpdatedAt          int64 `json:"updated_at"`
}