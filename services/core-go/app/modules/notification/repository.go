package notification

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, notif *Notification) error
	GetByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]Notification, error)
	MarkAsRead(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, notif *Notification) error {
	return r.db.WithContext(ctx).Create(notif).Error
}

func (r *repository) GetByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]Notification, error) {
	var notifs []Notification
	if limit <= 0 || limit > 100 {
		limit = 30
	}
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Find(&notifs).Error
	return notifs, err
}

func (r *repository) MarkAsRead(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("status", StatusRead).Error
}
