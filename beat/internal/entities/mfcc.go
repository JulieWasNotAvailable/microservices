package entities

import "github.com/google/uuid"

type MFCC struct {
	ID     uint
	BeatId uuid.UUID
	Col1   float64
	Col2   float64
}