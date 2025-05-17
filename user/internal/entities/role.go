package entities

import (
	"gorm.io/gorm"
)

type Role struct {
	ID  		uint     
	Rolename string 	`gorm:"unique"`
	User	 []User 	`gorm:"foreignKey:RoleID"`
}

func MigrateRole(db *gorm.DB) error {
	err := db.AutoMigrate(&Role{})
	return err
}