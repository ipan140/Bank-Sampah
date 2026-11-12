package wallet

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	GetOrCreateWallet(ctx context.Context, userID uuid.UUID) (*Wallet, error)
	GetWalletByUserID(ctx context.Context, userID uuid.UUID) (*Wallet, error)
	CreditBalance(ctx context.Context, userID uuid.UUID, amount float64, desc string, refID *uuid.UUID) (*LedgerTransaction, error)
	DebitBalance(ctx context.Context, userID uuid.UUID, amount float64, desc string, refID *uuid.UUID) (*LedgerTransaction, error)
	GetTransactions(ctx context.Context, userID uuid.UUID, page, limit int) ([]LedgerTransaction, int64, error)
	CreateRedeem(ctx context.Context, redeem *RedeemRequest) error
	GetRedeemByID(ctx context.Context, id uuid.UUID) (*RedeemRequest, error)
	ListRedeems(ctx context.Context, page, limit int, status string) ([]RedeemRequest, int64, error)
	UpdateRedeem(ctx context.Context, redeem *RedeemRequest) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetOrCreateWallet(ctx context.Context, userID uuid.UUID) (*Wallet, error) {
	var wallet Wallet
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&wallet).Error
	if err == nil {
		return &wallet, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		wallet = Wallet{
			ID:      uuid.New(),
			UserID:  userID,
			Balance: 0,
			Points:  0,
		}
		if err := r.db.WithContext(ctx).Create(&wallet).Error; err != nil {
			return nil, err
		}
		return &wallet, nil
	}
	return nil, err
}

func (r *repository) GetWalletByUserID(ctx context.Context, userID uuid.UUID) (*Wallet, error) {
	return r.GetOrCreateWallet(ctx, userID)
}

func (r *repository) CreditBalance(ctx context.Context, userID uuid.UUID, amount float64, desc string, refID *uuid.UUID) (*LedgerTransaction, error) {
	var txRecord *LedgerTransaction

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var wallet Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				wallet = Wallet{ID: uuid.New(), UserID: userID, Balance: 0, Points: 0}
				if err := tx.Create(&wallet).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		balanceBefore := wallet.Balance
		balanceAfter := balanceBefore + amount
		wallet.Balance = balanceAfter

		// Points reward calculation: 10 points per Rp 1.000
		pointsGained := int64(amount / 1000) * 10
		wallet.Points += pointsGained

		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		txRecord = &LedgerTransaction{
			WalletID:      wallet.ID,
			UserID:        userID,
			Type:          TxCreditDeposit,
			Amount:        amount,
			BalanceBefore: balanceBefore,
			BalanceAfter:  balanceAfter,
			Description:   desc,
			ReferenceID:   refID,
		}

		return tx.Create(txRecord).Error
	})

	return txRecord, err
}

func (r *repository) DebitBalance(ctx context.Context, userID uuid.UUID, amount float64, desc string, refID *uuid.UUID) (*LedgerTransaction, error) {
	var txRecord *LedgerTransaction

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var wallet Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			return err
		}

		if wallet.Balance < amount {
			return errors.New("insufficient balance")
		}

		balanceBefore := wallet.Balance
		balanceAfter := balanceBefore - amount
		wallet.Balance = balanceAfter

		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		txRecord = &LedgerTransaction{
			WalletID:      wallet.ID,
			UserID:        userID,
			Type:          TxDebitRedeem,
			Amount:        amount,
			BalanceBefore: balanceBefore,
			BalanceAfter:  balanceAfter,
			Description:   desc,
			ReferenceID:   refID,
		}

		return tx.Create(txRecord).Error
	})

	return txRecord, err
}

func (r *repository) GetTransactions(ctx context.Context, userID uuid.UUID, page, limit int) ([]LedgerTransaction, int64, error) {
	var items []LedgerTransaction
	var total int64

	query := r.db.WithContext(ctx).Model(&LedgerTransaction{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&items).Error
	return items, total, err
}

func (r *repository) CreateRedeem(ctx context.Context, redeem *RedeemRequest) error {
	return r.db.WithContext(ctx).Create(redeem).Error
}

func (r *repository) GetRedeemByID(ctx context.Context, id uuid.UUID) (*RedeemRequest, error) {
	var redeem RedeemRequest
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&redeem).Error
	if err != nil {
		return nil, err
	}
	return &redeem, nil
}

func (r *repository) ListRedeems(ctx context.Context, page, limit int, status string) ([]RedeemRequest, int64, error) {
	var items []RedeemRequest
	var total int64

	query := r.db.WithContext(ctx).Model(&RedeemRequest{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&items).Error
	return items, total, err
}

func (r *repository) UpdateRedeem(ctx context.Context, redeem *RedeemRequest) error {
	return r.db.WithContext(ctx).Save(redeem).Error
}
