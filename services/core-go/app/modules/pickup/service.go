package pickup

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	CreatePickup(ctx context.Context, nasabahID uuid.UUID, dto *CreatePickupDTO) (*PickupRequest, error)
	GetPickupByID(ctx context.Context, id uuid.UUID) (*PickupRequest, error)
	ListPickups(ctx context.Context, filter PickupFilter) ([]PickupRequest, int64, error)
	AssignCourier(ctx context.Context, id uuid.UUID, courierID uuid.UUID) (*PickupRequest, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status PickupStatus) (*PickupRequest, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreatePickup(ctx context.Context, nasabahID uuid.UUID, dto *CreatePickupDTO) (*PickupRequest, error) {
	if dto.ScheduledAt.Before(time.Now()) {
		return nil, errors.New("scheduled time must be in the future")
	}

	item := &PickupRequest{
		NasabahID:       nasabahID,
		Status:          StatusPending,
		PickupAddress:   dto.PickupAddress,
		Latitude:        dto.Latitude,
		Longitude:       dto.Longitude,
		ScheduledAt:     dto.ScheduledAt,
		EstimatedWeight: dto.EstimatedWeight,
		Notes:           dto.Notes,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *service) GetPickupByID(ctx context.Context, id uuid.UUID) (*PickupRequest, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) ListPickups(ctx context.Context, filter PickupFilter) ([]PickupRequest, int64, error) {
	return s.repo.FindAll(ctx, filter)
}

func (s *service) AssignCourier(ctx context.Context, id uuid.UUID, courierID uuid.UUID) (*PickupRequest, error) {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("pickup request not found")
	}

	item.CourierID = &courierID
	item.Status = StatusAssigned

	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *service) UpdateStatus(ctx context.Context, id uuid.UUID, status PickupStatus) (*PickupRequest, error) {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("pickup request not found")
	}

	item.Status = status
	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}
