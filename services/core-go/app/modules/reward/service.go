package reward

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Service interface {
	CreateReward(ctx context.Context, dto *CreateRewardDTO) (*RewardItem, error)
	ListRewards(ctx context.Context, activeOnly bool) ([]RewardItem, error)
	ClaimReward(ctx context.Context, userID uuid.UUID, rewardID uuid.UUID) (*RewardClaim, error)
	GetUserClaims(ctx context.Context, userID uuid.UUID) ([]RewardClaim, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateReward(ctx context.Context, dto *CreateRewardDTO) (*RewardItem, error) {
	if dto.PointsRequired <= 0 {
		return nil, errors.New("points required must be greater than zero")
	}

	item := &RewardItem{
		Title:          dto.Title,
		Description:    dto.Description,
		PointsRequired: dto.PointsRequired,
		Stock:          dto.Stock,
		ImageURL:       dto.ImageURL,
		IsActive:       true,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *service) ListRewards(ctx context.Context, activeOnly bool) ([]RewardItem, error) {
	return s.repo.FindAll(ctx, activeOnly)
}

func (s *service) ClaimReward(ctx context.Context, userID uuid.UUID, rewardID uuid.UUID) (*RewardClaim, error) {
	return s.repo.ClaimReward(ctx, userID, rewardID)
}

func (s *service) GetUserClaims(ctx context.Context, userID uuid.UUID) ([]RewardClaim, error) {
	return s.repo.GetUserClaims(ctx, userID)
}
