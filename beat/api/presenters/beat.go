package presenters

import (
	"github.com/google/uuid"
)

type Beat struct {
	ID          uuid.UUID    `json:"id" example:"019628ef-cd76-7d2d-bf80-48b8011fad40"`
	Name        string       `json:"name" validate:"required,min=2,max=60" example:"Summer Vibes"`
	Picture     string       `json:"picture" example:"https://storage.yandexcloud.net/imagesall/019623bd-3d0b-7dc2-8a1f-f782adeb42b4"`
	BeatmakerID uuid.UUID    `json:"beatmakerId" validate:"required" example:"019628ef-cd76-7d2d-bf80-48b8011fad40"`
	URL         string       `json:"url" validate:"required" example:"https://storage.yandexcloud.net/mp3beats/019623bd-3d0b-7dc2-8a1f-f782adeb42b4"`
	Price       int          `json:"price" validate:"required, gte=1" example:"2999"`
	Tags        []Tag        `json:"tags" validate:"required" gorm:"many2many:beat_tags;"` //many to many
	BPM         int          `json:"bpm" validate:"required,gte=20,lte=400" example:"120"`
	Description string       `json:"description" validate:"min=2,max=500" example:"Chill summer beat with tropical influences"`
	Genres      []Genre      `json:"genres" validate:"required" gorm:"many2many:beat_genres;"` //many to many
	Moods       []Mood       `json:"moods" validate:"required" gorm:"many2many:beat_moods;"`   //many to many
	KeynoteID   uint         `json:"keynoteId" validate:"required" example:"11"`               //keynote has many beats, but each beat has only one keynote`
	Timestamps  []Timestamp  `json:"timestamps" validate:"required" gorm:"foreignKey:BeatId"`  //a beat has many timestamps, but each timestamp has only one beat
	Instruments []Instrument `json:"instruments" gorm:"many2many:beat_instruments"`            //many to many
	Plays       int64        `json:"plays"`
	CreatedAt   int64        `json:"created_at"`
}

// @Description BeatSuccessResponse
type BeatSuccessResponse struct {
	Status bool        `json:"status" example:"true"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error" example:""`
}

// @Description BeatListResponse
type BeatListResponse struct {
	Status bool        `json:"status" example:"true"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error" example:""`
}

// @Description BeatErrorResponse
type BeatErrorResponse struct {
	Status bool   `json:"status" example:"false"`
	Data   string `json:"data" example:""`
	Error  string `json:"error" example:"error message"`
}

func CreateBeatSuccessResponse(data interface{}) *BeatSuccessResponse {
	return &BeatSuccessResponse{
		Status: true,
		Data:   data,
		Error:  "",
	}
}

func CreateBeatListResponse(data interface{}) *BeatListResponse {
	return &BeatListResponse{
		Status: true,
		Data:   data,
		Error:  "",
	}
}

func CreateBeatErrorResponse(err error) *BeatErrorResponse {
	return &BeatErrorResponse{
		Status: false,
		Data:   "",
		Error:  err.Error(),
	}
}
