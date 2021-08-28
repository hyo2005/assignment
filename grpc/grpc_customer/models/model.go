package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	Name       string    `gorm:"type:varchar(256)"`
	Address    string    `gorm:"type:varchar(256)"`
	Phone      string    `gorm:"type:varchar(256)"`
	License_id string    `gorm:"type:varchar(256)"`
	Active     bool      `gorm:"type:bool"`
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()`
	Email      string    `gorm:"type:varchar(256)"`
	Password   string    `gorm:"type:varchar(256)"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
