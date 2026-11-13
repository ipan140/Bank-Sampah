package ppob

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillType string

const (
	BillPLNToken BillType = "PLN_TOKEN"
	BillPulsa    BillType = "PULSA"
	BillData     BillType = "PAKET_DATA"
	BillBPJS     BillType = "BPJS"
	BillPDAM     BillType = "PDAM"
	BillPBB      BillType = "PBB"
)

type PPOBTransaction struct {
	ID                uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID            uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Type              BillType       `gorm:"type:varchar(30);not null" json:"type"`
	TargetNumber      string         `gorm:"type:varchar(50);not null" json:"target_number"` // No HP / No Meteran / No Pelanggan
	Amount            float64        `gorm:"type:decimal(12,2);not null" json:"amount"`
	AdminFee          float64        `gorm:"type:decimal(10,2);default:0" json:"admin_fee"`
	TotalDeducted     float64        `gorm:"type:decimal(12,2);not null" json:"total_deducted"`
	SerialNumberToken string         `gorm:"type:varchar(100)" json:"serial_number_token,omitempty"` // Token PLN / SN Pulsa
	Status            string         `gorm:"type:varchar(20);default:'SUCCESS'" json:"status"`
	CreatedAt         time.Time      `json:"created_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *PPOBTransaction) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}

type GoldSavingsAccount struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID           uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	BalanceGrams     float64        `gorm:"type:decimal(10,4);not null;default:0" json:"balance_grams"`
	TotalInvestedIDR float64        `gorm:"type:decimal(14,2);not null;default:0" json:"total_invested_idr"`
	LastUpdated      time.Time      `json:"last_updated"`
	CreatedAt        time.Time      `json:"created_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (g *GoldSavingsAccount) BeforeCreate(tx *gorm.DB) (err error) {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	g.LastUpdated = time.Now()
	return
}

type WasteDonation struct {
	ID                     uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID                 uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	CharityInstitutionName string         `gorm:"type:varchar(150);not null" json:"charity_institution_name"` // "Panti Asuhan Kasih Ibu", "Masjid Al-Ikhlas"
	Amount                 float64        `gorm:"type:decimal(12,2);not null" json:"amount"`
	Notes                  string         `gorm:"type:text" json:"notes"`
	IsAnonymous            bool           `gorm:"default:false" json:"is_anonymous"`
	CreatedAt              time.Time      `json:"created_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"-"`
}

func (d *WasteDonation) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return
}

type PayBillDTO struct {
	Type         BillType `json:"type" validate:"required"`
	TargetNumber string   `json:"target_number" validate:"required"`
	Amount       float64  `json:"amount" validate:"required,min=5000"`
}

type ConvertGoldDTO struct {
	AmountIDR float64 `json:"amount_idr" validate:"required,min=10000"`
}

type CreateDonationDTO struct {
	CharityInstitutionName string  `json:"charity_institution_name" validate:"required"`
	Amount                 float64 `json:"amount" validate:"required,min=5000"`
	Notes                  string  `json:"notes"`
	IsAnonymous            bool    `json:"is_anonymous"`
}
