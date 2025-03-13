package model

import (
	"time"
	"gorm.io/gorm"
)

//у одного бита может быть несколько лицензий
//и несколько лицензий у одного бита

type Beat struct{
	ID uint
	Author *string
	Title *string
	License *string
	Mood *string
	Date time.Time
	Genre *string
	Url *string
	FreeForNonProfit *uint
}

//for authomigration, because in postgres DB is not created authomatically
func MigrateBeats(db *gorm.DB) error {
	err := db.AutoMigrate(&Beat{})
	return err
}