package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cart struct {
	ID       uuid.UUID `json:"id"`
	Licenses []*License `json:"licenses" validate:"required" gorm:"many2many:cart_licenses;"`
	UserID   uuid.UUID `json:"userId"`
}

type CartLicenses struct {
	CartID    uuid.UUID `gorm:"primaryKey;constraint:OnDelete:CASCADE"`
	LicenseID uint      `gorm:"primaryKey"`
}

func MigrateAll(db *gorm.DB) error {
	err := db.AutoMigrate(
		&Cart{},
		&License{},
		&LicenseTemplate{},
	)
	return err
}
