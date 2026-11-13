package ppob

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"bank-sampah/core-go/app/modules/wallet"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	PayBillWithWallet(ctx context.Context, userID uuid.UUID, bill *PPOBTransaction) error
	GetOrCreateGoldAccount(ctx context.Context, userID uuid.UUID) (*GoldSavingsAccount, error)
	ConvertBalanceToGold(ctx context.Context, userID uuid.UUID, amountIDR float64, goldRatePerGram float64) (*GoldSavingsAccount, error)
	DonateFromWallet(ctx context.Context, donation *WasteDonation) error
	GetPPOBHistory(ctx context.Context, userID uuid.UUID, limit int) ([]PPOBTransaction, error)
	GetDonationHistory(ctx context.Context, userID uuid.UUID, limit int) ([]WasteDonation, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) PayBillWithWallet(ctx context.Context, userID uuid.UUID, bill *PPOBTransaction) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var userWallet wallet.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&userWallet).Error; err != nil {
			return errors.New("wallet not found")
		}

		if userWallet.Balance < bill.TotalDeducted {
			return errors.New("insufficient wallet balance to pay bill")
		}

		userWallet.Balance -= bill.TotalDeducted
		if err := tx.Save(&userWallet).Error; err != nil {
			return err
		}

		// Simulation SN or PLN Token
		if bill.Type == BillPLNToken {
			bill.SerialNumberToken = fmt.Sprintf("%04d-%04d-%04d-%04d-%04d", rand.Intn(9999), rand.Intn(9999), rand.Intn(9999), rand.Intn(9999), rand.Intn(9999))
		} else {
			bill.SerialNumberToken = fmt.Sprintf("SN%d%04d", time.Now().Unix(), rand.Intn(9999))
		}
		bill.Status = "SUCCESS"

		if err := tx.Create(bill).Error; err != nil {
			return err
		}

		// Record ledger
		ledger := wallet.LedgerTransaction{
			WalletID:      userWallet.ID,
			UserID:        userID,
			Type:          wallet.TxDebitRedeem,
			Amount:        bill.TotalDeducted,
			BalanceBefore: userWallet.Balance + bill.TotalDeducted,
			BalanceAfter:  userWallet.Balance,
			Description:   fmt.Sprintf("Pembayaran PPOB %s (%s)", bill.Type, bill.TargetNumber),
			ReferenceID:   &bill.ID,
			CreatedAt:     time.Now(),
		}
		return tx.Create(&ledger).Error
	})
}

func (r *repository) GetOrCreateGoldAccount(ctx context.Context, userID uuid.UUID) (*GoldSavingsAccount, error) {
	var account GoldSavingsAccount
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&account).Error
	if err == nil {
		return &account, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		account = GoldSavingsAccount{
			ID:               uuid.New(),
			UserID:           userID,
			BalanceGrams:     0,
			TotalInvestedIDR: 0,
			LastUpdated:      time.Now(),
		}
		if err := r.db.WithContext(ctx).Create(&account).Error; err != nil {
			return nil, err
		}
		return &account, nil
	}
	return nil, err
}

func (r *repository) ConvertBalanceToGold(ctx context.Context, userID uuid.UUID, amountIDR float64, goldRatePerGram float64) (*GoldSavingsAccount, error) {
	var goldAcc *GoldSavingsAccount

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var userWallet wallet.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&userWallet).Error; err != nil {
			return errors.New("wallet not found")
		}

		if userWallet.Balance < amountIDR {
			return errors.New("insufficient wallet balance to buy gold")
		}

		userWallet.Balance -= amountIDR
		if err := tx.Save(&userWallet).Error; err != nil {
			return err
		}

		var acc GoldSavingsAccount
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&acc).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				acc = GoldSavingsAccount{
					ID:               uuid.New(),
					UserID:           userID,
					BalanceGrams:     0,
					TotalInvestedIDR: 0,
				}
				if err := tx.Create(&acc).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		gramsPurchased := amountIDR / goldRatePerGram
		acc.BalanceGrams += gramsPurchased
		acc.TotalInvestedIDR += amountIDR
		acc.LastUpdated = time.Now()

		if err := tx.Save(&acc).Error; err != nil {
			return err
		}

		goldAcc = &acc

		ledger := wallet.LedgerTransaction{
			WalletID:      userWallet.ID,
			UserID:        userID,
			Type:          wallet.TxDebitRedeem,
			Amount:        amountIDR,
			BalanceBefore: userWallet.Balance + amountIDR,
			BalanceAfter:  userWallet.Balance,
			Description:   fmt.Sprintf("Konversi Saldo ke Tabungan Emas (%.4f gram)", gramsPurchased),
			ReferenceID:   &acc.ID,
			CreatedAt:     time.Now(),
		}
		return tx.Create(&ledger).Error
	})

	return goldAcc, err
}

func (r *repository) DonateFromWallet(ctx context.Context, donation *WasteDonation) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var userWallet wallet.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", donation.UserID).First(&userWallet).Error; err != nil {
			return errors.New("wallet not found")
		}

		if userWallet.Balance < donation.Amount {
			return errors.New("insufficient wallet balance for donation")
		}

		userWallet.Balance -= donation.Amount
		if err := tx.Save(&userWallet).Error; err != nil {
			return err
		}

		if err := tx.Create(donation).Error; err != nil {
			return err
		}

		ledger := wallet.LedgerTransaction{
			WalletID:      userWallet.ID,
			UserID:        donation.UserID,
			Type:          wallet.TxDebitRedeem,
			Amount:        donation.Amount,
			BalanceBefore: userWallet.Balance + donation.Amount,
			BalanceAfter:  userWallet.Balance,
			Description:   fmt.Sprintf("Sedekah/Donasi Sampah ke %s", donation.CharityInstitutionName),
			ReferenceID:   &donation.ID,
			CreatedAt:     time.Now(),
		}
		return tx.Create(&ledger).Error
	})
}

func (r *repository) GetPPOBHistory(ctx context.Context, userID uuid.UUID, limit int) ([]PPOBTransaction, error) {
	var list []PPOBTransaction
	if limit <= 0 || limit > 100 {
		limit = 30
	}
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Find(&list).Error
	return list, err
}

func (r *repository) GetDonationHistory(ctx context.Context, userID uuid.UUID, limit int) ([]WasteDonation, error) {
	var list []WasteDonation
	if limit <= 0 || limit > 100 {
		limit = 30
	}
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Find(&list).Error
	return list, err
}
