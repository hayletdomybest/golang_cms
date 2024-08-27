package entities

import (
	"time"
	"weex_admin/internal/domain/valueobjects"
)

type Role struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`

	Policies []valueobjects.Policy `gorm:"-"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
