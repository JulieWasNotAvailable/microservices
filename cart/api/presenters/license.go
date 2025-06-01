package presenters

import (
	"github.com/google/uuid"
)

type License struct {
	ID                uint            `json:"id"`
	BeatID            uuid.UUID       `json:"beatId"`
	LicenseTemplateID uint            `json:"licenseTemplateId"`
	LicenseTemplate   LicenseTemplate `json:"licenseTemplate" gorm:"foreignKey:LicenseTemplateID"`
	UserID            uuid.UUID       `json:"beatmakerId"`
	Price             int             `json:"price"`
}

type LicenseTemplate struct {
	ID                uint      `json:"id"`
	Name              string    `json:"name"`
	MP3               bool      `json:"mp3"`
	WAV               bool      `json:"wav"`
	ZIP               bool      `json:"zip"`
	Description       string    `json:"description"`
	MusicRecording    bool      `json:"musicRecording"`
	LiveProfit        bool      `json:"liveProfit"`
	DistributeCopies  int       `json:"distributeCopies"`
	AudioStreams      int       `json:"audioStreams"`
	RadioBroadcasting int       `json:"radioBroadcasting"`
	MusicVideos       int       `json:"musicVideos"`
	AvailableFilesID  int       `json:"availableFilesId"`
	Carts             []*Cart   `json:"-" gorm:"many2many:cart_licenses;joinForeignKey:LicenseID;joinReferences:CartID"`
	UserID            uuid.UUID `json:"userId"`
}
