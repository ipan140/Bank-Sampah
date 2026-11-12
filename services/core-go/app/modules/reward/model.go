package reward

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RewardItem struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Title          string         `gorm:"type:varchar(150);not null" json:"title"`
	Description    string         `gorm:"type:text" json:"description"`
	PointsRequired int64          `gorm:"not null" json:"points_required"`
	Stock          int            `gorm:"not null;default:0" json:"stock"`
	ImageURL       string         `gorm:"type:varchar(255)" json:"image_url"`
	IsActive       bool           `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (r *RewardItem) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return
}

type RewardClaim struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	RewardItemID uuid.UUID      `gorm:"type:uuid;not null;index" json:"reward_item_id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	PointsSpent  int64          `gorm:"not null" json:"points_spent"`
	ClaimCode    string         `gorm:"type:varchar(50);not null" json:"claim_code"`
	Status       string         `gorm:"type:varchar(20);default:'CLAIMED'" json:"status"` // CLAIMED, USED, EXPIRED
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c *RewardClaim) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}

type CreateRewardDTO struct {
	Title          string `json:"title" validate:"required"`
	Description    string `json:"description"`
	PointsRequired int64  `json:"points_required" validate:"required,min=1"`
	Stock          int    `json:"stock" validate:"required,min=1"`
	ImageURL       string `json:"image_url"`
}
