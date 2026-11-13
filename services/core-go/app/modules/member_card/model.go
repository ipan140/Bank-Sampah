package member_card

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MemberCard struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	CardNumber  string         `gorm:"type:varchar(30);uniqueIndex;not null" json:"card_number"`
	QRCodeToken string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"qr_code_token"` // Digunakan saat scan di HP / timbangan
	NFCUid      string         `gorm:"type:varchar(50);index" json:"nfc_uid,omitempty"`             // Opsional RFID/NFC fisik
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	IssuedAt    time.Time      `json:"issued_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *MemberCard) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	m.IssuedAt = time.Now()
	return
}

type ScanQRRequest struct {
	QRCodeToken string `json:"qr_code_token" validate:"required"`
}

type MemberInfoResponse struct {
	UserID     uuid.UUID `json:"user_id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	CardNumber string    `json:"card_number"`
	Balance    float64   `json:"balance"`
	Points     int64     `json:"points"`
}
