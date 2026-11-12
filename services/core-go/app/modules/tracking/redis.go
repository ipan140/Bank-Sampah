package tracking

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const CourierGeoKey = "couriers:geo"

type RedisRepository interface {
	UpdateCourierLocation(ctx context.Context, courierID uuid.UUID, lat, lng float64) error
	GetCourierLocation(ctx context.Context, courierID uuid.UUID) (*NearbyCourier, error)
	FindNearbyCouriers(ctx context.Context, lat, lng float64, radiusKm float64) ([]NearbyCourier, error)
}

type redisRepo struct {
	rdb *redis.Client
}

func NewRedisRepository(rdb *redis.Client) RedisRepository {
	return &redisRepo{rdb: rdb}
}

func (r *redisRepo) UpdateCourierLocation(ctx context.Context, courierID uuid.UUID, lat, lng float64) error {
	if r.rdb == nil {
		return nil
	}
	return r.rdb.GeoAdd(ctx, CourierGeoKey, &redis.GeoLocation{
		Name:      courierID.String(),
		Longitude: lng,
		Latitude:  lat,
	}).Err()
}

func (r *redisRepo) GetCourierLocation(ctx context.Context, courierID uuid.UUID) (*NearbyCourier, error) {
	if r.rdb == nil {
		return nil, fmt.Errorf("redis unavailable")
	}

	positions, err := r.rdb.GeoPos(ctx, CourierGeoKey, courierID.String()).Result()
	if err != nil {
		return nil, err
	}
	if len(positions) == 0 || positions[0] == nil {
		return nil, fmt.Errorf("courier location not found")
	}

	pos := positions[0]
	return &NearbyCourier{
		CourierID: courierID,
		Latitude:  pos.Latitude,
		Longitude: pos.Longitude,
	}, nil
}

func (r *redisRepo) FindNearbyCouriers(ctx context.Context, lat, lng float64, radiusKm float64) ([]NearbyCourier, error) {
	if r.rdb == nil {
		return []NearbyCourier{}, nil
	}

	query := &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Longitude:  lng,
			Latitude:   lat,
			Radius:     radiusKm,
			RadiusUnit: "km",
			Sort:       "ASC",
		},
		WithCoord: true,
		WithDist:  true,
	}

	locations, err := r.rdb.GeoSearchLocation(ctx, CourierGeoKey, query).Result()
	if err != nil {
		return nil, err
	}

	result := make([]NearbyCourier, 0, len(locations))
	for _, loc := range locations {
		cid, err := uuid.Parse(loc.Name)
		if err != nil {
			continue
		}
		dist, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", loc.Dist*1000), 64)
		result = append(result, NearbyCourier{
			CourierID: cid,
			Latitude:  loc.Latitude,
			Longitude: loc.Longitude,
			DistanceM: dist,
		})
	}

	return result, nil
}
