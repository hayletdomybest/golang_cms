package entities

import "time"

type User struct {
	ID        uint64 `gorm:"primaryKey"`
	Name      string `gorm:"unique;not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	RoleID    uint   `gorm:"not null"`
	Role      *Role  `gorm:"foreignKey:RoleID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
