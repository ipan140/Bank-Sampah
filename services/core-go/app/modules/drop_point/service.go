package drop_point

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	CreateDropPoint(ctx context.Context, dto *CreateDropPointDTO) (*DropPoint, error)
	ListDropPoints(ctx context.Context, city string) ([]DropPoint, error)
	GetDropPointByID(ctx context.Context, id uuid.UUID) (*DropPoint, error)
	GetSchedules(ctx context.Context) ([]MobilePickupSchedule, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateDropPoint(ctx context.Context, dto *CreateDropPointDTO) (*DropPoint, error) {
	item := &DropPoint{
		Name:           dto.Name,
		Address:        dto.Address,
		RT:             dto.RT,
		RW:             dto.RW,
		District:       dto.District,
		City:           dto.City,
		Latitude:       dto.Latitude,
		Longitude:      dto.Longitude,
		OperatingHours: dto.OperatingHours,
		ContactPerson:  dto.ContactPerson,
		ContactPhone:   dto.ContactPhone,
		IsActive:       true,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *service) ListDropPoints(ctx context.Context, city string) ([]DropPoint, error) {
	return s.repo.FindAll(ctx, city)
}

func (s *service) GetDropPointByID(ctx context.Context, id uuid.UUID) (*DropPoint, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) GetSchedules(ctx context.Context) ([]MobilePickupSchedule, error) {
	return s.repo.GetSchedules(ctx)
}
