package presenters

import (
	"github.com/JulieWasNotAvailable/microservices/unpublished/pkg/entities"
	"github.com/google/uuid"
)

//@Description presenters.UnpublishedBeat
type UnpublishedBeat struct {
	ID              uuid.UUID       		`json:"id" example:"019628ef-cd76-7d2d-bf80-48b8011fad40"`
	Name            string    				`json:"name" validate:"required,min=2,max=60" example:"Summer Vibes"`
	Picture         string    				`json:"picture" example:"https://storage.yandexcloud.net/imagesall/019623bd-3d0b-7dc2-8a1f-f782adeb42b4"`
	BeatmakerID     uuid.UUID       		`json:"beatmakerId" validate:"required" example:"019628ef-cd76-7d2d-bf80-48b8011fad40"`
	AvailableFiles 	entities.AvailableFiles `validate:"required" gorm:"foreignKey:UnpublishedBeatID;constraint:OnDelete:CASCADE;"`
	URL             string    				`json:"url" validate:"required" example:"https://storage.yandexcloud.net/mp3beats/019623bd-3d0b-7dc2-8a1f-f782adeb42b4"`
	Price           int       				`json:"price" validate:"required" example:"2999"`
	Tags            []entities.Tag      	`json:"tags" validate:"required" gorm:"many2many:tag_beats;"` //many to many
	BPM             int       				`json:"bpm" validate:"required,gte=20,lte=400" example:"120"`
	Description     string    				`json:"description" validate:"min=2,max=500" example:"Chill summer beat with tropical influences"`
	Genres         	[]entities.Genre      	`json:"genres" validate:"required" gorm:"many2many:genre_beats;"`       //many to many
	Moods          	[]entities.Mood      	`json:"moods" validate:"required" gorm:"many2many:mood_beats;"`       //many to many
	KeynoteID       uint       				`json:"keynoteId" validate:"required" example:"11"`    //keynote has many beats, but each beat has only one keynote`
	Timestamps    	[]entities.Timestamp    `json:"timestamps" validate:"required" gorm:"foreignKey:BeatID"` //a beat has many timestamps, but each timestamp has only one beat
	Instruments   	[]entities.Instrument   `json:"instruments" gorm:"many2many:instrument_beats"` //many to many
	Status          string    				`json:"status" example:"draft"`
	Err 			string
	SentToModerationAt int64 				`json:"sent_to_moderation_at"`
	CreatedAt      int64 					`json:"created_at"`
	UpdatedAt      int64 					`json:"updated_at"`
}

type BeatForPublishing struct {
	Beat UnpublishedBeat
	MFCC string
}

//@Description presenters.UnpublishedBeatSuccessResponse
type UnpublishedBeatSuccessResponse struct {
	Status 		bool   						`json:"status" example:"true"`
	Data   		UnpublishedBeat 			`json:"data"`
	Error 		string                     `json:"error" example:"null"`
}

//@Description presenters.UnpublishedBeatListResponse
type UnpublishedBeatListResponse struct {
	Status 	  bool   						`json:"status" example:"true"`
	Data      []UnpublishedBeat 			`json:"data"`
	Error 	  string                       `json:"error" example:"null"`
}

//@Description presenters.UnpublishedBeatErrorResponse
type UnpublishedBeatErrorResponse struct {
	Status 	  bool   						`json:"status" example:"true"`
	Data      string 						`json:"data"`
	Error 	  string                       `json:"error" example:"null"`
}

func CreateBeatSuccessResponse(data *entities.UnpublishedBeat) *UnpublishedBeatSuccessResponse {
    beat := UnpublishedBeat{
        ID:             data.ID,
        Name:           data.Name,
        Picture:        data.Picture,
        BeatmakerID:    data.BeatmakerID,
        AvailableFiles: data.AvailableFiles,
        URL:      		data.URL,
        Price:          data.Price,
        Tags:           data.Tags,
        BPM:            data.BPM,
        Description:    data.Description,
        Genres:         data.Genres,
        Moods:          data.Moods,
        KeynoteID:      data.KeynoteID,
        Instruments:    data.Instruments,
        Status:         string(data.Status),
		Err: 			data.Err,
        CreatedAt:      data.CreatedAt,
        UpdatedAt:      data.UpdatedAt,
    }

    return &UnpublishedBeatSuccessResponse{
        Status: true,
        Data:   beat,
        Error:  "",
    }
}

func CreateBeatSuccessResponse2(data UnpublishedBeat) *UnpublishedBeatSuccessResponse {
    return &UnpublishedBeatSuccessResponse{
        Status: true,
        Data:   data,
        Error:  "",
    }
}

func CreateBeatErrorResponse(err error) *UnpublishedBeatErrorResponse {
	return &UnpublishedBeatErrorResponse{
		Status: false,
		Data:   "",
		Error:  err.Error(),
	}
}

func CreateBeatListSuccessResponse(data []UnpublishedBeat) *UnpublishedBeatListResponse {
	return &UnpublishedBeatListResponse{
		Status: true,
		Data:  	data,
		Error:  "",
	}
}

func CreateBeatListErrorResponse(err error) *UnpublishedBeatListResponse {
	return &UnpublishedBeatListResponse{
		Status: false,
		Data:   nil,
		Error:  err.Error(),
	}
}