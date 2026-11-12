package wallet

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionType string

const (
	TxCreditDeposit TransactionType = "CREDIT_DEPOSIT" // Masuk dari setoran sampah
	TxDebitRedeem   TransactionType = "DEBIT_REDEEM"   // Keluar tarik tunai/e-wallet
	TxDebitReward   TransactionType = "DEBIT_REWARD"   // Tukar reward point
)

type RedeemStatus string

const (
	RedeemPending   RedeemStatus = "PENDING"
	RedeemApproved  RedeemStatus = "APPROVED"
	RedeemRejected  RedeemStatus = "REJECTED"
	RedeemCompleted RedeemStatus = "COMPLETED"
)

type Wallet struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	Balance   float64        `gorm:"type:decimal(14,2);not null;default:0" json:"balance"`
	Points    int64          `gorm:"default:0" json:"points"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (w *Wallet) BeforeCreate(tx *gorm.DB) (err error) {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return
}

type LedgerTransaction struct {
	ID            uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	WalletID      uuid.UUID       `gorm:"type:uuid;not null;index" json:"wallet_id"`
	UserID        uuid.UUID       `gorm:"type:uuid;not null;index" json:"user_id"`
	Type          TransactionType `gorm:"type:varchar(30);not null" json:"type"`
	Amount        float64         `gorm:"type:decimal(14,2);not null" json:"amount"`
	BalanceBefore float64         `gorm:"type:decimal(14,2);not null" json:"balance_before"`
	BalanceAfter  float64         `gorm:"type:decimal(14,2);not null" json:"balance_after"`
	Description   string          `gorm:"type:text" json:"description"`
	ReferenceID   *uuid.UUID      `gorm:"type:uuid;index" json:"reference_id,omitempty"` // ID Ticket Timbangan atau Redeem
	CreatedAt     time.Time       `json:"created_at"`
}

func (l *LedgerTransaction) BeforeCreate(tx *gorm.DB) (err error) {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return
}

type RedeemRequest struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Amount        float64        `gorm:"type:decimal(14,2);not null" json:"amount"`
	PaymentMethod string         `gorm:"type:varchar(50);not null" json:"payment_method"` // "BCA", "GOPAY", "OVO", "DANA"
	AccountName   string         `gorm:"type:varchar(100);not null" json:"account_name"`
	AccountNumber string         `gorm:"type:varchar(50);not null" json:"account_number"`
	Status        RedeemStatus   `gorm:"type:varchar(20);not null;default:'PENDING'" json:"status"`
	Notes         string         `gorm:"type:text" json:"notes"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (r *RedeemRequest) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return
}

type CreateRedeemDTO struct {
	Amount        float64 `json:"amount" validate:"required,min=10000"`
	PaymentMethod string  `json:"payment_method" validate:"required"`
	AccountName   string  `json:"account_name" validate:"required"`
	AccountNumber string  `json:"account_number" validate:"required"`
}

type ProcessRedeemDTO struct {
	Status RedeemStatus `json:"status" validate:"required"`
	Notes  string       `json:"notes"`
}
