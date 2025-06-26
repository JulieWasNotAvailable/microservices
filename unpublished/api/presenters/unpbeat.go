package presenters

import (
	"github.com/JulieWasNotAvailable/microservices/unpublished/internal/entities"
	"github.com/google/uuid"
)

// @Description presenters.UnpublishedBeatResponse
type UnpublishedBeatResponse struct {
	ID                 uuid.UUID               `json:"id" example:"019628ef-cd76-7d2d-bf80-48b8011fad40"`
	Name               string                  `json:"name" validate:"required,min=2,max=60" example:"Summer Vibes"`
	Picture            string                  `json:"picture" example:"https://storage.yandexcloud.net/imagesall/019623bd-3d0b-7dc2-8a1f-f782adeb42b4"`
	BeatmakerID        uuid.UUID               `json:"beatmakerId" validate:"required" example:"019628ef-cd76-7d2d-bf80-48b8011fad40"`
	BeatmakerName      string                  `json:"beatmakerName"`
	AvailableFiles     entities.AvailableFiles `json:"availableFiles" validate:"required" gorm:"foreignKey:UnpublishedBeatID;constraint:OnDelete:CASCADE;"`
	URL                string                  `json:"url" validate:"required" example:"https://storage.yandexcloud.net/mp3beats/019623bd-3d0b-7dc2-8a1f-f782adeb42b4"`
	Price              int                     `json:"price" validate:"required" example:"2999"`
	Tags               []entities.Tag          `json:"tags" validate:"required" gorm:"many2many:beat_tags;"` //many to many
	BPM                int                     `json:"bpm" validate:"required,gte=20,lte=400" example:"120"`
	Description        string                  `json:"description" validate:"min=2,max=500" example:"Chill summer beat with tropical influences"`
	Genres             []entities.Genre        `json:"genres" validate:"required" gorm:"many2many:beat_genres;"` //many to many
	Moods              []entities.Mood         `json:"moods" validate:"required" gorm:"many2many:beat_moods;"`   //many to many
	KeynoteID          uint                    `json:"keynoteId" validate:"required" example:"11"`               //keynote has many beats, but each beat has only one keynote`
	Keynote            Keynote                 `json:"keynote" gorm:"foreignKey:KeynoteID"`
	Timestamps         []entities.Timestamp    `json:"timestamps" validate:"required" gorm:"foreignKey:BeatID"` //a beat has many timestamps, but each timestamp has only one beat
	Instruments        []entities.Instrument   `json:"instruments" gorm:"many2many:beat_instruments"`           //many to many
	Status             string                  `json:"status" example:"draft"`                                  //status, err need to be in response, but not in request
	Err                string                  `json:"error"`
	SentToModerationAt int64                   `json:"sent_to_moderation_at"`
	CreatedAt          int64                   `json:"created_at"`
	UpdatedAt          int64                   `json:"updated_at"`
}

// @Description presenters.UnpublishedBeat
type UnpublishedBeatRequest struct {
	ID                 uuid.UUID               `json:"id" example:"019628ef-cd76-7d2d-bf80-48b8011fad40"`
	Name               string                  `json:"name" validate:"required,min=2,max=60" example:"Summer Vibes"`
	Picture            string                  `json:"picture" example:"https://storage.yandexcloud.net/imagesall/019623bd-3d0b-7dc2-8a1f-f782adeb42b4"`
	BeatmakerID        uuid.UUID               `json:"-" validate:"required" example:"019628ef-cd76-7d2d-bf80-48b8011fad40"`
	BeatmakerName      string                  `json:"-"`
	AvailableFiles     entities.AvailableFiles `json:"-" validate:"required" gorm:"foreignKey:UnpublishedBeatID;constraint:OnDelete:CASCADE;"`
	URL                string                  `json:"-" validate:"required" example:"https://storage.yandexcloud.net/mp3beats/019623bd-3d0b-7dc2-8a1f-f782adeb42b4"`
	Price              int                     `json:"-" validate:"required" example:"2999"`
	Tags               []entities.Tag          `json:"tags" validate:"required" gorm:"many2many:beat_tags;"` //many to many
	BPM                int                     `json:"bpm" validate:"required,gte=20,lte=400" example:"120"`
	Description        string                  `json:"description" validate:"min=2,max=500" example:"Chill summer beat with tropical influences"`
	Genres             []entities.Genre        `json:"genres" validate:"required" gorm:"many2many:beat_genres;"` //many to many
	Moods              []entities.Mood         `json:"moods" validate:"required" gorm:"many2many:beat_moods;"`   //many to many
	KeynoteID          uint                    `json:"keynoteId" validate:"required" example:"11"`               //keynote has many beats, but each beat has only one keynote`
	Keynote            Keynote                 `json:"-" gorm:"foreignKey:KeynoteID"`
	Timestamps         []entities.Timestamp    `json:"timestamps" validate:"required" gorm:"foreignKey:BeatID"` //a beat has many timestamps, but each timestamp has only one beat
	Instruments        []entities.Instrument   `json:"instruments" gorm:"many2many:beat_instruments"`           //many to many
	Status             string                  `json:"-" example:"draft"`                                       //status, err need to be in response, but not in request
	Err                string                  `json:"-"`
	SentToModerationAt int64                   `json:"-"`
	CreatedAt          int64                   `json:"-"`
	UpdatedAt          int64                   `json:"-"`
}

