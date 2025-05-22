package entities

type Role struct {
	ID       uint   `json:"id"`
	Rolename string `gorm:"unique"`
	User     []User `gorm:"foreignKey:RoleID"`
}
