package finance_report

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExpenseCategory string

const (
	ExpenseOperational ExpenseCategory = "OPERATIONAL" // BBM, servis kendaraan, timbangan
	ExpenseSalary      ExpenseCategory = "SALARY"      // Gaji petugas / kurir
	ExpensePackaging   ExpenseCategory = "PACKAGING"   // Karung, tali pres, kawat
	ExpenseOther       ExpenseCategory = "OTHER"
)

type OperationalExpense struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	Category    ExpenseCategory `gorm:"type:varchar(30);not null" json:"category"`
	Amount      float64         `gorm:"type:decimal(14,2);not null" json:"amount"`
	Description string          `gorm:"type:text;not null" json:"description"`
	ReceiptURL  string          `gorm:"type:varchar(255)" json:"receipt_url"`
	ExpenseDate time.Time       `json:"expense_date"`
	CreatedAt   time.Time       `json:"created_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (o *OperationalExpense) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return
}

type ProfitAndLossStatement struct {
	TotalB2BRevenue     float64 `json:"total_b2b_revenue"`      // Penjualan ke pabrik
	TotalCOGS           float64 `json:"total_cogs"`             // Biaya beli sampah dari nasabah
	GrossProfit         float64 `json:"gross_profit"`           // Laba Kotor (Revenue - COGS)
	GrossProfitMarginPct float64 `json:"gross_profit_margin_pct"`
	TotalExpenses       float64 `json:"total_operational_expenses"` // Beban operasional
	NetProfit           float64 `json:"net_profit"`             // Laba Bersih
	TotalNasabahPayouts float64 `json:"total_nasabah_payouts"`  // Total saldo ditarik / dicairkan
}

type CreateExpenseDTO struct {
	Category    ExpenseCategory `json:"category" validate:"required"`
	Amount      float64         `json:"amount" validate:"required,min=1000"`
	Description string          `json:"description" validate:"required"`
	ReceiptURL  string          `json:"receipt_url"`
	ExpenseDate time.Time       `json:"expense_date"`
}
