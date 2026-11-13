package compliance

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EPRCertificate struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CertificateNumber    string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"certificate_number"` // "CERT-EPR-2026-001"
	RecipientCompanyName string         `gorm:"type:varchar(150);not null" json:"recipient_company_name"`       // Restoran / Hotel / Kantor
	TaxNumberNPWP        string         `gorm:"type:varchar(50)" json:"tax_number_npwp"`
	TotalWasteKg         float64        `gorm:"type:decimal(12,2);not null" json:"total_waste_kg"`
	PlasticWasteKg       float64        `gorm:"type:decimal(12,2);default:0" json:"plastic_waste_kg"`
	PaperWasteKg         float64        `gorm:"type:decimal(12,2);default:0" json:"paper_waste_kg"`
	MetalWasteKg         float64        `gorm:"type:decimal(12,2);default:0" json:"metal_waste_kg"`
	Co2AvoidedKg         float64        `gorm:"type:decimal(12,2);not null" json:"co2_avoided_kg"`
	IssuedDate           time.Time      `json:"issued_date"`
	ValidUntil           time.Time      `json:"valid_until"`
	VerificationHash     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"verification_hash"`
	Status               string         `gorm:"type:varchar(20);default:'VALID'" json:"status"`
	CreatedAt            time.Time      `json:"created_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
}

func (e *EPRCertificate) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return
}

type SubscriptionFrequency string

const (
	FrequencyDaily       SubscriptionFrequency = "DAILY"
	FrequencyTwiceAWeek  SubscriptionFrequency = "TWICE_WEEKLY"
	FrequencyWeekly      SubscriptionFrequency = "WEEKLY"
	FrequencyBiWeekly    SubscriptionFrequency = "BI_WEEKLY"
)

type B2BSubscriptionContract struct {
	ID              uuid.UUID             `gorm:"type:uuid;primaryKey" json:"id"`
	ContractNumber  string                `gorm:"type:varchar(50);uniqueIndex;not null" json:"contract_number"`
	CompanyName     string                `gorm:"type:varchar(150);not null" json:"company_name"`
	Address         string                `gorm:"type:text;not null" json:"address"`
	ContactPerson   string                `gorm:"type:varchar(100);not null" json:"contact_person"`
	ContactPhone    string                `gorm:"type:varchar(30);not null" json:"contact_phone"`
	ContactEmail    string                `gorm:"type:varchar(100);not null" json:"contact_email"`
	PickupFrequency SubscriptionFrequency `gorm:"type:varchar(30);not null" json:"pickup_frequency"`
	MonthlyFeeIDR   float64               `gorm:"type:decimal(14,2);not null" json:"monthly_fee_idr"`
	StartDate       time.Time             `json:"start_date"`
	EndDate         time.Time             `json:"end_date"`
	IsActive        bool                  `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
	DeletedAt       gorm.DeletedAt        `gorm:"index" json:"-"`
}

func (b *B2BSubscriptionContract) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return
}

type IssueCertificateDTO struct {
	RecipientCompanyName string  `json:"recipient_company_name" validate:"required"`
	TaxNumberNPWP        string  `json:"tax_number_npwp"`
	PlasticWasteKg       float64 `json:"plastic_waste_kg"`
	PaperWasteKg         float64 `json:"paper_waste_kg"`
	MetalWasteKg         float64 `json:"metal_waste_kg"`
	ValidMonths          int     `json:"valid_months" validate:"required,min=1,max=36"`
}

type CreateSubscriptionDTO struct {
	CompanyName     string                `json:"company_name" validate:"required"`
	Address         string                `json:"address" validate:"required"`
	ContactPerson   string                `json:"contact_person" validate:"required"`
	ContactPhone    string                `json:"contact_phone" validate:"required"`
	ContactEmail    string                `json:"contact_email" validate:"required,email"`
	PickupFrequency SubscriptionFrequency `json:"pickup_frequency" validate:"required"`
	MonthlyFeeIDR   float64               `json:"monthly_fee_idr" validate:"required,min=100000"`
	DurationMonths  int                   `json:"duration_months" validate:"required,min=1"`
}
