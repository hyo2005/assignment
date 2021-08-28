package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()`
	Customer_id uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()`
	Flight_id   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()`
	Code        string    `gorm:"type:varchar(256)"`
	Status      string    `gorm:"type:varchar(256)"`
	Booked_date time.Time `gorm:"type:timestamp"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
