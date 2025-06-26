package producer

import (
	"github.com/JulieWasNotAvailable/microservices/unpublished/internal/entities"
)

type KafkaMessageToMFCC struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
}

type KafkaMessageBeatForPublishing struct {
	Beat entities.UnpublishedBeat
	MFCC []float64
}