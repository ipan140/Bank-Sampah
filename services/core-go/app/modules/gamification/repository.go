package gamification

import (
	"context"

	"bank-sampah/core-go/app/modules/weighbridge"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetLeaderboard(ctx context.Context, entityType EntityType, limit int) ([]LeaderboardEntry, error)
	GetUserRecycledWeight(ctx context.Context, userID uuid.UUID) (float64, error)
	GetUserBadges(ctx context.Context, userID uuid.UUID) ([]UserBadge, error)
	AwardBadge(ctx context.Context, userID uuid.UUID, code, title string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetLeaderboard(ctx context.Context, entityType EntityType, limit int) ([]LeaderboardEntry, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	var results []LeaderboardEntry
	type QueryResult struct {
		Name        string
		TotalKg     float64
		TotalTrx    int64
	}
	var qResults []QueryResult

	// Grouping based on users and their weighbridge net weight
	r.db.WithContext(ctx).Table("scale_tickets").
		Select("users.name as name, SUM(scale_tickets.net_weight) as total_kg, COUNT(scale_tickets.id) as total_trx").
		Joins("JOIN users ON users.id = scale_tickets.nasabah_id").
		Where("scale_tickets.deleted_at IS NULL").
		Group("users.name").
		Order("total_kg DESC").
		Limit(limit).
		Scan(&qResults)

	for i, q := range qResults {
		results = append(results, LeaderboardEntry{
			Rank:              i + 1,
			EntityName:        q.Name,
			EntityType:        entityType,
			TotalWeightKg:     q.TotalKg,
			TotalTransactions: q.TotalTrx,
		})
	}

	return results, nil
}

func (r *repository) GetUserRecycledWeight(ctx context.Context, userID uuid.UUID) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).Model(&weighbridge.ScaleTicket{}).
		Where("nasabah_id = ?", userID).
		Select("COALESCE(SUM(net_weight), 0)").
		Scan(&total).Error
	return total, err
}

func (r *repository) GetUserBadges(ctx context.Context, userID uuid.UUID) ([]UserBadge, error) {
	var badges []UserBadge
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("awarded_at ASC").Find(&badges).Error
	return badges, err
}

func (r *repository) AwardBadge(ctx context.Context, userID uuid.UUID, code, title string) error {
	var existing UserBadge
	err := r.db.WithContext(ctx).Where("user_id = ? AND badge_code = ?", userID, code).First(&existing).Error
	if err == nil {
		return nil // already awarded
	}

	badge := UserBadge{
		UserID:    userID,
		BadgeCode: code,
		Title:     title,
	}
	return r.db.WithContext(ctx).Create(&badge).Error
}
