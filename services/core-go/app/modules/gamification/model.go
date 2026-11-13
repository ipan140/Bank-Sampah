package gamification

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EntityType string

const (
	TypeRTRW     EntityType = "RT_RW"
	TypeSekolah  EntityType = "SEKOLAH"
	TypeIndividu EntityType = "INDIVIDU"
)

type LeaderboardEntry struct {
	Rank              int        `json:"rank"`
	EntityName        string     `json:"entity_name"`
	EntityType        EntityType `json:"entity_type"`
	TotalWeightKg     float64    `json:"total_weight_kg"`
	TotalTransactions int64      `json:"total_transactions"`
}

type BadgeDefinition struct {
	Code          string  `json:"code"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	MinKgRequired float64 `json:"min_kg_required"`
	IconURL       string  `json:"icon_url"`
}

type UserBadge struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	BadgeCode string         `gorm:"type:varchar(50);not null" json:"badge_code"`
	Title     string         `gorm:"type:varchar(100);not null" json:"title"`
	AwardedAt time.Time      `json:"awarded_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *UserBadge) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	u.AwardedAt = time.Now()
	return
}

type UserLevelProgress struct {
	LevelName       string   `json:"level_name"` // "Pemula Hijau", "Pahlawan Plastik", "Master Sirkular"
	TotalKgRecycled float64  `json:"total_kg_recycled"`
	NextLevelKg     float64  `json:"next_level_kg"`
	ProgressPct     float64  `json:"progress_pct"`
	UnlockedBadges  []string `json:"unlocked_badges"`
}