// @Description presenters.UnpublishedBeatSuccessResponse
type UnpublishedBeatSuccessResponse struct {
	Status bool                    `json:"status" example:"true"`
	Data   UnpublishedBeatResponse `json:"data"`
	Error  string                  `json:"error" example:"null"`
}

// @Description presenters.UnpublishedBeatListResponse
type UnpublishedBeatListResponse struct {
	Status bool                      `json:"status" example:"true"`
	Data   []UnpublishedBeatResponse `json:"data"`
	Error  string                    `json:"error" example:"null"`
}

// @Description presenters.UnpublishedBeatListResponse
type UnpublishedBeatErrorListResponse struct {
	Status bool     `json:"status" example:"true"`
	Data   string   `json:"data"`
	Error  []string `json:"error" example:"null"`
}

// @Description presenters.UnpublishedBeatErrorResponse
type UnpublishedBeatErrorResponse struct {
	Status bool   `json:"status" example:"true"`
	Data   string `json:"data"`
	Error  string `json:"error" example:"null"`
}

func CreateBeatSuccessResponse(data *entities.UnpublishedBeat) *UnpublishedBeatSuccessResponse {
	beat := UnpublishedBeatResponse{
		ID:             data.ID,
		Name:           data.Name,
		Picture:        data.Picture,
		BeatmakerID:    data.BeatmakerID,
		AvailableFiles: data.AvailableFiles,
		Price:          data.Price,
		Tags:           data.Tags,
		BPM:            data.BPM,
		Description:    data.Description,
		Genres:         data.Genres,
		Moods:          data.Moods,
		KeynoteID:      data.KeynoteID,
		Instruments:    data.Instruments,
		Status:         string(data.Status),
		Err:            data.Err,
		CreatedAt:      data.CreatedAt,
		UpdatedAt:      data.UpdatedAt,
	}

	return &UnpublishedBeatSuccessResponse{
		Status: true,
		Data:   beat,
		Error:  "",
	}
}

func CreateBeatSuccessResponse2(data UnpublishedBeatResponse) *UnpublishedBeatSuccessResponse {
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

func CreateBeatListSuccessResponse(data []UnpublishedBeatResponse) *UnpublishedBeatListResponse {
	return &UnpublishedBeatListResponse{
		Status: true,
		Data:   data,
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

func CreateBeatErrorArrayResponse(errs []error) *UnpublishedBeatErrorListResponse {
	stringErrs := []string{}
	for _, err := range errs {
		stringErrs = append(stringErrs, err.Error())
	}
	return &UnpublishedBeatErrorListResponse{
		Status: false,
		Data:   "",
		Error:  stringErrs,
	}
}

func RequestToEntity(data UnpublishedBeatRequest) entities.UnpublishedBeat {
	beat := entities.UnpublishedBeat{
		ID:             data.ID,
		Name:           data.Name,
		Picture:        data.Picture,
		BeatmakerID:    data.BeatmakerID,
		AvailableFiles: data.AvailableFiles,
		Price:          data.Price,
		Tags:           data.Tags,
		BPM:            data.BPM,
		Description:    data.Description,
		Genres:         data.Genres,
		Moods:          data.Moods,
		KeynoteID:      data.KeynoteID,
		Instruments:    data.Instruments,
		Status:         entities.Status(data.Status),
		Err:            data.Err,
		CreatedAt:      data.CreatedAt,
		UpdatedAt:      data.UpdatedAt,
	}
	return beat
}

func EntityToResponse(data entities.UnpublishedBeat) UnpublishedBeatResponse {
	beat := UnpublishedBeatResponse{
		ID:             data.ID,
		Name:           data.Name,
		Picture:        data.Picture,
		BeatmakerID:    data.BeatmakerID,
		AvailableFiles: data.AvailableFiles,
		Price:          data.Price,
		Tags:           data.Tags,
		BPM:            data.BPM,
		Description:    data.Description,
		Genres:         data.Genres,
		Moods:          data.Moods,
		KeynoteID:      data.KeynoteID,
		Instruments:    data.Instruments,
		Status:         string(data.Status),
		Err:            data.Err,
		CreatedAt:      data.CreatedAt,
		UpdatedAt:      data.UpdatedAt,
	}
	return beat
}

type BeatForPublishing struct {
	Beat UnpublishedBeatResponse
	MFCC []float64
}

func EntityListToResponse(beats []entities.UnpublishedBeat) []UnpublishedBeatResponse {
	beatResponses := []UnpublishedBeatResponse{}
	if len(beats) != 0 {
		for _, beat := range beats {
			beatResp := EntityToResponse(beat)
			beatResponses = append(beatResponses, beatResp)
		}
	}
	return beatResponses
}
