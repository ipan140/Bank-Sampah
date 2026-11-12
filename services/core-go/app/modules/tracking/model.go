package tracking

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourierLocation struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CourierID uuid.UUID      `gorm:"type:uuid;not null;index" json:"courier_id"`
	Latitude  float64        `gorm:"type:decimal(10,8);not null" json:"latitude"`
	Longitude float64        `gorm:"type:decimal(11,8);not null" json:"longitude"`
	Heading   float64        `gorm:"type:decimal(5,2)" json:"heading"` // Arah kompas derajat
	Speed     float64        `gorm:"type:decimal(5,2)" json:"speed"`   // km/h
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c *CourierLocation) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}

type UpdateLocationDTO struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
	Heading   float64 `json:"heading"`
	Speed     float64 `json:"speed"`
}

type NearbyCourier struct {
	CourierID uuid.UUID `json:"courier_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	DistanceM float64   `json:"distance_meters"`
}
