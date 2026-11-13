package notification

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationChannel string

const (
	ChannelInApp    NotificationChannel = "IN_APP"
	ChannelWhatsApp NotificationChannel = "WHATSAPP"
	ChannelPush     NotificationChannel = "PUSH"
)

type NotificationStatus string

const (
	StatusSent NotificationStatus = "SENT"
	StatusRead NotificationStatus = "READ"
)

type Notification struct {
	ID            uuid.UUID           `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        uuid.UUID           `gorm:"type:uuid;not null;index" json:"user_id"`
	Title         string              `gorm:"type:varchar(150);not null" json:"title"`
	Message       string              `gorm:"type:text;not null" json:"message"`
	Channel       NotificationChannel `gorm:"type:varchar(20);not null;default:'IN_APP'" json:"channel"`
	Status        NotificationStatus  `gorm:"type:varchar(20);not null;default:'SENT'" json:"status"`
	ReferenceType string              `gorm:"type:varchar(50)" json:"reference_type,omitempty"` // "WALLET_CREDIT", "PICKUP_DISPATCH"
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	DeletedAt     gorm.DeletedAt      `gorm:"index" json:"-"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) (err error) {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return
}

type SendNotificationDTO struct {
	UserID        uuid.UUID           `json:"user_id" validate:"required"`
	Title         string              `json:"title" validate:"required"`
	Message       string              `json:"message" validate:"required"`
	Channel       NotificationChannel `json:"channel"`
	ReferenceType string              `json:"reference_type"`
}
