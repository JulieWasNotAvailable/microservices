package model

import "gorm.io/gorm"

type License struct {
	ID           uint
	Author       *string
	Title        *string
	MP3Path      *string
	WAVPath      *string
	TRACKOUTPath *string
	Price        *uint
}

func MigrateLicenses(db *gorm.DB) error {
	err := db.AutoMigrate(&License{})
	return err
}