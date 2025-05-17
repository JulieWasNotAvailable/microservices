package entities

import "github.com/google/uuid"

type MFCC struct {
	ID     uint
	BeatId uuid.UUID
	Crm1   float64
	Crm2   float64
	Crm3   float64
	Crm4   float64
	Crm5   float64
	Crm6   float64
	Crm7   float64
	Crm8   float64
	Crm9   float64
	Crm10  float64
	Crm11  float64
	Crm12  float64
	Mlspc  float64
	Mfcc1  float64
	Mfcc2  float64
	Mfcc3  float64
	Mfcc4  float64
	Mfcc5  float64
	Mfcc6  float64
	Mfcc7  float64
	Mfcc8  float64
	Mfcc9  float64
	Mfcc10 float64
	Mfcc11 float64
	Mfcc12 float64
	Mfcc13 float64
	Mfcc14 float64
	Mfcc15 float64
	Mfcc16 float64
	Mfcc17 float64
	Mfcc18 float64
	Mfcc19 float64
	Mfcc20 float64
	Mfcc21 float64
	Mfcc22 float64
	Mfcc23 float64
	Mfcc24 float64
	Mfcc25 float64
	Mfcc26 float64
	Mfcc27 float64
	Mfcc28 float64
	Mfcc29 float64
	Mfcc30 float64
	Mfcc31 float64
	Mfcc32 float64
	Mfcc33 float64
	Mfcc34 float64
	Mfcc35 float64
	Mfcc36 float64
	Mfcc37 float64
	Mfcc38 float64
	Mfcc39 float64
	Mfcc40 float64
	Mfcc41 float64
	Mfcc42 float64
	Mfcc43 float64
	Mfcc44 float64
	Mfcc45 float64
	Mfcc46 float64
	Mfcc47 float64
	Mfcc48 float64
	Mfcc49 float64
	Mfcc50 float64
	Spc    float64
	Err    string
}

// ID uint
// BeatId string
// BPM 157
// Chroma1 int64
// Chroma2 int64
// Chroma3 int64
// ...
// Chroma12 int64
// Melspec int64
// Mfcc1 int64
// ...
// Mfcc50 int64
// Spectral int64
// Err string
