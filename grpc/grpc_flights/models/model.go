package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Flight struct {
	gorm.Model
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()`
	Name           string    `gorm:"type:varchar(256)"`
	From           string    `gorm:"type:varchar(256)"`
	To             string    `gorm:"type:varchar(256)"`
	Date           time.Time `gorm:"type:timestamp"`
	Status         string    `gorm:"type:varchar(256)"`
	Available_slot int       `gorm:"type:integer"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
