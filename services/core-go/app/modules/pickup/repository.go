package pickup

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, item *PickupRequest) error
	FindByID(ctx context.Context, id uuid.UUID) (*PickupRequest, error)
	FindAll(ctx context.Context, filter PickupFilter) ([]PickupRequest, int64, error)
	Update(ctx context.Context, item *PickupRequest) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, item *PickupRequest) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*PickupRequest, error) {
	var item PickupRequest
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *repository) FindAll(ctx context.Context, filter PickupFilter) ([]PickupRequest, int64, error) {
	var items []PickupRequest
	var total int64

	query := r.db.WithContext(ctx).Model(&PickupRequest{})
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.NasabahID != "" {
		query = query.Where("nasabah_id = ?", filter.NasabahID)
	}
	if filter.CourierID != "" {
		query = query.Where("courier_id = ?", filter.CourierID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Order("scheduled_at ASC").Limit(limit).Offset(offset).Find(&items).Error
	return items, total, err
}

func (r *repository) Update(ctx context.Context, item *PickupRequest) error {
	return r.db.WithContext(ctx).Save(item).Error
}
