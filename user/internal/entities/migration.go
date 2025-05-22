package entities

import "gorm.io/gorm"

func MigrateAll(db *gorm.DB) error {
	err := db.AutoMigrate(
		&Role{},
		&User{},
		&Metadata{},
		&User_Follows_Beatmakers{},
	)
	return err
}