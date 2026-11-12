package pickup

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PickupStatus string

const (
	StatusPending   PickupStatus = "PENDING"
	StatusAssigned  PickupStatus = "ASSIGNED"
	StatusInTransit PickupStatus = "IN_TRANSIT"
	StatusWeighed   PickupStatus = "WEIGHED"
	StatusCompleted PickupStatus = "COMPLETED"
	StatusCancelled PickupStatus = "CANCELLED"
)

type PickupRequest struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	NasabahID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"nasabah_id"`
	CourierID       *uuid.UUID     `gorm:"type:uuid;index" json:"courier_id,omitempty"`
	Status          PickupStatus   `gorm:"type:varchar(20);not null;default:'PENDING'" json:"status"`
	PickupAddress   string         `gorm:"type:text;not null" json:"pickup_address"`
	Latitude        float64        `gorm:"type:decimal(10,8);not null" json:"latitude"`
	Longitude       float64        `gorm:"type:decimal(11,8);not null" json:"longitude"`
	ScheduledAt     time.Time      `json:"scheduled_at"`
	EstimatedWeight float64        `gorm:"type:decimal(8,2)" json:"estimated_weight"`
	Notes           string         `gorm:"type:text" json:"notes"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *PickupRequest) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}

type CreatePickupDTO struct {
	PickupAddress   string    `json:"pickup_address" validate:"required"`
	Latitude        float64   `json:"latitude" validate:"required"`
	Longitude       float64   `json:"longitude" validate:"required"`
	ScheduledAt     time.Time `json:"scheduled_at" validate:"required"`
	EstimatedWeight float64   `json:"estimated_weight"`
	Notes           string    `json:"notes"`
}

type AssignCourierDTO struct {
	CourierID uuid.UUID `json:"courier_id" validate:"required"`
}

type UpdateStatusDTO struct {
	Status PickupStatus `json:"status" validate:"required"`
}

type PickupFilter struct {
	Status    string `query:"status"`
	NasabahID string `query:"nasabah_id"`
	CourierID string `query:"courier_id"`
	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
}
