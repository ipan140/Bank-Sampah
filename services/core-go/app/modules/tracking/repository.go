package tracking

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	LogLocation(ctx context.Context, loc *CourierLocation) error
	GetHistory(ctx context.Context, courierID uuid.UUID, limit int) ([]CourierLocation, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) LogLocation(ctx context.Context, loc *CourierLocation) error {
	return r.db.WithContext(ctx).Create(loc).Error
}

func (r *repository) GetHistory(ctx context.Context, courierID uuid.UUID, limit int) ([]CourierLocation, error) {
	var history []CourierLocation
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	err := r.db.WithContext(ctx).
		Where("courier_id = ?", courierID).
		Order("created_at DESC").
		Limit(limit).
		Find(&history).Error
	return history, err
}
