package weighbridge

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WasteCategory string

const (
	CategoryPlastic WasteCategory = "PLASTIC"
	CategoryPaper   WasteCategory = "PAPER"
	CategoryMetal   WasteCategory = "METAL"
	CategoryGlass   WasteCategory = "GLASS"
	CategoryOrganic WasteCategory = "ORGANIC"
)

type ScaleTicket struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	PickupID    *uuid.UUID     `gorm:"type:uuid;index" json:"pickup_id,omitempty"`
	NasabahID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"nasabah_id"`
	OperatorID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"operator_id"`
	Category    WasteCategory  `gorm:"type:varchar(30);not null" json:"category"`
	GrossWeight float64        `gorm:"type:decimal(10,2);not null" json:"gross_weight"` // Kg
	TareWeight  float64        `gorm:"type:decimal(10,2);not null" json:"tare_weight"`  // Kg
	NetWeight   float64        `gorm:"type:decimal(10,2);not null" json:"net_weight"`   // Kg
	PricePerKg  float64        `gorm:"type:decimal(12,2);not null" json:"price_per_kg"`
	TotalPrice  float64        `gorm:"type:decimal(14,2);not null" json:"total_price"`
	DeviceID    string         `gorm:"type:varchar(50)" json:"device_id"` // ESP32 MAC/ID
	IsSettled   bool           `gorm:"default:false" json:"is_settled"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *ScaleTicket) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	s.NetWeight = s.GrossWeight - s.TareWeight
	s.TotalPrice = s.NetWeight * s.PricePerKg
	return
}

type CreateTicketDTO struct {
	PickupID    *uuid.UUID    `json:"pickup_id"`
	NasabahID   uuid.UUID     `json:"nasabah_id" validate:"required"`
	Category    WasteCategory `json:"category" validate:"required"`
	GrossWeight float64       `json:"gross_weight" validate:"required"`
	TareWeight  float64       `json:"tare_weight"`
	PricePerKg  float64       `json:"price_per_kg" validate:"required"`
	DeviceID    string        `json:"device_id"`
}

type ScaleTicketFilter struct {
	NasabahID string `query:"nasabah_id"`
	Category  string `query:"category"`
	DeviceID  string `query:"device_id"`
	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
}
