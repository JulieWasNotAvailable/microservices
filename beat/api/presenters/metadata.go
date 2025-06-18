package presenters

import (
	"github.com/google/uuid"
)

type AvailableFiles struct {
	ID     uuid.UUID `json:"id"`
	MP3Url string    `json:"mp3url"`
	WAVUrl string    `json:"wavurl"`
	ZIPUrl string    `json:"zipurl"`
	BeatId uuid.UUID `json:"beatId"`
}

// @Description entities.Genre
type Genre struct {
	ID         uint    `json:"id" swaggerignore:"true"`
	Beat       []*Beat `json:"-" gorm:"many2many:beat_genres;" swaggerignore:"true"`
	PictureUrl string  `json:"picture_url"`
	Name       string  `json:"name" example:"Jerk"`
}

type Timestamp struct {
	ID        uint      `json:"id" swaggerignore:"true"`
	BeatId    uuid.UUID `json:"beatId" example:"01963e01-e46c-7996-996a-42ad3df115ac"`
	Name      string    `json:"title"`
	TimeStart int64     `json:"start_time" validate:"required,gte=1,lte=299"`
	TimeEnd   int64     `json:"end_time" validate:"required,gte=2,lte=300"`
}

type Tag struct {
	ID   uint    `json:"id"`
	Beat []*Beat `json:"-" gorm:"many2many:beat_tags;" swaggerignore:"true"`
	Name string  `json:"name"`
}

type TrendingTags []struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type TrendingGenres []struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type Mood struct {
	ID   uint    `json:"id"`
	Beat []*Beat `json:"-" gorm:"many2many:beat_moods;" swaggerignore:"true"`
	Name string  `json:"name"`
}

type Keynote struct {
	ID    uint   `json:"id"`
	Beats []Beat `json:"-" gorm:"foreignKey:KeynoteID" swaggerignore:"true"`
	Name  string `json:"name"`
}

type Instrument struct {
	ID   uint    `json:"id"`
	Beat []*Beat `json:"-" gorm:"many2many:beat_instruments;" swaggerignore:"true"`
	Name string  `json:"name"`
}

type MFCC struct {
	ID     uint      `json:"id"`
	BeatId uuid.UUID `json:"beatId"`
}

type Filters struct {
	Genres   []uint `json:"genres,omitempty"`
	Moods    []uint `json:"moods,omitempty"`
	Tags     []uint `json:"tags,omitempty"`
	MinPrice int    `json:"min_price,omitempty"`
	MaxPrice int    `json:"max_price,omitempty"`
	MinBPM   int    `json:"min_bpm,omitempty"`
	MaxBPM   int    `json:"max_bpm,omitempty"`
	Keynote  uint   `json:"keynote,omitempty"`
	ItemsNum int    `json:"items_num,omitempty"`
}

// @Description presenters.MetadataSuccessResponse
type MetadataSuccessResponse struct {
	Status bool        `json:"status" example:"true"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error" example:""`
}

// @Description presenters.MetadataListResponse
type MetadataListResponse struct {
	Status bool        `json:"status" example:"true"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error" example:""`
}

// @Description presenters.MetadataErrorResponse
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
