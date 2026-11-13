package drop_point

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, item *DropPoint) error
	FindAll(ctx context.Context, city string) ([]DropPoint, error)
	FindByID(ctx context.Context, id uuid.UUID) (*DropPoint, error)
	GetSchedules(ctx context.Context) ([]MobilePickupSchedule, error)
	CreateSchedule(ctx context.Context, s *MobilePickupSchedule) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, item *DropPoint) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *repository) FindAll(ctx context.Context, city string) ([]DropPoint, error) {
	var list []DropPoint
	query := r.db.WithContext(ctx).Where("is_active = ?", true)
	if city != "" {
		query = query.Where("city ILIKE ?", "%"+city+"%")
	}
	err := query.Order("city ASC, name ASC").Find(&list).Error
	return list, err
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*DropPoint, error) {
	var item DropPoint
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *repository) GetSchedules(ctx context.Context) ([]MobilePickupSchedule, error) {
	var list []MobilePickupSchedule
	err := r.db.WithContext(ctx).Order("day_of_week ASC, start_time ASC").Find(&list).Error
	return list, err
}

func (r *repository) CreateSchedule(ctx context.Context, s *MobilePickupSchedule) error {
	return r.db.WithContext(ctx).Create(s).Error
}
