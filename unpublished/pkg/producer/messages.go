package producer

import (
	"github.com/JulieWasNotAvailable/microservices/unpublished/api/presenters"
)

type KafkaMessageToMFCC struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
}

type KafkaMessageBeatForPublishing struct {
	Beat presenters.UnpublishedBeat
	MFCC []float64
}