package user

import (
	"context"
	"errors"

	"bank-sampah/core-go/app/modules/auth"
	"github.com/google/uuid"
)

type Service interface {
	ListUsers(ctx context.Context, filter UserFilter) ([]UserDetail, int64, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*UserDetail, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, req *UpdateProfileRequest) (*UserDetail, error)
	AssignRole(ctx context.Context, id uuid.UUID, role auth.Role) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ListUsers(ctx context.Context, filter UserFilter) ([]UserDetail, int64, error) {
	users, total, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	details := make([]UserDetail, len(users))
	for i, u := range users {
		details[i] = UserDetail{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Phone:     u.Phone,
			Role:      u.Role,
			Address:   u.Address,
			CreatedAt: u.CreatedAt,
		}
	}
	return details, total, nil
}

func (s *service) GetUserByID(ctx context.Context, id uuid.UUID) (*UserDetail, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &UserDetail{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      u.Role,
		Address:   u.Address,
		CreatedAt: u.CreatedAt,
	}, nil
}

func (s *service) UpdateProfile(ctx context.Context, id uuid.UUID, req *UpdateProfileRequest) (*UserDetail, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.Name != "" {
		u.Name = req.Name
	}
	if req.Phone != "" {
		u.Phone = req.Phone
	}
	if req.Address != "" {
		u.Address = req.Address
	}

	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}

	return &UserDetail{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      u.Role,
		Address:   u.Address,
		CreatedAt: u.CreatedAt,
	}, nil
}

func (s *service) AssignRole(ctx context.Context, id uuid.UUID, role auth.Role) error {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}

	u.Role = role
	return s.repo.Update(ctx, u)
}
