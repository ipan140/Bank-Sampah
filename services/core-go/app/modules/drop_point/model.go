package drop_point

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DropPoint struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name           string         `gorm:"type:varchar(150);not null" json:"name"`
	Address        string         `gorm:"type:text;not null" json:"address"`
	RT             string         `gorm:"type:varchar(10)" json:"rt"`
	RW             string         `gorm:"type:varchar(10)" json:"rw"`
	District       string         `gorm:"type:varchar(100)" json:"district"` // Kecamatan
	City           string         `gorm:"type:varchar(100);not null" json:"city"`
	Latitude       float64        `gorm:"type:decimal(10,8);not null" json:"latitude"`
	Longitude      float64        `gorm:"type:decimal(11,8);not null" json:"longitude"`
	OperatingHours string         `gorm:"type:varchar(100)" json:"operating_hours"` // "Senin - Sabtu: 08:00 - 16:00"
	ContactPerson  string         `gorm:"type:varchar(100)" json:"contact_person"`
	ContactPhone   string         `gorm:"type:varchar(30)" json:"contact_phone"`
	IsActive       bool           `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (d *DropPoint) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return
}

type MobilePickupSchedule struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	DropPointID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"drop_point_id"`
	DropPointName string         `gorm:"type:varchar(150)" json:"drop_point_name"`
	DayOfWeek     string         `gorm:"type:varchar(20);not null" json:"day_of_week"` // "Senin", "Rabu", "Jumat"
	StartTime     string         `gorm:"type:varchar(10);not null" json:"start_time"`  // "09:00"
	EndTime       string         `gorm:"type:varchar(10);not null" json:"end_time"`    // "12:00"
	CourierName   string         `gorm:"type:varchar(100)" json:"courier_name"`
	VehiclePlate  string         `gorm:"type:varchar(20)" json:"vehicle_plate"`
	CreatedAt     time.Time      `json:"created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *MobilePickupSchedule) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return
}

type CreateDropPointDTO struct {
	Name           string  `json:"name" validate:"required"`
	Address        string  `json:"address" validate:"required"`
	RT             string  `json:"rt"`
	RW             string  `json:"rw"`
	District       string  `json:"district"`
	City           string  `json:"city" validate:"required"`
	Latitude       float64 `json:"latitude" validate:"required"`
	Longitude      float64 `json:"longitude" validate:"required"`
	OperatingHours string  `json:"operating_hours"`
	ContactPerson  string  `json:"contact_person"`
	ContactPhone   string  `json:"contact_phone"`
}
