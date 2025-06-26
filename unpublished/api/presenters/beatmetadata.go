package presenters

import "github.com/google/uuid"

type AvailableFiles struct {
	ID                uuid.UUID `json:"id"`
	MP3Url            string    `json:"mp3url"`
	WAVUrl            string    `json:"wavurl"`
	ZIPUrl            string    `json:"zipurl"`
	UnpublishedBeatId uuid.UUID `json:"unpublishedBeatId"`
}

type Instrument struct {
	ID   uint   `json:"id" example:"1"`
	Name string `json:"name" example:"Piano"`
}

type Genre struct {
	ID   uint   `json:"id" example:"1"`
	Name string `json:"name" example:"Lo-Fi Hip Hop"`
}

type Timestamp struct {
	ID     uint      `json:"id" example:"1"`
	BeatID uuid.UUID `json:"beat_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name   string    `json:"name" example:"Intro"`
}

type Tag struct {
	ID   uint   `json:"id" example:"1"`
	UnpublishedBeat []*UnpublishedBeatResponse `json:"-"` //gorm:"many2many:beat_tags;"
	Name string `json:"name" example:"Chill"`
}

type Mood struct {
	ID   uint   `json:"id" example:"1"`
	Name string `json:"name" example:"Relaxed"`
}

type Keynote struct {
	ID   uint   `json:"id" example:"1"`
	Name string `json:"name" example:"C# Minor"`
}

//@Description presenters.MetadataSuccessResponse
type MetadataSuccessResponse struct {
	Status bool        `json:"status" example:"true"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error" example:""`
}

//@Description presenters.MetadataListResponse
type MetadataListResponse struct {
	Status bool        `json:"status" example:"true"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error" example:""`
}

//@Description presenters.MetadataErrorResponse
type MetadataErrorResponse struct {
	Status bool   `json:"status" example:"false"`
	Data   string `json:"data" example:""`
	Error  string `json:"error" example:"error message"`
}

func CreateMetadataSuccessResponse(data interface{}) *MetadataSuccessResponse {
	return &MetadataSuccessResponse{
		Status: true,
		Data:   data,
		Error:  "",
	}
}

func CreateMetadataListResponse(data interface{}) *MetadataListResponse {
	return &MetadataListResponse{
		Status: true,
		Data:   data,
		Error:  "",
	}
}

func CreateMetadataErrorResponse(err error) *MetadataErrorResponse {
	return &MetadataErrorResponse{
		Status: false,
		Data:   "",
		Error:  err.Error(),
	}
}