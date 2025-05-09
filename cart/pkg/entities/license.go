package entities

import "github.com/google/uuid"

type License struct {
	ID                uint      `json:"id"`
	BeatID            uuid.UUID `json:"beatId"`
	Cart              []*Cart   `gorm:"many2many:cart_licenses;"`
	LicenseTemplateID uint      `json:"licenseTemplateId"`
	MP3               bool      `json:"mp3"`
	WAV               bool      `json:"wav"`
	ZIP               bool      `json:"zip"`
	UserID            uuid.UUID `json:"beatmakerId"`
	Price             int       `json:"price"`
}

type LicenseTemplate struct {
	ID                uint       `json:"id"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	MusicRecording    bool       `json:"musicRecording"`
	LiveProfit        bool       `json:"liveProfit"`
	DistributeCopies  int        `json:"distributeCopies"`
	AudioStreams      int        `json:"audioStreams"`
	RadioBroadcasting int        `json:"radioBroadcasting"`
	MusicVideos       int        `json:"musicVideos"`
	AvailableFilesID  int        `json:"availableFilesId"`
	License           []*License `json:"-" gorm:"many2many:license_template_licenses;joinForeignKey:LicenseID;joinReferences:LicenseTemplateID"`
	UserID            uuid.UUID  `json:"userId"`
}
