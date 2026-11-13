package member_card

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	GetMyCard(ctx context.Context, userID uuid.UUID) (*MemberCard, error)
	ScanMemberQR(ctx context.Context, qrToken string) (*MemberInfoResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetMyCard(ctx context.Context, userID uuid.UUID) (*MemberCard, error) {
	return s.repo.GetOrCreateCard(ctx, userID)
}

func (s *service) ScanMemberQR(ctx context.Context, qrToken string) (*MemberInfoResponse, error) {
	return s.repo.GetMemberInfo(ctx, qrToken)
}
