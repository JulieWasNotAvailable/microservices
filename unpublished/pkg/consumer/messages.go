package consumer

import "github.com/google/uuid"

type KafkaMessageMFCC struct {
	ID       uuid.UUID `json:"beat_id"`
	Features []float64 `json:"features"`
	Error    string    `json:"error"`
}

type KafkaMessageURLUpdate struct {
	FileType string
	URL      string
}

type KafkaMessageLicenseCreated struct {
	BeatId uuid.UUID `json:"beatid"`
	Error  string    `json:"error"`
}

type KafkaMessageDeleteApprove struct {
	BeatId uuid.UUID `json:"beat_id"`
	Error  string    `json:"error"`
}
