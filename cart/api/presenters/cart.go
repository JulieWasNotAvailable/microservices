package presenters

import "github.com/google/uuid"

type Cart struct {
	ID       uuid.UUID `json:"id"`
	Licenses []License `json:"licenses" validate:"required" gorm:"many2many:cart_licenses;"`
	UserID    uuid.UUID      `json:"userId"`
}

type CartLicenses struct {
	CartID    uuid.Domain `gorm:"primaryKey;constraint:OnDelete:CASCADE"`
	LicenseID uint        `gorm:"primaryKey"`
}
