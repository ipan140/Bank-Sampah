package weighbridge

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, ticket *ScaleTicket) error
	FindByID(ctx context.Context, id uuid.UUID) (*ScaleTicket, error)
	FindAll(ctx context.Context, filter ScaleTicketFilter) ([]ScaleTicket, int64, error)
	MarkSettled(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, ticket *ScaleTicket) error {
	return r.db.WithContext(ctx).Create(ticket).Error
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*ScaleTicket, error) {
	var ticket ScaleTicket
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&ticket).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *repository) FindAll(ctx context.Context, filter ScaleTicketFilter) ([]ScaleTicket, int64, error) {
	var tickets []ScaleTicket
	var total int64

	query := r.db.WithContext(ctx).Model(&ScaleTicket{})
	if filter.NasabahID != "" {
		query = query.Where("nasabah_id = ?", filter.NasabahID)
	}
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	if filter.DeviceID != "" {
		query = query.Where("device_id = ?", filter.DeviceID)
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

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&tickets).Error
	return tickets, total, err
}

func (r *repository) MarkSettled(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&ScaleTicket{}).Where("id = ?", id).Update("is_settled", true).Error
}
