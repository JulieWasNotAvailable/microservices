package entities

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateLicense struct {
	BeatId      uuid.UUID
	UserId      uuid.UUID
	LicenseList []License
}

type License struct {
	ID                uint            `json:"id"`
	BeatID            uuid.UUID       `json:"beatId"`
	Cart              []*Cart         `json:"-" gorm:"many2many:cart_licenses;"`
	LicenseTemplateID uint            `json:"licenseTemplateId"`
	LicenseTemplate   LicenseTemplate `json:"licenseTemplate" gorm:"foreignKey:LicenseTemplateID"`
	UserID            uuid.UUID       `json:"-"`
	Price             int             `json:"price"`
}

type LicenseTemplate struct {
	ID                uint       `json:"id"`
	Name              string     `json:"name"`
	MP3               bool       `json:"mp3"`
	WAV               bool       `json:"wav"`
	ZIP               bool       `json:"zip"`
	Description       string     `json:"description"`
	MusicRecording    bool       `json:"musicRecording"`
	LiveProfit        bool       `json:"liveProfit"`
	DistributeCopies  int        `json:"distributeCopies"`
	AudioStreams      int        `json:"audioStreams"`
	RadioBroadcasting int        `json:"radioBroadcasting"`
	MusicVideos       int        `json:"musicVideos"`
	AvailableFilesID  int        `json:"availableFilesId"`
	License           []*License `json:"license" gorm:"foreignKey:LicenseTemplateID"`
	UserID            uuid.UUID  `json:"-"`
}

func (model *LicenseTemplate) BeforeSave(tx *gorm.DB) (err error) {
	if model.DistributeCopies < -1 {
		return errors.New("distribute copies value must be > -1 ")
	}
	if model.AudioStreams < -1 {
		return errors.New("audio streams value must be > -1 ")
	}
	if model.RadioBroadcasting < -1 {
		return errors.New("radio broadcasting value must be > -1 ")
	}
	if model.MusicVideos < -1 {
		return errors.New("music videos value must be > -1 ")
	}
	return nil
}
