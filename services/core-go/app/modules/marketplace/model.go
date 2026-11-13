package marketplace

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EcoCategory string

const (
	EcoCategoryKompos      EcoCategory = "KOMPOS_ORGANIK"
	EcoCategoryKerajinan   EcoCategory = "KERAJINAN_LIMBAH"
	EcoCategoryBijiPlastik EcoCategory = "PELET_BIJI_PLASTIK"
	EcoCategoryEcobrick    EcoCategory = "ECOBRICK"
	EcoCategoryOther       EcoCategory = "LAINNYA"
)

type EcoProduct struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(150);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Category    EcoCategory    `gorm:"type:varchar(50);not null" json:"category"`
	PriceIDR    float64        `gorm:"type:decimal(12,2);not null" json:"price_idr"`
	Stock       int            `gorm:"not null;default:0" json:"stock"`
	Unit        string         `gorm:"type:varchar(20);default:'pcs'" json:"unit"` // "kg", "karung", "pcs"
	ImageURL    string         `gorm:"type:varchar(255)" json:"image_url"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *EcoProduct) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}

type EcoOrder struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	OrderNumber     string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"order_number"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	TotalAmountIDR  float64        `gorm:"type:decimal(14,2);not null" json:"total_amount_idr"`
	ShippingAddress string         `gorm:"type:text;not null" json:"shipping_address"`
	Status          string         `gorm:"type:varchar(20);default:'PAID'" json:"status"` // PAID, PROCESSING, DELIVERED
	Items           []EcoOrderItem `gorm:"foreignKey:OrderID" json:"items"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (o *EcoOrder) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return
}

type EcoOrderItem struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	OrderID     uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`
	ProductID   uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	ProductName string    `gorm:"type:varchar(150);not null" json:"product_name"`
	Quantity    int       `gorm:"not null" json:"quantity"`
	UnitPrice   float64   `gorm:"type:decimal(12,2);not null" json:"unit_price"`
	Subtotal    float64   `gorm:"type:decimal(14,2);not null" json:"subtotal"`
}

func (i *EcoOrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	i.Subtotal = float64(i.Quantity) * i.UnitPrice
	return
}

type CreateProductDTO struct {
	Name        string      `json:"name" validate:"required"`
	Description string      `json:"description"`
	Category    EcoCategory `json:"category" validate:"required"`
	PriceIDR    float64     `json:"price_idr" validate:"required,min=100"`
	Stock       int         `json:"stock" validate:"required,min=1"`
	Unit        string      `json:"unit"`
	ImageURL    string      `json:"image_url"`
}

type CheckoutItemDTO struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
}

type CheckoutOrderDTO struct {
	ShippingAddress string            `json:"shipping_address" validate:"required"`
	Items           []CheckoutItemDTO `json:"items" validate:"required,min=1"`
}
