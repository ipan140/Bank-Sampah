package reward

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
	Create(ctx context.Context, item *RewardItem) error
	FindByID(ctx context.Context, id uuid.UUID) (*RewardItem, error)
	FindAll(ctx context.Context, activeOnly bool) ([]RewardItem, error)
	ClaimReward(ctx context.Context, userID uuid.UUID, rewardID uuid.UUID) (*RewardClaim, error)
	GetUserClaims(ctx context.Context, userID uuid.UUID) ([]RewardClaim, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, item *RewardItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*RewardItem, error) {
	var item RewardItem
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *repository) FindAll(ctx context.Context, activeOnly bool) ([]RewardItem, error) {
	var items []RewardItem
	query := r.db.WithContext(ctx)
	if activeOnly {
		query = query.Where("is_active = ? AND stock > 0", true)
	}
	err := query.Order("points_required ASC").Find(&items).Error
	return items, err
}

func (r *repository) ClaimReward(ctx context.Context, userID uuid.UUID, rewardID uuid.UUID) (*RewardClaim, error) {
	var claim *RewardClaim

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var item RewardItem
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", rewardID).First(&item).Error; err != nil {
			return err
		}

		if !item.IsActive || item.Stock <= 0 {
			return errors.New("reward out of stock or inactive")
		}

		var userWallet wallet.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&userWallet).Error; err != nil {
			return errors.New("user wallet not found")
		}

		if userWallet.Points < item.PointsRequired {
			return errors.New("insufficient points to claim reward")
		}

		// Deduct points and stock
		userWallet.Points -= item.PointsRequired
		item.Stock -= 1

		if err := tx.Save(&userWallet).Error; err != nil {
			return err
		}
		if err := tx.Save(&item).Error; err != nil {
			return err
		}

		code := fmt.Sprintf("RWD-%d-%04d", time.Now().Unix(), rand.Intn(9999))
		claim = &RewardClaim{
			RewardItemID: item.ID,
			UserID:       userID,
			PointsSpent:  item.PointsRequired,
			ClaimCode:    code,
			Status:       "CLAIMED",
		}

		return tx.Create(claim).Error
	})

	return claim, err
}

func (r *repository) GetUserClaims(ctx context.Context, userID uuid.UUID) ([]RewardClaim, error) {
	var claims []RewardClaim
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&claims).Error
	return claims, err
}
