package producer

import "github.com/google/uuid"

type KafkaMessageCreateLicense struct {
	BeatId uuid.UUID `json:"beatid"`
	Error  string    `json:"error"`
}
