package gamification

import (
	"context"

	"github.com/google/uuid"
)

var StandardBadges = []BadgeDefinition{
	{Code: "ECO_BRONZE", Title: "Bronze Recycler", Description: "Setor sampah minimal 10 kg", MinKgRequired: 10.0, IconURL: "/badges/bronze.png"},
	{Code: "ECO_SILVER", Title: "Silver Eco-Hero", Description: "Setor sampah minimal 50 kg", MinKgRequired: 50.0, IconURL: "/badges/silver.png"},
	{Code: "ECO_GOLD", Title: "Gold Climate Champion", Description: "Setor sampah minimal 100 kg", MinKgRequired: 100.0, IconURL: "/badges/gold.png"},
	{Code: "ECO_PLATINUM", Title: "Master Sirkular", Description: "Setor sampah minimal 500 kg", MinKgRequired: 500.0, IconURL: "/badges/platinum.png"},
}

type Service interface {
	GetLeaderboard(ctx context.Context, entityType EntityType, limit int) ([]LeaderboardEntry, error)
	GetUserLevelAndBadges(ctx context.Context, userID uuid.UUID) (*UserLevelProgress, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetLeaderboard(ctx context.Context, entityType EntityType, limit int) ([]LeaderboardEntry, error) {
	if entityType == "" {
		entityType = TypeIndividu
	}
	return s.repo.GetLeaderboard(ctx, entityType, limit)
}

func (s *service) GetUserLevelAndBadges(ctx context.Context, userID uuid.UUID) (*UserLevelProgress, error) {
	totalKg, err := s.repo.GetUserRecycledWeight(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Auto-award badges if threshold reached
	var unlocked []string
	for _, b := range StandardBadges {
		if totalKg >= b.MinKgRequired {
			_ = s.repo.AwardBadge(ctx, userID, b.Code, b.Title)
			unlocked = append(unlocked, b.Title)
		}
	}

	// Calculate user level & next tier
	levelName := "Warga Peduli (Pemula)"
	nextLevelKg := 10.0
	progressPct := (totalKg / nextLevelKg) * 100

	if totalKg >= 500.0 {
		levelName = "Master Sirkular (Level Maksimal)"
		nextLevelKg = 500.0
		progressPct = 100.0
	} else if totalKg >= 100.0 {
		levelName = "Gold Climate Champion"
		nextLevelKg = 500.0
		progressPct = ((totalKg - 100.0) / 400.0) * 100
	} else if totalKg >= 50.0 {
		levelName = "Silver Eco-Hero"
		nextLevelKg = 100.0
		progressPct = ((totalKg - 50.0) / 50.0) * 100
	} else if totalKg >= 10.0 {
		levelName = "Bronze Recycler"
		nextLevelKg = 50.0
		progressPct = ((totalKg - 10.0) / 40.0) * 100
	}

	if progressPct > 100.0 {
		progressPct = 100.0
	}

	return &UserLevelProgress{
		LevelName:       levelName,
		TotalKgRecycled: totalKg,
		NextLevelKg:     nextLevelKg,
		ProgressPct:     progressPct,
		UnlockedBadges:  unlocked,
	}, nil
}
