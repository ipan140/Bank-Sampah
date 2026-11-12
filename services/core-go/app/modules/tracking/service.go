package tracking

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	UpdateLocation(ctx context.Context, courierID uuid.UUID, dto *UpdateLocationDTO) error
	GetCourierLiveLocation(ctx context.Context, courierID uuid.UUID) (*NearbyCourier, error)
	FindNearbyCouriers(ctx context.Context, lat, lng float64, radiusKm float64) ([]NearbyCourier, error)
	GetLocationHistory(ctx context.Context, courierID uuid.UUID, limit int) ([]CourierLocation, error)
}

type service struct {
	repo      Repository
	redisRepo RedisRepository
}

func NewService(repo Repository, redisRepo RedisRepository) Service {
	return &service{
		repo:      repo,
		redisRepo: redisRepo,
	}
}

func (s *service) UpdateLocation(ctx context.Context, courierID uuid.UUID, dto *UpdateLocationDTO) error {
	// 1. Update in Redis for real-time low-latency lookups
	_ = s.redisRepo.UpdateCourierLocation(ctx, courierID, dto.Latitude, dto.Longitude)

	// 2. Persist to Postgres for historical breadcrumbs
	loc := &CourierLocation{
		CourierID: courierID,
		Latitude:  dto.Latitude,
		Longitude: dto.Longitude,
		Heading:   dto.Heading,
		Speed:     dto.Speed,
	}
	return s.repo.LogLocation(ctx, loc)
}

func (s *service) GetCourierLiveLocation(ctx context.Context, courierID uuid.UUID) (*NearbyCourier, error) {
	return s.redisRepo.GetCourierLocation(ctx, courierID)
}

func (s *service) FindNearbyCouriers(ctx context.Context, lat, lng float64, radiusKm float64) ([]NearbyCourier, error) {
	if radiusKm <= 0 {
		radiusKm = 10.0 // 10km default radius
	}
	return s.redisRepo.FindNearbyCouriers(ctx, lat, lng, radiusKm)
}

func (s *service) GetLocationHistory(ctx context.Context, courierID uuid.UUID, limit int) ([]CourierLocation, error) {
	return s.repo.GetHistory(ctx, courierID, limit)
}
