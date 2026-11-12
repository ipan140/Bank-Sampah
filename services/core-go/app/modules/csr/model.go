package csr

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CsrProgram struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	PartnerName string         `gorm:"type:varchar(150);not null" json:"partner_name"`
	Title       string         `gorm:"type:varchar(200);not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	FundAmount  float64        `gorm:"type:decimal(14,2);not null" json:"fund_amount"`
	TargetKg    float64        `gorm:"type:decimal(10,2);not null" json:"target_kg"`
	CollectedKg float64        `gorm:"type:decimal(10,2);default:0" json:"collected_kg"`
	StartDate   time.Time      `json:"start_date"`
	EndDate     time.Time      `json:"end_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *CsrProgram) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}

type ImpactSummary struct {
	TotalWasteKg       float64 `json:"total_waste_kg"`
	PlasticWasteKg     float64 `json:"plastic_waste_kg"`
	PaperWasteKg       float64 `json:"paper_waste_kg"`
	MetalWasteKg       float64 `json:"metal_waste_kg"`
	Co2SavedKg         float64 `json:"co2_saved_kg"`         // Approx 1 kg plastic recycled saves 1.5 kg CO2
	TreesEquivalent    float64 `json:"trees_saved_equiv"`    // Approx 50 kg paper recycled saves 1 tree
	TotalNasabahHelped int64   `json:"total_nasabah_helped"`
	TotalFundsDistributed float64 `json:"total_funds_distributed"`
}

type CreateProgramDTO struct {
	PartnerName string    `json:"partner_name" validate:"required"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	FundAmount  float64   `json:"fund_amount" validate:"required,min=100000"`
	TargetKg    float64   `json:"target_kg" validate:"required,min=1"`
	StartDate   time.Time `json:"start_date" validate:"required"`
	EndDate     time.Time `json:"end_date" validate:"required"`
}
