package csr

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Service interface {
	CreateProgram(ctx context.Context, dto *CreateProgramDTO) (*CsrProgram, error)
	GetProgramByID(ctx context.Context, id uuid.UUID) (*CsrProgram, error)
	ListPrograms(ctx context.Context) ([]CsrProgram, error)
	GetImpactMetrics(ctx context.Context) (*ImpactSummary, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateProgram(ctx context.Context, dto *CreateProgramDTO) (*CsrProgram, error) {
	if dto.FundAmount < 100000 {
		return nil, errors.New("fund amount must be at least Rp 100.000")
	}

	program := &CsrProgram{
		PartnerName: dto.PartnerName,
		Title:       dto.Title,
		Description: dto.Description,
		FundAmount:  dto.FundAmount,
		TargetKg:    dto.TargetKg,
		CollectedKg: 0,
		StartDate:   dto.StartDate,
		EndDate:     dto.EndDate,
	}

	if err := s.repo.CreateProgram(ctx, program); err != nil {
		return nil, err
	}

	return program, nil
}

func (s *service) GetProgramByID(ctx context.Context, id uuid.UUID) (*CsrProgram, error) {
	return s.repo.FindProgramByID(ctx, id)
}

func (s *service) ListPrograms(ctx context.Context) ([]CsrProgram, error) {
	return s.repo.ListPrograms(ctx)
}

func (s *service) GetImpactMetrics(ctx context.Context) (*ImpactSummary, error) {
	return s.repo.CalculateImpact(ctx)
}
