package inventory

import (
	"time"

	"bank-sampah/core-go/app/modules/weighbridge"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StockMovementType string

const (
	MovementInWeighbridge StockMovementType = "STOCK_IN_WEIGHBRIDGE" // Masuk dari timbangan nasabah
	MovementOutSale       StockMovementType = "STOCK_OUT_SALE"        // Keluar dijual ke pabrik offtaker
	MovementAdjustment    StockMovementType = "STOCK_ADJUSTMENT"      // Koreksi fisik gudang
)

type WasteStock struct {
	ID          uuid.UUID                 `gorm:"type:uuid;primaryKey" json:"id"`
	Category    weighbridge.WasteCategory `gorm:"type:varchar(30);uniqueIndex;not null" json:"category"`
	CurrentKg   float64                   `gorm:"type:decimal(12,2);not null;default:0" json:"current_kg"`
	AvgCostPerKg float64                  `gorm:"type:decimal(12,2);not null;default:0" json:"avg_cost_per_kg"` // Rata-rata harga beli dari nasabah
	LastUpdated time.Time                 `json:"last_updated"`
}

func (s *WasteStock) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	s.LastUpdated = time.Now()
	return
}

type StockMutation struct {
	ID          uuid.UUID         `gorm:"type:uuid;primaryKey" json:"id"`
	Category    weighbridge.WasteCategory `gorm:"type:varchar(30);not null;index" json:"category"`
	Type        StockMovementType `gorm:"type:varchar(30);not null" json:"type"`
	QuantityKg  float64           `gorm:"type:decimal(12,2);not null" json:"quantity_kg"`
	StockBefore float64           `gorm:"type:decimal(12,2);not null" json:"stock_before"`
	StockAfter  float64           `gorm:"type:decimal(12,2);not null" json:"stock_after"`
	ReferenceID *uuid.UUID        `gorm:"type:uuid;index" json:"reference_id,omitempty"` // ID Ticket Timbangan atau ID Penjualan Pabrik
	Notes       string            `gorm:"type:text" json:"notes"`
	CreatedAt   time.Time         `json:"created_at"`
}

func (m *StockMutation) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return
}

type B2BSaleOrder struct {
	ID             uuid.UUID                 `gorm:"type:uuid;primaryKey" json:"id"`
	BuyerFactory   string                    `gorm:"type:varchar(150);not null" json:"buyer_factory"` // Nama Pabrik Daur Ulang / Offtaker
	InvoiceNumber  string                    `gorm:"type:varchar(50);uniqueIndex;not null" json:"invoice_number"`
	Category       weighbridge.WasteCategory `gorm:"type:varchar(30);not null" json:"category"`
	WeightSoldKg   float64                   `gorm:"type:decimal(12,2);not null" json:"weight_sold_kg"`
	CostPricePerKg float64                   `gorm:"type:decimal(12,2);not null" json:"cost_price_per_kg"` // Modal beli dari nasabah
	SalePricePerKg float64                   `gorm:"type:decimal(12,2);not null" json:"sale_price_per_kg"` // Harga jual ke pabrik
	TotalRevenue   float64                   `gorm:"type:decimal(14,2);not null" json:"total_revenue"`
	TotalCost      float64                   `gorm:"type:decimal(14,2);not null" json:"total_cost"`
	GrossProfit    float64                   `gorm:"type:decimal(14,2);not null" json:"gross_profit"` // Total Revenue - Total Cost
	SaleDate       time.Time                 `json:"sale_date"`
	CreatedAt      time.Time                 `json:"created_at"`
}

func (b *B2BSaleOrder) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	b.TotalRevenue = b.WeightSoldKg * b.SalePricePerKg
	b.TotalCost = b.WeightSoldKg * b.CostPricePerKg
	b.GrossProfit = b.TotalRevenue - b.TotalCost
	return
}

type CreateB2BSaleDTO struct {
	BuyerFactory   string                    `json:"buyer_factory" validate:"required"`
	Category       weighbridge.WasteCategory `json:"category" validate:"required"`
	WeightSoldKg   float64                   `json:"weight_sold_kg" validate:"required,min=1"`
	SalePricePerKg float64                   `json:"sale_price_per_kg" validate:"required,min=100"`
	Notes          string                    `json:"notes"`
}
